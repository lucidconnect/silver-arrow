package merchant

import (
	"encoding/base64"
	"fmt"
	"strings"

	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/graphql/merchant/graph/model"
	"github.com/lucidconnect/silver-arrow/repository/models"
)

// ProductId is base64 encoded

func (m *MerchantService) CreateProduct(input model.NewProduct) (*model.Product, error) {
	productID := uuid.New()

	merchant, err := m.repository.FetchMerchantByAddress(input.Owner)
	if err != nil {
		return nil, err
	}
	chainId := int64(input.Chain)
	product := &models.Product{
		ID:             productID,
		Name:           input.Name,
		Chain:          chainId,
		Owner:          input.Owner,
		Token:          input.Token,
		DepositAddress: input.ReceivingAddress,
		MerchantID:     merchant.ID,
		CreatedAt:      time.Now(),
		Mode:           model.ModeTest,
	}
	if err := m.repository.CreateProduct(product); err != nil {
		log.Err(err).Send()
		return nil, err
	}

	productObj := &model.Product{
		Name:             input.Name,
		Chain:            input.Chain,
		Owner:            input.Owner,
		Token:            input.Token,
		ProductID:        productID.String(),
		ReceivingAddress: input.ReceivingAddress,
	}
	return productObj, nil
}

func (m *MerchantService) FetchProductsByOwner(owner string) ([]*model.Product, error) {
	var products []*model.Product
	ms, err := m.repository.FetchProductsByOwner(owner)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	for _, v := range ms {
		subscriptions, err := parseMerchantSubscriptions(v.Subscriptions)
		if err != nil {
			log.Err(err).Send()
			continue
		}
		product := &model.Product{
			Name:             v.Name,
			Owner:            v.Owner,
			Chain:            int(v.Chain),
			ProductID:        v.ID.String(),
			MerchantID:       v.MerchantID.String(),
			ReceivingAddress: v.DepositAddress,
			Subscriptions:    subscriptions,
		}
		products = append(products, product)
	}
	return products, nil
}

func parseMerchantSubscriptions(subs []models.Subscription) ([]*model.Sub, error) {
	var subscriptions []*model.Sub

	for _, sub := range subs {
		subscription := &model.Sub{
			Chain:         int(sub.Chain),
			Token:         sub.Token,
			Amount:        int(sub.Amount),
			Active:        sub.Active,
			Interval:      fmt.Sprintf("%v days", sub.Interval),
			WalletAddress: sub.WalletAddress,
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (m *MerchantService) FetchProduct(pid string) (*model.Product, error) {
	id, err := uuid.Parse(pid)
	if err != nil {
		id, err = parseUUID(pid)
		if err != nil {
			log.Err(err).Send()
			return nil, errors.New("invalid product id")
		}
	}
	v, _ := m.repository.FetchProduct(id)

	subscriptions, err := parseMerchantSubscriptions(v.Subscriptions)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}
	createdAt := v.CreatedAt.Format(time.RFC3339)
	product := &model.Product{
		Name:             v.Name,
		Mode:             v.Mode,
		Owner:            v.Owner,
		Chain:            int(v.Chain),
		ProductID:        pid,
		MerchantID:       v.MerchantID.String(),
		ReceivingAddress: v.DepositAddress,
		CreatedAt:        &createdAt,
		Subscriptions:    subscriptions,
	}

	return product, nil
}

func (m *MerchantService) UpdateProductMode(merchantId uuid.UUID, productId, mode string) error {
	id := ParseUUID(productId)

	update := map[string]interface{}{
		"mode": mode,
	}
	err := m.repository.UpdateProduct(id, merchantId, update)
	if err != nil {
		log.Err(err).Send()
		return err
	}
	return nil
}

func Base64EncodeUUID(id uuid.UUID) (string, error) {
	b, err := id.MarshalBinary()
	if err != nil {
		err = errors.Wrap(err, "marshalling uuid failed")
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func parseUUID(mid string) (uuid.UUID, error) {
	b, err := base64.RawURLEncoding.DecodeString(mid)
	if err != nil {
		return uuid.Nil, err
	}
	id, err := uuid.FromBytes(b)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func removeUnderscore(s string) string {
	strArr := strings.Split(s, "_")
	return strings.ToTitle(strings.Join(strArr, ""))
}
