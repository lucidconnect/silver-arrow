package erc4337

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lucidconnect/silver-arrow/abi/EntryPoint"
	"github.com/lucidconnect/silver-arrow/abi/erc20"
	"github.com/pkg/errors"
)

type Client struct {
	ctx  context.Context
	c, p *ethclient.Client
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

func InitialiseBundler(chain int64) (*ERCBundler, error) {
	network, err := GetNetwork(chain)
	if err != nil {
		return nil, err
	}
	rpc := os.Getenv(fmt.Sprintf("%s_NODE_URL", network))
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
		log.Fatal().Err(err).Msg("Oops! Something went wrong ")
		return nil, err
	}

	if paymasterUrl != "" {
		fmt.Println("initialising with a paymaster")
		paymasterClient, err = ethclient.Dial(paymasterUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("Oops! Something went wrong ")
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

func (nc *Client) GetAccountNonce(entryPoint, address common.Address) (*big.Int, error) {
	e, err := EntryPoint.NewEntryPoint(entryPoint, nc.c)
	if err != nil {
		err = errors.Wrap(err, "error initialising entrypoint instance")
		return nil, err
	}
	opts := &bind.CallOpts{
		Pending: false,
		Context: nc.ctx,
		From:    address,
	}

	nonce, err := e.GetNonce(opts, address, common.Big0)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func (nc *Client) GetErc20TokenBalance(tokenAddress, walletAddress common.Address) (*big.Int, error) {
	ercToken, err := erc20.NewErc20(tokenAddress, nc.c)
	if err != nil {
		err = errors.Wrap(err, "error initialising erc20 token instance")
		return nil, err
	}

	opts := &bind.CallOpts{
		Pending: true,
		Context: nc.ctx,
	}

	balance, err := ercToken.BalanceOf(opts, walletAddress)
	if err != nil {
		err = errors.Wrapf(err, "error occured fetchin erc20 token balance for wallet at %v", walletAddress.Hex())
		return nil, err
	}

	return balance, nil
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
	preVerificationGas, _ := userop["preVerificationGas"].(string)

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

	// add a 3s delay (not ideal should implement a more elegant way to prevent rate limiting)

	time.Sleep(3 * time.Second)
	// fmt.Println("payload",request)
	err := nc.c.Client().CallContext(nc.ctx, &result, "eth_sendUserOperation", request, entryPoint)
	if err != nil {
		if err.Error() == "AA13 initCode failed or OOG" {
			// increase verification gas limit

		}
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
	callGasLimit, _ := userop["callGasLimit"].(string)
	verificationGasLimit, _ := userop["verificationGasLimit"].(string)
	preVerificationGas, _ := userop["preVerificationGas"].(string)

	maxFeePerGas, _ := userop["maxFeePerGas"].(string)
	maxPriorityFeePerGas, _ := userop["maxPriorityFeePerGas"].(string)
	paymasterAndData, _ := userop["paymasterAndData"].(string)
	signature, _ := userop["signature"].(string)
	callData, _ := userop["callData"].(string)

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
	err := nc.c.Client().CallContext(nc.ctx, result, "eth_estimateUserOperationGas", request, entrypointAddress)
	if err != nil {
		err = errors.Wrap(err, "eth_estimateUserOperationGas call error")
		return nil, err
	}

	fmt.Println("eth_estimateUserOperationGas - ", result)
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

	// fmt.Println("user operation - ", result)
	return result, nil
}

// eth_getUserOperationReciept

func GetNetwork(chainId int64) (string, error) {
	switch chainId {
	case 1:
		return ETHEREUM, nil
	case 5:
		return GOERLI, nil
	case 137:
		return POLYGON, nil
	case 84531:
		return BASE_GOERLI, nil
	case 8453:
		return BASE, nil
	case 80001:
		return MUMBAI, nil
	default:
		return "NOT SUPPORTED", errors.New("Unsupported chain")
	}
}
