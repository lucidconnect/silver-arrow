// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Account struct {
	Email   *string `json:"email,omitempty"`
	Address string  `json:"address"`
	Signer  *string `json:"signer,omitempty"`
}

type NewTransferRequest struct {
	Chain  int     `json:"chain"`
	Token  string  `json:"token"`
	Amount float64 `json:"amount"`
	Sender string  `json:"sender"`
	Target string  `json:"target"`
}

type PaymentIntent struct {
	Type           PaymentType `json:"type"`
	Email          *string     `json:"email,omitempty"`
	Chain          int         `json:"chain"`
	Token          string      `json:"token"`
	Amount         int         `json:"amount"`
	Interval       int         `json:"interval"`
	ProductID      string      `json:"productId"`
	OwnerAddress   string      `json:"ownerAddress"`
	WalletAddress  string      `json:"walletAddress"`
	FirstChargeNow bool        `json:"firstChargeNow"`
}

type RequestValidation struct {
	Chain         int    `json:"chain"`
	UserOpHash    string `json:"userOpHash"`
	SignedMessage string `json:"signedMessage"`
}

type SubscriptionData struct {
	ID                  string  `json:"id"`
	Token               string  `json:"token"`
	Amount              int     `json:"amount"`
	Interval            int     `json:"interval"`
	ProductID           string  `json:"productId"`
	WalletAddress       string  `json:"walletAddress"`
	SubscriptionKey     string  `json:"subscriptionKey"`
	CreatedAt           *string `json:"createdAt,omitempty"`
	TransactionHash     *string `json:"transactionHash,omitempty"`
	TransactionExplorer *string `json:"transactionExplorer,omitempty"`
}

type SubscriptionMod struct {
	SubscriptionID string       `json:"subscriptionId"`
	Toggle         StatusToggle `json:"toggle"`
}

type TransactionData struct {
	ID                  *string     `json:"id,omitempty"`
	Type                PaymentType `json:"type"`
	Chain               int         `json:"chain"`
	Token               string      `json:"token"`
	Amount              int         `json:"amount"`
	Interval            *int        `json:"interval,omitempty"`
	Reference           string      `json:"reference"`
	ProductID           *string     `json:"productId,omitempty"`
	WalletAddress       string      `json:"walletAddress"`
	SubscriptionKey     *string     `json:"subscriptionKey,omitempty"`
	CreatedAt           *string     `json:"createdAt,omitempty"`
	TransactionHash     *string     `json:"transactionHash,omitempty"`
	TransactionExplorer *string     `json:"transactionExplorer,omitempty"`
}

type ValidationData struct {
	UserOpHash string `json:"userOpHash"`
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

type StatusToggle string

const (
	StatusToggleCancel  StatusToggle = "cancel"
	StatusToggleDisable StatusToggle = "disable"
	StatusToggleEnable  StatusToggle = "enable"
)

var AllStatusToggle = []StatusToggle{
	StatusToggleCancel,
	StatusToggleDisable,
	StatusToggleEnable,
}

func (e StatusToggle) IsValid() bool {
	switch e {
	case StatusToggleCancel, StatusToggleDisable, StatusToggleEnable:
		return true
	}
	return false
}

func (e StatusToggle) String() string {
	return string(e)
}

func (e *StatusToggle) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StatusToggle(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StatusToggle", str)
	}
	return nil
}

func (e StatusToggle) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
