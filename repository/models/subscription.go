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
	Status                 string
	Interval               int64     `gorm:"not null"`
	UserOpHash             string    `gorm:"index"`
	MerchantId             string    `gorm:"index"`
	ProductID              uuid.UUID `gorm:"index"`
	CheckoutSessionID      uuid.UUID
	ProductName            string
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
	DisabledAt             time.Time
	CancelledAt            time.Time
	Payments               []Payment
	TransactionHash        string
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

type Payment struct {
	ID                    uuid.UUID   `gorm:"primaryKey"`
	Type                  PaymentType `gorm:"not null"`
	Chain                 int64       `gorm:"not null"`
	Token                 string      `gorm:"not null"`
	TokenAddress          string
	Status                PaymentStatus `gorm:"not null"`
	Amount                int64         `gorm:"not null"`
	Source                string
	WalletID              uuid.UUID `gorm:"not null"`
	ProductID             uuid.UUID `gorm:"not null"`
	Sponsored             bool
	Reference             uuid.UUID `gorm:"index;not null"`
	UserOpHash            string
	Destination           string
	Acknowledged          bool
	SubscriptionID        uuid.UUID
	TransactionHash       string
	BlockExplorerTx       string
	SubscriptionPublicKey string
	CreatedAt             time.Time
	WebhookAcknowledgedAt time.Time
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

type PaymentStatus string
type PaymentType string

const (
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"

	PaymentTypeSingle    PaymentType = "single"
	PaymentTypeRecurring PaymentType = "recurring"
)
