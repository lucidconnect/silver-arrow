package models

import (
	"time"

	"github.com/google/uuid"
)

type Price struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	Active        bool
	Type          string
	Interval      string
	IntervalCount int64
	TrialPeriod   int64
	ProductID     uuid.UUID
	MerchantID    uuid.UUID
	Product       Product
	Amount        int64
	Token         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
