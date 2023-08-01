package models

import "time"

type Subscription struct {
	Token          string    `bson:"token"`  // really the token contract
	Amount         int64     `bson:"amount"` // amount in wei
	Active         bool      `bson:"active"`
	Interval       int64     `bson:"interval"`
	UserOpHash     string    `bson:"userop_hash"`
	SigningKey     string    `bson:"signing_key"`
	MerchantId     string    `bson:"merchant_id"`
	NextChargeAt   time.Time `bson:"next_charge_at"`
	WalletAddress  string    `bson:"wallet_address"`
	SubscriptionId string    `bson:"subscription_id"`
}
