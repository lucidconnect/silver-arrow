package erc4337

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lucidconnect/silver-arrow/abi/KernelFactory"
	Kernel "github.com/lucidconnect/silver-arrow/abi/kernel"
	KernelStorage "github.com/lucidconnect/silver-arrow/abi/kernelStorage"
	"github.com/pkg/errors"
)

/**
Types
*/

var (
	SUDO_MODE      = "0x00000000"
	VALIDATOR_MODE = "0x00000001"
	ENABLE_MODE    = "0x00000002"
)

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type GasEstimateResult struct {
	PreVerificationGas   string `json:"preVerificationGas"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	CallGasLimit         string `json:"callGasLimit"`
}

type UserOperation struct {
	Sender               string `json:"sender"`
	Nonce                string `json:"nonce"`
	InitCode             string `json:"initCode"`
	CallData             string `json:"callData"`
	CallGasLimit         string `json:"callGasLimit"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	PreVerificationGas   string `json:"preVerificationGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	PaymasterAndData     string `json:"paymasterAndData"`
	Signature            string `json:"signature"`
}

type PaymasterResult struct {
	PaymasterAndData     string `json:"paymasterAndData"`
	PreVerificationGas   string `json:"preVerificationGas"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	CallGasLimit         string `json:"callGasLimit"`
}

type AlchemyPaymasterRequest struct {
	PolicyId       string `json:"policyId"`
	EntryPoint     string `json:"entryPoint"`
	DummySignature string `json:"dummySignature"`
	UserOperation  any    `json:"userOperation"`
	FeeOverride    any    `json:"feeOverride"`
}
type AlchemyPaymasterResult struct {
	PaymasterAndData     string      `json:"paymasterAndData"`
	PreVerificationGas   string      `json:"preVerificationGas"`
	VerificationGasLimit string      `json:"verificationGasLimit"`
	CallGasLimit         string      `json:"callGasLimit"`
	MaxFeePerGas         string      `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string      `json:"maxPriorityFeePerGas"`
	Error                ErrorObject `json:"error"`
}

type ErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func newUserOp(userop map[string]any) UserOperation {
	sender, _ := userop["sender"].(string)
	nonce, _ := userop["nonce"].(string)
	initCode, _ := userop["initCode"].(string)
	callGasLimit, _ := userop["callGasLimit"].(string)
	verificationGasLimit, _ := userop["verificationGasLimit"].(string)
	preVerificationGas, _ := userop["preVerificationGas"].(string)

	maxFeePerGas, _ := userop["maxFeePerGas"].(string)
	maxPriorityFeePerGas, _ := userop["maxPriorityFeePerGas"].(string)
	paymasterAndData, _ := userop["paymasterAndData"].(string)
	signature, _ := userop["signature"].(string)
	callData, _ := userop["callData"].(string)

	return UserOperation{
		Sender:               sender,
		Nonce:                nonce,
		InitCode:             initCode,
		CallGasLimit:         callGasLimit,
		VerificationGasLimit: verificationGasLimit,
		PreVerificationGas:   preVerificationGas,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
		PaymasterAndData:     paymasterAndData,
		Signature:            signature,
		CallData:             callData,
	}
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

func GetExecuteFnData(to string, amount *big.Int, callData []byte) ([]byte, error) {
	dest := common.HexToAddress(to)

	contractABI, err := Kernel.KernelMetaData.GetAbi()
	// contractABI, err := abi.JSON(strings.NewReader(accountABI))
	if err != nil {
		err = errors.Wrap(err, "abi.JSON() unable to read contract abi")
		return nil, err
	}
	// op := vm.CALL
	payload, err := contractABI.Pack("execute", dest, amount, callData, uint8(0))
	if err != nil {
		err = errors.Wrap(err, "PacK() unable to prepare tx payload")
		return nil, err
	}
	return payload, nil
}

func GetCreateAccountFnData(accountImplementation common.Address, enableData []byte, index *big.Int) ([]byte, error) {
	factory := KernelFactory.KernelFactoryABI
	factoryAbi, err := abi.JSON(strings.NewReader(factory))
	if err != nil {
		err = errors.Wrap(err, "abi.JSON() unable to read contract abi")
		return nil, err
	}

	payload, err := factoryAbi.Pack("createAccount", accountImplementation, enableData, index)
	if err != nil {
		return nil, err
	}
	fmt.Println("acc payload ", payload)
	return payload, nil
}

func EncodeKernelStorageWithSelector(selector string, args ...interface{}) ([]byte, error) {
	kernelStorage := KernelStorage.KernelStorageABI
	kernelStorageAbi, err := abi.JSON(strings.NewReader(kernelStorage))
	if err != nil {
		err = errors.Wrap(err, "abi.JSON() unable to read kernelStorage abi")
		return nil, err
	}
	payload, err := kernelStorageAbi.Pack(selector, args...)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func getErc20TokenABI() string {
	tokenABI := `[{
        "constant": false,
        "inputs": [
            {
                "name": "to",
                "type": "address"
            },
            {
                "name": "value",
                "type": "uint256"
            }
        ],
        "name": "transfer",
        "outputs": [
            {
                "name": "",
                "type": "bool"
            }
        ],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }]`
	return tokenABI
}

func getAccountABI() string {
	accountABI := `[{
		"inputs": [
			{
				"internalType": "address",
				"name": "to",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "value",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "data",
				"type": "bytes"
			},
			{
				"internalType": "bytes",
				"name": "operation",
				"type": "uint8"
			}
		],
		"name": "execute",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}]`

	return accountABI
}

// func getERC20TokenAddress(token string) string {
// 	return ""
// }

func GetEntryPointAddress() common.Address {
	entrypointAddress := "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"
	return common.HexToAddress(entrypointAddress)
}

func getVerificationGasLimit() *big.Int {
	return big.NewInt(300000)
}

func getMaxFeePerGas() *big.Int {
	return big.NewInt(3079999999)
}

func getMaxPriorityFeePerGas() *big.Int {
	return big.NewInt(3079999999)
}
