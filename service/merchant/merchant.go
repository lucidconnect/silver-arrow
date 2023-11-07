package merchant

import (
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/api/graphql/merchant/graph/model"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type MerchantService struct {
	repository repository.Database
}

func NewMerchantService(r repository.Database) *MerchantService {
	return &MerchantService{
		repository: r,
	}
}

func (m *MerchantService) CreateAccessKeys(owner string) (*model.AccessKey, error) {
	pk, sk, err := auth.CreateAccessKey()
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	accessKey := &model.AccessKey{
		PublicKey:  pk,
		PrivateKey: sk,
	}
	// check if the merchant exists, if not create a new entry
	merchant, err := m.repository.FetchMerchantByAddress(owner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// create a new entry
			merchant = &models.Merchant{
				ID:           uuid.New(),
				PublicKey:    pk,
				OwnerAddress: owner,
			}
			err := m.repository.AddMerchant(merchant)
			if err != nil {
				log.Err(err).Send()
				return nil, err
			}
		} else {
			log.Err(err).Send()
			return nil, err
		}
	} else {
		err = m.repository.UpdateMerchantKey(merchant.ID, pk)
		if err != nil {
			log.Err(err).Send()
			return nil, err
		}
	}
	return accessKey, nil
}

func (m *MerchantService) FetchMerchantKey(owner string) (string, error) {
	merchant, err := m.repository.FetchMerchantByAddress(owner)
	if err != nil {
		log.Err(err).Send()
		return "", err
	}

	return merchant.PublicKey, nil
}

func (m *MerchantService) SummarizeMerchant(owner string) (*model.MerchantStats, error) {
	var totalSubscriptions, totalUsers int
	products, err := m.repository.FetchProductsByOwner(owner)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	nProduct := len(products)

	for _, product := range products {
		subs := len(product.Subscriptions)
		totalSubscriptions += subs
	}
	totalUsers = totalSubscriptions
	stats := &model.MerchantStats{
		Users:         totalUsers,
		Products:      nProduct,
		Subscriptions: totalSubscriptions,
	}
	return stats, nil
}
