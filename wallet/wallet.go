package wallet

import (
	"fmt"
	"log"
	"math/big"
	"os"
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

func (ws *WalletService) ValidateSubscription(userop map[string]any) (*model.SubscriptionData, error) {
	// k, _ := Kernel.NewKernel(common.HexToAddress(input.WalletAddress), nil)
	// // k.ValidateUserOp()
	opHash, err := ws.bundler.SendUserOp(userop)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "userop_hash", Value: opHash}}
	results, err := ws.repository.FindSubscriptionsByFilter(filter)
	result := results[0]
	subData := &model.SubscriptionData{
		ID:            result.SubscriptionId,
		Token:         result.Token,
		Amount:        float64(result.Amount),
		Interval:      int(result.Interval),
		MerchantID:    result.MerchantId,
		WalletAddress: result.WalletAddress,
	}
	fmt.Println("userop hash ", opHash)
	return subData, err
}

func (ws *WalletService) AddSubscription(input model.NewSubscription) (*model.ValidationData, map[string]any, error) {
	var nextChargeAt time.Time
	publicKey, signingKey, err := CreateAccessKey()
	if err != nil {
		return nil, nil, err
	}

	amount := amountToWei(input.Amount)
	// if input.NextChargeAt != nil {
	// 	// check if the date is today and execute the charge
	// 	nextCharge := input.NextChargeAt.Format(time.DateOnly)
	// 	if nextCharge == time.Now().Format(time.DateOnly) {
	// 		// execute first
	// 		// set the next charge date
	// 	}
	// }

	// secondsInt := input.Interval * 60 * 60 * 24
	// seconds := time.Duration(secondsInt) * time.Second
	// _ = seconds / time.Nanosecond
	// // hours := 24 * time.Hour * time.Duration(days)
	// // minute := hours * time.Minute
	interval := daysToNanoSeconds(int64(input.Interval))

	nextChargeAt = time.Now().Add(interval)

	// strconv.
	// CreateaWhitelistData(in)
	callData, err := createValidatorEnableData(publicKey, input.MerchantID)
	if err != nil {
		err = errors.Wrap(err, "error creating validator data")
		return nil, nil, err
	}

	
	initCode, err := ws.getContractInitCode(common.HexToAddress(input.OwnerAddress))
	if err != nil {
		return nil, nil, err
	}
	op, err := ws.bundler.CreateUnsignedUserOperation(input.WalletAddress, input.WalletAddress, initCode, callData, common.Big0, common.Big0, false, 80001)
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
		Amount:         amount,
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

	// return &model.SubscriptionData{
	// 	ID:            fmt.Sprint(id),
	// 	Token:         input.Token,
	// 	Amount:        input.Amount,
	// 	Interval:      input.Interval,
	// 	MerchantID:    input.MerchantID,
	// 	WalletAddress: input.WalletAddress,
	// 	// ValidatorData: hexutil.Encode(enableData),
	// }
	return &model.ValidationData{
		UserOpHash: opHash.Hex(),
	}, op, nil
}

func amountToWei(a float64) int64 {
	return int64(1000000000000000000 * a)
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
	enableData = append(enableData, m...)
	return enableData, nil
}

func (ws *WalletService) getContractInitCode(accountAddress common.Address) ([]byte, error) {
	initCode := []byte{}
	factoryAddress := os.Getenv("KERNEL_FACTORY_ADDRESS")

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
	return initCode, nil
}
