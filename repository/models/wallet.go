package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID                   uuid.UUID `gorm:"primaryKey"`
	Email                string    `gorm:"index"`
	UserID               uuid.UUID
	SignerAddress        string `gorm:"index"`
	WalletAddress        string `gorm:"index"`
	TurnkeySubOrgID      string
	TurnkeySubOrgName    string
	TurnkeyPrivateKeyTag string
	Subscriptions        []Subscription
	Keys                 []Key
	Payments             []Payment
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New()
	return
}

type User struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	MerchantID     uuid.UUID `gorm:"index"`
	ClientUniqueID string
	WalletParams   []Wallet
}
