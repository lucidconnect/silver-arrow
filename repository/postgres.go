package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	if err = db.AutoMigrate(models.Wallet{}, models.Key{}, models.Subscription{}); err != nil {
		log.Fatal("Error migrating database models")
	}
	// db.Model(&models.Subscription{}).
	// 	Exec(createForeignKeyIfNotExistsQuery("subscriptions", "wallets", "wallet_address", "wallet_address"))

	return db, nil
}

type DB struct {
	Db *gorm.DB
}

func NewDB(db *gorm.DB) *DB {
	return &DB{
		Db: db,
	}
}

func (p *DB) AddAccount(addressData *models.Wallet) error {
	return p.Db.Create(addressData).Error
}

func (p *DB) FetchAccountByAddress(address string) (*models.Wallet, error) {
	var wallet *models.Wallet
	err := p.Db.Where("wallet_address = ?", address).Find(&wallet).Error
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (p *DB) AddSubscription(subscriptionData *models.Subscription, key *models.Key) error {
	tx := p.Db.Begin()
	subscriptionData.Key = *key
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

func (p *DB) FetchWalletSubscriptions(address string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := p.Db.Where("wallet_address = ?", address).Preload("Key").Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (p *DB) DeactivateSubscription(id uint) error {
	return p.Db.Where("id = ?", id).UpdateColumn("active", false).Error
}

func (p *DB) UpdateSubscription(id uint) error {
	return errors.New("unimplemented")
}

func (p *DB) FetchDueSubscriptions(days int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	startInterval := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	endInterval := startInterval.Add(24 * time.Hour)

	err := p.Db.Where("expires_at >= ? AND expires_at <= ?", startInterval, endInterval).Preload("Key").Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (p *DB) AddSubscriptionKey(key *models.Key) error {
	return p.Db.Create(key).Error
}

func (p *DB) GetSecretKey(publicKey string) (string, error) {
	var key *models.Key
	if err := p.Db.Where("private_key_id = ?", publicKey).Find(&key).Error; err != nil {
		return "", err
	}

	return key.PrivateKeyId, nil
}

func (p *DB) FindSubscriptionByHash(hash string) (*models.Subscription, error) {
	var subscription *models.Subscription
	if err := p.Db.Where("user_op_hash = ?", hash).Preload("Key").Find(&subscription).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

// returns the private key ID
func (p *DB) GetSubscriptionKey(publicKey string) (string, error) {

	return "", nil
}

func (p *DB) CreateMerchant(m *models.Merchant) error {
	return p.Db.Create(m).Error
}

func (p *DB) FetchMerchant(id uuid.UUID) (*models.Merchant, error) {
	var merchant *models.Merchant
	if err := p.Db.Where("id = ?", id).Find(&merchant).Error; err != nil {
		return nil, err
	}
	return merchant, nil
}

func (p *DB) FetchMerchanstByOwner(owner string) ([]models.Merchant, error) {
	var merchants []models.Merchant
	if err := p.Db.Where("owner = ?", owner).Find(&merchants).Error; err != nil {
		return nil, err
	}
	return merchants, nil
}

func (p *DB) FindSubscriptionByMerchant(merchantId string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	if err := p.Db.Where("merchant_id = ?", merchantId).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func createForeignKeyIfNotExistsQuery(fromTable, targetTable, fromCol, targetCol string) string {
	return fmt.Sprintf(`
		DO $$
		BEGIN
			IF NOT (
				SELECT
					COUNT(1) >= 1
				FROM 
					information_schema.table_constraints AS tc 
					JOIN information_schema.key_column_usage AS kcu
						ON tc.constraint_name = kcu.constraint_name
						AND tc.table_schema = kcu.table_schema
					JOIN information_schema.constraint_column_usage AS ccu
						ON ccu.constraint_name = tc.constraint_name
						AND ccu.table_schema = tc.table_schema
				WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name='%s' AND ccu.column_name = '%s'
			) THEN
					ALTER TABLE %s
					ADD CONSTRAINT %s_%s_%s_%s_foreign
					FOREIGN KEY (%s) REFERENCES %s(%s)
					ON UPDATE RESTRICT
					ON DELETE RESTRICT;
			END IF;
		END;
		$$;`,
		fromTable, targetCol, fromTable, fromTable, fromCol, targetTable, targetCol, fromCol, targetTable, targetCol,
	)
}
