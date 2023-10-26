package erc4337

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/lucidconnect/silver-arrow/erc20"
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
		policyId := os.Getenv("POLICY_ID")
		paymaster, err = b.client.RequestGasAndPaymasterAndData(policyId, b.EntryPoint, hexutil.Encode(tok), o)
		if err != nil {
			err = errors.Wrap(err, "call to sponsor user op failed")
			return nil, err
		}
		callGasLimit = paymaster.CallGasLimit
		verificationGas = paymaster.VerificationGasLimit
		maxPriorityFeePerGas = paymaster.MaxPriorityFeePerGas
		preVerificationGas = paymaster.PreVerificationGas
		maxFeePerGas = paymaster.MaxFeePerGas
		paymasterAndData = paymaster.PaymasterAndData
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
		maxPriorityFeePerGas, err = b.client.GetMaxPriorityFee()
		if err != nil {
			return nil, err
		}
		// maxPriorityFeePerGas = hexutil.EncodeBig(getMaxPriorityFeePerGas())
		maxFeePerGas = calcMaxFeePerGas(maxPriorityFeePerGas)
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

func calcMaxFeePerGas(maxPriorityFee string) string {
	maxPriorityFeeBig, _ := new(big.Int).SetString(maxPriorityFee, 0)
	maxFeePerGasBig := new(big.Int).Add(maxPriorityFeeBig, big.NewInt(22))

	return hexutil.EncodeBig(maxFeePerGasBig)
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
		// log.Err(err).Msg("SignUserOp() - ")
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

	tokenAddress := erc20.GetTokenAddress(token, 80001)
	callData, err := GetExecuteFnData(accountABI, tokenAddress, common.Big0, erc20TransferData)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}
