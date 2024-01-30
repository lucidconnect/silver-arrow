package api

import (
	"github.com/lucidconnect/silver-arrow/core"
)

/*
* NewProduct example object

	{
		"name": "xyz",
		"receivingAddress": "0xabc123",
		"firstChargeNow": true,
		"priceData": {
			"type": "recurring",
			"token": "usdc",
			"chain": 80001,
			"amount": 15,
			"interval": "day",
			"intervalCount": 1,
			"productId": "1234-abc-xyz",
		}
		"depositAddressData": [{
			"productId": "1234-abc-xyz",
			"merchantId": "1234-abc-xyz",
			"walletAddress": "0xabc123",
			"active": true
		},
		{
			"productId": "1234-abc-xyz",
			"merchantId": "1234-abc-xyz",
			"walletAddress": "0xabc123",
			"active": true
		}]
	}
*/
type NewProduct struct {
	Name             string   `json:"name"`
	ReceivingAddress string   `json:"receivingAddress"`
	FirstChargeNow   bool     `json:"firstChargeNow"`
	PriceData        NewPrice `json:"priceData"`
}

type NewPrice struct {
	Type          core.PriceType        `json:"type"`
	Token         string                `json:"token"`
	Chain         int64                 `json:"chain"`
	Amount        float64               `json:"amount"`
	Interval      core.RecuringInterval `json:"interval"`
	IntervalCount int64                 `json:"intervalCount"`
	ProductID     string                `json:"productId,omitempty"`
	TrialPeriod   int                   `json:"trialPeriod,omitempty"`
}

type PriceDataResponse struct {
	ID            string                `json:"id"`
	Active        bool                  `json:"active"`
	Type          core.PriceType        `json:"type"`
	Token         string                `json:"token"`
	Chain         int64                 `json:"chain"`
	Amount        int64                 `json:"amount"`
	Interval      core.RecuringInterval `json:"interval"`
	IntervalCount int64                 `json:"intervalCount"`
	ProductID     string                `json:"productId,omitempty"`
	TrialPeriod   int                   `json:"trialPeriod,omitempty"`
}

type ProductResponse struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	ReceivingAddress string            `json:"receivingAddress"`
	FirstChargeNow   bool              `json:"firstChargeNow"`
	DefaultPriceData PriceDataResponse `json:"priceData"`
}
