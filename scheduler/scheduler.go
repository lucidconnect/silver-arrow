package scheduler

import (
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	LucidMerchant "github.com/helicarrierstudio/silver-arrow/abi/LucidMerchant"
	"github.com/helicarrierstudio/silver-arrow/erc4337"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type Scheduler struct {
	bundler   *erc4337.ERCBundler
	queue     repository.Queuer
	datastore repository.WalletRepository
}

func NewScheduler(data repository.WalletRepository, bundler *erc4337.ERCBundler) *Scheduler {
	queue := repository.NewDeque()
	return &Scheduler{
		bundler:   bundler,
		queue:     queue,
		datastore: data,
	}
}

// create a valid the user op and add it to a queue
func (s *Scheduler) ScheduleUserOp() {
	// read from the database and fetch subscriptions with nextChargeAt = today
	now := time.Now().UTC().Truncate(24 * time.Hour)
	filter := bson.D{{Key: "next_charge_at", Value: now.String()}}
	result, err := s.datastore.FindSubscriptionsByFilter(filter)
	if err != nil {
		log.Println(err)
	}

	for _, v := range result {
		s.queue.Write(v)
	}
}

// func (s *Scheduler) SendUserOp(sub models.Subscription) error {
// 	token := sub.Token
// 	amount, _ := new(big.Int).SetString(sub.Amount,10)
// 	address, err := s.fetchMerchantAddress(sub.MerchantId)
// 	if err != nil {
// 		return err
// 	}

// 	data, err := erc4337.CreateTransferCallData(address, token, amount)
// 	if err != nil {
// 		return err
// 	}
// 	nonce := getAccountNonce(sub.WalletAddress)
// 	signingKey := sub.SigningKey
// 	op, err := s.bundler.CreateUserOperation(sub.WalletAddress, address, data, nonce, amount, true, signingKey, 0)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = s.bundler.SendUserOp(op)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func getAccountNonce(address string) *big.Int {
	return big.NewInt(0)
}

// FetchMerchantAddress call's the merchant contract
// and fetches the address for the given MerchantId
func (s *Scheduler) fetchMerchantAddress(merchantId string) (string, error) {
	contractAddress := os.Getenv("MERCHANT_CONTRACT")
	backend := s.bundler.GetClient()

	l, err := LucidMerchant.NewLucidMerchant(common.HexToAddress(contractAddress), backend)
	if err != nil {
		return "", err
	}

	sbyte := make([]byte, 32)
	copy(sbyte, []byte(merchantId))

	m, err := l.GetMerchant(nil, [32]byte(sbyte))
	if err != nil {
		return "", err
	}
	return m.ReceivingAddress.Hex(), nil
}
