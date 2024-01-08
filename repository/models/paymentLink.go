package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentLink struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	MerchantID   uuid.UUID
	MerchantName string
	CallbackURL  string
	ProductID    uuid.UUID
	Product      Product
	CreatedAt    time.Time
	DeletedAt    time.Time
}
