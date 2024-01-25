package erc4337_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/joho/godotenv"
	"github.com/lucidconnect/silver-arrow/core/gateway"
	"github.com/lucidconnect/silver-arrow/core/service/erc4337"
	"github.com/stretchr/testify/assert"
)

func Test_EstimateUserOperationGas(t *testing.T) {
	if err := godotenv.Load("../../../.env.test"); err != nil {
		t.Error(err)
	}
	entryPoint := os.Getenv("ENTRY_POINT")
	alchemy, _ := erc4337.NewAlchemyService(80001)
	userop, err := packTestUserOperation()
	if !assert.NoError(t, err) {
		t.Error(err)
	}
	sig := userop["signature"].(string)
	result, err := alchemy.EstimateUserOperationGas(entryPoint, sig, userop)
	if !assert.NoError(t, err) {
		t.Error(err)
	}

	fmt.Println(result)
}

func packTestUserOperation() (map[string]any, error) {
	key := "0xe81f9f7146470e1e728cc44d22089098de6be6ebe3ca39f21b7b092f09b10cf5"
	sender := "0xb96442F14ac82E21c333A8bB9b03274Ae26eb79D"
	target := "0xB77ce6ec08B85DcC468B94Cea7Cc539a3BbF9510"
	// token := "USDC"
	chainId := 80001

	alchemy, err := erc4337.NewAlchemyService(int64(chainId))
	if err != nil {
		return nil, err
	}

	// amount := big.NewInt(0)

	validatorAddress := os.Getenv("VALIDATOR_ADDRESS")
	executorAddress := os.Getenv("EXECUTOR_ADDRESS")

	data, err := gateway.SetValidatorExecutor(target, validatorAddress, executorAddress, sender, int64(chainId))
	if err != nil {
		return nil, err
	}
	nonce, _ := alchemy.GetAccountNonce(common.HexToAddress(sender))

	op, err := alchemy.CreateUnsignedUserOperation(sender, nil, data, nonce, true, int64(chainId))
	if err != nil {
		return nil, err
	}

	fmt.Println("user op", op)

	sig, _, err := erc4337.SignUserOp(op, key, erc4337.SUDO_MODE, nil, int64(chainId))
	if err != nil {
		return nil, err
	}

	op["signature"] = hexutil.Encode(sig)

	return op, nil
}
