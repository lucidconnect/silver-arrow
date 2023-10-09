package merchant

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/google/uuid"
	"github.com/helicarrierstudio/silver-arrow/graphql/merchant/graph/model"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

func (m *MerchantService) CreateProduct(input model.NewProduct) (*model.Product, error) {
	id := uuid.New()

	merchant, err := m.repository.FetchMerchantByAddress(input.Owner)
	if err != nil {
		return nil, err
	}
	chainId := int64(input.Chain)
	product := &models.Product{
		ID:             id,
		Name:           input.Name,
		Chain:          chainId,
		Owner:          input.Owner,
		Token:          input.Token,
		DepositAddress: input.ReceivingAddress,
		MerchantID:     merchant.ID,
		CreatedAt:      time.Now(),
	}
	if err := m.repository.CreateProduct(product); err != nil {
		log.Err(err).Send()
		return nil, err
	}

	productID, _ := Base64EncodeUUID(id)
	merchantObj := &model.Product{
		Name:             input.Name,
		Chain:            input.Chain,
		Owner:            input.Owner,
		Token:            input.Token,
		ProductID:        productID,
		ReceivingAddress: input.ReceivingAddress,
	}
	return merchantObj, nil
}

func (m *MerchantService) FetchProductsByOwner(owner string) ([]*model.Product, error) {
	var products []*model.Product
	ms, err := m.repository.FetchProductsByOwner(owner)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	for _, v := range ms {
		ProductID, _ := Base64EncodeUUID(v.ID)
		subscriptions, err := fetchMerchantSubscriptions(m.repository, ProductID)
		if err != nil {
			log.Err(err).Send()
			continue
		}
		product := &model.Product{
			Name:             v.Name,
			Owner:            v.Owner,
			Chain:            int(v.Chain),
			ProductID:        ProductID,
			ReceivingAddress: v.DepositAddress,
			Subscriptions:    subscriptions,
		}
		products = append(products, product)
	}
	return products, nil
}

func fetchMerchantSubscriptions(repo repository.Database, merchant string) ([]*model.Sub, error) {
	var subscriptions []*model.Sub
	subs, err := repo.FindSubscriptionByProduct(merchant)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

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
	id := ParseUUID(pid)
	v, _ := m.repository.FetchProduct(id)

	createdAt := v.CreatedAt.Format(time.RFC3339)
	merchant := &model.Product{
		Name:             v.Name,
		Owner:            v.Owner,
		Chain:            int(v.Chain),
		ProductID:        pid,
		ReceivingAddress: v.DepositAddress,
		CreatedAt:        &createdAt,
	}

	return merchant, nil
}

func Base64EncodeUUID(id uuid.UUID) (string, error) {
	b, err := id.MarshalBinary()
	if err != nil {
		err = errors.Wrap(err, "marshalling uuid failed")
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func ParseUUID(mid string) uuid.UUID {
	b, _ := base64.RawURLEncoding.DecodeString(mid)
	id, _ := uuid.FromBytes(b)
	return id
}

func removeUnderscore(s string) string {
	strArr := strings.Split(s, "_")
	return strings.ToTitle(strings.Join(strArr, ""))
}
