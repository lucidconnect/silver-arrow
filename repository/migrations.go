package repository

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/lucidconnect/silver-arrow/repository/migrations"
	"github.com/rs/zerolog/log"
)

func (db *PostgresDB) RunMigrations() {
	migrator := gormigrate.New(db.Db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.DropPublicKeysColumnFromProductsTable,
		migrations.DropDepositAddressColumnFromProductsTable,
		migrations.DropDepositWalletIdColumnFromProductsTable,
	})

	if err := migrator.Migrate(); err != nil {
		log.Err(err).Msg("Unable to run migrations")
	}
}
