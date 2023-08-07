package erc4337

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
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

func (b *ERCBundler) GetClient() *ethclient.Client {
	return b.client.GetEthClient()
}

func (b *ERCBundler) AccountNonce(sender string) (*big.Int, error) {
	senderAddress := common.HexToAddress(sender)

	nonce, err := b.client.GetAccountNonce(senderAddress)
	if err != nil {
		err = errors.Wrap(err, "AccountNonce() -")
		return nil, err
	} // pimlico
	// fmt.Println("nonce:", nonce)
	return nonce, nil
}

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
// CreateUserOperation returns a signed useroperation
func (b *ERCBundler) CreateUserOperation(sender, target string, callData []byte, nonce, amount *big.Int, sponsored bool, key string, chain int64) (map[string]any, error) {
	var paymasterResult *PaymasterResult
	var err error
	var callGasLimit, verificationGas, preVerificationGas string

	senderAddress := common.HexToAddress(sender)
	tok := make([]byte, 65)
	rand.Read(tok)

	o := map[string]any{
		"sender":               senderAddress.Hex(),
		"nonce":                hexutil.EncodeBig(nonce),
		"initCode":             "0x",
		"callData":             hexutil.Encode(callData),
		"callGasLimit":         hexutil.EncodeBig(big.NewInt(0)),
		"verificationGasLimit": hexutil.EncodeBig(big.NewInt(0)),
		"preVerificationGas":   hexutil.EncodeBig(big.NewInt(0)),
		"maxFeePerGas":         getMaxFeePerGas().String(),
		"maxPriorityFeePerGas": getMaxPriorityFeePerGas().String(),
		"signature":            hexutil.Encode(tok),
		"paymasterAndData":     "0x",
	}
	fmt.Println(o)
	paymasterContext := map[string]any{
		"type": "payg",
	}

	if sponsored {
		paymasterResult, err = b.client.SponsorUserOperation(b.EntryPoint, o, paymasterContext)
		if err != nil {
			err = errors.Wrap(err, "call to sponsor user op failed")
			return nil, err
		}

		// callGasLimit, err = hexutil.DecodeBig(paymasterResult.CallGasLimit)
		// if err != nil {
		// 	err = errors.Wrapf(err, "decoding gas limit - %v failed", paymasterResult.CallGasLimit)
		// 	return nil, err
		// }

		// verificationGas, err = hexutil.DecodeBig(paymasterResult.VerificationGasLimit)
		// if err != nil {
		// 	err = errors.Wrapf(err, "decoding verification gas limit - %v failed", paymasterResult.VerificationGasLimit)
		// 	return nil, err
		// }

		// preVerificationGas, err = hexutil.DecodeBig(paymasterResult.PreVerificationGas)
		// if err != nil {
		// 	err = errors.Wrapf(err, "decoding pre verification gas limit - %v failed", paymasterResult.PreVerificationGas)
		// 	return nil, err
		// }

		o["paymasterAndData"] = paymasterResult.PaymasterAndData
	} else {
		fmt.Println("not using paymaster")
		result, err := b.client.EstimateUserOperationGas(b.EntryPoint, o)
		if err != nil {
			return nil, err
		}
		callGasLimit = result.CallGasLimit
		verificationGas = result.VerificationGasLimit
		preVerificationGas = result.PreVerificationGas
	}

	o["callGasLimit"] = callGasLimit
	o["verificationGasLimit"] = verificationGas
	o["preVerificationGas"] = preVerificationGas

	// TODO: Do I need a merchantId here?
	sig, err := SignUserOp(o, key, SUDO_MODE, nil, chain)
	if err != nil {
		err = errors.Wrap(err, "call to sign user op failed")
		return nil, err
	}

	o["signature"] = hexutil.Encode(sig)
	fmt.Println(o)
	return o, nil
}

func (b *ERCBundler) CreateUnsignedUserOperation(sender, target string, initCode, callData []byte, nonce *big.Int, sponsored bool, chain int64) (map[string]any, error) {
	// var paymasterResult *PaymasterResult
	var paymaster *AlchemyPaymasterResult
	var err error
	var callGasLimit, maxFeePerGas, maxPriorityFeePerGas, verificationGas, preVerificationGas string

	senderAddress := common.HexToAddress(sender)
	tok := make([]byte, 65)
	rand.Read(tok)
	fmt.Println("senderAddress", senderAddress)
	o := map[string]any{
		"sender":   senderAddress.Hex(),
		"nonce":    hexutil.EncodeBig(nonce),
		"initCode": hexutil.Encode(initCode),
		"callData": hexutil.Encode(callData),
		// "callGasLimit":         hexutil.EncodeBig(big.NewInt(0)),
		// "verificationGasLimit": hexutil.EncodeBig(big.NewInt(0)),
		// "preVerificationGas":   hexutil.EncodeBig(big.NewInt(0)),
		// "maxFeePerGas":         hexutil.EncodeBig(getMaxFeePerGas()),
		// "maxPriorityFeePerGas": hexutil.EncodeBig(getMaxPriorityFeePerGas()),
		// "signature":            hexutil.Encode(tok),
		// "paymasterAndData":     "0x",
	}
	fmt.Println(o)

	// paymasterContext := map[string]any{
	// 	"type": "payg",
	// }

	if sponsored {
		policyId := "2e865ced-98e5-4265-a20e-b46c695a28bd"
		paymaster, err = b.client.RequestGasAndPaymasterAndData(policyId, b.EntryPoint, hexutil.Encode(tok), o)
		if err != nil {
			err = errors.Wrap(err, "call to sponsor user op failed")
			return nil, err
		}

		callGasLimit = paymaster.CallGasLimit
		verificationGas = paymaster.VerificationGasLimit
		preVerificationGas = paymaster.PreVerificationGas
		maxPriorityFeePerGas = paymaster.MaxPriorityFeePerGas
		maxFeePerGas = paymaster.MaxFeePerGas
	} else {
		fmt.Println("not using paymaster")
		o := map[string]any{
			"sender":               senderAddress.Hex(),
			"nonce":                hexutil.EncodeBig(nonce),
			"initCode":             hexutil.Encode(initCode),
			"callData":             hexutil.Encode(callData),
			"callGasLimit":         "0xec00",
			"verificationGasLimit": "0x9d5e",
			"preVerificationGas":   "0xab90",
			"maxFeePerGas":         hexutil.EncodeBig(getMaxFeePerGas()),
			"maxPriorityFeePerGas": hexutil.EncodeBig(getMaxPriorityFeePerGas()),
			"signature":            hexutil.Encode(tok),
			"paymasterAndData":     "0x",
		}

		if initCode == nil {
			result, err := b.client.EstimateUserOperationGas(b.EntryPoint, o)
			if err != nil {
				return nil, err
			}
			fmt.Println("user op gas", result)

			callGasLimit = result.CallGasLimit
		} else {
			callGasLimit = "0xec00"
		}

		verificationGas = hexutil.EncodeBig(getVerificationGasLimit())
		preVerificationGas = hexutil.EncodeBig(getPreVerificationGas())

		maxPriorityFeePerGas = hexutil.EncodeBig(getMaxPriorityFeePerGas())
		maxFeePerGas = hexutil.EncodeBig(getMaxFeePerGas())
	}

	o["paymasterAndData"] = "0x"
	o["callGasLimit"] = callGasLimit
	o["verificationGasLimit"] = verificationGas
	o["preVerificationGas"] = preVerificationGas
	o["maxPriorityFeePerGas"] = maxPriorityFeePerGas
	o["maxFeePerGas"] = maxFeePerGas
	o["signature"] = hexutil.Encode(tok)
	// operation, err := userop.New(o)
	// if err != nil {
	// 	return nil, err
	// }
	// entrypoint := GetEntryPointAddress()
	// opHash := operation.GetUserOpHash(entrypoint, big.NewInt(chain))
	// opHash.Hex()
	// sig, err := signUserOp(o, key, chain)
	// if err != nil {
	// 	err = errors.Wrap(err, "call to sign user op failed")
	// 	return nil, err
	// }

	// o["signature"] = hexutil.Encode(sig)
	// fmt.Println(o)
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

func SignUserOp(op map[string]any, key, mode string, merchantId []byte, chain int64) ([]byte, error) {
	chainId := big.NewInt(chain)
	entrypoint := GetEntryPointAddress()

	operation, err := userop.New(op)
	if err != nil {
		return nil, err
	}

	opHash := operation.GetUserOpHash(entrypoint, chainId)

	fmt.Println("userop hash - ", opHash)
	fmt.Println("userop hash bytes - ", opHash.Bytes())
	// privKey, err := getSigningKey(key)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	kb, _ := hexutil.Decode(key)
	privKey, err := crypto.ToECDSA(kb)
	if err != nil {
		fmt.Println("key err", err)
		return nil, err
	}
	// Kernel has a specific convention for encoding signatures in order to determing the mode see (https://github.com/stackup-wallet/userop.js/blob/main/src/preset/builder/kernel.ts#L114-L123)
	sig, _ := hexutil.Decode(mode)

	if merchantId != nil {
		sig = append(sig, merchantId...)
	}

	signatureBytes, err := crypto.Sign(opHash[:], privKey)
	if err != nil {
		err = errors.Wrap(err, "signUserOp() failure - ")
		return nil, err
	}
	signatureBytes[64] += 27
	sig = append(sig, signatureBytes...)
	signature := hexutil.Encode(sig)

	fmt.Println("signature length - ", len(signatureBytes))
	fmt.Println("signature - ", hexutil.Encode(signatureBytes))
	fmt.Println("signature - ", signature)
	return sig, nil
}

/*
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
		fmt.Println("call data ", hexutil.Encode(callData))
		return callData, nil
	}

	erc20Token := getErc20TokenABI()

	erc20TransferData, err := GetTransferFnData(erc20Token, toAddress, amount)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create erc20 call data")
		return nil, err
	}

	tokenAddress := getTokenAddres(token)
	callData, err := GetExecuteFnData(accountABI, tokenAddress, common.Big0, erc20TransferData)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

func CreateSetExecutionCallData(enableData []byte) ([]byte, error) {
	kernelAbi := getKernelStorageAbi()
	lucidValidator := os.Getenv("LUCID_VALIDATOR")
	callData, err := GetSetExecutionFnData(kernelAbi, lucidValidator, enableData)
	if err != nil {
		err = errors.Wrap(err, "CreateSetExecutionCallData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

func CreateFactoryFnData(enableData []byte) ([]byte, error) {
	factoryAbi := getAccountFactoryAbi()
	defaultValidatorAddress := os.Getenv("DEFAULT_VALIDATOR")
	callData, err := GetCreateAccountFnData(factoryAbi, defaultValidatorAddress, enableData)
	if err != nil {
		err = errors.Wrap(err, "CreateFactoryFnData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

func getTokenAddres(token string) string {
	return "0x0fa8781a83e46826621b3bc094ea2a0212e71b23"
}
