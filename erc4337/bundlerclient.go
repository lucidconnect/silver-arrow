package erc4337

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	KernelStorage "github.com/helicarrierstudio/silver-arrow/abi/kernelStorage"
	"github.com/pkg/errors"
)

type Client struct {
	ctx  context.Context
	c, p *ethclient.Client
}

type GasEstimateResult struct {
	PreVerificationGas   int `json:"preVerificationGas"`
	VerificationGasLimit int `json:"verificationGasLimit"`
	CallGasLimit         int `json:"callGasLimit"`
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

func InitialiseBundler() (*ERCBundler, error) {
	rpc := os.Getenv("NODE_URL")
	paymaster := os.Getenv("PAYMASTER_URL")
	entryPoint := os.Getenv("ENTRY_POINT")

	node, err := Dial(rpc, paymaster)
	if err != nil {
		return nil, err
	}
	// time.DateOnly
	bundler := NewERCBundler(entryPoint, node)
	if bundler == nil {
		return nil, errors.New("bundler was not initialised")
	}

	return bundler, nil
}

func Dial(url, paymasterUrl string) (*Client, error) {
	var client, paymasterClient *ethclient.Client
	ctx := context.Background()

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("Oops! Something went wrong ", err)
		return nil, err
	}

	if paymasterUrl != "" {
		fmt.Println("initialising with a paymaster")
		paymasterClient, err = ethclient.Dial(paymasterUrl)
		if err != nil {
			log.Fatal("Oops! Something went wrong ", err)
			return nil, err
		}
	}

	return newClient(ctx, client, paymasterClient), nil
}

func newClient(ctx context.Context, client, paymasterClient *ethclient.Client) *Client {
	return &Client{
		ctx: ctx,
		c:   client,
		p:   paymasterClient,
	}
}

func (nc *Client) GetEthClient() *ethclient.Client {
	return nc.c
}

// GetAccountCode is a wrapper around ethclient.CodeAt()
func (nc *Client) GetAccountCode(address common.Address) ([]byte, error) {
	return nc.c.CodeAt(nc.ctx, address, nil)
}

func (nc *Client) GetAccountNonce(address common.Address) (*big.Int, error) {
	k, err := KernelStorage.NewKernelStorage(address, nc.c)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{
		Pending: true,
		Context: nil,
	}

	return k.GetNonce0(opts)
}

// SendUserOperation sends a user operation to an alt mempool
func (nc *Client) SendUserOperation(entryPoint string, userop map[string]any) (string, error) {
	var result string
	// fmt.Println("user op", userop)
	// obj, _ := json.Marshal(userop)

	sender, _ := userop["sender"].(string)
	nonce, _ := userop["nonce"].(string)
	initCode, _ := userop["initCode"].(string)
	callGasLimit, _ := userop["callGasLimit"].(string)
	verificationGasLimit, _ := userop["verificationGasLimit"].(string)
	preVerificationGas, ok := userop["preVerificationGas"].(string)
	if !ok {
		log.Panic(userop["preVerificationGas"])
	}

	maxFeePerGas, _ := userop["maxFeePerGas"].(string)
	maxPriorityFeePerGas, _ := userop["maxPriorityFeePerGas"].(string)
	paymasterAndData, _ := userop["paymasterAndData"].(string)
	signature, _ := userop["signature"].(string)
	callData, _ := userop["callData"].(string)

	// for k, v := range userop {
	// 	v.(string)
	// }

	request := UserOperation{
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
	fmt.Println(request)
	err := nc.c.Client().CallContext(nc.ctx, &result, "eth_sendUserOperation", request, entryPoint)
	if err != nil {
		err = errors.Wrap(err, "eth_sendUserOperation call error")
		return result, err
	}

	return result, nil
}

// Estimate the gas parameters required for the user operation
func (nc *Client) EstimateUserOperationGas(entrypointAddress string, userop map[string]any) (*GasEstimateResult, error) {
	result := &GasEstimateResult{}

	sender, _ := userop["sender"].(string)
	nonce, _ := userop["nonce"].(string)
	initCode, _ := userop["initCode"].(string)
	// callGasLimit, _ := userop["callGasLimit"].(string)
	// verificationGasLimit, _ := userop["verificationGasLimit"].(string)
	// preVerificationGas, ok := userop["preVerificationGas"].(string)
	// if !ok {
	// 	log.Panic(userop["preVerificationGas"])
	// }

	// maxFeePerGas, _ := userop["maxFeePerGas"].(string)
	// maxPriorityFeePerGas, _ := userop["maxPriorityFeePerGas"].(string)
	// paymasterAndData, _ := userop["paymasterAndData"].(string)
	signature, _ := userop["signature"].(string)
	callData, _ := userop["callData"].(string)

	request := UserOperation{
		Sender:               sender,
		Nonce:                nonce,
		InitCode:             initCode,
		// CallGasLimit:         callGasLimit,
		// VerificationGasLimit: verificationGasLimit,
		// PreVerificationGas:   preVerificationGas,
		// MaxFeePerGas:         maxFeePerGas,
		// MaxPriorityFeePerGas: maxPriorityFeePerGas,
		// PaymasterAndData:     paymasterAndData,
		Signature:            signature,
		CallData:             callData,
	}
	err := nc.c.Client().CallContext(nc.ctx, result, "eth_estimateUserOperationGas", request, entrypointAddress)
	if err != nil {
		err = errors.Wrap(err, "eth_estimateUserOperationGas call error")
		return nil, err
	}

	return result, nil
}

// Estimate the gas parameters required for the user operation
func (nc *Client) GetBalance(address string) (*big.Int, error) {
	return nc.c.BalanceAt(nc.ctx, common.HexToAddress(address), nil)
}

// GetUserOperationByHash calls eth_getUserOperationByHash bundler rpc method
func (nc *Client) GetUserOperationByHash(userophash string) (map[string]any, error) {
	var result map[string]any

	err := nc.c.Client().CallContext(nc.ctx, &result, "eth_getUserOperationByHash", userophash)
	if err != nil {
		err = errors.Wrap(err, "eth_getUserOperationByHash call error")
		return nil, err
	}

	fmt.Println("user operation - ", result)
	return result, nil
}

// eth_getUserOperationReciept
