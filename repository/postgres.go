package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(dbconn *sql.DB) (*gorm.DB, error) {
	// ...
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("Connecting to database")
	dialector := postgres.New(postgres.Config{
		DSN:        dsn,
		DriverName: "postgres",
		Conn:       dbconn,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database")
		return nil, err
	}

	// ...

	db.Model(&models.Subscription{}).
		Exec(createForeignKeyIfNotExistsQuery("subscriptions", "accounts", "account_address", "account_address"))

	return db, nil
}

type PostgresDb struct {
	Db *gorm.DB
}

func NewPostgresDb(db *gorm.DB) *PostgresDb {
	return &PostgresDb{
		Db: db,
	}
}

func (p *PostgresDb) SetAddress(addressData models.Wallet) error {
	return p.Db.Create(addressData).Error
}

func (p *PostgresDb) AddSubscription(subscriptionData models.Subscription) error {
	return p.Db.Create(subscriptionData).Error
}

func (p *PostgresDb) FetchWalletSubscriptions(address string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := p.Db.Where("wallet_address = ?", address).Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (p *PostgresDb) DeactivateSubscription(id uint) error {
	return p.Db.Where("id = ?", id).UpdateColumn("active", false).Error
}

func (p *PostgresDb) UpdateSubscription(id uint) error {
	return errors.New("unimplemented")
}

func (p *PostgresDb) FetchDueSubscriptions(days int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription

	startInterval := time.Now().Add(time.Duration(days) * 24 * time.Hour)
	endInterval := startInterval.Add(24 * time.Hour)

	err := p.Db.Where("expires_at >= ? AND expires_at <= ?", startInterval, endInterval).Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (p *PostgresDb) SetKey(key models.Key) error {
	return p.Db.Create(key).Error
}

func (p *PostgresDb) GetSecretKey(publicKey string) (string, error) {
	var key *models.Key
	if err := p.Db.Where("subscription_key = ?", publicKey).Find(&key).Error; err != nil {
		return "", err
	}

	return key.SecretKey, nil
}

func (p *PostgresDb) FindSubscriptionByHash(hash string) (*models.Subscription, error) {
	var subscription *models.Subscription
	if err := p.Db.Where("user_op_hash = ?", hash).Find(&subscription).Error; err != nil {
		return nil, err
	}
	return subscription, nil
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
