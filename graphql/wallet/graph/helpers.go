package graph

import (
	"context"
	"os"
	"strconv"

	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"github.com/helicarrierstudio/silver-arrow/service/merchant"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func getAuthenticatedAndActiveMerchant(ctx context.Context) (*models.Merchant, error) {
	useAuthStr := os.Getenv("USE_AUTH")
	useAuth, _ := strconv.ParseBool(useAuthStr)
	if !useAuth {
		return &models.Merchant{}, nil
	}

	merchant, err := merchant.ForContext(ctx)
	if err != nil {
		err = errors.Wrapf(err, "merchant authorization failed %v", ctx)
		log.Err(err).Send()
		return nil, err
	}

	return merchant, nil
}
