package erc4337

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog/log"
)

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

func (nc *Client) SponsorUserOperation(entryPoint string, userop, pc interface{}) (*PaymasterResult, error) {
	result := &PaymasterResult{}

	err := nc.p.Client().CallContext(nc.ctx, result, "pm_sponsorUserOperation", userop, entryPoint, pc)
	if err != nil {
		log.Err(err).Msg("pm_sponsorUserOperation")
		return nil, err
	}

	fmt.Println("pm_sponsorUserOperation - ", result)
	return result, nil
}

// paymaster data request for alchemy 2e865ced-98e5-4265-a20e-b46c695a28bd
func (nc *Client) RequestGasAndPaymasterAndData(policyId, entryPoint, dummySignature string, userop any) (*AlchemyPaymasterResult, error) {
	result := &AlchemyPaymasterResult{}

	feeOverride := map[string]string{
		"maxFeePerGas":         "0x29260CA6A",
		"maxPriorityFeePerGas": "0x29260CA6A",
	}
	request := AlchemyPaymasterRequest{
		PolicyId:       policyId,
		EntryPoint:     entryPoint,
		DummySignature: dummySignature,
		UserOperation:  userop,
		FeeOverride:    feeOverride,
	}

	time.Sleep(3 * time.Second)
	err := nc.p.Client().CallContext(nc.ctx, result, "alchemy_requestGasAndPaymasterAndData", request)
	if err != nil {
		log.Err(err).Msg("alchemy_requestGasAndPaymasterAndData")
		return nil, err
	}

	// add 130098856
	var maxFee, maxPriorityFee *big.Int
	maxFeex := new(big.Int).Mul(new(big.Int).SetBytes(hexutil.MustDecode(result.MaxFeePerGas)), big.NewInt(10))
	maxFee = new(big.Int).Div(maxFeex, big.NewInt(7))

	maxPriorityFeex := new(big.Int).Mul(new(big.Int).SetBytes(hexutil.MustDecode(result.MaxPriorityFeePerGas)), big.NewInt(10))
	maxPriorityFee = new(big.Int).Div(maxPriorityFeex, big.NewInt(7))
	
	// maxFee := new(big.Int).SetBytes(hexutil.MustDecode(result.MaxFeePerGas)).Add(big.NewInt(130098856))
	// maxPriorityFee := new(big.Int).SetBytes(hexutil.MustDecode(result.MaxPriorityFeePerGas)).Add(big.NewInt(130098856))

	result.MaxFeePerGas = hexutil.EncodeBig(maxFee)
	result.MaxPriorityFeePerGas = hexutil.EncodeBig(maxPriorityFee)
	fmt.Println("alchemy_requestGasAndPaymasterAndData - ", result)
	return result, nil
}
