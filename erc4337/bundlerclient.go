package erc4337

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type Client struct {
	ctx  context.Context
	c, p *ethclient.Client
}

type GasEstimateResult struct {
	PreVerificationGas *big.Int `json:"PreVerificationGas"`
	VerificationGas    *big.Int `json:"VerificationGas"`
	CallGasLimit       *big.Int `json:"CallGasLimit"`
}

// func initialiseBundler() (*useroperation.ERCBundler, error) {
// 	rpc := os.Getenv("NODE_URL")
// 	paymaster := os.Getenv("PAYMASTER_URL")
// 	entryPoint := os.Getenv("ENTRY_POINT")

// 	node, err := Dial(rpc, paymaster)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// time.DateOnly
// 	bundler := useroperation.NewERCBundler(entryPoint, node)
// 	if bundler == nil {
// 		return nil, errors.New("bundler was not initialised")
// 	}

// 	return bundler, nil
// }

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

func (nc *Client) GetAccountNonce(address common.Address) (uint64, error) {
	return nc.c.NonceAt(nc.ctx, address, nil)
}

// SendUserOperation sends a user operation to an alt mempool
func (nc *Client) SendUserOperation(entryPoint string, userop interface{}) (string, error) {
	var result string

	err := nc.c.Client().CallContext(nc.ctx, &result, "eth_sendUserOperation", userop, entryPoint)
	if err != nil {
		err = errors.Wrap(err, "eth_sendUserOperation call error")
		return result, err
	}

	return result, nil
}

// Estimate the gas parameters required for the user operation
func (nc *Client) EstimateUserOperationGas(entrypointAddress string, userop interface{}) (*GasEstimateResult, error) {
	result := &GasEstimateResult{}

	err := nc.c.Client().CallContext(nc.ctx, result, "eth_estimateUserOperationGas", userop, entrypointAddress)
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
func (nc *Client) GetUserOperationByHash(userophash string) error {
	var result map[string]interface{}

	err := nc.c.Client().CallContext(nc.ctx, &result, "eth_getUserOperationByHash", userophash)
	if err != nil {
		err = errors.Wrap(err, "eth_getUserOperationByHash call error")
		return err
	}

	fmt.Println("user operation - ", result)
	return nil
}

// eth_getUserOperationReciept
