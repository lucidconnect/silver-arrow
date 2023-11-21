package merchant

import (
	"crypto/rand"
	"errors"
	"math/big"
	"os"
	"strings"

	convoy "github.com/frain-dev/convoy-go"
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/graphql/merchant/graph/model"
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var alphaNumericRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(length int) string {
	runeLength := len(alphaNumericRunes)
	var b strings.Builder
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(runeLength)))
		b.WriteRune(letterRunes[n.Int64()])
	}
	str := b.String()
	return str
}

func (m *MerchantService) CreateMerchant(input model.NewMerchant) (*model.Merchant, error) {
	id := uuid.New()

	token := RandStringRunes(32)
	convoyClient := convoy.New(convoy.Options{
		APIKey:    os.Getenv("CONVOY_API_KEY"),
		ProjectID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	endpoint, err := convoyClient.Endpoints.Create(&convoy.CreateEndpointRequest{
		Secret:       token,
		URL:          *input.WebHookURL,
		Description:  input.Name + "'s default endpoint",
		SupportEmail: input.Email,
	}, &convoy.EndpointQueryParam{
		GroupID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	if err != nil {
		log.Err(err).Msg("failed to create app endpoint")
		return nil, errors.New("failed to create merchant's endpoint on Convoy")
	}

	_, err = convoyClient.Subscriptions.Create(&convoy.CreateSubscriptionRequest{
		Name:       input.Name + "'s default subscription",
		EndpointID: endpoint.UID,
		FilterConfig: &convoy.FilterConfiguration{
			EventTypes: []string{"*"},
		},
	})

	if err != nil {
		log.Err(err).Msgf("failed to create convoy subscription for merchant with convoy endpoint id %v.", endpoint.UID)
	}

	merchant := &models.Merchant{
		ID:               id,
		Name:             input.Name,
		Email:            input.Email,
		OwnerAddress:     input.Owner,
		WebhookToken:     token,
		ConvoyEndpointID: endpoint.UID,
		WebhookUrl:       *input.WebHookURL,
	}

	if err := m.repository.AddMerchant(merchant); err != nil {
		log.Err(err).Msg("creating merchant failed")
		return nil, errors.New("merchant creation failed")
	}

	merchantObj := &model.Merchant{
		ID:    id.String(),
		Name:  input.Name,
		Email: input.Email,
	}
	return merchantObj, nil
}

func (m *MerchantService) UpdateMerchantWebhook(merchant models.Merchant, webHookUrl string) (*model.Merchant, error) {
	convoyClient := convoy.New(convoy.Options{
		APIKey:    os.Getenv("CONVOY_API_KEY"),
		ProjectID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	_, err := convoyClient.Endpoints.Update(merchant.ConvoyEndpointID, &convoy.CreateEndpointRequest{
		Secret:       merchant.WebhookToken,
		URL:          webHookUrl,
		Description:  merchant.Name + "'s default endpoint",
		SupportEmail: merchant.Email,
	}, &convoy.EndpointQueryParam{
		GroupID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	if err != nil {
		log.Err(err).Msg("failed to update app endpoint")
		return nil, errors.New("failed to update merchant's endpoint on Convoy")
	}

	if err := m.repository.UpdateMerchantWebhookUrl(merchant.ID, webHookUrl); err != nil {
		log.Err(err).Send()
		return nil, errors.New("updating merchant webhook url failed")
	}
	return &model.Merchant{
		ID:         merchant.ID.String(),
		Name:       merchant.Name,
		PublicKey:  merchant.PublicKey,
		WebHookURL: webHookUrl,
	}, nil
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
