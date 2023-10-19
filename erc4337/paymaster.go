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

	// feeOverride := map[string]string{
	// 	"maxFeePerGas":         "0x29260CA6A",
	// 	"maxPriorityFeePerGas": "0x29260CA6A",
	// }
	request := AlchemyPaymasterRequest{
		PolicyId:       policyId,
		EntryPoint:     entryPoint,
		DummySignature: dummySignature,
		UserOperation:  userop,
		// FeeOverride:    feeOverride,
	}

	time.Sleep(3 * time.Second)
	err := nc.p.Client().CallContext(nc.ctx, result, "alchemy_requestGasAndPaymasterAndData", request)
	if err != nil {
		log.Err(err).Msg("alchemy_requestGasAndPaymasterAndData")
		return nil, err
	}

	// add 130098856
	// var maxFee, maxPriorityFee *big.Int
	// maxFeePerGas, err := hexutil.DecodeUint64(result.MaxFeePerGas)
	// if err != nil {
	// 	log.Err(err).Msgf("maxFeePerGas returned from Alchemy %v", result.MaxFeePerGas)
	// 	return result, err
	// }
	// maxFeePerGas = (maxFeePerGas * 2) / 7

	// maxPriorityFee, err := hexutil.DecodeUint64(result.MaxPriorityFeePerGas)
	// if err != nil {
	// 	log.Err(err).Msgf("maxFeePerGas returned from Alchemy %v", result.MaxFeePerGas)
	// 	return result, err
	// }
	// maxPriorityFee = (maxPriorityFee * 2) / 7

	// result.MaxFeePerGas = hexutil.EncodeUint64(maxFeePerGas)
	// result.MaxPriorityFeePerGas = hexutil.EncodeUint64(maxPriorityFee)
	fmt.Println("alchemy_requestGasAndPaymasterAndData - ", result)
	return result, nil
}

func (nc *Client) RequestPaymasterAndData(policyId, entryPoint, dummySignature string, userop any) (*AlchemyPaymasterResult, error) {
	result := &AlchemyPaymasterResult{}

	// feeOverride := map[string]string{
	// 	"maxFeePerGas":         "0x29260CA6A",
	// 	"maxPriorityFeePerGas": "0x29260CA6A",
	// }
	request := AlchemyPaymasterRequest{
		PolicyId:       policyId,
		EntryPoint:     entryPoint,
		DummySignature: dummySignature,
		UserOperation:  userop,
		// FeeOverride:    feeOverride,
	}

	time.Sleep(3 * time.Second)
	err := nc.p.Client().CallContext(nc.ctx, result, "alchemy_requestPaymasterAndData", request)
	if err != nil {
		log.Err(err).Msg("alchemy_requestPaymasterAndData")
		return nil, err
	}

	// add 130098856
	// var maxFee, maxPriorityFee *big.Int

	maxFeePerGasBig, _ := new(big.Int).SetString(result.MaxFeePerGas, 0)
	_maxFeePerGas := new(big.Int).Mul(maxFeePerGasBig, big.NewInt(10))
	maxFeePerGas := new(big.Int).Div(_maxFeePerGas, big.NewInt(7))

	maxPriorityFeeBig, _ := new(big.Int).SetString(result.MaxPriorityFeePerGas, 0)
	_maxPriorityFee := new(big.Int).Mul(maxPriorityFeeBig, big.NewInt(10))
	maxPriorityFee := new(big.Int).Div(_maxPriorityFee, big.NewInt(7))

	result.MaxFeePerGas = hexutil.EncodeBig(maxFeePerGas)
	result.MaxPriorityFeePerGas = hexutil.EncodeBig(maxPriorityFee)
	fmt.Println("alchemy_PaymasterAndData - ", result)
	return result, nil
}
