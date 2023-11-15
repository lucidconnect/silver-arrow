package server

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/common"
	LucidMerchant "github.com/lucidconnect/silver-arrow/abi/LucidMerchant"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/service/wallet"
	"github.com/pkg/errors"
)

var defaultChain int64

// create a valid the user op and add it to a queue
func (s *Server) SubscriptionJob() {
	// read from the database and fetch subscriptions expiring in 3 days
	subsDueIn3, err := s.database.FetchDueSubscriptions(3)
	if err != nil {
		log.Err(err).Send()
	}

	dueToday, err := s.database.FetchDueSubscriptions(0)
	if err != nil {
		log.Err(err).Send()
	}

	// client := s.bundler.GetClient()

	// loop through subsDueIn3 and check if the account has enough to cover the sub
	for _, sub := range subsDueIn3 {
		amount := big.NewInt(sub.Amount)
		wallet := common.HexToAddress(sub.WalletAddress)
		token := common.HexToAddress(sub.TokenAddress)
		// chain := sub.Chain

		balance, err := s.bundler.GetErc20TokenBalance(token, wallet)
		if err != nil {
			log.Err(err).Send()
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
		tokenAddress := common.HexToAddress(sub.TokenAddress)

		balance, err := s.bundler.GetErc20TokenBalance(tokenAddress, wallet)
		if err != nil {
			log.Err(err).Send()
			continue
		}

		if balance.CmpAbs(amount) < 1 {
			fmt.Printf("Wallet at address %v does not have enough balance to complete the transaction", wallet)
			// should send an email notification to the user's email
			continue
		} else {
			// initiate user operation
			time.Sleep(15 * time.Second)
			// get the account
			// ws := wallet.NewWalletService(s.database)
			s.initialisePayment(sub)
		}
	}
}

func (s *Server) initialisePayment(sub models.Subscription) {
	var sponsored bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		sponsored = true
	default:
		sponsored = false
	}

	merchantId := uuid.MustParse(sub.MerchantId)
	merchant, err := s.database.FetchMerchantById(merchantId)
	if err != nil {
		log.Err(err).Msg("merchant not found")
		return
	}

	walletService := wallet.NewWalletService(s.database, sub.Chain)
	reference := uuid.New()

	payment := &models.Payment{
		Type:                  wallet.PaymentTypeRecurring.String(),
		Chain:                 sub.Chain,
		Token:                 sub.Token,
		Amount:                sub.Amount,
		Source:                sub.WalletAddress,
		WalletID:              sub.WalletID,
		ProductID:             sub.ProductID,
		Sponsored:             sponsored,
		Reference:             reference,
		Destination:           sub.MerchantDepositAddress,
		SubscriptionID:        sub.ID,
		SubscriptionPublicKey: sub.Key.PublicKey,
	}

	userop, useropHash, err := walletService.CreatePayment(payment)
	if err != nil {
		err = errors.Wrap(err, "creating payment operation failed")
		log.Err(err).Caller().Send()
	}

	signature, err := walletService.SignPaymentOperation(userop, useropHash)
	if err != nil {
		err = errors.Wrap(err, "signing payment operation failed")
		log.Err(err).Caller().Send()
	}
	userop["signature"] = signature

	_, err = walletService.ExecutePaymentOperation(userop, payment.Chain)
	if err != nil {
		log.Err(err).Send()
	}

	nextChargeAt := time.Now().Add((time.Duration(sub.Interval)))

	update := map[string]interface{}{
		"expires_at":     nextChargeAt,
		"next_charge_at": nextChargeAt,
	}
	err = s.database.UpdateSubscription(sub.ID, update)
	if err != nil {
		log.Err(err).Send()
	}

	s.TriggerWebhook(*merchant, reference.String())
}

func getAccountNonce(address string) *big.Int {
	return big.NewInt(0)
}

// FetchMerchantAddress call's the merchant contract
// and fetches the address for the given MerchantId
func (s *Server) fetchMerchantAddress(merchantId string) (string, error) {
	contractAddress := os.Getenv("MERCHANT_CONTRACT")
	backend := s.bundler.GetEthBackend()

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
