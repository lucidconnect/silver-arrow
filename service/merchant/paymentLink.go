package merchant

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/conversions"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/graphql/merchant/graph/model"
	"github.com/rs/zerolog/log"
)

type PaymentLinkQueryParams struct {
	PaymentLinkId *string
	ProductId     *string
}

func (m *MerchantService) CreatePaymentLink(input model.NewPaymentLink) (string, error) {
	id := uuid.New()
	productId, err := uuid.Parse(input.ProductID)
	if err != nil {
		log.Err(err).Msg("invalid product id")
		return "", err
	}
	product, err := m.repository.FetchProduct(productId)
	if err != nil {
		log.Err(err).Msgf("failed to fetch product with id %v", productId)
		return "", err
	}

	newPaymentLink := &models.PaymentLink{
		ID:          id,
		MerchantID:  product.MerchantID,
		CallbackURL: input.CallbackURL,
		ProductID:   productId,
		Product:     *product,
		CreatedAt:   time.Now(),
	}

	err = m.repository.CreatePaymentLink(newPaymentLink)
	if err != nil {
		log.Err(err).Send()
		return "", err
	}

	return id.String(), nil
}

func (m *MerchantService) FetchPaymentLink(queryParam PaymentLinkQueryParams) (*model.PaymentLinkDetails, error) {
	var paymentLink *models.PaymentLink
	if queryParam.PaymentLinkId != nil {
		id, err := uuid.Parse(*queryParam.PaymentLinkId)
		if err != nil {
			log.Err(err).Msg("parsing uuid failed")
			return nil, err
		}
		paymentLink, err = m.repository.FetchPaymentLink(id)
		if err != nil {
			log.Err(err).Msgf("failed to fetch payment link with id %v", id)
			return nil, err
		}
	} else {
		id, err := uuid.Parse(*queryParam.ProductId)
		if err != nil {
			log.Err(err).Msg("parsing uuid failed")
			return nil, err
		}
		paymentLink, err = m.repository.FetchPaymentLink(id)
		if err != nil {
			log.Err(err).Msgf("failed to fetch payment link with id %v", id)
			return nil, err
		}
	}

	interval := conversions.ParseNanoSecondsToDay(paymentLink.Product.Interval)
	amount := conversions.ParseTransferAmountFloat(paymentLink.Product.Token, paymentLink.Product.Amount)

	paymentLinkDetalais := &model.PaymentLinkDetails{
		ID: paymentLink.ID.String(),
		Mode: paymentLink.Product.Mode,
		ProductID:    paymentLink.ProductID.String(),
		MerchantID:   paymentLink.MerchantID.String(),
		MerchantName: paymentLink.MerchantName,
		ProductName:  paymentLink.Product.Name,
		Interval:     int(interval),
		CallbackURL:  paymentLink.CallbackURL,
		Amount:       amount,
		Token:        paymentLink.Product.Token,
		Chain:        int(paymentLink.Product.Chain),
	}

	return paymentLinkDetalais, nil
}

// func nanoSecondsToDay(ns int64) int64 {
// 	interval := time.Duration(ns)
// 	hours := interval.Hours()

// 	days := hours / 24
// 	return int64(days)
// }

// func parseTransferAmountFloat(token string, amount int64) float64 {
// 	var divisor int
// 	if token == "USDC" || token == "USDT" {
// 		divisor = 6
// 	} else {
// 		divisor = 18
// 	}
// 	minorFactor := math.Pow10(divisor)
// 	parsedAmount := float64(amount) / minorFactor

// 	return parsedAmount
// }
