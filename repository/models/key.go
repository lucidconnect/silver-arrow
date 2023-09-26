package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Key struct {
	ID             uuid.UUID `gorm:"primarykey"`
	CreatedAt      time.Time
	PublicKey      string `gorm:"not null"`
	PrivateKeyId   string `gorm:"not null"`
	SubscriptionID uuid.UUID
	WalletID       uuid.UUID
}

func (k *Key) BeforeCreate(tx *gorm.DB) (err error) {
	k.ID = uuid.New()
	return
}
