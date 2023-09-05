package scheduler

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	LucidMerchant "github.com/helicarrierstudio/silver-arrow/abi/LucidMerchant"
	"github.com/helicarrierstudio/silver-arrow/erc4337"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/wallet"
	"github.com/pkg/errors"
)

type Scheduler struct {
	bundler       *erc4337.ERCBundler
	queue         repository.Queuer
	datastore     repository.WalletRepository
	walletService *wallet.WalletService
}

func NewScheduler(data repository.WalletRepository, bundler *erc4337.ERCBundler, wallet *wallet.WalletService) *Scheduler {
	queue := repository.NewDeque()
	return &Scheduler{
		bundler:       bundler,
		queue:         queue,
		datastore:     data,
		walletService: wallet,
	}
}

// create a valid the user op and add it to a queue
func (s *Scheduler) SubscriptionJob() {
	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}
	// read from the database and fetch subscriptions expiring in 3 days
	subsDueIn3, err := s.datastore.FetchDueSubscriptions(3)
	if err != nil {
		log.Println(err)
	}
	
	dueToday, err := s.datastore.FetchDueSubscriptions(0)
	if err != nil {
		log.Println(err)
	}

	client := s.bundler.GetClient()

	// loop through subsDueIn3 and check if the account has enough to cover the sub
	for _, sub := range subsDueIn3 {
		amount := big.NewInt(sub.Amount)
		wallet := common.HexToAddress(sub.WalletAddress)
		token := common.HexToAddress(sub.TokenAddress)

		balance, err := client.GetErc20TokenBalance(token, wallet)
		if err != nil {
			log.Println(err)
			continue
		}

		if balance.CmpAbs(amount) < 1 {
			fmt.Printf("Wallet at address %v does not have enough balance to complete the transaction", wallet)
			// should send an email notification to the user's email
			continue
		} else {
			continue
		}
	}

	for _, sub := range dueToday {
		amount := big.NewInt(sub.Amount)
		wallet := common.HexToAddress(sub.WalletAddress)
		token := common.HexToAddress(sub.TokenAddress)

		balance, err := client.GetErc20TokenBalance(token, wallet)
		if err != nil {
			log.Println(err)
			continue
		}

		if balance.CmpAbs(amount) < 1 {
			fmt.Printf("Wallet at address %v does not have enough balance to complete the transaction", wallet)
			// should send an email notification to the user's email
			continue
		} else {
			// initiate user operation
			time.Sleep(15 * time.Second)
			err = s.walletService.ExecuteCharge(sub.WalletAddress, sub.MerchantDepositAddress, sub.MerchantId, sub.Token, sub.Key.SecretKey, sub.Amount, usePaymaster)
			if err != nil {
				err = errors.Wrapf(err, "ExecuteCharge() - error occurred during charge execution for subscription %v - ", sub.ID)
				log.Println(err)
				continue
			}
		}
	}
}

func getAccountNonce(address string) *big.Int {
	return big.NewInt(0)
}

// FetchMerchantAddress call's the merchant contract
// and fetches the address for the given MerchantId
func (s *Scheduler) fetchMerchantAddress(merchantId string) (string, error) {
	contractAddress := os.Getenv("MERCHANT_CONTRACT")
	backend := s.bundler.GetClient().GetEthClient()

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
