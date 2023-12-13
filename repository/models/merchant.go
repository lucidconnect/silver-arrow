package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/graphql/merchant/graph/model"
	"gorm.io/gorm"
)

type Product struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	Name             string    `gorm:"unique"`
	Chain            int64     `gorm:"not null"`
	Owner            string    `gorm:"index;not null"`
	Token            string    `gorm:"not null"`
	DepositAddress   string    `gorm:"not null"`
	MerchantID       uuid.UUID
	CreatedAt        time.Time
	Subscriptions    []Subscription
	CheckoutSessions []CheckoutSession
	Payments         []Payment
	Mode             model.Mode
}

type Merchant struct {
	ID                 uuid.UUID `gorm:"primaryKey"`
	Name               string
	Email              string
	WebhookUrl         string
	WebhookToken       string
	OwnerAddress       string `gorm:"unique"` // web3 wallet that owns this account
	Products           []Product
	ConvoyEndpointID   string
	MerchantAccessKeys []MerchantAccessKey
}

type MerchantAccessKey struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Mode       string
	PublicKey  string    `gorm:"index"`
	MerchantID uuid.UUID `gorm:"index"`
	CreatedAt  time.Time
}

func (m *MerchantAccessKey) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
