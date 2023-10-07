package repository

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/helicarrierstudio/silver-arrow/repository/migrations"
)

func (db *DB) RunMigrations() {
	migrator := gormigrate.New(db.Db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.DropPublicKeysColumnFromProductsTable,
	})

	if err := migrator.Migrate(); err != nil {
		log.Fatalf("Unable to run migrations: %v", err)
	}
}
