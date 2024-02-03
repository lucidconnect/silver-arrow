package erc4337

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lucidconnect/silver-arrow/abi/lucidTokenActions"
	"github.com/lucidconnect/silver-arrow/abi/tokenActions"
)

type DebitInstruction struct {
	Token,
	Destination common.Address
	Amount *big.Int
}

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

func BatchTransferErc20Action(instructions []DebitInstruction) ([]byte, error) {
	var inputs []lucidTokenActions.TransferInput
	tokenActionAbi, err := lucidTokenActions.LucidTokenActionsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	for _, instruction := range instructions {
		in := lucidTokenActions.TransferInput{
			Token:  instruction.Token,
			To:     instruction.Destination,
			Amount: instruction.Amount,
		}

		inputs = append(inputs, in)
	}
	callData, err := tokenActionAbi.Pack("batchTransfer20Action", inputs)
	if err != nil {
		return nil, err
	}

	return callData, nil
}
