package erc4337

import (
	"fmt"
	"log"
)

type PaymasterResult struct {
	PaymasterAndData     string `json:"paymasterAndData"`
	PreVerificationGas   string `json:"preVerificationGas"`
	VerificationGasLimit string `json:"verificationGasLimit"`
	CallGasLimit         string `json:"callGasLimit"`
}

func (nc *Client) SponsorUserOperation(entryPoint string, userop, pc interface{}) (*PaymasterResult, error) {
	result := &PaymasterResult{}

	err := nc.p.Client().CallContext(nc.ctx, result, "pm_sponsorUserOperation", userop, entryPoint)
	if err != nil {
		log.Printf("pm_sponsorUserOperation -  message: %v \n", err)
		return nil, err
	}

	fmt.Println("pm_sponsorUserOperation - ", result)
	return result, nil
}
