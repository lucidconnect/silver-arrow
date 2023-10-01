package merchant

import (
	"github.com/google/uuid"
	"github.com/helicarrierstudio/silver-arrow/auth"
	"github.com/helicarrierstudio/silver-arrow/graphql/merchant/graph/model"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
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
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		err = m.repository.UpdateMerchantKey(merchant.ID, pk)
		if err != nil {
			return nil, err
		}
	}
	return accessKey, nil
}

