package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	Name             string    `gorm:"unique"`
	Chain            int64     `gorm:"not null"`
	Owner            string    `gorm:"index;not null"`
	Token            string    `gorm:"not null"`
	DepositAddress   string    `gorm:"not null"`
	DepositWalletID  uuid.UUID
	MerchantID       uuid.UUID
	DepositWallet    DepositWallet
	Merchant         Merchant
	Subscriptions    []Subscription
	CheckoutSessions []CheckoutSession
	Payments         []Payment
	Mode             string
	Amount           int64
	Interval         int64
	InstantCharge    bool
	PaymentType      string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}
