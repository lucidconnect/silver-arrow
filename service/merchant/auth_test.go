package merchant

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stretchr/testify/assert"
)

func TestSignatureVerification(t *testing.T) {
	// xc, err := hex.DecodeString("a379c53749b89b570912c33788588b54c5539e5ba5ecd748fca6279eb8b27831")
	// if !assert.Nil(t, err) {
	// 	t.Fail()
	// }

	// fmt.Println(xc)

	dec, err := hex.DecodeString("716dcb74a5b473599a55a80e063d0de8091605c96813cc8542df78e2b4c8075d")
	if !assert.Nil(t, err) {
		t.Fail()
	}

	fmt.Println(dec)
	testHash := hexutil.MustDecode("0x716dcb74a5b473599a55a80e063d0de8091605c96813cc8542df78e2b4c8075d")
	fmt.Println(testHash)
	testSig := hexutil.MustDecode("0x162c9dff76986aa795d464e6cd76d3d26864c3fe58896e96af746bec76af6db94e14b1a91110327ca726803acc5317c0f8eaa9df4afa04892590015aace2bd321b")
	
	addr := common.HexToAddress("0x31De6d0c0d1ad1223081d72903b9C14773D7857b")
	fmt.Println(addr)
	// testSig[64] += 27
	// fmt.Println(testSig)

	ethSignedMsg := ecrecover.ToEthSignedMessageHash(testHash)
	// ad, err := ecrecover.Recover(ethSignedMsg, testSig)
	recoveredPubKey, err := ecrecover.Recover(ethSignedMsg, testSig)
	if !assert.Nil(t, err) {
		t.Fail()
	}

	// pub, _ := crypto.SigToPub(ethSignedMsg, testSig)
	// fmt.Println(crypto.CompressPubkey(pub))
	fmt.Printf("recovered public key: %v", (recoveredPubKey))
	// ok := crypto.VerifySignature(pubKey, ethSignedMsg, testSig)
	// if !ok {
	// 	t.Fail()
	// }
	// fmt.Println("valid? ", ok)
}

func ToEthSignedMessageHash(message []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	return ecrecover.Keccak256([]byte(msg))
}
