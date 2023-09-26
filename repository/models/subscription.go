package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// `Subscription` belongs to `Wallet`, `WalletID` is the foreign key
// `Subscription` has one Key, `SubscriptionID` is the foreign key
type Subscription struct {
	gorm.Model
	ID                     uuid.UUID
	Token                  string `gorm:"not null"` // really the token contract
	Amount                 int64  `gorm:"not null"` // amount in decimal precision
	Active                 bool   `gorm:"not null"`
	Interval               int64  `gorm:"not null"`
	UserOpHash             string `gorm:"index"`
	MerchantId             string `gorm:"index"`
	MerchantDepositAddress string
	ExpiresAt              time.Time `gorm:"index;type:timestamptz"`
	NextChargeAt           time.Time `gorm:"index;type:timestamptz"`
	WalletAddress          string    `gorm:"index"`
	TokenAddress           string
	WalletID               uuid.UUID
	Key                    Key   `gorm:"foreignKey:SubscriptionID"`
	Chain                  int64 //
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
