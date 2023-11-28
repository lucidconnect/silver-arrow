package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SetupDatabase(dbconn *sql.DB) (*gorm.DB, error) {
	// ...
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("Connecting to database")
	dialector := postgres.Open(dsn)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database")
		return nil, err
	}

	// ...
	if err = db.AutoMigrate(models.Payment{}, models.Wallet{}, models.Merchant{}, models.Key{}, models.Subscription{}, models.Product{}); err != nil {
		log.Fatal().Err(err).Msg("Error migrating database models")
	}
	// db.Model(&models.Subscription{}).
	// 	Exec(createForeignKeyIfNotExistsQuery("subscriptions", "wallets", "wallet_address", "wallet_address"))

	return db, nil
}

type PostgresDB struct {
	Db *gorm.DB
}

func NewPostgresDB(db *gorm.DB) *PostgresDB {
	return &PostgresDB{
		Db: db,
	}
}

func (p *PostgresDB) AddAccount(addressData *models.Wallet) error {
	return p.Db.Create(addressData).Error
}

func (p *PostgresDB) FetchAccountByAddress(address string) (*models.Wallet, error) {
	var wallet *models.Wallet
	err := p.Db.Where("wallet_address = ?", address).First(&wallet).Error
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (p *PostgresDB) AddSubscription(subscriptionData *models.Subscription, key *models.Key) error {
	tx := p.Db.Begin()
	// subscriptionData.Key = *key
	if err := tx.Create(subscriptionData).Error; err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	// if err := tx.Create(key).Error; err != nil {
	// 	fmt.Println(err)
	// 	tx.Rollback()
	// 	return err
	// }

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) FetchWalletSubscriptions(address string, status *string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	var err error
	if status != nil {
		err = p.Db.Where("wallet_address = ? AND status =  ?", address, *status).Preload(clause.Associations).Find(&subscriptions).Error
	} else {
		err = p.Db.Where("wallet_address = ? AND status = ?", address, "active").Preload(clause.Associations).Find(&subscriptions).Error
	}
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (p *PostgresDB) DeactivateSubscription(id uint) error {
	return p.Db.Where("id = ?", id).UpdateColumn("active", false).Error
}

func (p *PostgresDB) UpdateSubscription(id uuid.UUID, update map[string]interface{}) error {
	var subscription *models.Subscription

	if err := p.Db.Model(&subscription).Where("id = ?", id).Updates(update).Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) FetchDueSubscriptions(days int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	startInterval := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	endInterval := startInterval.Add(24 * time.Hour)

	err := p.Db.Where("active = ? AND expires_at >= ? AND expires_at <= ?", true, startInterval, endInterval).Preload("Key").Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (p *PostgresDB) AddSubscriptionKey(key *models.Key) error {
	return p.Db.Create(key).Error
}

func (p *PostgresDB) GetSecretKey(publicKey string) (string, error) {
	var key *models.Key
	if err := p.Db.Where("private_key_id = ?", publicKey).Find(&key).Error; err != nil {
		return "", err
	}

	return key.PrivateKeyId, nil
}

func (p *PostgresDB) FindSubscriptionByHash(hash string) (*models.Subscription, error) {
	var subscription *models.Subscription
	if err := p.Db.Where("user_op_hash = ?", hash).Preload("Key").Find(&subscription).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

func (p *PostgresDB) FindSubscriptionById(id uuid.UUID) (*models.Subscription, error) {
	var subscription *models.Subscription
	if err := p.Db.Where("id = ?", id).Preload("Key").Find(&subscription).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

func (p *PostgresDB) FindSubscriptionByProductId(id uuid.UUID) (*models.Subscription, error) {
	var subscription *models.Subscription
	if err := p.Db.Where("product_id = ? AND active = ?", id, true).First(&subscription).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

// returns the private key ID
func (p *PostgresDB) GetSubscriptionKey(publicKey string) (string, error) {
	var key models.Key
	if err := p.Db.Where("public_key = ?", publicKey).Find(&key).Error; err != nil {
		return "", err
	}
	return key.PrivateKeyId, nil
}

func (p *PostgresDB) GetWalletMetadata(address string) (string, string, uuid.UUID, error) {
	var wallet models.Wallet
	if err := p.Db.Where("wallet_address = ?", address).Find(&wallet).Error; err != nil {
		return "", "", uuid.Nil, err
	}
	keyTag := wallet.TurnkeyPrivateKeyTag
	orgId := wallet.TurnkeySubOrgID
	walletId := wallet.ID

	return keyTag, orgId, walletId, nil
}

func (p *PostgresDB) CreateProduct(m *models.Product) error {
	return p.Db.Create(m).Error
}

func (p *PostgresDB) FetchProduct(id uuid.UUID) (*models.Product, error) {
	var product *models.Product
	if err := p.Db.Where("id = ?", id).Preload("Subscriptions").Find(&product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (p *PostgresDB) FetchProductsByOwner(owner string) ([]models.Product, error) {
	var products []models.Product
	if err := p.Db.Where("owner = ?", owner).Preload("Subscriptions").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (p *PostgresDB) UpdateProduct(id, merchantId uuid.UUID, update map[string]interface{}) error {
	var product *models.Product

	if err := p.Db.Model(&product).Where("id = ? AND merchant_id = ?", id, merchantId).Updates(update).Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) AddMerchant(merchant *models.Merchant) error {
	return p.Db.Create(merchant).Error
}

func (p *PostgresDB) FetchMerchantByAddress(address string) (*models.Merchant, error) {
	var merchant *models.Merchant
	if err := p.Db.Where("owner_address = ?", address).Preload("MerchantAccessKeys").First(&merchant).Error; err != nil {
		return nil, err
	}
	return merchant, nil
}

func (p *PostgresDB) FetchMerchantByPublicKey(key string) (*models.Merchant, error) {
	var merchant *models.Merchant

	if err := p.Db.Joins("JOIN access_keys ON access_keys.merchant_id = merchants.id").Where("access_keys.public_key = ?", key).First(&merchant).Error; err != nil {
		return nil, err
	}

	return merchant, nil
}

func (p *PostgresDB) UpdateMerchantKey(id uuid.UUID, key, mode string) error {
	var accessKey *models.MerchantAccessKey

	err := p.Db.Model(&accessKey).Where("merchant_id = ? AND mode = ?", id, mode).Update("public_key", key).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) FetchMerchantKey(key string) (*models.MerchantAccessKey, error) {
	var accessKey *models.MerchantAccessKey
	if err := p.Db.Where("public_key = ?", key).First(&accessKey).Error; err != nil {
		return nil, err
	}
	return accessKey, nil
}

func (p *PostgresDB) CreateMerchantAccessKeys(key *models.MerchantAccessKey) error {
	return p.Db.Create(key).Error
}

func (p *PostgresDB) DeleteMerchantAccessKey(id uuid.UUID, key *models.MerchantAccessKey) error {
	return p.Db.Where("id = ?", id).Delete(key).Error
}

func (p *PostgresDB) UpdateMerchantWebhookUrl(id uuid.UUID, webhookUrl string) error {
	var merchant *models.Merchant

	if err := p.Db.Model(&merchant).Where("id = ?", id).Update("webhook_url", webhookUrl).Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) CreatePayment(payment *models.Payment) error {
	return p.Db.Create(payment).Error
}

func (p *PostgresDB) UpdatePayment(id uuid.UUID, paymentUpdate map[string]any) error {
	var payment *models.Payment
	if err := p.Db.Model(&payment).Where("id = ?", id).Updates(paymentUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) FindPaymentById(id uuid.UUID) (*models.Payment, error) {
	var payment *models.Payment
	if err := p.Db.Where("id = ? ", id).First(&payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PostgresDB) FindPaymentByReference(reference uuid.UUID) (*models.Payment, error) {
	var payment *models.Payment
	if err := p.Db.Where("reference = ? ", reference).First(&payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PostgresDB) FindPaymentByUseropHash(hash string) (*models.Payment, error) {
	var payment *models.Payment
	if err := p.Db.Where("user_op_hash = ? ", hash).First(&payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (p *PostgresDB) FindAllPaymentsByWallet(address string) ([]models.Payment, error) {
	var wallet models.Wallet

	if err := p.Db.Where("wallet_address = ?", address).Preload("Payments").First(&wallet).Error; err != nil {
		return nil, err
	}
	payments := wallet.Payments
	return payments, nil
}

func (p *PostgresDB) FetchAllPaymentsByProduct(productId uuid.UUID) ([]models.Payment, error) {
	var subscriptions []models.Subscription
	var payments []models.Payment
	if err := p.Db.Where("product_id = ?", productId).Preload("Payments").Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	for _, subscription := range subscriptions {
		payments = append(payments, subscription.Payments...)
	}
	return payments, nil
}

func (p *PostgresDB) FetchMerchantById(id uuid.UUID) (*models.Merchant, error) {
	var merchant *models.Merchant
	if err := p.Db.Where("id = ?", id).Preload("MerchantAccessKeys").First(&merchant).Error; err != nil {
		return nil, err
	}

	return merchant, nil
}

func (p *PostgresDB) CreateWebhookEvent(webhookEvent *models.WebhookEvent) error {
	return p.Db.Create(webhookEvent).Error
}

func (p *PostgresDB) UpdateWebhookEvent(webhookEvent *models.WebhookEvent) error {
	return p.Db.Save(webhookEvent).Error
}
