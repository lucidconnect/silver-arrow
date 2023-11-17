package graph

import (
	"context"

	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func getAuthenticatedAndActiveMerchant(ctx context.Context) (*models.Merchant, error) {
	merchant, err := auth.AuthMerchantContext(ctx)
	if err != nil {
		err = errors.Wrapf(err, "merchant authorization failed %v", ctx)
		log.Err(err).Send()
		return nil, err
	}

	return merchant, nil
}
