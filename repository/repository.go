package repository

import "github.com/helicarrierstudio/silver-arrow/repository/models"

type SchedulerRepository interface {
	SetAddress(interface{}) error
	AddSubscription(interface{}) error
	ListSubscriptions(address string) ([]models.Subscription, error)
	RemoveSubscription(id int64) error
}
