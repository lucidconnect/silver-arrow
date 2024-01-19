package models

import (
	"time"

	"github.com/google/uuid"
)

// PaymentLink belongs to a product
type PaymentLink struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	MerchantID   uuid.UUID
	MerchantName string
	CallbackURL  string
	ProductID    uuid.UUID
	PriceID      uuid.UUID
	Product      Product
	Price        Price
	CreatedAt    time.Time
	DeletedAt    time.Time
}
