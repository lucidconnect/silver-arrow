package repository

import (
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
)

type Database interface {
	// Wallet
	AddAccount(*models.Wallet) error
	FetchAccountByAddress(address string) (*models.Wallet, error)
	AddSubscription(*models.Subscription, *models.Key) error
	FetchWalletSubscriptions(address string) ([]models.Subscription, error)
	FetchDueSubscriptions(days int) ([]models.Subscription, error)
	FindSubscriptionByHash(hash string) (*models.Subscription, error)
	FindSubscriptionById(id uuid.UUID) (*models.Subscription, error) 
	UpdateSubscription(uuid.UUID, map[string]interface{}) error
	DeactivateSubscription(id uint) error
	AddSubscriptionKey(*models.Key) error
	GetSubscriptionKey(publicKey string) (string, error)
	GetWalletMetadata(wallet string) (string, string, uuid.UUID, error)

	// Product
	CreateProduct(*models.Product) error
	FetchProduct(uuid.UUID) (*models.Product, error)
	FetchProductsByOwner(string) ([]models.Product, error)
	FindSubscriptionByProduct(string) ([]models.Subscription, error)

	// Merchant
	AddMerchant(*models.Merchant) error
	FetchMerchantByAddress(string) (*models.Merchant, error)
	FetchMerchantByPublicKey(string) (*models.Merchant, error)
	UpdateMerchantKey(uuid.UUID, string) error
}

type Queuer interface {
	Read() (models.Subscription, error)
	Write(models.Subscription)
}

type CacheWizard interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
