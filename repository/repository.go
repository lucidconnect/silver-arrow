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
	FetchWalletSubscriptions(address string, status *string) ([]models.Subscription, error)
	FetchDueSubscriptions(days int) ([]models.Subscription, error)
	FindSubscriptionByHash(hash string) (*models.Subscription, error)
	FindSubscriptionById(id uuid.UUID) (*models.Subscription, error)
	FindSubscriptionByProductId(id uuid.UUID, wallet string) (*models.Subscription, error)
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
	UpdateProduct(productId uuid.UUID, merchantId uuid.UUID, update map[string]interface{}) error

	// Price
	CreatePrice(*models.Price) error
	FetchPrice(uuid.UUID) (*models.Price, error)
	FetchAllPrices(merchantId uuid.UUID) ([]models.Price, error)
	FetchAllPricesByProduct(productId uuid.UUID) ([]models.Price, error)
	UpdatePrice(priceId uuid.UUID, update map[string]interface{}) error

	// DepositWallet
	AddDepositWallet(*models.DepositWallet) error
	FetchDepositWallet(uuid.UUID) (*models.DepositWallet, error)
	FetchDepositWalletByMerchant(uuid.UUID) ([]models.DepositWallet, error)
	UpdateDepositWallet(walletId uuid.UUID, update *models.DepositWallet) error
	DeleteDepositWallet(id uuid.UUID) error

	// Merchant
	AddMerchant(*models.Merchant) error
	FetchMerchantById(uuid.UUID) (*models.Merchant, error)
	FetchMerchantByAddress(string) (*models.Merchant, error)
	FetchMerchantByPublicKey(string) (*models.Merchant, error)
	UpdateMerchantKey(id uuid.UUID, key, mode string) error
	UpdateMerchantWebhookUrl(uuid.UUID, string) error
	CreateMerchantAccessKeys(*models.MerchantAccessKey) error
	FetchMerchantKey(string) (*models.MerchantAccessKey, error)
	DeleteMerchantAccessKey(id uuid.UUID, key *models.MerchantAccessKey) error

	// Webhook
	CreateWebhookEvent(*models.WebhookEvent) error
	UpdateWebhookEvent(*models.WebhookEvent) error

	// Checkout Sesison
	CreateCheckoutSession(*models.CheckoutSession) error
	FetchCheckoutSession(id uuid.UUID) (*models.CheckoutSession, error)
	UpdateCheckoutSession(*models.CheckoutSession) error

	// PaymentLink
	CreatePaymentLink(*models.PaymentLink) error
	FetchPaymentLink(id uuid.UUID) (*models.PaymentLink, error)
	FetchPaymentLinkByProduct(productId uuid.UUID) (*models.PaymentLink, error)
	FetchPaymentLinkByMerchant(merchantId uuid.UUID) ([]models.PaymentLink, error)
	DeletePaymentLink(id uuid.UUID) error

	// Tokens
	AddToken(*models.Token) error
	FetchAllTokens(chain int64) ([]models.Token, error)
	FetchOneToken(name string, chain int64) (*models.Token, error)
}

type Queuer interface {
	Read() (models.Subscription, error)
	Write(models.Subscription)
}

type CacheWizard interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
