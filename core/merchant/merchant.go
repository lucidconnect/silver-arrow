package merchant

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	convoy "github.com/frain-dev/convoy-go"
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/graphql/merchant/graph/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type MerchantService struct {
	repository repository.Database
	merchant   uuid.UUID // id of the active merchant
}

type DepositAddress struct {
	ID            uuid.UUID
	MerchantID    string
	WalletAddress string
	Percentage    float64
	Active        bool
	Note          string
}

func NewMerchantService(r repository.Database, merchant uuid.UUID) *MerchantService {
	return &MerchantService{
		repository: r,
		merchant:   merchant,
	}
}

var alphaNumericRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(length int) string {
	runeLength := len(alphaNumericRunes)
	var b strings.Builder
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(runeLength)))
		b.WriteRune(alphaNumericRunes[n.Int64()])
	}
	str := b.String()
	return str
}

func (m *MerchantService) CreateMerchant(input model.NewMerchant) (*model.Merchant, error) {
	id := uuid.New()

	token := RandStringRunes(32)
	// convoyClient := convoy.New(convoy.Options{
	// 	APIKey:    os.Getenv("CONVOY_API_KEY"),
	// 	ProjectID: os.Getenv("CONVOY_PROJECT_ID"),
	// })

	// endpoint, err := convoyClient.Endpoints.Create(&convoy.CreateEndpointRequest{
	// 	Secret:       token,
	// 	Description:  input.Name + "'s default endpoint",
	// 	SupportEmail: input.Email,
	// }, &convoy.EndpointQueryParam{
	// 	GroupID: os.Getenv("CONVOY_PROJECT_ID"),
	// })

	// if err != nil {
	// 	log.Err(err).Msg("failed to create app endpoint")
	// 	return nil, errors.New("failed to create merchant's endpoint on Convoy")
	// }

	// _, err = convoyClient.Subscriptions.Create(&convoy.CreateSubscriptionRequest{
	// 	Name:       input.Name + "'s default subscription",
	// 	EndpointID: endpoint.UID,
	// 	FilterConfig: &convoy.FilterConfiguration{
	// 		EventTypes: []string{"*"},
	// 	},
	// })

	// if err != nil {
	// 	log.Err(err).Msgf("failed to create convoy subscription for merchant with convoy endpoint id %v.", endpoint.UID)
	// }

	merchant := &models.Merchant{
		ID:           id,
		Name:         input.Name,
		Email:        input.Email,
		OwnerAddress: input.Owner,
		WebhookToken: token,
		CreatedAt:    time.Now(),
		// ConvoyEndpointID: endpoint.UID,
		// TestPublicKey: key.PublicKey,
	}

	if err := m.repository.AddMerchant(merchant); err != nil {
		log.Err(err).Msg("creating merchant failed")
		return nil, errors.New("merchant creation failed")
	}

	key, err := m.CreateAccessKeys(input.Owner, model.ModeTest.String())
	if err != nil {
		log.Err(err).Msg("creating merchant test keys failed")
		return nil, errors.New("creating merchant test keys failed")
	}

	merchantObj := &model.Merchant{
		ID:        id.String(),
		Name:      input.Name,
		Email:     input.Email,
		AccessKey: key,
	}
	return merchantObj, nil
}

func (m *MerchantService) UpdateMerchantWebhook(merchant models.Merchant, webHookUrl string) (*model.Merchant, error) {
	convoyClient := convoy.New(convoy.Options{
		APIKey:    os.Getenv("CONVOY_API_KEY"),
		ProjectID: os.Getenv("CONVOY_PROJECT_ID"),
	})

	if merchant.ConvoyEndpointID != "" {
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
	} else {
		endpoint, err := convoyClient.Endpoints.Create(&convoy.CreateEndpointRequest{
			Name:         merchant.Name,
			Secret:       merchant.WebhookToken,
			Description:  merchant.Name + "'s default endpoint",
			SupportEmail: merchant.Email,
			URL:          webHookUrl,
			OwnerID:      merchant.ID.String(),
		}, &convoy.EndpointQueryParam{
			GroupID: os.Getenv("CONVOY_PROJECT_ID"),
		})

		if err != nil {
			log.Err(err).Msg("failed to create app endpoint")
			return nil, errors.New("failed to create merchant's endpoint on Convoy")
		}

		_, err = convoyClient.Subscriptions.Create(&convoy.CreateSubscriptionRequest{
			Name:       merchant.Name + "'s default subscription",
			EndpointID: endpoint.UID,
			FilterConfig: &convoy.FilterConfiguration{
				EventTypes: []string{"*"},
			},
		})

		if err != nil {
			log.Err(err).Msgf("failed to create convoy subscription for merchant with convoy endpoint id %v.", endpoint.UID)
		}
	}

	if err := m.repository.UpdateMerchantWebhookUrl(merchant.ID, webHookUrl); err != nil {
		log.Err(err).Send()
		return nil, errors.New("updating merchant webhook url failed")
	}
	return &model.Merchant{
		ID:           merchant.ID.String(),
		Name:         merchant.Name,
		WebHookURL:   webHookUrl,
		WebhookToken: merchant.WebhookToken,
	}, nil
}

func (m *MerchantService) CreateAccessKeys(owner, mode string) (*model.MerchantAccessKey, error) {
	pk, sk, err := auth.CreateAccessKey()
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	accessKey := &model.MerchantAccessKey{
		Mode:       model.Mode(mode),
		PublicKey:  pk,
		PrivateKey: sk,
	}
	merchant, err := m.repository.FetchMerchantByAddress(owner)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("merchant not found")
		} else {
			log.Err(err).Send()
			return nil, err
		}
	} else {
		// if an existing key exists for the mode, delete it and create a new one
		keys := merchant.MerchantAccessKeys
		var targetKey models.MerchantAccessKey
		for _, key := range keys {
			if key.Mode == mode {
				targetKey = key
			}
		}

		err = m.repository.DeleteMerchantAccessKey(targetKey.ID, &targetKey)
		if err != nil {
			log.Err(err).Send()
			return nil, err
		}
		newKey := &models.MerchantAccessKey{
			Mode:       mode,
			PublicKey:  pk,
			MerchantID: merchant.ID,
			CreatedAt:  time.Now(),
		}
		err = m.repository.CreateMerchantAccessKeys(newKey)
		if err != nil {
			log.Err(err).Send()
			return nil, err
		}
	}
	return accessKey, nil
}

func (m *MerchantService) FetchMerchantKey(owner, mode string) (string, error) {
	merchant, err := m.repository.FetchMerchantByAddress(owner)
	if err != nil {
		log.Err(err).Send()
		return "", err
	}
	keys := merchant.MerchantAccessKeys
	var targetKey models.MerchantAccessKey
	for _, key := range keys {
		if key.Mode == mode {
			targetKey = key
		}
	}
	key := targetKey.PublicKey
	return key, nil
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

func (m *MerchantService) AddDepositAddress(address *DepositAddress) (*DepositAddress, error) {
	id := uuid.New()
	depositWallet := &models.DepositWallet{
		ID:            id,
		MerchantID:    m.merchant,
		WalletAddress: address.WalletAddress,
		Percentage:    address.Percentage,
		Active:        address.Active,
		Note:          address.Note,
	}

	if err := m.repository.AddDepositWallet(depositWallet); err != nil {
		log.Err(err).Caller().Send()
		return nil, fmt.Errorf("failed to add merchant deposit wallet")
	}

	address.ID = id

	return address, nil
}

func (m *MerchantService) UpdateDepositWallet(id string, wallet *DepositAddress) error {
	walletId, _ := uuid.Parse(id)
	update := &models.DepositWallet{
		WalletAddress: wallet.WalletAddress,
	}

	if err := m.repository.UpdateDepositWallet(walletId, update); err != nil {
		log.Err(err).Caller().Send()
		return fmt.Errorf("failed to update merchant deposit detail")
	}

	return nil
}

func (m *MerchantService) ListDepositWallets(merchantId string) ([]*model.DepositWallet, error) {
	var depositAddresses []*model.DepositWallet
	mid, _ := uuid.Parse(merchantId)
	depositWallets, err := m.repository.FetchDepositWalletByMerchant(mid)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	for _, wallet := range depositWallets {
		depositAddress := &model.DepositWallet{
			ID: wallet.ID.String(),
		}

		depositAddresses = append(depositAddresses, depositAddress)
	}

	return depositAddresses, nil
}

func (m *MerchantService) DeleteDepositWallet(walletId string) error {
	wid, _ := uuid.Parse(walletId)

	if err := m.repository.DeleteDepositWallet(wid); err != nil {
		log.Err(err).Caller().Send()
		return err
	}

	return nil
}
