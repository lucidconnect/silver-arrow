package core

import "github.com/google/uuid"

type SessionState string

const (
	SessionExpiry  int64        = 1.8e12 // time in nanoseconds for a session to expire (30 minutes)
	SessionActive  SessionState = "active"
	SessionExpired SessionState = "expired"
)

// CheckoutSession describes a session that triggers a payment (one-time) or initiates a subscription
// A checkout session can be created by the merchant on behalf of a customer to start a checkout processs
// or created whenever a customer visits a payment link.
type CheckoutSession struct {
	Id          uuid.UUID
	ProductId   uuid.UUID
	MerchantId  uuid.UUID
	CustomerId  string
	CallbackURL string
	PaymentLink string
	State       SessionState
	CreatedAt   int64
}

// create session
