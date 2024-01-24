package merchant

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/conversions"
	"github.com/lucidconnect/silver-arrow/core"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/graphql/merchant/graph/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Price struct {
	ID           uuid.UUID
	Active       bool
	Amount       int64
	Token        string
	Chain        int64
	Type         core.PriceType
	IntervalUnit core.RecuringInterval
	Interval     int64
	TrialPeriod  int64
	ProductId    string
	MerchantId   string
	CreatedAt    int64
}

func (m *MerchantService) CreatePrice(price *Price, productId string) (*Price, error) {
	id := uuid.New()

	pid, err := uuid.Parse(productId)
	if err != nil {
		return nil, err
	}
	priceObject := &models.Price{
		ID:           id,
		Active:       price.Active,
		Type:         string(price.Type),
		Interval:     price.Interval,
		IntervalUnit: string(price.IntervalUnit),
		TrialPeriod:  price.TrialPeriod,
		ProductID:    pid,
		MerchantID:   m.merchant,
		Amount:       price.Amount,
		CreatedAt:    time.Now(),
	}
	if err = m.repository.CreatePrice(priceObject); err != nil {
		log.Err(err).Caller().Send()
		return nil, fmt.Errorf("failed to create price for product [%v]", price.ProductId)
	}

	return price, nil
}

func (m *MerchantService) RetrievePriceData(priceId string) (*model.PriceData, error) {
	pid, err := uuid.Parse(priceId)
	if err != nil {
		errMsg := fmt.Errorf("invalid uuid [%v] supplied", priceId)
		log.Err(err).Caller().Msg(errMsg.Error())
		return nil, errMsg
	}

	price, err := m.repository.FetchPrice(pid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errMsg := fmt.Errorf("price [%v] not found", pid)
			return nil, errMsg
		} else {
			errMsg := fmt.Errorf("error occured fetching price data")
			log.Err(err).Msg(errMsg.Error())
			return nil, errMsg
		}
	}
	// priceData := &Price{
	// 	ID:            price.ID,
	// 	Active:        price.Active,
	// 	Token:         price.Token,
	// 	Amount:        price.Amount,
	// 	Type:          PriceType(price.Type),
	// 	Interval:      RecuringInterval(price.Interval),
	// 	IntervalCount: price.IntervalCount,
	// 	TrialPeriod:   price.TrialPeriod,
	// 	ProductId:     price.ProductID.String(),
	// 	CreatedAt:     price.CreatedAt.Unix(),
	// }
	amount := conversions.ParseTransferAmountFloat(price.Token, price.Amount)

	priceData := &model.PriceData{
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
	return priceData, nil
}

func (m *MerchantService) RetrieveProductPriceData(productId string) ([]*model.PriceData, error) {
	pid, err := uuid.Parse(productId)
	if err != nil {
		errMsg := fmt.Errorf("invalid uuid [%v] supplied", pid)
		log.Err(err).Caller().Msg(errMsg.Error())
		return nil, errMsg
	}

	priceObjects, err := m.repository.FetchAllPricesByProduct(pid)
	if err != nil {
		log.Err(err).Send()
		return nil, err
	}

	return ParsePriceObjects(priceObjects)
}
