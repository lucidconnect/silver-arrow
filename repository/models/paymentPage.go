package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentLink struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	MerhantId uuid.UUID
	ProductId uuid.UUID
	Product   Product
	CreatedAt time.Time
	DeletedAt time.Time
}
