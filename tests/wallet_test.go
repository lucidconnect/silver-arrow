package tests

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/helicarrierstudio/silver-arrow/erc4337"
	"github.com/helicarrierstudio/silver-arrow/graph/model"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/wallet"
	"github.com/stretchr/testify/assert"
)

func TestAddSubscription(t *testing.T) {
	r := repository.NewWalletRepo(db)
	ercBundler := erc4337.NewERCBundler(entrypointAddress, nodeClient)

	ws := wallet.NewWalletService(r, ercBundler)
	// mId := randKey()
	mId := "3838hr8hud9dijh3j"
	key := "0xe81f9f7146470e1e728cc44d22089098de6be6ebe3ca39f21b7b092f09b10cf5"
	p, _ := crypto.HexToECDSA(key[2:])
	owner := crypto.PubkeyToAddress(p.PublicKey).Hex()
	fmt.Println("owner", owner)
	newSub := model.NewSubscription{
		Chain:        80001,
		NextChargeAt: nil,
		Token:        "USDC",
		Amount:       1,
		Interval:     30,
		MerchantID:   mId,
		// WalletAddress: "0x14De44b6100dE479655D752ECD2230D10F8fA061",
		WalletAddress: "0x6a6F07c5c32F5fb20393a2110B2Bf0925e59571b",
		OwnerAddress:  "0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2",
	}

	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}

	_, op, err := ws.AddSubscription(newSub, usePaymaster, big.NewInt(0))
	assert.NoError(t, err)

	sig, _, err := erc4337.SignUserOp(op, key, erc4337.SUDO_MODE, nil, 80001)
	assert.NotEmpty(t, sig)
	assert.NoError(t, err)

	op["signature"] = hexutil.Encode(sig)
	fmt.Println(op["signature"])

	data, sKey, err := ws.ValidateSubscription(op)
	assert.NotEmpty(t, data)
	assert.NotEmpty(t, sKey)
	assert.NoError(t, err)

	// target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"

	// err = ws.ExecuteCharge(newSub.WalletAddress, target, mId, "USDC", key, int64(newSub.Amount), usePaymaster)
	// assert.NoError(t, err)

	// fmt.Println("Data", data)
	// t.Fail()
}

func randKey() string {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}
	// fmt.Println(key)
	return hexutil.Encode(key)
}

func TestSignature(t *testing.T) {
	sender := "0x14De44b6100dE479655D752ECD2230D10F8fA061"
	target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"
	token := "USDC"

	ercBundler := erc4337.NewERCBundler(entrypointAddress, nodeClient)

	amount := big.NewInt(1000000)
	data, err := erc4337.CreateTransferCallData(target, token, amount)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	nonce, _ := ercBundler.AccountNonce(sender)
	key := "0xe81f9f7146470e1e728cc44d22089098de6be6ebe3ca39f21b7b092f09b10cf5"

	chainId := 80001
	op, err := ercBundler.CreateUnsignedUserOperation(sender, nil, data, nonce, true, int64(chainId))
	assert.NoError(t, err)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	sig, hash, err := erc4337.SignUserOp(op, key, erc4337.SUDO_MODE, nil, int64(chainId))
	assert.NoError(t, err)
	pubkey, err := crypto.Ecrecover(hash, sig[4:])
	assert.NoError(t, err)
	fmt.Println(pubkey)
	p, err := crypto.UnmarshalPubkey(pubkey)
	assert.NoError(t, err)

	address := crypto.PubkeyToAddress(*p).Hex()

	fmt.Println(address)
	t.Fail()
}
