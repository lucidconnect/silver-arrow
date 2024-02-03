package merchant

import (
	"encoding/base64"
	"fmt"

	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/conversions"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/api"
	"github.com/lucidconnect/silver-arrow/server/graphql/merchant/graph/model"
)

// Product describes a good/service offered by a merchant/
// Each version of a good/service will be a separate product.
// Product can be used in conjuction with Price to configure pricing options
type Product struct {
	ID             uuid.UUID
	Name           string
	Chain          int64
	Token          string
	Active         bool
	CreatedAt      int64
	Mode           string
	Price          string // (optional) id for the price object
	InstantCharge  bool
	PaymentType    string
	Owner          string
	DepositAddress []*models.DepositWallet
}

func (m *MerchantService) CreateProduct(product *Product) (*Product, error) {
	productID := uuid.New()

	chainId := int64(product.Chain)
	depositWallets := product.DepositAddress
	productObj := &models.Product{
		ID:    productID,
		Name:  product.Name,
		Chain: chainId,
		Owner: product.Owner,
		MerchantID:     m.merchant,
		CreatedAt:      time.Now(),
		Mode:           model.ModeTest.String(),
		DepositWallets: depositWallets,
	}
	if err := m.repository.CreateProduct(productObj); err != nil {
		log.Err(err).Send()
		return nil, err
	}

	product.ID = productID
	return product, nil
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
		}
		priceObjects, err := m.repository.FetchAllPricesByProduct(v.ID)
		if err != nil {
			log.Err(err).Send()
		}
		priceData, _ := ParsePriceObjects(priceObjects)
		// interval := conversions.ParseNanoSecondsToDay(v.Interval)
		var depositWallets []*model.DepositWallet
		for _, wallet := range v.DepositWallets {
			w := ParseDepositWalletToGraphqlObject(*wallet)

			depositWallets = append(depositWallets, &w)
		}

		product := &model.Product{
			Name:             v.Name,
			Mode:             model.Mode(v.Mode),
			Owner:            v.Owner,
			ProductID:        v.ID.String(),
			MerchantID:       v.MerchantID.String(),
			ReceivingAddress: depositWallets,
			Subscriptions:    subscriptions,
			PriceData:        priceData,
		}
		products = append(products, product)
	}
	return products, nil
}

func (m *MerchantService) FetchProduct(pid string) (*model.Product, error) {
	id, err := uuid.Parse(pid)
	if err != nil {
		id, err = parseUUID(pid)
		if err != nil {
			log.Err(err).Msg("invalid product id")
			return nil, errors.New("invalid product id")
		}
	}
	v, err := m.repository.FetchProduct(id)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, errors.New("product not found")
	}

	subscriptions, err := parseMerchantSubscriptions(v.Subscriptions)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}
	createdAt := v.CreatedAt.Format(time.RFC3339)
	var depositWallets []*model.DepositWallet
	for _, wallet := range v.DepositWallets {
		w := ParseDepositWalletToGraphqlObject(*wallet)

		depositWallets = append(depositWallets, &w)
	}
	product := &model.Product{
		Name:             v.Name,
		Mode:             model.Mode(v.Mode),
		Owner:            v.Owner,
		ProductID:        pid,
		MerchantID:       v.MerchantID.String(),
		DefaultPrice:     v.DefaultPriceID.String(),
		ReceivingAddress: depositWallets,
		CreatedAt:        &createdAt,
		Subscriptions:    subscriptions,
	}

	return product, nil
}

func (m *MerchantService) UpdateProductMode(merchantId uuid.UUID, productId, mode string) error {
	id, err := uuid.Parse(productId)
	if err != nil {
		return err
	}

	var chainId int

	switch mode {
	case model.ModeLive.String():
		chainId = 10
	case model.ModeTest.String():
		chainId = 80001
	}
	update := map[string]interface{}{
		"mode":  mode,
		"chain": chainId,
	}

	err = m.repository.UpdateProduct(id, merchantId, update)
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

func ParseGraphqlInput(gqlInput model.NewProduct) *Product {
	var depositWallets []*models.DepositWallet
	depositAddresses := gqlInput.ReceivingAddress
	for _, address := range depositAddresses {
		wallet := models.DepositWallet{
			WalletAddress: address.Address,
			Percentage:    address.Percentage,
			Note:          address.Note,
		}
		depositWallets = append(depositWallets, &wallet)
	}
	p := &Product{
		Name: gqlInput.Name,
		Owner:          gqlInput.Owner,
		DepositAddress: depositWallets,
		Active:         true,
		InstantCharge: gqlInput.FirstChargeNow,
		PaymentType:   gqlInput.PriceData.Type.String(),
	}
	return p
}

func ParseNewApiProduct(input api.NewProduct) *Product {
	var depositWallets []*models.DepositWallet
	depositAddresses := input.ReceivingAddress
	for _, address := range depositAddresses {
		wallet := models.DepositWallet{
			WalletAddress: address.Address,
			Percentage:    address.Percentage,
			Note:          address.Note,
		}
		depositWallets = append(depositWallets, &wallet)
	}
	p := &Product{
		Name: input.Name,
		Owner:          input.Owner,
		DepositAddress: depositWallets,
		Active:         true,
		InstantCharge: input.FirstChargeNow,
		PaymentType:   string(input.PriceData.Type),
	}
	return p
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

func ParsePriceObjects(prices []models.Price) ([]*model.PriceData, error) {
	var priceData []*model.PriceData

	for _, price := range prices {
		amount := conversions.ParseTransferAmountFloat(price.Token, price.Amount)

		p := &model.PriceData{
			ID:           price.ID.String(),
			Type:         model.PaymentType(price.Type),
			Active:       price.Active,
			Amount:       amount,
			Token:        price.Token,
			Chain:        int(price.Chain),
			IntervalUnit: model.IntervalType(price.IntervalUnit),
			Interval:     int(price.Interval),
			ProductID:    price.ProductID.String(),
			TrialPeriod:  int(price.TrialPeriod),
		}
		priceData = append(priceData, p)
	}
	return priceData, nil
}

func ParseDepositWalletToGraphqlObject(wallet models.DepositWallet) model.DepositWallet {
	return model.DepositWallet{
		ID:         wallet.ID.String(),
		Address:    wallet.WalletAddress,
		Percentage: wallet.Percentage,
		Merchant:   wallet.MerchantID.String(),
		Note:       &wallet.Note,
	}
}
