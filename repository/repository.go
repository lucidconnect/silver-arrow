package repository

import (
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

type SchedulerRepository interface{}

type WalletRepository interface {
	SetAddress(models.Wallet) error
	AddSubscription(models.Subscription) (any, error)
	ListSubscriptions(address string) ([]models.Subscription, error)
	RemoveSubscription(id int64) error
	FindSubscriptionsByFilter(filter any) ([]models.Subscription, error)
}

type Queuer interface {
	Read() (models.Subscription, error)
	Write(models.Subscription)
}

type CacheWizard interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
