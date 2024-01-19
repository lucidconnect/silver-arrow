package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Each product created will have a deposit wallet attached to it
// A deposit wallet can be used for more than one product
// A product can only have one deposit wallet attached to it
type DepositWallet struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	MerchantID    uuid.UUID
	Merchant      Merchant
	WalletAddress string
}

type Merchant struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	Name         string
	Email        string
	WebhookUrl   string
	WebhookToken string
	OwnerAddress string `gorm:"unique"` // web3 wallet that owns this account
	// Products           []Product
	ConvoyEndpointID   string
	MerchantAccessKeys []MerchantAccessKey
	PaymentLinks       []PaymentLink
	CreatedAt          time.Time
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
