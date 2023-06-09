package repository

import (
	"database/sql"
	"fmt"

	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(dsn string, dbconn *sql.DB) (*gorm.DB, error) {
	// ...
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
	return db, nil
}

type PostgresDb struct {
	Db *gorm.DB
}

func (p *PostgresDb) SetAddress(addressData interface{}) error {
	return p.Db.Create(addressData).Error
}

func (p *PostgresDb) AddSubscription(subscriptionData interface{}) error {
	return p.Db.Create(subscriptionData).Error
}

func (p *PostgresDb) ListSubscriptions(address string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	err := p.Db.Where("account_address = ?", address).Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (p *PostgresDb) RemoveSubscription(id uint64) error {
	return p.Db.Where("id = ?", id).UpdateColumn("active", false).Error
}
