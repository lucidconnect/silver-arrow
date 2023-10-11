package auth

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
}

func TestAccessKeyConsistency(t *testing.T) {
	address, privKey, _ := CreateAccessKey()
	fmt.Println("Address:", address)

	msg := []byte("hello")
	hashed := sha256.New().Sum(msg)
	msgHash := ecrecover.ToEthSignedMessageHash(crypto.Keccak256(hashed))
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
	privateKey := "0xe81f9f7146470e1e728cc44d22089098de6be6ebe3ca39f21b7b092f09b10cf5"
	publicKey := "0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2"

	p, err := crypto.HexToECDSA(privateKey[2:])
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	p.Public()

	address := crypto.PubkeyToAddress(p.PublicKey)
	fmt.Println("address", address)
	hash := common.HexToHash("0x5782249a8bfdb910171fb10b976bcf9b1c2aa4cfc619c6e670002cd06a9f5ea4")
	sig, err := secp256k1.Sign((hash.Bytes()), hexutil.MustDecode(privateKey))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	sig[64] += 27
	fmt.Println("raw sig - ", hexutil.Encode(sig))


	sig[64] -= 27
	pub, err := crypto.SigToPub((hash.Bytes()), sig)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	fmt.Println(crypto.PubkeyToAddress(*pub))

	assert.Equal(t, crypto.PubkeyToAddress(*pub).Hex(), publicKey)
	t.Fail()

}
