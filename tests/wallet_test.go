package tests

import (
	"crypto/rand"
	"fmt"
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
	r := repository.NewMongoDb(mongoClient)
	ercBundler := erc4337.NewERCBundler(entrypointAddress, nodeClient)

	ws := wallet.NewWalletService(r, ercBundler)
	mId := randKey()

	key := "0xe81f9f7146470e1e728cc44d22089098de6be6ebe3ca39f21b7b092f09b10cf5"
	p, _ := crypto.HexToECDSA(key[2:])
	owner := crypto.PubkeyToAddress(p.PublicKey).Hex()
	fmt.Println("owner", owner)
	newSub := model.NewSubscription{
		Chain:         80001,
		NextChargeAt:  nil,
		Token:         "USDC",
		Amount:        1,
		Interval:      30,
		MerchantID:    mId,
		WalletAddress: "0x14De44b6100dE479655D752ECD2230D10F8fA061",
		OwnerAddress:  "0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2",
	}

	_, op, err := ws.AddSubscription(newSub)
	assert.NoError(t, err)

	fmt.Println(op)
	sig, err := erc4337.SignUserOp(op, key, erc4337.SUDO_MODE, nil, 80001)
	assert.NoError(t, err)
	op["signature"] = hexutil.Encode(sig)
	data, sKey, err := ws.ValidateSubscription(op)
	assert.NoError(t, err)

	target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"
	err = ws.ExecuteCharge(newSub.WalletAddress, target, mId, "USDC", sKey, int64(newSub.Amount))
	assert.NoError(t, err)

	fmt.Println("Data", data)
	t.Fail()
}

func randKey() string {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}
	fmt.Println(key)
	return hexutil.Encode(key)
}
