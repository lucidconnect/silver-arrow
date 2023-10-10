package erc4337

import (
	"fmt"

	"github.com/pkg/errors"
)

func (nc *Client) GetMaxPriorityFee() (string, error) {
	var result string

	err := nc.c.Client().CallContext(nc.ctx, result, "rundler_maxPriorityFeePerGas")
	if err != nil {
		err = errors.Wrap(err, "rundler_maxPriorityFeePerGas call error")
		return "", err
	}

	fmt.Println("rundler_maxPriorityFeePerGas - ", result)
	return result, nil
}
