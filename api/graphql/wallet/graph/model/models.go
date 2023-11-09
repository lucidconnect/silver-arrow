package model

import "time"

type NewSubscription struct {
	Chain          int
	Email          string
	Token          string
	Amount         int
	Interval       int
	ProductID      string
	OwnerAddress   string
	WalletAddress  string
	NextChargeDate *time.Time
}
