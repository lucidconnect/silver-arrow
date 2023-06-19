package useroperation

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
func CreateUserOp(sender, to, token string, amount *big.Int) ([]byte, error) {
	senderAddress := common.HexToAddress(sender)
	callData, err := CreateTransferCallData(to, token, amount)
	if err != nil {
		err = errors.Wrap(err, "CreateUserOp(): failed to create call data")
		return nil, err
	}
	operation := map[string]any{
		"sender":               senderAddress,
		"nonce":                0,
		"initCode":             getContractInitCode(),
		"callData":             callData,
		"callGasLimit":         getCallGasLimit(),
		"verificationGasLimit": getVerificationGasLimit(),
		"preVerificationGas":   getPreVerificationGas(),
		"maxFeePerGas":         getMaxFeePerGas(),
		"paymasterAndData":     nil,
	}

	op, err := userop.New(operation)
	if err != nil {
		err = errors.Wrap(err, "CreateUserOp(): failed to create userop")
		return nil, err
	}

	op.Signature = signUserOp(operation)
	fmt.Println("user operation object", op)
	return op.PackForSignature(), nil
}

func signUserOp(operation map[string]any) []byte {
	// not implemented
	return nil
}

func getContractInitCode() []byte {
	return nil
}

func getCallGasLimit() *big.Int {
	return big.NewInt(0)
}

func getVerificationGasLimit() *big.Int {
	return big.NewInt(0)
}

func getPreVerificationGas() *big.Int {
	return big.NewInt(0)
}

func getMaxFeePerGas() *big.Int {
	return big.NewInt(0)
}

func getMaxPriorityFeePerGas() *big.Int {
	return big.NewInt(0)
}

/*
const accountABI = ["function execute(address to, uint256 value, bytes data)"];
// An ABI can be fragments and does not have to include the entire interface.
// As long as it includes the parts we want to use.
const partialERC20TokenABI = ["function transfer(address to, uint amount) returns (bool)",];

encodeFunctionCallData("execute",[tokenAddress,])

There are two possible scenerios here:
1. The call data is for an erc20 token transfer
2. The call data is for an Eth transfer
*/
func CreateTransferCallData(toAddress, token string, amount *big.Int) ([]byte, error) {
	accountABI := getAccountABI()

	if token == "ETH" {
		callData, err := GetExecuteFnData(accountABI, toAddress, amount, nil)
		if err != nil {
			err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
			return nil, err
		}
		return callData, nil
	}

	erc20Token := getErc20TokenABI()

	erc20TransferData, err := GetTransferFnData(erc20Token, toAddress, amount)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create erc20 call data")
		return nil, err
	}

	callData, err := GetExecuteFnData(accountABI, "", common.Big0, erc20TransferData)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

// This wil be encoded as the data passed into the execute function
// to is the destination address for the transaction
// amount is the value to be transferred
func GetTransferFnData(erc20TokenABI string, to string, amount *big.Int) ([]byte, error) {
	dest := common.HexToAddress(to)
	contractABI, err := abi.JSON(strings.NewReader(erc20TokenABI))
	if err != nil {
		err = errors.Wrap(err, "GetTransferFnData(): unable to read contract abi")
		return nil, err
	}

	payload, err := contractABI.Pack("transfer", dest, amount)
	if err != nil {
		err = errors.Wrap(err, "GetTransferFnData(): unable to prepare tx payload")
		return nil, err
	}

	return payload, nil
}

func GetExecuteFnData(accountABI, to string, amount *big.Int, callData []byte) ([]byte, error) {
	dest := common.HexToAddress(to)

	contractABI, err := abi.JSON(strings.NewReader(accountABI))
	if err != nil {
		err = errors.Wrap(err, "abi.JSON() unable to read contract abi")
		return nil, err
	}

	payload, err := contractABI.Pack("execute", dest, amount, callData)
	if err != nil {
		err = errors.Wrap(err, "PacK() unable to prepare tx payload")
		return nil, err
	}
	return payload, nil
}

func getErc20TokenABI() string {
	tokenABI := ""
	return tokenABI
}

func getAccountABI() string {
	accountABI := ""

	return accountABI
}

func getERC20TokenAddress(token string) string {
	return ""
}
