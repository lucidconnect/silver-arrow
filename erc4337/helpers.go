package erc4337

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lucidconnect/silver-arrow/abi/KernelFactory"
	Kernel "github.com/lucidconnect/silver-arrow/abi/kernel"
	KernelStorage "github.com/lucidconnect/silver-arrow/abi/kernelStorage"
	"github.com/pkg/errors"
)

/**

 */

func getKernelStorageAbi() string {
	kernelABI := `[{
		"inputs": [
			{
				"internalType": "bytes4",
				"name": "_selector",
				"type": "bytes4"
			},
			{
				"internalType": "address",
				"name": "_executor",
				"type": "address"
			},
			{
				"internalType": "contract IKernelValidator",
				"name": "_validator",
				"type": "address"
			},
			{
				"internalType": "uint48",
				"name": "_validUntil",
				"type": "uint48"
			},
			{
				"internalType": "uint48",
				"name": "_validAfter",
				"type": "uint48"
			},
			{
				"internalType": "bytes",
				"name": "_enableData",
				"type": "bytes"
			}
		],
		"name": "setExecution",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}]`
	return kernelABI
}

func getAccountFactoryAbi(factoryAddress string) string {
	factory := KernelFactory.KernelFactoryABI
	return factory
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

func GetSetExecutionFnData(accountABI, validator, kernel string, enableData []byte) ([]byte, error) {
	contractABI, err := abi.JSON(strings.NewReader(accountABI))
	if err != nil {
		err = errors.Wrap(err, "abi.JSON() unable to read contract abi")
		return nil, err
	}

	// kernel execute fn selector: 0x51945447
	selector, err := hexutil.Decode("0x51945447")
	if err != nil {
		err = errors.Wrap(err, "invalid selector hex")
		return nil, err
	}
	fnSelector := [4]byte{}
	copy(fnSelector[:], selector)
	executorAddress := common.HexToAddress(kernel)
	validatorAddress := common.HexToAddress(validator)

	fmt.Println("enable data -", hexutil.Encode(enableData))
	payload, err := contractABI.Pack("setExecution", fnSelector, executorAddress, validatorAddress, big.NewInt(99999999999), big.NewInt(0), enableData)
	if err != nil {
		err = errors.Wrap(err, "abi.Pack() - ")
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
	payload, err := kernelStorageAbi. Pack(selector, args...)
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

func (b *ERCBundler) getContractInitCode(address common.Address) []byte {
	code, err := b.client.GetAccountCode(address)
	if err != nil {
		log.Err(err).Msg("Oops! could not fetch account code")
		return nil
	}
	fmt.Println("Account code - ", code)
	return code
}

func getCallGasLimit() *big.Int {
	return big.NewInt(70000)
}

func getVerificationGasLimit() *big.Int {
	return big.NewInt(300000)
}

func getPreVerificationGas() *big.Int {
	return big.NewInt(89925)
}

func getMaxFeePerGas() *big.Int {
	return big.NewInt(3079999999)
}

func getMaxPriorityFeePerGas() *big.Int {
	return big.NewInt(3079999999)
}

func getSigningKey(privateKey string) (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.HexToECDSA(privateKey[2:])
	if err != nil {
		err = errors.Wrapf(err, "private key parse failure, %v", privateKey)
		return nil, err
	}
	return privKey, nil
}
