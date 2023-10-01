package graph

import (
	"context"

	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"github.com/helicarrierstudio/silver-arrow/service/merchant"
	"github.com/pkg/errors"
)

func getAuthenticatedAndActiveMerchant(ctx context.Context) (*models.Merchant, error) {
	merchant, err := merchant.ForContext(ctx)
	if err != nil {
		err = errors.Wrapf(err, "merchant authorization failed %v", ctx)
		return nil, err
	}

	return merchant, nil
}
