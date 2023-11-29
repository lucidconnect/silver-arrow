package repository

import (
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
)

type Database interface {
	// Wallet
	AddAccount(*models.Wallet) error
	FetchAccountByAddress(address string) (*models.Wallet, error)
	GetWalletMetadata(wallet string) (string, string, uuid.UUID, error)
	FindAllPaymentsByWallet(address string) ([]models.Payment, error)

	// Subscriptions
	AddSubscription(*models.Subscription, *models.Key) error
	AddSubscriptionKey(*models.Key) error
	FetchWalletSubscriptions(address string) ([]models.Subscription, error)
	FetchDueSubscriptions(days int) ([]models.Subscription, error)
	FindSubscriptionByHash(hash string) (*models.Subscription, error)
	FindSubscriptionById(id uuid.UUID) (*models.Subscription, error)
	FindSubscriptionByProductId(id uuid.UUID) (*models.Subscription, error)
	UpdateSubscription(uuid.UUID, map[string]interface{}) error
	DeactivateSubscription(id uint) error
	GetSubscriptionKey(publicKey string) (string, error)

	// Payments
	CreatePayment(*models.Payment) error
	UpdatePayment(id uuid.UUID, update map[string]any) error
	FindPaymentById(id uuid.UUID) (*models.Payment, error)
	FindPaymentByReference(uuid.UUID) (*models.Payment, error)
	FindPaymentByUseropHash(hash string) (*models.Payment, error)

	// Product
	CreateProduct(*models.Product) error
	FetchProduct(uuid.UUID) (*models.Product, error)
	FetchProductsByOwner(string) ([]models.Product, error)
	FetchAllPaymentsByProduct(productId uuid.UUID) ([]models.Payment, error)

	// Merchant
	AddMerchant(*models.Merchant) error
	FetchMerchantById(uuid.UUID) (*models.Merchant, error)
	FetchMerchantByAddress(string) (*models.Merchant, error)
	FetchMerchantByPublicKey(string) (*models.Merchant, error)
	UpdateMerchantKey(uuid.UUID, string) error
	UpdateMerchantWebhookUrl(uuid.UUID, string) error

	// Webhook
	CreateWebhookEvent(*models.WebhookEvent) error
	UpdateWebhookEvent(*models.WebhookEvent) error
}

type Queuer interface {
	Read() (models.Subscription, error)
	Write(models.Subscription)
}

type CacheWizard interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
