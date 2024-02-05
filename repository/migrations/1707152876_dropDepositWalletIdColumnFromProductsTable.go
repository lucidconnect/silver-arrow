package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"gorm.io/gorm"
)

var DropDepositWalletIdColumnFromProductsTable *gormigrate.Migration = &gormigrate.Migration{
	ID: "1707152876",
	Migrate: func(db *gorm.DB) error {
		if db.Migrator().HasTable("products") {
			if db.Migrator().HasColumn(&models.Product{}, "deposit_wallet_id") {
				return db.Migrator().DropColumn(&models.Product{}, "deposit_wallet_id")
			}
		}
		return nil
	},
}
