package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"gorm.io/gorm"
)

var DropPublicKeysColumnFromProductsTable *gormigrate.Migration = &gormigrate.Migration{
	ID: "1696672968",
	Migrate: func(db *gorm.DB) error {
		if db.Migrator().HasTable("products") {
			if db.Migrator().HasColumn(&models.Product{}, "public_key") {
				return db.Migrator().DropColumn(&models.Product{}, "public_key")
			}
		}
		return nil
	},
}
