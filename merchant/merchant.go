package merchant

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/helicarrierstudio/silver-arrow/graphql/merchant/graph/model"
	"github.com/helicarrierstudio/silver-arrow/repository"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
)

type MerchantService struct {
	repository      repository.Database
	client          *ethclient.Client
	ContractAddress common.Address
}

func NewMerchantService(c *ethclient.Client, r repository.Database, address string) *MerchantService {
	return &MerchantService{
		repository:      r,
		client:          c,
		ContractAddress: common.HexToAddress(address),
	}
}

func (m *MerchantService) CreateMerchant(input model.NewMerchant) (*model.Merchant, error) {
	id := uuid.New()

	chainId := int64(input.Chain)
	merchant := &models.Merchant{
		ID:             id,
		Name:           input.Name,
		Chain:          chainId,
		Owner:          input.Owner,
		Token:          input.Token,
		DepositAddress: input.ReceivingAddress,
	}
	if err := m.repository.CreateMerchant(merchant); err != nil {
		return nil, err
	}
	merchantId, _ := EncodeUUIDToMerchantId(id)

	merchantObj := &model.Merchant{
		Name:             input.Name,
		Chain:            input.Chain,
		Owner:            input.Owner,
		Token:            input.Token,
		MerchantID:       merchantId,
		ReceivingAddress: input.ReceivingAddress,
	}
	return merchantObj, nil
}

func (m *MerchantService) UpdateMerchant() {}

func (m *MerchantService) FetchMerchantsByOwner(owner string) ([]*model.Merchant, error) {
	var merchants []*model.Merchant
	ms, err := m.repository.FetchMerchanstByOwner(owner)
	if err != nil {
		return nil, err
	}

	for _, v := range ms {
		merchantId, _ := EncodeUUIDToMerchantId(v.ID)
		subscriptions, err := fetchMerchantSubscriptions(m.repository, merchantId)
		if err != nil {
			log.Println(err)
		}
		merchant := &model.Merchant{
			Name:             v.Name,
			Owner:            v.Owner,
			Chain:            int(v.Chain),
			MerchantID:       merchantId,
			ReceivingAddress: v.DepositAddress,
			Subscriptions:    subscriptions,
		}
		merchants = append(merchants, merchant)
	}
	return merchants, nil
}

func fetchMerchantSubscriptions(repo repository.Database, merchant string) ([]*model.Sub, error) {
	var subscriptions []*model.Sub
	subs, err := repo.FindSubscriptionByMerchant(merchant)
	if err != nil {
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

func (m *MerchantService) FetchMerchant(mid string) (*model.Merchant, error) {
	id := ParseMerchantIdtoUUID(mid)
	v, _ := m.repository.FetchMerchant(id)

	merchant := &model.Merchant{
		Name:             v.Name,
		Owner:            v.Owner,
		Chain:            int(v.Chain),
		MerchantID:       mid,
		ReceivingAddress: v.DepositAddress,
	}

	return merchant, nil
}

func EncodeUUIDToMerchantId(id uuid.UUID) (string, error) {
	b, err := id.MarshalBinary()
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func ParseMerchantIdtoUUID(mid string) uuid.UUID {
	b, _ := base64.RawStdEncoding.DecodeString(mid)
	id, _ := uuid.FromBytes(b)
	return id
}

func removeUnderscore(s string) string {
	strArr := strings.Split(s, "_")
	return strings.ToTitle(strings.Join(strArr, ""))
}
