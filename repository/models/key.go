package models

type Key struct {
	ID              uint `gorm:"primarykey"`
	SubscriptionKey string `gorm:"index"`
	SecretKey       string `gorm:"not null"`
	SubscriptionID  uint
}
