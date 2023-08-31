package models

type Wallet struct {
	ID             uint         `gorm:"primaryKey"`
	Email          string         `gorm:"index"`
	SignerAddress  string         `gorm:"index"`
	AccountAddress string         `gorm:"index"`
}
