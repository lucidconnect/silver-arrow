package erc4337

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lucidconnect/silver-arrow/abi/EntryPoint"
	"github.com/lucidconnect/silver-arrow/abi/erc20"
	"github.com/pkg/errors"

	// "github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type AlchemyService struct {
	EntryPoint string

	ctx     context.Context
	backend *ethclient.Client
}

func NewAlchemyService(chain int64) (*AlchemyService, error) {
	return initialiseAlchemyService(chain)
}

func initialiseAlchemyService(chain int64) (*AlchemyService, error) {
	ctx := context.Background()
	network, err := GetNetwork(chain)
	if err != nil {
		log.Err(err).Msgf("chain %v not supported", chain)
		return nil, err
	}

	rpc := os.Getenv(fmt.Sprintf("%s_NODE_URL", network))
	entryPoint := os.Getenv("ENTRY_POINT")

	node, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	return &AlchemyService{entryPoint, ctx, node}, nil
}

func (bs *AlchemyService) GetEthBackend() *ethclient.Client {
	return bs.backend
}

func (bs *AlchemyService) GetAccountNonce(address common.Address) (*big.Int, error) {
	entryPoint := common.HexToAddress(bs.EntryPoint)
	e, err := EntryPoint.NewEntryPoint(entryPoint, bs.backend)
	if err != nil {

		err = errors.Wrap(err, "error initialising entrypoint instance")
		return nil, err
	}
	opts := &bind.CallOpts{
		Pending: false,
		Context: bs.ctx,
		From:    address,
	}

	nonce, err := e.GetNonce(opts, address, common.Big0)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

// GetAccountCode is a wrapper around ethclient.CodeAt()
func (bs *AlchemyService) GetAccountCode(address common.Address) ([]byte, error) {
	return bs.backend.CodeAt(bs.ctx, address, nil)
}

func (bs *AlchemyService) GetNativeTokenBalance(walletAddress common.Address) (*big.Int, error) {
	return bs.backend.BalanceAt(bs.ctx, walletAddress, nil)
}

func (bs *AlchemyService) GetErc20TokenBalance(tokenAddress, walletAddress common.Address) (*big.Int, error) {
	ercToken, err := erc20.NewErc20(tokenAddress, bs.backend)
	if err != nil {
		err = errors.Wrap(err, "error initialising erc20 token instance")
		return nil, err
	}

	opts := &bind.CallOpts{
		Pending: true,
		Context: bs.ctx,
	}

	balance, err := ercToken.BalanceOf(opts, walletAddress)
	if err != nil {
		err = errors.Wrapf(err, "error occured fetchin erc20 token balance for wallet at %v", walletAddress.Hex())
		return nil, err
	}

	return balance, nil
}

// erc-4337 specific calls

func (a *AlchemyService) CreateUnsignedUserOperation(sender string, initCode, callData []byte, nonce *big.Int, sponsored bool, chain int64) (map[string]any, error) {
	var paymaster *AlchemyPaymasterResult
	var gasEstimate *GasEstimateResult
	var err error
	var paymasterAndData, maxFeePerGas, maxPriorityFeePerGas string
	tok := make([]byte, 65)
	rand.Read(tok)
	dummySignature := "0x00000001d16b1fd2d2703b7214bee4a66979f386e1f1af9cd48629d3c5a436b567941cc43125a5d744bae0fe0fbf32144f6c2f9ccdc25c6e28d464953297cea608c748cc1c"
	o := map[string]any{
		"sender":   sender,
		"nonce":    hexutil.EncodeBig(nonce),
		"initCode": hexutil.Encode(initCode),
		"callData": hexutil.Encode(callData),
	}

	// gas estimation
	gasEstimate, err = a.EstimateUserOperationGas(a.EntryPoint, dummySignature, o)
	if err != nil {
		err = errors.Wrap(err, "gas estimation failed")
		return nil, err
	}

	maxPriorityFeePerGas, err = a.GetMaxPriorityFee()
	if err != nil {
		return nil, err
	}
	maxFeePerGas = calcMaxFeePerGas(maxPriorityFeePerGas)
	o["preVerificationGas"] = gasEstimate.PreVerificationGas
	o["verificationGasLimit"] = gasEstimate.VerificationGasLimit
	o["callGasLimit"] = gasEstimate.CallGasLimit
	o["maxFeePerGas"] = maxFeePerGas
	o["maxPriorityFeePerGas"] = maxPriorityFeePerGas

	if sponsored {
		policyId := os.Getenv("POLICY_ID")

		paymaster, err = a.RequestPaymasterAndData(policyId, a.EntryPoint, dummySignature, o)
		if err != nil {
			err = errors.Wrap(err, "call to sponsor user op failed")
			return nil, err
		}
		paymasterAndData = paymaster.PaymasterAndData
	} else {
		fmt.Println("not using paymaster")
		paymasterAndData = "0x"
	}

	o["signature"] = dummySignature
	o["paymasterAndData"] = paymasterAndData

	return o, nil
}

// Estimate the gas parameters required for the user operation
func (bs *AlchemyService) EstimateUserOperationGas(entrypointAddress, dummySignature string, userop map[string]any) (*GasEstimateResult, error) {
	result := &GasEstimateResult{}

	userop["callGasLimit"] = "0x16710"
	userop["verificationGasLimit"] = hexutil.EncodeBig(getVerificationGasLimit())
	userop["preVerificationGas"] = hexutil.EncodeBig(getVerificationGasLimit())
	userop["maxPriorityFeePerGas"] = hexutil.EncodeBig(getMaxPriorityFeePerGas())
	userop["maxFeePerGas"] = hexutil.EncodeBig(getMaxFeePerGas())
	userop["signature"] = dummySignature
	userop["paymasterAndData"] = "0x"

	userOperation := newUserOp(userop)
	backoffOperation := func() error {
		err := bs.backend.Client().CallContext(bs.ctx, result, "eth_estimateUserOperationGas", userOperation, entrypointAddress)
		if err != nil {
			err = parseRpcError(err)
		}

		// rpcErr := err.(rpc.Error)
		// if rpcErr.ErrorCode() != 429 {
		// 	return &backoff.PermanentError{
		// 		Err: rpcErr,
		// 	}
		// }
		return err
	}
	err := backoff.Retry(backoffOperation, backoff.NewExponentialBackOff())
	if err != nil {
		err = errors.Wrap(err, "eth_estimateUserOperationGas call error")
		return nil, err
	}
	pvg, _ := hexutil.DecodeUint64(result.PreVerificationGas)
	vgl, _ := hexutil.DecodeUint64(result.VerificationGasLimit)
	result.PreVerificationGas = hexutil.EncodeUint64(pvg * 2)
	result.VerificationGasLimit = hexutil.EncodeUint64(vgl * 2)
	log.Debug().Msgf("eth_estimateUserOperationGas - %v", result)
	return result, nil
}

func (bs *AlchemyService) RequestPaymasterAndData(policyId, entryPoint, dummySignature string, userop any) (*AlchemyPaymasterResult, error) {
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

	backoffOperation := func() error {
		err := bs.backend.Client().CallContext(bs.ctx, result, "alchemy_requestPaymasterAndData", request)
		if err != nil {
			err = parseRpcError(err)
		}
		return err
	}
	err := backoff.Retry(backoffOperation, backoff.NewExponentialBackOff())
	if err != nil {
		log.Err(err).Msg("alchemy_requestPaymasterAndData")
		return nil, err
	}

	fmt.Println("alchemy_PaymasterAndData - ", result)
	return result, nil
}

// paymaster data request for alchemy 2e865ced-98e5-4265-a20e-b46c695a28bd
func (bs *AlchemyService) RequestGasAndPaymasterAndData(policyId, entryPoint, dummySignature string, userop any) (*AlchemyPaymasterResult, error) {
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

	backoffOperation := func() error {
		err := bs.backend.Client().CallContext(bs.ctx, result, "alchemy_requestGasAndPaymasterAndData", request)
		if err != nil {
			err = parseRpcError(err)
		}
		return err
	}
	err := backoff.Retry(backoffOperation, backoff.NewExponentialBackOff())
	if err != nil {
		log.Err(err).Msg("alchemy_requestGasAndPaymasterAndData")
		return nil, err
	}
	fmt.Println("alchemy_requestGasAndPaymasterAndData - ", result)
	return result, nil
}

func (bs *AlchemyService) GetMaxPriorityFee() (string, error) {
	var result string

	err := bs.backend.Client().CallContext(bs.ctx, &result, "rundler_maxPriorityFeePerGas")
	if err != nil {
		err = errors.Wrap(err, "rundler_maxPriorityFeePerGas call error")
		return "", err
	}

	log.Debug().Msgf("rundler_maxPriorityFeePerGas - %v", result)

	resultUint, err := hexutil.DecodeUint64(result)
	if err != nil {
		return result, err
	}

	maxPriorityFee := resultUint * 2
	// maxPriorityFeeBig, _ := new(big.Int).SetString(result, 0)
	// maxPriorityFee := new(big.Int).Mul(maxPriorityFeeBig, big.NewInt(2))
	log.Debug().Msgf("adjusted maxPriorityFeePerGas - %v", maxPriorityFee)
	return hexutil.EncodeUint64(maxPriorityFee), nil
}

// SendUserOperation sends a user operation to an alt mempool and returns the userop hash if call is successful
func (bs *AlchemyService) SendUserOperation(userop map[string]any) (string, error) {
	var result string
	// log.Debug().Msgf("user op", userop)
	// obj, _ := json.Marshal(userop)
	userOperation := newUserOp(userop)

	backOffOperation := func() error {
		err := bs.backend.Client().CallContext(bs.ctx, &result, "eth_sendUserOperation", userOperation, bs.EntryPoint)
		if err != nil {
			err = parseRpcError(err)
			log.Err(err).Msg("rpc error")
		}
		return err
	}

	err := backoff.Retry(backOffOperation, backoff.NewExponentialBackOff())
	if err != nil {
		// if err.Error() == "AA13 initCode failed or OOG" {
		// 	// increase verification gas limit and retry
		// }
		err = errors.Wrap(err, "eth_sendUserOperation call error")
		return result, err
	}

	return result, nil
}

func (bs *AlchemyService) GetUserOperationByHash(userophash string) (map[string]any, error) {
	var result map[string]any

	backOffOperation := func() error {
		err := bs.backend.Client().CallContext(bs.ctx, &result, "eth_getUserOperationByHash", userophash)
		if err != nil {
			err = errors.Wrap(err, "eth_getUserOperationByHash call error")
			log.Err(err).Send()
		}

		if result == nil {
			err = errors.New("null rpc result")
			log.Debug().Err(err).Send()
		}

		return err
	}

	err := backoff.Retry(backOffOperation, backoff.NewExponentialBackOff())
	if err != nil {
		err = errors.Wrap(err, "eth_sendUserOperation call error")
		return result, err
	}
	log.Debug().Msgf("user operation - %v", result)
	return result, nil
}

func calcMaxFeePerGas(maxPriorityFee string) string {
	maxPriorityFeeBig, _ := new(big.Int).SetString(maxPriorityFee, 0)
	maxFeePerGasBig := new(big.Int).Add(maxPriorityFeeBig, big.NewInt(22))

	return hexutil.EncodeBig(maxFeePerGasBig)
}

func parseRpcError(err error) error {
	if !strings.Contains(err.Error(), "429") {
		return &backoff.PermanentError{
			Err: err,
		}
	}
	return err
}

func (bs *AlchemyService) IsAccountDeployed(address string, chain int64) bool {
	// bundler, err := erc4337.NewAlchemyService(chain)
	// if err != nil {
	// 	err = errors.Wrap(err, "failed to initialise bundler")
	// 	log.Panic().Err(err).Send()
	// 	return false
	// }

	code, err := bs.GetAccountCode(common.HexToAddress(address))
	if err != nil {
		log.Err(err).Caller().Send()
		return false
	}
	fmt.Println("Code ", code)
	if len(code) == 0 {
		log.Info().Msg("account not deployed, should be deployed first!")
		return false
	}
	return true
}

func (bs *AlchemyService) GetTransactionHash(useropHash string) (string, error) {
	useropResult, err := bs.GetUserOperationByHash(useropHash)
	if err != nil {
		err = errors.Wrap(err, "fetching the transction hash failed")
		log.Err(err).Caller().Send()
		return "", err
	}

	transactionHash := useropResult["transactionHash"].(string)
	return transactionHash, nil
}
