package testhelpers

import (
	"context"
	"fmt"
	"time"

	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	Db *gorm.DB
}

func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		// postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	dialector := gpostgres.Open(connStr)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database")
		return nil, err
	}

	if err = db.AutoMigrate(models.Payment{}, models.Wallet{}, models.Merchant{}, models.Key{}, models.Subscription{}, models.Product{}); err != nil {
		log.Fatal().Err(err).Msg("Error migrating database models")
	}
	return &PostgresContainer{
		PostgresContainer: pgContainer,
		Db:  db,
	}, nil
}
