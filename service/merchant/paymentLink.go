package merchant

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/graphql/merchant/graph/model"
	"github.com/lucidconnect/silver-arrow/repository/models"
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
		ID:        id,
		MerhantId: product.MerchantID,
		ProductId: productId,
		Product:   *product,
		CreatedAt: time.Now(),
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

	paymentLinkDetalais := &model.PaymentLinkDetails{
		ID: paymentLink.ID.String(),
		// Mode: paymentLink.Product.Mode,
		ProductID:  paymentLink.ProductId.String(),
		MerchantID: paymentLink.MerhantId.String(),
		Amount:     int(paymentLink.Product.Amount),
		Token:      paymentLink.Product.Token,
		Chain:      int(paymentLink.Product.Chain),
	}

	return paymentLinkDetalais, nil
}
