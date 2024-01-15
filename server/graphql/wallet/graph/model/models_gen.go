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

type CheckoutSession struct {
	ID          string      `json:"id"`
	Chain       int         `json:"chain"`
	Token       string      `json:"token"`
	Amount      int         `json:"amount"`
	Interval    int         `json:"interval"`
	ProductID   string      `json:"productId"`
	MerchantID  string      `json:"merchantId"`
	PaymentType PaymentType `json:"paymentType"`
	ChargeLater bool        `json:"chargeLater"`
}

type NewTransferRequest struct {
	Chain  int     `json:"chain"`
	Token  string  `json:"token"`
	Amount float64 `json:"amount"`
	Sender string  `json:"sender"`
	Target string  `json:"target"`
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
	Mode           Mode        `json:"mode"`
	Type           PaymentType `json:"type"`
	Email          *string     `json:"email,omitempty"`
	Chain          int         `json:"chain"`
	Token          string      `json:"token"`
	Amount         float64     `json:"amount"`
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
	ID                  string     `json:"id"`
	Token               string     `json:"token"`
	Amount              int        `json:"amount"`
	Interval            int        `json:"interval"`
	ProductID           string     `json:"productId"`
	MerchantID          string     `json:"merchantId"`
	ProductName         string     `json:"productName"`
	WalletAddress       string     `json:"walletAddress"`
	SubscriptionKey     string     `json:"subscriptionKey"`
	CreatedAt           string     `json:"createdAt"`
	NextChargeDate      time.Time  `json:"nextChargeDate"`
	TransactionHash     string     `json:"transactionHash"`
	TransactionExplorer string     `json:"transactionExplorer"`
	Payments            []*Payment `json:"payments,omitempty"`
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
	Interval            int         `json:"interval"`
	Reference           string      `json:"reference"`
	ProductID           string      `json:"productId"`
	WalletAddress       string      `json:"walletAddress"`
	SubscriptionKey     string      `json:"subscriptionKey"`
	CreatedAt           string      `json:"createdAt"`
	TransactionHash     string      `json:"transactionHash"`
	TransactionExplorer string      `json:"transactionExplorer"`
}

type ValidationData struct {
	UserOpHash string `json:"userOpHash"`
}

type Mode string

const (
	ModeTest Mode = "test"
	ModeLive Mode = "live"
)

var AllMode = []Mode{
	ModeTest,
	ModeLive,
}

func (e Mode) IsValid() bool {
	switch e {
	case ModeTest, ModeLive:
		return true
	}
	return false
}

func (e Mode) String() string {
	return string(e)
}

func (e *Mode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Mode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Mode", str)
	}
	return nil
}

func (e Mode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
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

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusDisabled  SubscriptionStatus = "disabled"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

var AllSubscriptionStatus = []SubscriptionStatus{
	SubscriptionStatusActive,
	SubscriptionStatusDisabled,
	SubscriptionStatusCancelled,
}

func (e SubscriptionStatus) IsValid() bool {
	switch e {
	case SubscriptionStatusActive, SubscriptionStatusDisabled, SubscriptionStatusCancelled:
		return true
	}
	return false
}

func (e SubscriptionStatus) String() string {
	return string(e)
}

func (e *SubscriptionStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SubscriptionStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SubscriptionStatus", str)
	}
	return nil
}

func (e SubscriptionStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
