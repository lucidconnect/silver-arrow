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
	DepositAddress string
}

func (m *MerchantService) CreateProduct(product *Product) (*Product, error) {
	productID := uuid.New()

	// merchant, err := m.repository.FetchMerchantByAddress(input.Owner)
	// if err != nil {
	// 	return nil, err
	// }
	chainId := int64(product.Chain)
	// amount := conversions.ParseFloatAmountToInt(input.Token, input.Amount)
	// interval := conversions.ParseDaysToNanoSeconds(int64(input.Interval))

	productObj := &models.Product{
		ID:             productID,
		Name:           product.Name,
		Chain:          chainId,
		Owner:          product.Owner,
		// Token:          product.Token,
		DepositAddress: product.DepositAddress,
		MerchantID:     m.merchant,
		CreatedAt:      time.Now(),
		Mode:           model.ModeTest.String(),
		// InstantCharge: product.InstantCharge,
		// PaymentType:   product.PaymentType,
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
			continue
		}
		// interval := conversions.ParseNanoSecondsToDay(v.Interval)

		product := &model.Product{
			Name:             v.Name,
			Mode:             model.Mode(v.Mode),
			Owner:            v.Owner,
			Chain:            int(v.Chain),
			// Token:            v.Token,
			// Interval:         int(interval),
			ProductID:        v.ID.String(),
			MerchantID:       v.MerchantID.String(),
			ReceivingAddress: v.DepositAddress,
			Subscriptions:    subscriptions,
		}
		defaultPrice := v.DefaultPriceID
		if defaultPrice != uuid.Nil {
			priceData, err := m.RetrievePriceData(defaultPrice.String())
			if err != nil {
				log.Err(err).Send()
			}
			amount := conversions.ParseTransferAmountFloat(priceData.Token, priceData.Amount)
			gqlPriceData := &model.PriceData{
				ID:            priceData.ID.String(),
				Type:          model.PaymentType(priceData.Type),
				Active:        priceData.Active,
				Amount:        amount,
				Token:         priceData.Token,
				Interval:      model.IntervalType(priceData.Interval),
				IntervalCount: int(priceData.IntervalCount),
				ProductID:     priceData.ProductId,
				TrialPeriod:   int(priceData.TrialPeriod),
			}
			product.PriceData = gqlPriceData
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
	// interval := conversions.ParseNanoSecondsToDay(v.Interval)
	// amount := conversions.ParseTransferAmountFloat(v.Token, v.Amount)
	product := &model.Product{
		Name:             v.Name,
		Mode:             model.Mode(v.Mode),
		Owner:            v.Owner,
		Chain:            int(v.Chain),
		// Token:            v.Token,
		// Interval:         int(interval),
		ProductID:        pid,
		MerchantID:       v.MerchantID.String(),
		ReceivingAddress: v.DepositAddress,
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

// func removeUnderscore(s string) string {
// 	strArr := strings.Split(s, "_")
// 	return strings.ToTitle(strings.Join(strArr, ""))
// }

// func parseFloatAmountToInt(token string, amount float64) int64 {
// 	var divisor int
// 	if token == "USDC" || token == "USDT" {
// 		divisor = 6
// 	} else {
// 		divisor = 18
// 	}
// 	minorFactor := math.Pow10(divisor)
// 	parsedAmount := int64(amount * minorFactor)

// 	return parsedAmount
// }

// func daysToNanoSeconds(days int64) time.Duration {
// 	nanoSsecondsInt := days * 24 * 60 * 60 * 1e9
// 	return time.Duration(nanoSsecondsInt)
// }

func ParseGraphqlInput(gqlInput model.NewProduct) *Product {
	p := &Product{
		Name:  gqlInput.Name,
		Chain: int64(gqlInput.Chain),
		// Mode:  gqlInput.Mode.String(),
		// Price: ,
		Owner:          gqlInput.Owner,
		DepositAddress: gqlInput.ReceivingAddress,
		Active:         true,
		// Mode: gqlInput.T,
		InstantCharge: gqlInput.FirstChargeNow,
		PaymentType:   gqlInput.PaymentType.String(),
	}
	return p
}
