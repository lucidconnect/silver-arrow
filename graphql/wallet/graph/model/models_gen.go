// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Account struct {
	Email   *string `json:"email,omitempty"`
	Address string  `json:"address"`
	Signer  *string `json:"signer,omitempty"`
}

type NewSubscription struct {
	Chain         int        `json:"chain"`
	NextChargeAt  *time.Time `json:"nextChargeAt,omitempty"`
	Token         string     `json:"token"`
	Amount        int        `json:"amount"`
	Interval      int        `json:"interval"`
	ProductID     string     `json:"productId"`
	MerchantID    string     `json:"merchantId"`
	WalletAddress string     `json:"walletAddress"`
	OwnerAddress  string     `json:"ownerAddress"`
}

type SubscriptionData struct {
	ID              string  `json:"id"`
	Token           string  `json:"token"`
	Amount          int     `json:"amount"`
	Interval        int     `json:"interval"`
	MerchantID      string  `json:"merchantId"`
	WalletAddress   string  `json:"walletAddress"`
	SubscriptionKey string  `json:"subscriptionKey"`
	CreatedAt       *string `json:"createdAt,omitempty"`
}

type SubscriptionValidation struct {
	Chain         int    `json:"chain"`
	UserOpHash    string `json:"userOpHash"`
	SignedMessage string `json:"signedMessage"`
}

type ValidationData struct {
	UserOpHash string `json:"userOpHash"`
}
