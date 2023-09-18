package repository

import (
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

type WalletRepository interface {
	AddAccount(*models.Wallet) error
	AddSubscription(*models.Subscription, *models.Key) error
	FetchWalletSubscriptions(address string) ([]models.Subscription, error)
	FetchDueSubscriptions(days int) ([]models.Subscription, error)
	FindSubscriptionByHash(hash string) (*models.Subscription, error)
	UpdateSubscription(id uint) error
	DeactivateSubscription(id uint) error
	AddSubscriptionKey(*models.Key) error
	GetSubscriptionKey(publicKey string) (string, error)
}

type Queuer interface {
	Read() (models.Subscription, error)
	Write(models.Subscription)
}

type CacheWizard interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
