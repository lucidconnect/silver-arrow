package model

import (
	"time"

	"github.com/google/uuid"
)

type NewSubscription struct {
	Chain          int
	Email          string
	Token          string
	Amount         float64
	Interval       int
	ProductID      uuid.UUID
	ProductName    string
	OwnerAddress   string
	WalletAddress  string
	DepositAddress string
	NextChargeDate *time.Time
}
