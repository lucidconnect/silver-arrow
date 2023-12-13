package api

type NewCheckoutSession struct {
	Amount            int64  `json:"amount"` // Amount intended to be collected by this PaymentIntent. A positive integer representing how much to charge in the smallest currency unit
	Chain             int64  `json:"chain"`
	Token             string `json:"token"` // Token intended for this Payment intent e.g USDC
	Customer          string `json:"customer"`
	Interval          int64  `json:"interval"`
	ProductId         string `json:"product_id"`
	ChargeLater       bool   `json:"charge_later"`
	PaymentMode       Mode   `json:"payment_type"`
	ClientReferenceId string `json:"client_reference_id"`
}

type CheckoutSessiontResponse struct {
	Id          string `json:"id"`
	Amount      int64  `json:"amount"`
	Chain       int64  `json:"chain"`
	Token       string `json:"token"` // Token intended for this Payment intent e.g USDC
	Customer    string `json:"customer"`
	Interval    int64  `json:"interval"`
	ProductId   string `json:"product_id"`
	ChargeLater bool   `json:"charge_later"`
	PaymentMode string `json:"payment_type"`
	Mode        string `json:"mode"`
}

type Mode string

const (
	SinglePaymentMode    Mode = "single"
	RecurringPaymentMode Mode = "recurring"
)

func (e Mode) IsValid() bool {
	switch e {
	case SinglePaymentMode, RecurringPaymentMode:
		return true
	}
	return false
}

func (e Mode) String() string {
	return string(e)
}
