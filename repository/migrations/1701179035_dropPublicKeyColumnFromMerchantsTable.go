package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"gorm.io/gorm"
)

var DropPublicKeyColumnFromMerchantsTable *gormigrate.Migration = &gormigrate.Migration{
	ID: "1701179035",
	Migrate: func(db *gorm.DB) error {
		if db.Migrator().HasTable("products") {
			if db.Migrator().HasColumn(&models.Merchant{}, "public_key") {
				return db.Migrator().DropColumn(&models.Merchant{}, "public_key")
			}
		}
		return nil
	},
}
