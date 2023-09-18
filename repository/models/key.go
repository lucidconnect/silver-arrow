package models

import "gorm.io/gorm"

type Key struct {
	gorm.Model
	PublicKey      string `gorm:"not null"`
	PrivateKeyId   string `gorm:"not null"`
	SubscriptionID uint
}
