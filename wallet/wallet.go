package wallet

import (
	"log"

	"github.com/helicarrierstudio/silver-arrow/graph/model"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

type WalletService struct {
	repository repository.WalletRepository
}

func NewWalletService(r repository.WalletRepository) *WalletService {
	return &WalletService{
		repository: r,
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

func (ws *WalletService) AddSubscription(input model.NewSubscription) (*model.SubscriptionData, error) {
	CreateAccessKey()
	sub := models.Subscription{
		Amount:     input.Amount,
		Active:     true,
		Interval:   input.Interval,
		Token:      input.Token,
		MerchantId: input.MerchantID,
	}

	// strconv.
	// CreateaWhitelistData(in)
	err := ws.repository.AddSubscription(sub)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return nil, nil
}
