package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"gorm.io/gorm"
)

var DropDepositAddressColumnFromProductsTable *gormigrate.Migration = &gormigrate.Migration{
	ID: "1707154122",
	Migrate: func(db *gorm.DB) error {
		if db.Migrator().HasTable("products") {
			if db.Migrator().HasColumn(&models.Product{}, "deposit_address") {
				return db.Migrator().DropColumn(&models.Product{}, "deposit_address")
			}
		}
		return nil
	},
}
