package models

import "github.com/google/uuid"

type CheckoutSession struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	Customer      string
	ProductID     uuid.UUID
	MerchantID    uuid.UUID
	PaymentLinkID uuid.UUID
	PaymentLink   PaymentLink // multiple checkout_sesisons can belong to a single payment_link and subscription
	CallbackURL   string
}
