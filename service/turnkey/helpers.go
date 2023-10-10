package turnkey

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type SignRawPayloadResult struct {
	R string `json:"r"`
	S string `json:"s"`
	V string `json:"v"`
}

func (r *SignRawPayloadResult) GetResult() string {
	return ""
}

type TurnkeyTimestamp struct {
	Seconds string `json:"seconds"`
	Nano    string `json:"nano"`
}

type TurnkeyResponse struct {
	Activity TurnkeyActivity `json:"activity"`
}

type TurnkeyActivity struct {
	Id             string         `json:"id"`
	OrganizationID string         `json:"organizationId"`
	Status         string         `json:"status"`
	Type           string         `json:"type"`
	Intent         any            `json:" intent"`
	Result         map[string]any `json:"result"`
	// Votes          []string         `json:"votes"`
	Fingerprint string           `json:"fingerprint"`
	CanApprove  bool             `json:"canApprove"`
	CanReject   bool             `json:"canReject"`
	CreatedAt   TurnkeyTimestamp `json:"createdAt"`
	UpdatedAt   TurnkeyTimestamp `json:"updated"`
}

// ResponseError is a structured format to display an HTTP error response.
type ResponseError struct {
	Code int    `json:"responseCode"`
	Text string `json:"responseBody"`
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Text)
}

func GetPrivateKeyIdFromResult(result map[string]any) (privateKeyId, address string, err error) {
	var privateKeyResult CreatePrivateKeyResult

	privateKeyI := result["createPrivateKeysResultV2"]

	privateKeyB, err := json.Marshal(privateKeyI)
	if err != nil {
		err = errors.Wrap(err, "marshalling turnkey createPrivateKeysResult interface failed")
		return "", "", err
	}

	err = json.Unmarshal(privateKeyB, &privateKeyResult)
	if err != nil {
		err = errors.Wrap(err, "failed to unmarshall turnkey createPrivateKeysResult interface into struct")
		return "", "", err
	}

	privateKeyId, address = privateKeyResult.getResult()

	fmt.Printf("privateKeyId: %s \n address: %s \n", privateKeyId, address)
	return
}

func ExtractSubOrganizationIdFromResult(result map[string]any) string {
	m := result["createSubOrganizationResult"]
	internalMap, ok := m.(map[string]any)
	if !ok {
		return ""
	}

	orgId, ok := internalMap["subOrganizationId"].(string)
	if !ok {
		return ""
	}

	return orgId
}

func ExctractTurnkeySignatureFromResult(result map[string]any) (*TurnkeySignature, error) {
	m := result["signRawPayloadResult"]
	internalMap, ok := m.(map[string]any)
	if !ok {
		fmt.Println("map", m)
		return nil, errors.New("malformed object")
	}

	r := fmt.Sprintf("0x%v", internalMap["r"])
	s := fmt.Sprintf("0x%v", internalMap["s"])
	v := fmt.Sprintf("0x%v", internalMap["v"])

	return &TurnkeySignature{
		R: r,
		S: s,
		V: v,
	}, nil
}

func ExtractPrivateKeyTagIdFromResult(result map[string]any) string {
	m := result["createPrivateKeyTagResult"]
	internalMap, ok := m.(map[string]any)
	if !ok {
		return ""
	}

	tagId, ok := internalMap["privateKeyTagId"].(string)
	if !ok {
		return ""
	}

	return tagId
}
