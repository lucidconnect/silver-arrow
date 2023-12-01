package repository

import (
	"context"
	"testing"

	"github.com/lucidconnect/silver-arrow/testhelpers"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

func Test_WalletRepo_AddSubscription(t *testing.T) {
	// type fields struct {
	// 	db *gorm.DB
	// }
	// type args struct {
	// 	ctx   context.Context
	// 	input models.Subscription
	// }

	// tests := []struct {
	// 	name       string
	// 	args       args
	// 	beforeTest func(sqlmock.Sqlmock)
	// 	wantErr bool
	// }{
	// 	{
	// 		name: "fail create subscription",
	// 		args: args{
	// 			ctx: context.TODO(),
	// 			input: models.Subscription{Token: "USDC", Amount: 10, Active: true, Interval: 1, UserOpHash: "opHash", MerchantId: "mId",},
	// 		},
	// 	},

	// }
}

type PostgresRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *PostgresDB
	ctx         context.Context
}

func (suite *PostgresRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	suite.pgContainer = pgContainer
	suite.repository = NewPostgresDB(pgContainer.Db)
}

func (suite *PostgresRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatal().Err(err).Msg("error terminating postgres container")
	}
}