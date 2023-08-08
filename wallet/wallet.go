package wallet

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/helicarrierstudio/silver-arrow/erc4337"
	"github.com/helicarrierstudio/silver-arrow/graph/model"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"github.com/pkg/errors"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
	"go.mongodb.org/mongo-driver/bson"
)

type WalletService struct {
	repository repository.WalletRepository
	bundler    *erc4337.ERCBundler
}

func NewWalletService(r repository.WalletRepository, b *erc4337.ERCBundler) *WalletService {
	return &WalletService{
		repository: r,
		bundler:    b,
	}
}

func (ws *WalletService) AddAccount(input model.Account) error {
	walletAddress := input.Address
	// email := input.Email
	wallet := models.Wallet{
		AccountAddress: walletAddress,
	}
	err := ws.repository.SetAddress(wallet)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type userOp struct {
	CallData             string `json:"callData"`
	CallGasLimit         string `json:"callGasLimit"`
	InitCode             string `json:"initCode"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	Nonce                string `json:"nonce"`
	PaymasterAndData     string `json:"paymasterAndData"`
	PreVerificationGas   string `json:"preVerificationGas"`
	Sender               string `json:"sender"`
	Signature            string `json:"signature"`
	VerificationGasLimit string `json:"verificationGasLimit"`
}

func convertMapToStruct(m map[string]interface{}, s interface{}) error {
	stValue := reflect.ValueOf(s).Elem()
	sType := stValue.Type()
	for i := 0; i < sType.NumField(); i++ {
		field := sType.Field(i)
		if value, ok := m[field.Name]; ok {
			stValue.Field(i).Set(reflect.ValueOf(value))
		}
	}
	return nil
}

func (ws *WalletService) ValidateSubscription(userop map[string]any) (*model.SubscriptionData, string, error) {
	// k, _ := Kernel.NewKernel(common.HexToAddress(input.WalletAddress), nil)
	// k.ValidateUserOp()
	// var opStruct userOp
	// err := mapstructure.Decode(userop, &opStruct)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, "", err
	// }
	// fmt.Println("user operation struct: ", opStruct)
	opHash, err := ws.bundler.SendUserOp(userop)
	if err != nil {
		fmt.Println("userop hash ", opHash)
		return nil, "", err
	}
	filter := bson.D{{Key: "userop_hash", Value: opHash}}
	results, err := ws.repository.FindSubscriptionsByFilter(filter)
	if err != nil {
		return nil, "", err
	}
	result := results[0]
	// a, _ := new(big.Int).SetString(result.Amount, 10)
	token := result.Token
	// var amount int64
	// if token == "USDC" {
	// 	amount, _  = strconv.Atoi(result.Amount)
	// } else {
	// 	amount = weiToAmount(a)
	// }

	amount, _ := strconv.Atoi(result.Amount)
	subData := &model.SubscriptionData{
		ID:            result.SubscriptionId,
		Token:         token,
		Amount:        amount,
		Interval:      int(result.Interval),
		MerchantID:    result.MerchantId,
		WalletAddress: result.WalletAddress,
	}
	fmt.Println("userop hash ", opHash)

	return subData, result.SigningKey, nil
}

func (ws *WalletService) AddSubscription(input model.NewSubscription) (*model.ValidationData, map[string]any, error) {
	var nextChargeAt time.Time
	var initCode []byte
	var nonce, amount *big.Int
	publicKey, signingKey, err := CreateAccessKey()
	if err != nil {
		return nil, nil, err
	}
	amount = big.NewInt(int64(input.Amount)) // This will cause a bug for amounts that are fractional

	interval := daysToNanoSeconds(int64(input.Interval))

	nextChargeAt = time.Now().Add(interval)

	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}

	isAccountDeployed := ws.isAccountDeployed(input.WalletAddress)
	if !isAccountDeployed {
		initCode, err = getContractInitCode(common.HexToAddress(input.OwnerAddress))
		if err != nil {
			return nil, nil, err
		}
		// usePaymaster = false
	}

	callData, err := createValidatorEnableData(publicKey, input.MerchantID)
	if err != nil {
		err = errors.Wrap(err, "error creating validator data")
		return nil, nil, err
	}

	if !isAccountDeployed {
		nonce = common.Big0
	} else {
		nonce, err = ws.bundler.AccountNonce(input.WalletAddress)
		if err != nil {
			log.Println(err)
			return nil, nil, err
		}
	}

	op, err := ws.bundler.CreateUnsignedUserOperation(input.WalletAddress, input.WalletAddress, initCode, callData, nonce, usePaymaster, int64(input.Chain))
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	entrypoint := erc4337.GetEntryPointAddress()
	operation, err := userop.New(op)
	if err != nil {
		return nil, nil, err
	}
	opHash := operation.GetUserOpHash(entrypoint, big.NewInt(int64(input.Chain)))

	sub := models.Subscription{
		Token:          input.Token,
		Amount:         amount.String(),
		Active:         false,
		Interval:       interval.Nanoseconds(),
		UserOpHash:     opHash.Hex(),
		SigningKey:     signingKey,
		MerchantId:     input.MerchantID,
		NextChargeAt:   nextChargeAt,
		WalletAddress:  input.WalletAddress,
		SubscriptionId: publicKey,
	}
	_, err = ws.repository.AddSubscription(sub)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return &model.ValidationData{
		UserOpHash: opHash.Hex(),
	}, op, nil
}

func amountToWei(amount any) (*big.Int, error) {
	etherInWei := new(big.Int)
	etherInWei.SetString("1000000000000000000", 10)

	switch v := amount.(type) {
	case *big.Int:
		weiAmount := new(big.Int).Mul(v, etherInWei)
		return weiAmount, nil
	case *big.Float:
		weiAmount := new(big.Int)
		weiAmountFloat := new(big.Float).Mul(v, big.NewFloat(1e18))
		weiAmountFloat.Int(weiAmount)
		return weiAmount, nil
	default:
		return nil, fmt.Errorf("unsupported input type: %T", amount)
	}
}

func amountToMwei(amount int64) (*big.Int, error) {
	etherInMWei := new(big.Int)
	etherInMWei.SetString("1000000", 10)
	v := big.NewInt(int64(amount))
	mWeiAmount := new(big.Int).Mul(v, etherInMWei)
	return mWeiAmount, nil
	// switch v := amount.(type) {
	// case *big.Int:
	// 	mWeiAmount := new(big.Int).Mul(v, etherInMWei)
	// 	return mWeiAmount, nil
	// case *big.Float:
	// 	mWeiAmount := new(big.Int)
	// 	mWeiAmountFloat := new(big.Float).Mul(v, big.NewFloat(1e6))
	// 	mWeiAmountFloat.Int(mWeiAmount)
	// 	return mWeiAmount, nil
	// default:
	// 	return nil, fmt.Errorf("unsupported input type: %T", amount)
	// }
}

func mWeiToAmount(amt *big.Int) int64 {
	etherInMWei := new(big.Int)
	etherInMWei.SetString("1000000", 10)

	result := new(big.Int)
	result.Div(amt, etherInMWei)
	return result.Int64()
}

func weiToAmount(amt *big.Int) int64 {
	etherInWei := new(big.Int)
	etherInWei.SetString("1000000000000000000", 10)

	result := new(big.Int)
	result.Div(amt, etherInWei)
	return result.Int64()
}

func (ws *WalletService) ExecuteCharge(sender, target, mId, token, key string, amount int64) error {
	// if token != "USDC" {
	// 	//	should convert to wei
	// }

	actualAmount, err := amountToMwei(amount)
	if err != nil {
		return err
	}
	data, err := erc4337.CreateTransferCallData(target, token, actualAmount)
	if err != nil {
		err = errors.Wrap(err, "CreateTransferCallData() - ")
		return err
	}

	nonce, err := ws.bundler.AccountNonce(sender)
	if err != nil {
		log.Println(err)
		return err
	}
	merchant, _ := hexutil.Decode(mId)
	chainId := int64(80001)

	op, err := ws.bundler.CreateUnsignedUserOperation(sender, target, nil, data, nonce, true, chainId)
	if err != nil {
		err = errors.Wrap(err, "CreateUnsignedUserOperation() - ")
		return err
	}
	// fmt.Println("user operation", op)

	sig, _ := erc4337.SignUserOp(op, key, erc4337.VALIDATOR_MODE, merchant, int64(chainId))

	op["signature"] = hexutil.Encode(sig)

	opHash, err := ws.bundler.SendUserOp(op)
	if err != nil {
		err = errors.Wrap(err, "SendUserOp() - ")
		return err
	}

	fmt.Println("user operation hash -", opHash) // 0x28b45cf378c23fbdbbcb4f4c4d085791eb6d660214ff4a2402e40fd1c73751c6 0xfcd3b481cc3ba345fcf24c777463baf60dbb1f7475ca297b9259d020044565be
	return nil

}

func daysToNanoSeconds(days int64) time.Duration {
	secondsInt := days * 24 * 60 * 60 * 1e9
	return time.Duration(secondsInt)
}

func createValidatorEnableData(publicKey, merchantId string) ([]byte, error) {
	enableData, err := hexutil.Decode(publicKey)
	if err != nil {
		err = errors.Wrap(err, "failed to decode public key hex")
		return nil, err
	}

	m := []byte(merchantId)
	mid := common.LeftPadBytes(m, 32)
	fmt.Println("length of merchant id", mid)
	enableData = append(enableData, mid...)

	data, err := erc4337.CreateSetExecutionCallData(enableData)
	if err != nil {
		err = errors.Wrap(err, "CreateSetExecutionCallData() - ")
		return nil, err
	}

	return data, nil
}

func getContractInitCode(accountAddress common.Address) ([]byte, error) {
	initCode := []byte{}
	factoryAddress := os.Getenv("KERNEL_FACTORY_ADDRESS")
	fmt.Println("accountAddress ", accountAddress)
	data := accountAddress.Bytes()
	fmt.Println("enable data ", hexutil.Encode(data))
	code, err := erc4337.CreateFactoryFnData(accountAddress.Bytes())
	if err != nil {
		return nil, err
	}
	factoryAddressToBytes := common.FromHex(factoryAddress)
	initCode = append(initCode, factoryAddressToBytes...)
	initCode = append(initCode, code...)

	fmt.Println("initcode: ", (initCode))
	fmt.Println(hexutil.Encode(initCode))
	return initCode, nil
}

func (ws *WalletService) isAccountDeployed(address string) bool {
	code, err := ws.bundler.GetClient().CodeAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		fmt.Println("An error occured")
		return false
	}
	fmt.Println("Code ", code)
	if len(code) == 0 {
		fmt.Println("account not deployed, should be deployed first!")
		return false
	}
	return true
}
