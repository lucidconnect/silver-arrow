package erc4337

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
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

func (b *ERCBundler) AccountNonce(sender string) uint64 {
	senderAddress := common.HexToAddress(sender)
	nonce, _ := b.client.GetAccountNonce(senderAddress)
	fmt.Println("nonce:", nonce)
	return nonce
}

// UserOperation represents an EIP-4337 style transaction for a smart contract account.
// CreateUserOperation returns a signed useroperation
func (b *ERCBundler) CreateUserOperation(sender, target, token string, callData []byte, nonce, amount *big.Int, sponsored bool, key string, chain int64) (map[string]any, error) {
	var paymasterResult *PaymasterResult
	var err error
	var callGasLimit, verificationGas, preVerificationGas *big.Int

	senderAddress := common.HexToAddress(sender)
	tok := make([]byte, 65)
	rand.Read(tok)

	o := map[string]any{
		"sender":               senderAddress,
		"nonce":                nonce,
		"initCode":             "0x",
		"callData":             hexutil.Encode(callData),
		"callGasLimit":         big.NewInt(0),
		"verificationGasLimit": big.NewInt(0),
		"preVerificationGas":   big.NewInt(0),
		"maxFeePerGas":         getMaxFeePerGas(),
		"maxPriorityFeePerGas": getMaxPriorityFeePerGas(),
		"signature":            hexutil.Encode(tok),
		"paymasterAndData":     "0x",
	}

	paymasterContext := map[string]any{
		"type": "payg",
	}

	if sponsored {
		paymasterResult, err = b.client.SponsorUserOperation(b.EntryPoint, o, paymasterContext)
		if err != nil {
			err = errors.Wrap(err, "call to sponsor user op failed")
			return nil, err
		}

		callGasLimit, err = hexutil.DecodeBig(paymasterResult.CallGasLimit)
		if err != nil {
			err = errors.Wrapf(err, "decoding gas limit - %v failed", paymasterResult.CallGasLimit)
			return nil, err
		}

		verificationGas, err = hexutil.DecodeBig(paymasterResult.VerificationGasLimit)
		if err != nil {
			err = errors.Wrapf(err, "decoding verification gas limit - %v failed", paymasterResult.VerificationGasLimit)
			return nil, err
		}

		preVerificationGas, err = hexutil.DecodeBig(paymasterResult.PreVerificationGas)
		if err != nil {
			err = errors.Wrapf(err, "decoding pre verification gas limit - %v failed", paymasterResult.PreVerificationGas)
			return nil, err
		}

		o["paymasterAndData"] = paymasterResult.PaymasterAndData
	} else {
		fmt.Println("not using paymaster")
		result, err := b.client.EstimateUserOperationGas(b.EntryPoint, o)
		if err != nil {
			return nil, err
		}
		callGasLimit = result.CallGasLimit
		verificationGas = result.VerificationGas
		preVerificationGas = result.PreVerificationGas
	}

	o["callGasLimit"] = callGasLimit
	o["verificationGasLimit"] = verificationGas
	o["preVerificationGas"] = preVerificationGas

	sig, err := signUserOp(o, key, chain)
	if err != nil {
		err = errors.Wrap(err, "call to sign user op failed")
		return nil, err
	}

	o["signature"] = hexutil.Encode(sig)
	fmt.Println(o)
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

func (b *ERCBundler) GetUserOp(userophash string) error {
	return b.client.GetUserOperationByHash(userophash)
}

func signUserOp(op map[string]any, key string, chain int64) ([]byte, error) {
	chainId := big.NewInt(chain)
	entrypoint := getEntryPointAddress()

	operation, err := userop.New(op)
	if err != nil {
		return nil, err
	}

	opHash := operation.GetUserOpHash(entrypoint, chainId)

	fmt.Println("userop hash - ", opHash)
	fmt.Println("userop hash bytes - ", opHash.Bytes())
	privKey, err := getSigningKey(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Kernel has a specific convention for encoding signatures in order to determing the mode see (https://github.com/stackup-wallet/userop.js/blob/main/src/preset/builder/kernel.ts#L114-L123)
	sig, _ := hexutil.Decode("0x00000000")

	signatureBytes, err := crypto.Sign(opHash[:], privKey)
	if err != nil {
		err = errors.Wrap(err, "signUserOp() failure - ")
		return nil, err
	}
	signatureBytes[64] += 27
	sig = append(sig, signatureBytes...)
	signature := hexutil.Encode(sig)

	fmt.Println("signature - ", sig)
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

	callData, err := GetExecuteFnData(accountABI, "0x0000000000000000000000000000000000001010", common.Big0, erc20TransferData)
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
	// op := vm.CALL
	payload, err := contractABI.Pack("execute", dest, amount, callData, uint8(0))
	if err != nil {
		err = errors.Wrap(err, "PacK() unable to prepare tx payload")
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

func getEntryPointAddress() common.Address {
	entrypointAddress := "0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"
	return common.HexToAddress(entrypointAddress)
}

func (b *ERCBundler) getContractInitCode(address common.Address) []byte {
	code, err := b.client.GetAccountCode(address)
	if err != nil {
		log.Println("Oops! could not fetch account code", err)
		return nil
	}
	fmt.Println("Account code - ", code)
	return code
}

func getCallGasLimit() *big.Int {
	return big.NewInt(60000)
}

func getVerificationGasLimit() *big.Int {
	return big.NewInt(60624)
}

func getPreVerificationGas() *big.Int {
	return big.NewInt(59925)
}

func getMaxFeePerGas() *big.Int {
	return big.NewInt(2400000018)
}

func getMaxPriorityFeePerGas() *big.Int {
	return big.NewInt(2400000018)
}

func getSigningKey(privateKey string) (*ecdsa.PrivateKey, error) {
	privKey, err := crypto.HexToECDSA(privateKey[2:])
	if err != nil {
		err = errors.Wrap(err, "private key parse failure")
		return nil, err
	}
	return privKey, nil
}
