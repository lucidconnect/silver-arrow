package models

import (
	"github.com/google/uuid"
)

type Merchant struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	Name           string    `gorm:"unique"`
	Chain          int64     `gorm:"not null"`
	Owner          string    `gorm:"index;not null"`
	Token          string    `gorm:"not null"`
	DepositAddress string    `gorm:"not null"`
}
