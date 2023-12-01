package tests

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	merchant_model "github.com/lucidconnect/silver-arrow/graphql/merchant/graph/model"
	"github.com/lucidconnect/silver-arrow/graphql/wallet/graph/model"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/lucidconnect/silver-arrow/service/merchant"
	"github.com/lucidconnect/silver-arrow/service/wallet"
	"github.com/stretchr/testify/assert"
)

func TestAddSubscription(t *testing.T) {
	r := repository.NewPostgresDB(db)
	ms := merchant.NewMerchantService(r)
	ws := wallet.NewWalletService(r, defaultChain)
	// mId := randKey()
	newProduct := merchant_model.NewProduct{
		Name:             "test xyz",
		Owner:            "0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2",
		Chain:            int(defaultChain),
		Token:            "USDC",
		ReceivingAddress: "0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2",
	}
	product, _ := ms.CreateProduct(newProduct)

	key := "0xe81f9f7146470e1e728cc44d22089098de6be6ebe3ca39f21b7b092f09b10cf5"
	p, _ := crypto.HexToECDSA(key[2:])
	owner := crypto.PubkeyToAddress(p.PublicKey).Hex()
	fmt.Println("owner", owner)
	newSub := model.NewSubscription{
		Chain:     int(defaultChain),
		Token:     "USDC",
		Amount:    1,
		Interval:  30,
		ProductID: uuid.MustParse(product.ProductID),
		// WalletAddress: "0x14De44b6100dE479655D752ECD2230D10F8fA061",
		WalletAddress: "0xb96442F14ac82E21c333A8bB9b03274Ae26eb79D",
		OwnerAddress:  "0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2",
	}

	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}

	chain := int64(newSub.Chain)
	_, op, err := ws.AddSubscription(uuid.MustParse("e0b3849f-5870-4ee1-ab1a-3882c0da7903"), newSub, usePaymaster, big.NewInt(0), chain)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	sig, _, err := erc4337.SignUserOp(op, key, erc4337.SUDO_MODE, nil, 80001)
	assert.NotEmpty(t, sig)
	assert.NoError(t, err)

	op["signature"] = hexutil.Encode(sig)
	fmt.Println(op["signature"])

	data, err := ws.ValidateSubscription(op, chain)
	assert.NotEmpty(t, data)
	assert.NoError(t, err)

	// target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"

	// err = ws.ExecuteCharge(newSub.WalletAddress, target, mId, "USDC", key, int64(newSub.Amount), usePaymaster)
	// assert.NoError(t, err)

	// fmt.Println("Data", data)
	// t.Fail()
}

func TestSubscriptionIsUnique(t *testing.T) {
	r := repository.NewPostgresDB(db)
	ws := wallet.NewWalletService(r, defaultChain)

	newSub := model.NewSubscription{
		Chain:     int(defaultChain),
		Token:     "USDC",
		Amount:    1,
		Interval:  30,
		ProductID: uuid.MustParse("aad69be2-8513-4fe1-b5df-63720630ae6b"),
		// WalletAddress: "0x14De44b6100dE479655D752ECD2230D10F8fA061",
		WalletAddress: "0xb96442F14ac82E21c333A8bB9b03274Ae26eb79D",
		OwnerAddress:  "0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2",
	}

	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}

	chain := int64(newSub.Chain)
	_, _, err := ws.AddSubscription(uuid.MustParse("e0b3849f-5870-4ee1-ab1a-3882c0da7903"), newSub, usePaymaster, big.NewInt(0), chain)
	if !assert.Error(t, err) {
		t.FailNow()
	}
}

func randKey(length int) string {
	key := make([]byte, length)

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

	ercBundler, _ := erc4337.NewAlchemyService(defaultChain)

	amount := big.NewInt(1000000)
	data, err := erc4337.CreateTransferCallData(target, token, defaultChain, amount)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	nonce, _ := ercBundler.GetAccountNonce(common.HexToAddress(sender))
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
