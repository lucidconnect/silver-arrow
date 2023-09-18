package models

import (
	"time"

	"gorm.io/gorm"
)

// `Subscription` belongs to `Wallet`, `WalletID` is the foreign key
// `Subscription` has one Key, `SubscriptionID` is the foreign key
type Subscription struct {
	gorm.Model
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
	WalletID               uint
	Key                    Key `gorm:"foreignKey:SubscriptionID"`
}
