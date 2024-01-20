package models

import (
	"time"

	"github.com/google/uuid"
)

type Price struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	Active        bool
	Amount        int64
	Token         string
	Chain         int64
	Type          string
	Interval      string
	IntervalCount int64
	TrialPeriod   int64
	ProductID     uuid.UUID
	MerchantID    uuid.UUID
	Product       Product
	CreatedAt time.Time
	UpdatedAt time.Time
}
