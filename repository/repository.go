package repository

import (
	"github.com/google/uuid"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

type Database interface {
	// Wallet
	AddAccount(*models.Wallet) error
	FetchAccountByAddress(address string) (*models.Wallet, error)
	AddSubscription(*models.Subscription, *models.Key) error
	FetchWalletSubscriptions(address string) ([]models.Subscription, error)
	FetchDueSubscriptions(days int) ([]models.Subscription, error)
	FindSubscriptionByHash(hash string) (*models.Subscription, error)
	UpdateSubscription(id uint) error
	DeactivateSubscription(id uint) error
	AddSubscriptionKey(*models.Key) error
	GetSubscriptionKey(publicKey string) (string, error)
	GetWalletMetadata(wallet string) (string, string, uuid.UUID, error)

	// Merchant
	CreateMerchant(*models.Merchant) error
	FetchMerchant(uuid.UUID) (*models.Merchant, error)
	FetchMerchanstByOwner(string) ([]models.Merchant, error)
	FindSubscriptionByMerchant(string) ([]models.Subscription, error)
}

type Queuer interface {
	Read() (models.Subscription, error)
	Write(models.Subscription)
}

type CacheWizard interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
