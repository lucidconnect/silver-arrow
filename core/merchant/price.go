package merchant

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type PriceType string
type RecuringInterval string

const (
	PriceTypeOneTime      PriceType        = "one-time"
	PriceTypeRecurring    PriceType        = "recurring"
	RecuringIntervalDay   RecuringInterval = "day"
	RecuringIntervalWeek  RecuringInterval = "week"
	RecuringIntervalMonth RecuringInterval = "month"
	RecuringIntervalYear  RecuringInterval = "year"
)

type Price struct {
	ID            uuid.UUID
	Active        bool
	Amount        int64
	Token         string
	Type          PriceType
	Interval      RecuringInterval
	IntervalCount int64
	TrialPeriod   int64
	ProductId     string
	MerchantId    string
	CreatedAt     int64
}

func (m *MerchantService) CreatePrice(price *Price, productId string) (*Price, error) {
	id := uuid.New()

	pid, err := uuid.Parse(productId)
	if err != nil {
		return nil, err
	}
	priceObject := &models.Price{
		ID:            id,
		Active:        price.Active,
		Type:          string(price.Type),
		Interval:      string(price.Interval),
		IntervalCount: price.IntervalCount,
		TrialPeriod:   price.TrialPeriod,
		ProductID:     pid,
		MerchantID:    m.merchant,
		Amount:        price.Amount,
		CreatedAt:     time.Now(),
	}
	if err = m.repository.CreatePrice(priceObject); err != nil {
		log.Err(err).Caller().Send()
		return nil, fmt.Errorf("failed to create price for product [%v]", price.ProductId)
	}

	return price, nil
}

func (m *MerchantService) RetrievePriceData(priceId string) (*Price, error) {
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
	priceData := &Price{
		ID:            price.ID,
		Active:        price.Active,
		Token:         price.Token,
		Amount:        price.Amount,
		Type:          PriceType(price.Type),
		Interval:      RecuringInterval(price.Interval),
		IntervalCount: price.IntervalCount,
		TrialPeriod:   price.TrialPeriod,
		ProductId:     price.ProductID.String(),
		CreatedAt:     price.CreatedAt.Unix(),
	}
	return priceData, nil
}
