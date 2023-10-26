package erc4337

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/lucidconnect/silver-arrow/erc20"
	"github.com/pkg/errors"
	"github.com/rmanzoku/ethutils/ecrecover"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

func CreateTransferCallData(toAddress, token string, amount *big.Int) ([]byte, error) {
	accountABI := getAccountABI()

	if token == "ETH" {
		callData, err := GetExecuteFnData(accountABI, toAddress, amount, nil)
		if err != nil {
			err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
			return nil, err
		}
		return callData, nil
	}

	erc20Token := getErc20TokenABI()

	erc20TransferData, err := GetTransferFnData(erc20Token, toAddress, amount)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create erc20 call data")
		return nil, err
	}

	tokenAddress := erc20.GetTokenAddress(token, 80001)
	callData, err := GetExecuteFnData(accountABI, tokenAddress, common.Big0, erc20TransferData)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData(): failed to create final call data")
		return nil, err
	}

	return callData, nil
}

func SignUserOp(op map[string]any, key, mode string, merchantId []byte, chain int64) ([]byte, []byte, error) {
	chainId := big.NewInt(chain)
	entrypoint := GetEntryPointAddress()

	operation, err := userop.New(op)
	if err != nil {
		return nil, nil, err
	}

	opHash := operation.GetUserOpHash(entrypoint, chainId)
	hash := opHash.Bytes()

	kb, err := hexutil.Decode(key)
	if err != nil {
		err = errors.Wrap(err, "SignUserOp() - ")
		return nil, nil, err
	}
	// fmt.Println("kb", kb)
	// common.FromHex(key)

	// Kernel has a specific convention for encoding signatures in order to determing the mode see (https://github.com/stackup-wallet/userop.js/blob/main/src/preset/builder/kernel.ts#L114-L123)
	signature, _ := hexutil.Decode(mode)

	// if merchantId != nil {
	// 	fmt.Println("merchant id - ", merchantId)
	// 	signature = append(signature, merchantId...)
	// }

	sig, err := secp256k1.Sign(ecrecover.ToEthSignedMessageHash(hash), kb)
	if err != nil {
		err = errors.Wrap(err, "generating signature failed.")
		return nil, nil, err
	}
	sig[64] += 27
	signature = append(signature, sig...)

	fmt.Println("raw sig - ", hexutil.Encode(sig))
	fmt.Println("hash - ", hexutil.Encode(opHash[:]))

	return signature, ecrecover.ToEthSignedMessageHash(hash), nil
}
