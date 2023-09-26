package erc4337

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pkg/errors"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

var (
	SUDO_MODE      = "0x00000000"
	VALIDATOR_MODE = "0x00000001"
	ENABLE_MODE    = "0x00000002"
)

type ERCBundler struct {
	EntryPoint string
	client     *Client
}

func NewERCBundler(entrypoint string, client *Client) *ERCBundler {
	return &ERCBundler{
		EntryPoint: entrypoint,
		client:     client,
	}
}

func (b *ERCBundler) GetClient() *Client {
	return b.client
}

func (b *ERCBundler) AccountNonce(sender string) (*big.Int, error) {
	senderAddress := common.HexToAddress(sender)

	nonce, err := b.client.GetAccountNonce(common.HexToAddress(b.EntryPoint), senderAddress)
	if err != nil {
		if err.Error() == "no contract code at given address" {
			return common.Big0, nil
		}
		err = errors.Wrap(err, "AccountNonce() -")
		return nil, err
	}
	return nonce, nil
}

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
// CreateUserOperation returns a signed useroperation
// func (b *ERCBundler) CreateUserOperation(sender, target string, callData []byte, nonce, amount *big.Int, sponsored bool, key string, chain int64) (map[string]any, error) {
// 	var paymasterResult *PaymasterResult
// 	var err error
// 	var callGasLimit, verificationGas, preVerificationGas string

// 	senderAddress := common.HexToAddress(sender)
// 	tok := make([]byte, 65)
// 	rand.Read(tok)

// 	o := map[string]any{
// 		"sender":               senderAddress.Hex(),
// 		"nonce":                hexutil.EncodeBig(nonce),
// 		"initCode":             "0x",
// 		"callData":             hexutil.Encode(callData),
// 		"callGasLimit":         hexutil.EncodeBig(big.NewInt(0)),
// 		"verificationGasLimit": hexutil.EncodeBig(big.NewInt(0)),
// 		"preVerificationGas":   hexutil.EncodeBig(big.NewInt(0)),
// 		"maxFeePerGas":         getMaxFeePerGas().String(),
// 		"maxPriorityFeePerGas": getMaxPriorityFeePerGas().String(),
// 		"signature":            hexutil.Encode(tok),
// 		"paymasterAndData":     "0x",
// 	}
// 	// fmt.Println(o)
// 	paymasterContext := map[string]any{
// 		"type": "payg",
// 	}

// 	if sponsored {
// 		paymasterResult, err = b.client.SponsorUserOperation(b.EntryPoint, o, paymasterContext)
// 		if err != nil {
// 			err = errors.Wrap(err, "call to sponsor user op failed")
// 			return nil, err
// 		}

// 		// callGasLimit, err = hexutil.DecodeBig(paymasterResult.CallGasLimit)
// 		// if err != nil {
// 		// 	err = errors.Wrapf(err, "decoding gas limit - %v failed", paymasterResult.CallGasLimit)
// 		// 	return nil, err
// 		// }

// 		// verificationGas, err = hexutil.DecodeBig(paymasterResult.VerificationGasLimit)
// 		// if err != nil {
// 		// 	err = errors.Wrapf(err, "decoding verification gas limit - %v failed", paymasterResult.VerificationGasLimit)
// 		// 	return nil, err
// 		// }

// 		// preVerificationGas, err = hexutil.DecodeBig(paymasterResult.PreVerificationGas)
// 		// if err != nil {
// 		// 	err = errors.Wrapf(err, "decoding pre verification gas limit - %v failed", paymasterResult.PreVerificationGas)
// 		// 	return nil, err
// 		// }

// 		o["paymasterAndData"] = paymasterResult.PaymasterAndData
// 	} else {
// 		// fmt.Println("not using paymaster")
// 		result, err := b.client.EstimateUserOperationGas(b.EntryPoint, o)
// 		if err != nil {
// 			return nil, err
// 		}
// 		callGasLimit = result.CallGasLimit
// 		verificationGas = result.VerificationGasLimit
// 		preVerificationGas = result.PreVerificationGas
// 	}

// 	o["callGasLimit"] = callGasLimit
// 	o["verificationGasLimit"] = verificationGas
// 	o["preVerificationGas"] = preVerificationGas

// 	// TODO: Do I need a merchantId here?
// 	sig, _, err := SignUserOp(o, key, SUDO_MODE, nil, chain)
// 	if err != nil {
// 		err = errors.Wrap(err, "CreateUserOperation() call to sign user op failed")
// 		return nil, err
// 	}

// 	o["signature"] = hexutil.Encode(sig)
// 	// fmt.Println(o)
// 	return o, nil
// }

func (b *ERCBundler) CreateUnsignedUserOperation(sender string, initCode, callData []byte, nonce *big.Int, sponsored bool, chain int64) (map[string]any, error) {
	// var paymasterResult *PaymasterResult
	var paymaster *AlchemyPaymasterResult
	var gasEstimate *GasEstimateResult
	var err error
	var paymasterAndData, callGasLimit, maxFeePerGas, maxPriorityFeePerGas, verificationGas, preVerificationGas string
	senderAddress := common.HexToAddress(sender)
	tok := make([]byte, 65)
	rand.Read(tok)
	// fmt.Println("senderAddress", senderAddress)
	o := map[string]any{
		"sender":   senderAddress.Hex(),
		"nonce":    hexutil.EncodeBig(nonce),
		"initCode": hexutil.Encode(initCode),
		"callData": hexutil.Encode(callData),
	}

	if sponsored {
		policyId := "2e865ced-98e5-4265-a20e-b46c695a28bd"
		paymaster, err = b.client.RequestGasAndPaymasterAndData(policyId, b.EntryPoint, hexutil.Encode(tok), o)
		if err != nil {
			err = errors.Wrap(err, "call to sponsor user op failed")
			return nil, err
		}

		callGasLimit = paymaster.CallGasLimit
		verificationGas = hexutil.EncodeBig(getVerificationGasLimit())
		xy := hexutil.MustDecode(paymaster.VerificationGasLimit)
		fmt.Printf("paymaster returned verification gas limit - %v", new(big.Int).SetBytes(xy) )
		// verificationGas = paymaster.VerificationGasLimit
		// preVerificationGas = hexutil.EncodeBig(getPreVerificationGas())
		maxPriorityFeePerGas = paymaster.MaxPriorityFeePerGas
		preVerificationGas = paymaster.PreVerificationGas
		maxFeePerGas = paymaster.MaxFeePerGas
		paymasterAndData = "0x"
	} else {
		fmt.Println("not using paymaster")
		o["callGasLimit"] = "0x16710"
		o["verificationGasLimit"] = hexutil.EncodeBig(getVerificationGasLimit())
		o["preVerificationGas"] = hexutil.EncodeBig(getVerificationGasLimit())
		o["maxPriorityFeePerGas"] = hexutil.EncodeBig(getMaxPriorityFeePerGas())
		o["maxFeePerGas"] = hexutil.EncodeBig(getMaxFeePerGas())
		o["signature"] = hexutil.Encode(tok)
		o["paymasterAndData"] = "0x"
		// if initCode == nil {
		gasEstimate, err = b.client.EstimateUserOperationGas(b.EntryPoint, o)
		if err != nil {
			return nil, err
		} else {
			callGasLimit = "0x16710"
		}
		paymasterAndData = "0x"
		verificationGas = hexutil.EncodeBig(getVerificationGasLimit())
		// preVerificationGas = hexutil.EncodeBig(getVerificationGasLimit())

		maxPriorityFeePerGas = hexutil.EncodeBig(getMaxPriorityFeePerGas())
		maxFeePerGas = hexutil.EncodeBig(getMaxFeePerGas())
		preVerificationGas = gasEstimate.PreVerificationGas
	}

	o["callGasLimit"] = callGasLimit
	o["verificationGasLimit"] = verificationGas
	o["preVerificationGas"] = preVerificationGas
	o["maxPriorityFeePerGas"] = maxPriorityFeePerGas
	o["maxFeePerGas"] = maxFeePerGas
	o["signature"] = hexutil.Encode(tok)
	o["paymasterAndData"] = paymasterAndData

	return o, nil
}

// SendUserOp uses the necessary inputs to send a useroperation to the smart account
func (b *ERCBundler) SendUserOp(op map[string]any) (string, error) {
	return b.client.SendUserOperation(b.EntryPoint, op)
}

func (b *ERCBundler) GetBalance(address string) error {
	bal, err := b.client.GetBalance(address)
	fmt.Printf("%s balance: %s \n", address, bal)
	return err
}

func (b *ERCBundler) GetUserOp(userophash string) (map[string]any, error) {
	return b.client.GetUserOperationByHash(userophash)
}

func SignUserOp(op map[string]any, key, mode string, merchantId []byte, chain int64) ([]byte, []byte, error) {
	chainId := big.NewInt(chain)
	entrypoint := GetEntryPointAddress()

	operation, err := userop.New(op)
	if err != nil {
		return nil, nil, err
	}

	opHash := operation.GetUserOpHash(entrypoint, chainId)
	hash := opHash.Bytes()

	kb, err := hexutil.Decode(key)
	if err != nil {
		log.Println("SignUserOp() - ", err)
		err = errors.Wrap(err, "SignUserOp() - ")
		return nil, nil, err
	}
	fmt.Println("kb", kb)
	common.FromHex(key)

	// Kernel has a specific convention for encoding signatures in order to determing the mode see (https://github.com/stackup-wallet/userop.js/blob/main/src/preset/builder/kernel.ts#L114-L123)
	signature, _ := hexutil.Decode(mode)

	// if merchantId != nil {
	// 	fmt.Println("merchant id - ", merchantId)
	// 	signature = append(signature, merchantId...)
	// }

	sig, err := secp256k1.Sign(ecrecover.ToEthSignedMessageHash(hash), kb)
	if err != nil {
		err = errors.Wrap(err, "generating signature failed.")
		return nil, nil, err
	}
	sig[64] += 27
	signature = append(signature, sig...)

	fmt.Println("raw sig - ", hexutil.Encode(sig))
	fmt.Println("hash - ", hexutil.Encode(opHash[:]))

	return signature, ecrecover.ToEthSignedMessageHash(hash), nil
}

/*
There are two possible scenerios here:
1. The call data is for an erc20 token transfer
2. The call data is for an Eth transfer
*/
// TODO: Create a Token type
func CreateTransferCallData(toAddress, token string, amount *big.Int) ([]byte, error) {
	accountABI := getAccountABI()

	if token == "ETH" {
		callData, err := GetExecuteFnData(accountABI, toAddress, amount, nil)
		if err != nil {
			err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
			return nil, err
		}
		// fmt.Println("call data ", hexutil.Encode(callData))
		return callData, nil
	}

	erc20Token := getErc20TokenABI()

	erc20TransferData, err := GetTransferFnData(erc20Token, toAddress, amount)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create erc20 call data")
		return nil, err
	}

	tokenAddress := GetTokenAddres(token)
	callData, err := GetExecuteFnData(accountABI, tokenAddress, common.Big0, erc20TransferData)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

func CreateSetExecutionCallData(enableData []byte, kernel string) ([]byte, error) {
	kernelAbi := getKernelStorageAbi()
	lucidValidator := os.Getenv("LUCID_VALIDATOR")
	callData, err := GetSetExecutionFnData(kernelAbi, lucidValidator, kernel, enableData)
	if err != nil {
		err = errors.Wrap(err, "CreateSetExecutionCallData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

func CreateFactoryFnData(enableData []byte, index *big.Int) ([]byte, error) {
	factoryAbi := getAccountFactoryAbi()
	defaultValidatorAddress := os.Getenv("DEFAULT_VALIDATOR")
	callData, err := GetCreateAccountFnData(factoryAbi, defaultValidatorAddress, enableData, index)
	if err != nil {
		err = errors.Wrap(err, "CreateFactoryFnData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

func GetTokenAddres(token string) string {
	return "0x0fa8781a83e46826621b3bc094ea2a0212e71b23"
}
