package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// `Subscription` belongs to `Wallet`, `WalletID` is the foreign key
// `Subscription` has one Key, `SubscriptionID` is the foreign key
type Subscription struct {
	ID                     uuid.UUID `gorm:"primaryKey"`
	Token                  string    `gorm:"not null"` // really the token contract
	Amount                 int64     `gorm:"not null"` // amount in decimal precision
	Active                 bool      `gorm:"not null"`
	Interval               int64     `gorm:"not null"`
	UserOpHash             string    `gorm:"index"`
	MerchantId             string    `gorm:"index"`
	ProductID              uuid.UUID `gorm:"index"`
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
	Payments               []Payment
	TransactionHash        string
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

type Payment struct {
	ID              uuid.UUID `gorm:"primaryKey"`
	Type            string    `gorm:"not null"`
	Chain           int64     `gorm:"not null"`
	Token           string    `gorm:"not null"`
	Status          string    `gorm:"not null"`
	Amount          int64     `gorm:"not null"`
	Source          string
	WalletID        uuid.UUID `gorm:"not null"`
	ProductID       uuid.UUID `gorm:"not null"`
	Reference       uuid.UUID `gorm:"not null"`
	UserOpHash      string
	Destination     string
	SubscriptionID  uuid.UUID
	TransactionHash string
	BlockExplorerTx string
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
