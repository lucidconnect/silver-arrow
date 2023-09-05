package repository

import (
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

type SchedulerRepository interface{}

// type WalletRepository interface {
// 	SetAddress(models.Wallet) error
// 	AddSubscription(models.Subscription) (any, error)
// 	ListSubscriptions(address string) ([]models.Subscription, error)
// 	RemoveSubscription(id int64) error
// 	FindSubscriptionsByFilter(filter any) ([]models.Subscription, error)
// }

type WalletRepository interface {
	KeyManager
	AddressBook
	Subscriptions
}

type AddressBook interface {
	SetAddress(models.Wallet) error
}

type Subscriptions interface {
	AddSubscription(models.Subscription) error
	FetchWalletSubscriptions(address string) ([]models.Subscription, error)
	FetchDueSubscriptions(days int) ([]models.Subscription, error)
	FindSubscriptionByHash(hash string) (*models.Subscription, error)
	UpdateSubscription(id uint) error
	DeactivateSubscription(id uint) error
}

type KeyManager interface {
	SetKey(models.Key) error
	GetSecretKey(publicKey string) (string, error)
}

type Queuer interface {
	Read() (models.Subscription, error)
	Write(models.Subscription)
}

type CacheWizard interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
