package tests

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lucidconnect/silver-arrow/core/service/erc4337"
	"github.com/lucidconnect/silver-arrow/core/service/turnkey"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessKey(t *testing.T) {
	tk, _ := turnkey.NewTurnKeyService()
	if !assert.NotEmpty(t, tk) {
		t.FailNow()
	}

	orgId := "0714b26b-dfab-4de9-be01-0c2908025916"
	keyTag := "33c9bf06-f1c2-4dd1-909d-b426338f2553"

	// orgId := ""
	name := fmt.Sprintf("test-key-%s", randKey(6))
	activityId, err := tk.CreatePrivateKey(orgId, name, keyTag)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	result, err := tk.GetActivity(orgId, activityId)
	if !assert.NoError(t, err) || !assert.NotEmpty(t, result) {
		t.FailNow()
	}

	privateKeyId, address, err := turnkey.GetPrivateKeyIdFromResult(result)
	if !assert.NoError(t, err) || !assert.NotEmpty(t, privateKeyId) || !assert.NotEmpty(t, address) {
		t.FailNow()
	}

	fmt.Printf("privateKeyId - %v \n address - %v \n", privateKeyId, address)
}

func TestCreateSubOrganization(t *testing.T) {
	tk, _ := turnkey.NewTurnKeyService()
	if !assert.NotEmpty(t, tk) {
		t.FailNow()
	}

	orgId := tk.TurnkeyClient.DefaultOrganization()

	name := fmt.Sprintf("test-org-%s", randKey(6))
	fmt.Println("organization id - ", *orgId)
	activityId, err := tk.CreateSubOrganization(*orgId, name)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	result, err := tk.GetActivity(*orgId, activityId)
	if !assert.NoError(t, err) || !assert.NotEmpty(t, result) {
		t.FailNow()
	}
	fmt.Println(result)
}

func TestCreatePrivateKeyTag(t *testing.T) {
	tk, _ := turnkey.NewTurnKeyService()
	if !assert.NotEmpty(t, tk) {
		t.FailNow()
	}

	orgId := "0714b26b-dfab-4de9-be01-0c2908025916"

	// orgId := ""
	tag := fmt.Sprintf("test-key-%s", randKey(6))
	activityId, err := tk.CreatePrivateKeyTag(orgId, tag)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	result, err := tk.GetActivity(orgId, activityId)
	if !assert.NoError(t, err) || !assert.NotEmpty(t, result) {
		t.FailNow()
	}

	fmt.Println(result)
}

func TestSignMessage(t *testing.T) {
	tk, _ := turnkey.NewTurnKeyService()
	if !assert.NotEmpty(t, tk) {
		t.FailNow()
	}
	orgId := "0714b26b-dfab-4de9-be01-0c2908025916"
	privateKeyId := "001547ac-6758-419b-944b-cb8f94e7792e"
	address := "0xEED4276D700B776DC315496109ac2b1b483CaCa9"
	message := []byte("hello")
	messageHash := hexutil.Encode(crypto.Keccak256(message))
	activityId, err := tk.SignMessage(orgId, privateKeyId, messageHash)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println(activityId)

	result, err := tk.GetActivity(orgId, activityId)
	if !assert.NoError(t, err) || !assert.NotEmpty(t, result) {
		t.FailNow()
	}

	sig, err := turnkey.ExctractTurnkeySignatureFromResult(result)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	signature := sig.ParseSignature(erc4337.VALIDATOR_MODE)
	fmt.Printf("Signature - %v", signature)

	recoveredAddress, err := ecrecover.Recover(crypto.Keccak256(message), hexutil.MustDecode(signature))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	fmt.Println("recovered address - ", recoveredAddress.Hex())

	if !assert.Equal(t, address, recoveredAddress.Hex()) {
		t.FailNow()
	}
}
