package wallet

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessKey(t *testing.T) {
	pubKey, privKey, err := CreateAccessKey()
	assert.NoError(t, err)
	fmt.Println(pubKey)
	fmt.Println(privKey)
	t.Fail()
}

func TestAccessKeyConsistency(t *testing.T) {
	address, privKey, _ := CreateAccessKey()
	fmt.Println("Address:", address)

	msg := []byte("hello")
	msgHash := ecrecover.ToEthSignedMessageHash(crypto.Keccak256(msg))
	key, _ := hexutil.Decode(privKey)
	sig, err := secp256k1.Sign(msgHash, key)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	sig[64] += 27

	recovered, err := ecrecover.Recover(msgHash, sig)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	fmt.Println("recovered address", recovered)
	if !assert.Equal(t, address, recovered.Hex()) {
		t.FailNow()
	}
}

func TestSignature(t *testing.T) {
	privateKey := "0xc1fce60cfb4b32bf4584e577904d806f8c5af28104d34e9923466eb8ca6faeff"
	publicKey := "0x6574f281AAaA788cf89e5269E9c842E50c5713fe"

	// fmt.Println(publicKey)
	p, err := crypto.HexToECDSA(privateKey[2:])
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	p.Public()

	address := crypto.PubkeyToAddress(p.PublicKey)
	fmt.Println("address", address)

	hash := crypto.Keccak256([]byte("stuff"))
	signature, err := crypto.Sign(hash, p)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	x := "0x00000002ffffffffffff00000000000040acee1113697bdee3077493896fa759d1b3e25500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020b77ce6ec08b85dcc468b94cea7cc539a3bbf9510ffffffffffff0000000000000000000000000000000000000000000000000000000000000000000000000041169aebb266dfebe8fb1501580331fb20806d92b925a887fdfeed400b1522b73b6f21d9b69c570c7211e80f09f14216a74be72bce0ff4c1c4771fb9cbd195f1c51b13ea99673f85abe74c79ff548dd015d34ffd260a2f0aa1e80218b3d5eca013f147f54168a7f1468800b25cf8bc7195de91ee8a6d8fbf3c9005d80dc53b95f36800"

	fmt.Println(hexutil.Decode(x))

	// recovered, err := crypto.Ecrecover(hash, signature)
	// if !assert.NoError(t, err) {
	// 	t.FailNow()
	// }
	// pb, _ :=	crypto.UnmarshalPubkey(recovered)

	// fmt.Println(len(recovered))
	pub, _ := crypto.SigToPub(hash, signature)
	fmt.Println(crypto.PubkeyToAddress(*pub))
	// fmt.Println("recovered address", common.BytesToAddress(recovered[0:20]))

	assert.Equal(t, crypto.PubkeyToAddress(*pub).Hex(), publicKey)
	t.Fail()

}
