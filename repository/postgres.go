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

	db.Model(&models.Subscription{}).
		Exec(createForeignKeyIfNotExistsQuery("subscriptions", "accounts", "account_address", "account_address"))

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
