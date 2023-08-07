package erc4337

import (
	"fmt"
	"log"
	"time"
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
		log.Printf("pm_sponsorUserOperation -  message: %v \n", err)
		return nil, err
	}

	fmt.Println("pm_sponsorUserOperation - ", result)
	return result, nil
}

// paymaster data request for alchemy 2e865ced-98e5-4265-a20e-b46c695a28bd
func (nc *Client) RequestGasAndPaymasterAndData(policyId, entryPoint, dummySignature string, userop any) (*AlchemyPaymasterResult, error) {
	result := &AlchemyPaymasterResult{}

	request := AlchemyPaymasterRequest{
		PolicyId:       policyId,
		EntryPoint:     entryPoint,
		DummySignature: dummySignature,
		UserOperation:  userop,
	}

	time.Sleep(3 * time.Second)
	err := nc.p.Client().CallContext(nc.ctx, result, "alchemy_requestGasAndPaymasterAndData", request)
	if err != nil {
		log.Printf("alchemy_requestGasAndPaymasterAndData -  message: %v \n", err)
		return nil, err
	}

	fmt.Println("alchemy_requestGasAndPaymasterAndData - ", result)
	return result, nil
}
