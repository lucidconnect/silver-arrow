package erc4337

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (nc *Client) GetMaxPriorityFee() (string, error) {
	var result string

	err := nc.c.Client().CallContext(nc.ctx, &result, "rundler_maxPriorityFeePerGas")
	if err != nil {
		err = errors.Wrap(err, "rundler_maxPriorityFeePerGas call error")
		return "", err
	}

	log.Debug().Msgf("rundler_maxPriorityFeePerGas - %v", result)
	maxPriorityFeeBig, _ := new(big.Int).SetString(result, 0)
	maxPriorityFee := new(big.Int).Mul(maxPriorityFeeBig, big.NewInt(2))
	log.Debug().Msgf("adjusted maxPriorityFeePerGas - %v", maxPriorityFee)
	return hexutil.EncodeBig(maxPriorityFee), nil
}
