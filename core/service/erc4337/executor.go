package erc4337

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lucidconnect/silver-arrow/abi/tokenActions"
)

func TransferErc20Action(token, to common.Address, amount *big.Int) ([]byte, error) {
	tokenActionAbi, err := tokenActions.TokenActionsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	callData, err := tokenActionAbi.Pack("transfer20Action", token, amount, to)
	if err != nil {
		return nil, err
	}

	return callData, nil
}
