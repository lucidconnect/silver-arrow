package api

type NewCheckoutSession struct {
	Customer    string `json:"customer"`
	ProductId   string `json:"product_id"`
	CallbackUrl string `json:"callback_url"`
}

type CheckoutSessionResponse struct {
	SessionId string `json:"session_id"`
	Url       string `json:"url"`
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
