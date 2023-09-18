package turnkey

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type TurnkeyRequest struct {
	Type           string `json:"type"`
	TimestampMs    string `json:"timestampMs"`
	OrganizationId string `json:"organizationId"`
	Parameters     any    `json:"parameters"`
}

func newEthereumPrivateKeyRequest(orgId, keyName, keyTag string) TurnkeyRequest {
	activity := "ACTIVITY_TYPE_CREATE_PRIVATE_KEYS_V2"
	timestamp := time.Now().UnixMilli()
	keyParam := createPrivateKeyParam(keyName, keyTag, SECP256K1, ETHEREUM)
	return TurnkeyRequest{
		Type:           activity,
		TimestampMs:    strconv.FormatInt(timestamp, 10),
		OrganizationId: orgId,
		Parameters:     keyParam,
	}
}

func newSubOrganizationRequest(orgId, subOrgName string) TurnkeyRequest {
	activity := "ACTIVITY_TYPE_CREATE_SUB_ORGANIZATION_V2"
	// turnkeyUsername := os.Getenv("TURNKEY_USERNAME")\
	turnkeyUsername := "x"
	timestamp := time.Now().UnixMilli()
	subOrgParams := newSubOrganizationIntentV3(subOrgName, turnkeyUsername)
	return TurnkeyRequest{
		Type:           activity,
		TimestampMs:    strconv.FormatInt(timestamp, 10),
		OrganizationId: orgId,
		Parameters:     subOrgParams,
	}
}

func newSignPayloadRequest(orgId, privateKeyId, payload string) TurnkeyRequest {
	activity := "ACTIVITY_TYPE_SIGN_RAW_PAYLOAD"
	timestamp := time.Now().UnixMilli()
	params := newV1SignPayloadIntent(privateKeyId, payload)

	return TurnkeyRequest{
		Type:           activity,
		TimestampMs:    strconv.FormatInt(timestamp, 10),
		OrganizationId: orgId,
		Parameters:     params,
	}
}

func newPrivateKeyTagRequest(orgId, keyTag string, keyIds []string) TurnkeyRequest {
	activity := "ACTIVITY_TYPE_CREATE_PRIVATE_KEY_TAG"
	timestamp := time.Now().UnixMilli()
	params := newV1CreatePrivateKeyTagIntent(keyTag, keyIds)

	return TurnkeyRequest{
		Type:           activity,
		TimestampMs:    strconv.FormatInt(timestamp, 10),
		OrganizationId: orgId,
		Parameters:     params,
	}
}

type GetActivityRequest struct {
	ActivityId     string `json:"activityId"`
	OrganizationId string `json:"organizationId"`
}

func newActivityPollRequest(orgId, activityId string) GetActivityRequest {
	return GetActivityRequest{
		ActivityId:     activityId,
		OrganizationId: orgId,
	}
}

type RootUserParams struct {
	UserName       string        `json:"userName"`
	ApiKeys        []ApiKeyParam `json:"apiKeys"`
	Authenticators any           `json:"authenticators"`
	UserEmail      string        `json:"userEmail"`
}

func newRootUserParam(username, email string) RootUserParams {
	apiKeyName := os.Getenv("TURNKEY_API_KEY_NAME")
	publicKey := os.Getenv("TURNKEY_KEY")
	apiKey := ApiKeyParam{
		ApiKeyName: apiKeyName,
		PublicKey:  publicKey,
	}

	apiKeys := []ApiKeyParam{apiKey}

	fmt.Println("api keys - ", apiKeys)
	apiKeys = append(apiKeys, newApiKeyParam())
	return RootUserParams{
		UserName:       username,
		UserEmail:      email,
		Authenticators: []any{},
		ApiKeys:        apiKeys,
	}
}

/**
Intents start here
*/

type CreateSubOrganizationIntentV3 struct {
	SubOrganizationName string           `json:"subOrganizationName"`
	RootUsers           []RootUserParams `json:"rootUsers"`
	RootQuorumThreshold int32            `json:"rootQuorumThreshold"`
	// PrivateKeys         any              `json:"privateKeys"`
}

func newSubOrganizationIntentV3(orgName, username string) CreateSubOrganizationIntentV3 {
	rootUsers := []RootUserParams{}
	rootUsers = append(rootUsers, newRootUserParam(username, ""))
	return CreateSubOrganizationIntentV3{
		SubOrganizationName: orgName,
		RootUsers:           rootUsers,
		RootQuorumThreshold: 1,
		// PrivateKeys:         []any{},
	}
}

type V1SignRawPayloadIntent struct {
	PrivateKeyId string `json:"privateKeyId"`
	Payload      string `json:"payload"`
	Encoding     string `json:"encoding"`
	HashFunction string `json:"hashFunction"`
}

func newV1SignPayloadIntent(privateKeyId, payload string) V1SignRawPayloadIntent {
	return V1SignRawPayloadIntent{
		PrivateKeyId: privateKeyId,
		Payload:      payload,
		Encoding:     "PAYLOAD_ENCODING_HEXADECIMAL",
		HashFunction: "HASH_FUNCTION_NO_OP",
	}
}

type V1CreatePrivateKeyTagIntent struct {
	PrivateKeyTagName string   `json:"privateKeyTagName"`
	PrivateKeyIds     []string `json:"privateKeyIds"`
}

func newV1CreatePrivateKeyTagIntent(keyTag string, privateKeyIds []string) V1CreatePrivateKeyTagIntent {
	return V1CreatePrivateKeyTagIntent{
		PrivateKeyTagName: keyTag,
		PrivateKeyIds:     privateKeyIds,
	}
}

/**
Params start here
*/

type ApiKeyParam struct {
	ApiKeyName string `json:"apiKeyName"`
	PublicKey  string `json:"publicKey"`
}

func newApiKeyParam() ApiKeyParam {
	keyName := os.Getenv("TURNKEY_API_NAME")
	publicKey := os.Getenv("TURNKEY_PUBLIC_KEY")
	return ApiKeyParam{
		ApiKeyName: keyName,
		PublicKey:  publicKey,
	}
}

// PrivateKey holds the fields for the privateKey object field sent to turnkey
type PrivateKey struct {
	PrivateKeyName string   `json:"privateKeyName"`
	Curve          string   `json:"curve"`
	PrivateKeyTags []string `json:"privateKeyTags"`
	AddressFormats []string `json:"addressFormats"`
}

type PrivateKeyParameter struct {
	PrivateKeys []PrivateKey `json:"privateKeys"`
}

type CreatePrivateKeyResult struct {
	PrivateKeys []privateKey `json:"privateKeys"`
}

func (r *CreatePrivateKeyResult) getResult() (privateKeyId, address string) {
	pk := r.PrivateKeys[0]
	address = pk.Addresses[0].Address
	privateKeyId = pk.PrivateKeyId
	return
}

type privateKey struct {
	Addresses    []address `json:"addresses"`
	PrivateKeyId string    `json:"privateKeyId"`
}

type address struct {
	Address string `json:"address"`
}

func createPrivateKeyParam(privateKeyName, privateKeyTag string, curveType Curve, format AddressFormat) PrivateKeyParameter {
	privateKey := PrivateKey{
		PrivateKeyName: privateKeyName,
		Curve:          string(curveType),
		PrivateKeyTags: []string{privateKeyTag},
		AddressFormats: []string{string(format)},
	}
	return PrivateKeyParameter{
		PrivateKeys: []PrivateKey{privateKey},
	}
}

type TurnkeySignature struct {
	R string `json:"r"`
	S string `json:"s"`
	V string `json:"v"`
}

func (sig *TurnkeySignature) ParseSignature(mode string) string {
	signature, _ := hexutil.Decode(mode)

	signatureBytes := hexutil.MustDecode(sig.R)
	fmt.Println(len(signatureBytes))
	signatureBytes = append(signatureBytes, hexutil.MustDecode(sig.S)...)
	signatureBytes = append(signatureBytes, hexutil.MustDecode(sig.V)...)
	fmt.Println(len(signatureBytes))
	signatureBytes[64] += 27
	signature = append(signature, signatureBytes...)

	return hexutil.Encode(signature)
}
