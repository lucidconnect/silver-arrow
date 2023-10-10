package merchant

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stretchr/testify/assert"
)

func TestSignatureVerification(t *testing.T) {
	// xc, err := hex.DecodeString("a379c53749b89b570912c33788588b54c5539e5ba5ecd748fca6279eb8b27831")
	// if !assert.Nil(t, err) {
	// 	t.Fail()
	// }

	// fmt.Println(xc)

	testHash := hexutil.MustDecode("0xa379c53749b89b570912c33788588b54c5539e5ba5ecd748fca6279eb8b27831")
	// fmt.Println(testHash)
	testSig := hexutil.MustDecode("0x1142d7b29272ad3d7e28928d55dd9a168adbf2b8cc66820cb0ef785741f9bf6b7745232d9241d8147fee1b11c4243dfdaa0f81c3b3927991cb34f5c401da01821b")
	pubKey := hexutil.MustDecode("0x02e9ec28a584b0ff8206bec029c1fdcc40688f0f6310d6888dec28e8849afcb9c7")
	// fmt.Println("",pubKey)
	pub, err := crypto.DecompressPubkey(pubKey)
	if !assert.Nil(t, err) {
		t.Fail()
	}
	addr := crypto.PubkeyToAddress(*pub)
	fmt.Println(addr)
	// testSig[64] += 27
	// fmt.Println(testSig)

	ethSignedMsg := ecrecover.ToEthSignedMessageHash(testHash)
	// ad, err := ecrecover.Recover(ethSignedMsg, testSig)
	recoveredPubKey, err := ecrecover.Recover(ethSignedMsg, testSig)
	if !assert.Nil(t, err) {
		t.Fail()
	}

	ecrecover.Example()

	// pub, _ := crypto.SigToPub(ethSignedMsg, testSig)
	// fmt.Println(crypto.CompressPubkey(pub))
	fmt.Printf("recovered public key: %v", (recoveredPubKey))
	// ok := crypto.VerifySignature(pubKey, ethSignedMsg, testSig)
	// if !ok {
	// 	t.Fail()
	// }
	// fmt.Println("valid? ", ok)
}
