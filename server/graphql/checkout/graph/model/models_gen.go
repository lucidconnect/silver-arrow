// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Account struct {
	Email   *string `json:"email,omitempty"`
	Address string  `json:"address"`
	Signer  *string `json:"signer,omitempty"`
}

type BillingHistory struct {
	Date        time.Time `json:"date"`
	Amount      float64   `json:"amount"`
	ExplorerURL string    `json:"explorerUrl"`
}

type Payment struct {
	Chain     int           `json:"chain"`
	Token     string        `json:"token"`
	Status    PaymentStatus `json:"status"`
	Amount    float64       `json:"amount"`
	Source    string        `json:"source"`
	ProductID string        `json:"productId"`
	Reference string        `json:"reference"`
}

type PaymentIntent struct {
	Type              PaymentType `json:"type"`
	Email             *string     `json:"email,omitempty"`
	Chain             int         `json:"chain"`
	Token             string      `json:"token"`
	Amount            float64     `json:"amount"`
	CheckoutSessionID *string     `json:"checkoutSessionId,omitempty"`
	ProductID         string      `json:"productId"`
	PriceID           string      `json:"priceId"`
	OwnerAddress      string      `json:"ownerAddress"`
	WalletAddress     string      `json:"walletAddress"`
	FirstChargeNow    bool        `json:"firstChargeNow"`
}

type RequestValidation struct {
	Chain         int    `json:"chain"`
	UserOpHash    string `json:"userOpHash"`
	SignedMessage string `json:"signedMessage"`
}

type TransactionData struct {
	ID                  *string     `json:"id,omitempty"`
	Type                PaymentType `json:"type"`
	Chain               int         `json:"chain"`
	Token               string      `json:"token"`
	Amount              int         `json:"amount"`
	Interval            int         `json:"interval"`
	Reference           string      `json:"reference"`
	ProductID           string      `json:"productId"`
	WalletAddress       string      `json:"walletAddress"`
	SubscriptionKey     string      `json:"subscriptionKey"`
	CreatedAt           string      `json:"createdAt"`
	TransactionHash     string      `json:"transactionHash"`
	TransactionExplorer string      `json:"transactionExplorer"`
}

type PaymentStatus string

const (
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
)

var AllPaymentStatus = []PaymentStatus{
	PaymentStatusFailed,
	PaymentStatusPending,
	PaymentStatusSuccess,
}

func (e PaymentStatus) IsValid() bool {
	switch e {
	case PaymentStatusFailed, PaymentStatusPending, PaymentStatusSuccess:
		return true
	}
	return false
}

func (e PaymentStatus) String() string {
	return string(e)
}

func (e *PaymentStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PaymentStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PaymentStatus", str)
	}
	return nil
}

func (e PaymentStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PaymentType string

const (
	PaymentTypeSingle    PaymentType = "single"
	PaymentTypeRecurring PaymentType = "recurring"
)

var AllPaymentType = []PaymentType{
	PaymentTypeSingle,
	PaymentTypeRecurring,
}

func (e PaymentType) IsValid() bool {
	switch e {
	case PaymentTypeSingle, PaymentTypeRecurring:
		return true
	}
	return false
}

func (e PaymentType) String() string {
	return string(e)
}

func (e *PaymentType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PaymentType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PaymentType", str)
	}
	return nil
}

func (e PaymentType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
