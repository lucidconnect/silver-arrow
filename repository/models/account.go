package models

type Account struct {
	ID             uint64 `gorm:"primaryKey"`
	SignerAddress  string `gorm:"not null"`
	AccountAddress string `gorm:"not null"`
}
