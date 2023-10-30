package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/lucidconnect/silver-arrow/service/wallet"
	"github.com/stretchr/testify/assert"
)

func Test_SendUseropWithPaymaster(t *testing.T) {
	sender := "0xb96442F14ac82E21c333A8bB9b03274Ae26eb79D"
	target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"
	token := "ETH"
	key := "0xe81f9f7146470e1e728cc44d22089098de6be6ebe3ca39f21b7b092f09b10cf5"
	chainId := 80001
	node, err := erc4337.NewAlchemyService(int64(chainId))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	nonce, _ := node.GetAccountNonce(common.HexToAddress(sender))

	amount := big.NewInt(0)
	data, err := erc4337.CreateTransferCallData(target, token, int64(chainId), amount)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	_, err = wallet.GetContractInitCode(common.HexToAddress("0x85fc2E4425d0DAba7426F50091a384ee05D37Cd2"), common.Big0)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	op, err := node.CreateUnsignedUserOperation(sender, nil, data, nonce, true, int64(chainId))
	assert.NoError(t, err)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	fmt.Println("user op", op)
	fmt.Println("max fee per gas", op["maxFeePerGas"])
	sig, _, err := erc4337.SignUserOp(op, key, erc4337.SUDO_MODE, nil, int64(chainId))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	// sig[4:][64] -= 27
	// fmt.Println(sig)
	// fmt.Println(sig[0:4])
	// secp256k1.RecoverPubkey(hash, sig[4:])
	// pubkey, err := crypto.Ecrecover(hash, sig[4:])
	// if !assert.NoError(t, err) {
	// 	t.FailNow()
	// }
	// if !assert.NotEmpty(t, pubkey) {
	// 	t.FailNow()
	// }
	// fmt.Println(pubkey)
	// p, err := crypto.UnmarshalPubkey(pubkey)
	// if !assert.NoError(t, err) {
	// 	t.FailNow()
	// }
	// address := crypto.PubkeyToAddress(*p).Hex()

	// fmt.Println(address)

	op["signature"] = hexutil.Encode(sig)

	useropHash, err := node.SendUserOperation(op)
	if !assert.NoError(t, err) {
		t.Fail()
	}

	fmt.Println(useropHash)
	// t.Fail()
	// set validator

}
