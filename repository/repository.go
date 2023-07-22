package repository

import (
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

type SchedulerRepository interface{}

type WalletRepository interface {
	SetAddress(models.Wallet) error
	AddSubscription(models.Subscription) error
	ListSubscriptions(address string) ([]models.Subscription, error)
	RemoveSubscription(id int64) error
}
