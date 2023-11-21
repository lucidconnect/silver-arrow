package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	Name           string    `gorm:"unique"`
	Chain          int64     `gorm:"not null"`
	Owner          string    `gorm:"index;not null"`
	Token          string    `gorm:"not null"`
	DepositAddress string    `gorm:"not null"`
	MerchantID     uuid.UUID
	CreatedAt      time.Time
	Subscriptions  []Subscription
	Payments       []Payment
}

type Merchant struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	Name             string
	Email            string
	PublicKey        string // lucid public key for authenticating requests
	WebhookUrl       string
	WebhookToken     string
	OwnerAddress     string `gorm:"unique"` // web3 wallet that owns this account
	Products         []Product
	ConvoyEndpointID string
}
