package turnkey

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	tk "github.com/tkhq/go-sdk"
	"github.com/tkhq/go-sdk/pkg/api/client"
	"github.com/tkhq/go-sdk/pkg/api/client/who_am_i"
	"github.com/tkhq/go-sdk/pkg/api/models"
	"github.com/tkhq/go-sdk/pkg/apikey"
	"github.com/tkhq/go-sdk/pkg/store/local"
)

const TURNKEY_API_SIGNATURE_SCHEME = "SIGNATURE_SCHEME_TK_API_P256"

type Curve string
type AddressFormat string

const (
	// curve types
	SECP256K1 Curve = "CURVE_SECP256K1"
	ED25519   Curve = "CURVE_ED25519"

	// address formats
	UNCOMPRESSED AddressFormat = "ADDRESS_FORMAT_UNCOMPRESSED"
	COMPRESSED   AddressFormat = "ADDRESS_FORMAT_COMPRESSED"
	ETHEREUM     AddressFormat = "ADDRESS_FORMAT_ETHEREUM"
)

type TurnkeyService struct {
	Host          string
	Protocol      string
	TurnkeyClient *tk.Client
}

func NewTurnKeyService() (*TurnkeyService, error) {
	host := os.Getenv("TURNKEY_HOST")
	client, err := initTurnkeyClient()
	if err != nil {
		log.Println(err, "unable to initialise Turnkey client")
		return nil, err
	}
	return &TurnkeyService{
		Protocol:      "https",
		Host:          host,
		TurnkeyClient: client,
	}, nil
}

func initKeys(keyPath, keysDir string) (*tk.Client, error) {
	store := local.New()
	err := store.SetKeysDirectory(keysDir)
	if err != nil {
		return nil, err
	}

	key, err := store.Load(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load API key")
	}

	return &tk.Client{
		Client:        client.NewHTTPClient(nil),
		Authenticator: &tk.Authenticator{Key: key},
		APIKey:        key,
	}, nil
}
func initTurnkeyClient() (*tk.Client, error) {
	keyPath := os.Getenv("TK_KEYS_NAME")
	keyDir := os.Getenv("TK_KEYS_DIR")

	fmt.Println("path ", keyPath)
	// client, err := tk.New(keyPath)
	client, err := initKeys(keyPath, keyDir)
	if err != nil {
		err = errors.Wrap(err, "tk.New()")
		log.Println(err)
		return nil, err
	}

	fmt.Println("org ", *client.APIKey)
	p := who_am_i.NewPublicAPIServiceGetWhoamiParams().WithBody(&models.V1GetWhoamiRequest{
		OrganizationID: client.DefaultOrganization(),
	})

	resp, err := client.V0().WhoAmi.PublicAPIServiceGetWhoami(p, client.Authenticator)
	if err != nil {
		err = errors.Wrap(err, "PublicAPIServiceGetWhoami()")
		return nil, err
	}

	fmt.Println("turnkey user id - ", *resp.Payload.UserID)
	return client, nil
}

func (tk *TurnkeyService) CreatePrivateKeyTag(orgId, keyTag string) (string, error) {
	var activityId string
	if orgId == "" {
		orgId = *tk.TurnkeyClient.DefaultOrganization()
	}

	privateKeyIds := []string{}
	privateKeyTagRequest := newPrivateKeyTagRequest(orgId, keyTag, privateKeyIds)
	apiKey := tk.TurnkeyClient.APIKey
	path := "/public/v1/submit/create_private_key_tag"

	payload, err := json.Marshal(privateKeyTagRequest)
	if err != nil {
		return "", err
	}

	stamp, err := apikey.Stamp(payload, apiKey)
	if err != nil {
		return "", err
	}

	response, err := tk.makeRequest(path, stamp, payload)
	if err != nil {
		return "", err
	}

	activityId = response.Activity.Id
	return activityId, nil
}

func (tk *TurnkeyService) CreatePrivateKey(orgId, name, tag string) (string, error) {
	log.Println(orgId)

	if orgId == "" {
		orgId = *tk.TurnkeyClient.DefaultOrganization()
	}
	keyRequest := newEthereumPrivateKeyRequest(orgId, name, tag)
	apiKey := tk.TurnkeyClient.APIKey
	path := "/public/v1/submit/create_private_keys"

	payload, err := json.Marshal(keyRequest)
	if err != nil {
		return "", err
	}

	stamp, err := apikey.Stamp(payload, apiKey)
	if err != nil {
		return "", err
	}

	response, err := tk.makeRequest(path, stamp, payload)
	if err != nil {
		return "", err
	}

	return response.Activity.Id, nil
}

func (tk *TurnkeyService) CreateSubOrganization(orgId, subOrgName string) (string, error) {
	if orgId == "" {
		orgId = *tk.TurnkeyClient.DefaultOrganization()
	}
	apiKey := tk.TurnkeyClient.APIKey
	path := "/public/v1/submit/create_sub_organization"
	subOrganizationRequest := newSubOrganizationRequest(orgId, subOrgName)

	payload, err := json.Marshal(subOrganizationRequest)
	if err != nil {
		return "", err
	}

	stamp, err := apikey.Stamp(payload, apiKey)
	if err != nil {
		return "", err
	}

	response, err := tk.makeRequest(path, stamp, payload)
	if err != nil {
		return "", err
	}

	return response.Activity.Id, nil
}

func (tk *TurnkeyService) SignMessage(orgId, privateKeyId, message string) (string, error) {
	if orgId == "" {
		orgId = *tk.TurnkeyClient.DefaultOrganization()
	}
	signatureRequest := newSignPayloadRequest(orgId, privateKeyId, message)
	apiKey := tk.TurnkeyClient.APIKey
	path := "/public/v1/submit/sign_raw_payload"

	payload, err := json.Marshal(signatureRequest)
	if err != nil {
		return "", err
	}

	stamp, err := apikey.Stamp(payload, apiKey)
	if err != nil {
		return "", err
	}

	response, err := tk.makeRequest(path, stamp, payload)
	if err != nil {
		return "", err
	}

	return response.Activity.Id, nil
}

func (tk *TurnkeyService) GetActivity(orgId, activityId string) (map[string]any, error) {
	if orgId == "" {
		orgId = *tk.TurnkeyClient.DefaultOrganization()
	}
	activityRequest := newActivityPollRequest(orgId, activityId)
	apiKey := tk.TurnkeyClient.APIKey
	path := "/public/v1/query/get_activity"

	payload, err := json.Marshal(activityRequest)
	if err != nil {
		return nil, err
	}

	stamp, err := apikey.Stamp(payload, apiKey)
	if err != nil {
		return nil, err
	}

	response, err := tk.makeRequest(path, stamp, payload)
	if err != nil {
		return nil, err
	}

	return response.Activity.Result, nil
}

func (tk *TurnkeyService) makeRequest(path, stamp string, payload []byte) (*TurnkeyResponse, error) {
	response, err := post(context.Background(), tk.Protocol, tk.Host, path, payload, stamp)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		errMsg := &ResponseError{
			Code: response.StatusCode,
			Text: response.Status,
		}
		return nil, errMsg
	}

	if response.StatusCode != http.StatusOK {
		errMsg := &ResponseError{
			Code: response.StatusCode,
			Text: string(responseBodyBytes),
		}
		return nil, errMsg
	}
	var parsedResponse *TurnkeyResponse

	err = json.Unmarshal(responseBodyBytes, &parsedResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return parsedResponse, nil
}

func post(ctx context.Context, protocol string, host string, path string, body []byte, stamp string) (*http.Response, error) {
	url := fmt.Sprintf("%s://%s%s", protocol, host, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "error while creating HTTP POST request")
	}

	req.Header.Set("X-Stamp", stamp)

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error while sending HTTP POST request")
	}

	return response, nil
}

func maybeParseJSON(payload any) any {
	bytes, ok := payload.([]byte)
	if ok {
		var decoded any
		err := json.Unmarshal(bytes, &decoded)
		if err == nil {
			return decoded
		}

	}
	return payload
}
