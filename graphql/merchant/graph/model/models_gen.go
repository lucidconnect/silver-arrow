// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AccessKey struct {
	Mode       Mode   `json:"mode"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type Merchant struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	PublicKey  string     `json:"publicKey"`
	WebHookURL string     `json:"webHookUrl"`
	AccessKey  *AccessKey `json:"accessKey"`
}

type MerchantAccessKeyQuery struct {
	Mode            Mode   `json:"mode"`
	MerchantAddress string `json:"merchantAddress"`
}

type MerchantStats struct {
	Users         int `json:"users"`
	Products      int `json:"products"`
	Subscriptions int `json:"subscriptions"`
}

type MerchantUpdate struct {
	MerchantID string  `json:"merchantId"`
	Name       *string `json:"name,omitempty"`
	Email      *string `json:"email,omitempty"`
	IsActive   *bool   `json:"isActive,omitempty"`
	WebHookURL *string `json:"webHookUrl,omitempty"`
}

type NewMerchant struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Owner string `json:"owner"`
	// This would be the url where payment status event would be delivered to
	WebHookURL *string `json:"webHookUrl,omitempty"`
}

type NewMerchantKey struct {
	MerchantAddress string `json:"merchantAddress"`
	Mode            Mode   `json:"mode"`
}

type NewProduct struct {
	Name             string `json:"name"`
	Owner            string `json:"owner"`
	Chain            int    `json:"chain"`
	Token            string `json:"token"`
	ReceivingAddress string `json:"receivingAddress"`
}

type Product struct {
	Name             string  `json:"name"`
	Mode             Mode    `json:"mode"`
	Owner            string  `json:"owner"`
	Chain            int     `json:"chain"`
	Token            string  `json:"token"`
	ProductID        string  `json:"productId"`
	MerchantID       string  `json:"merchantId"`
	ReceivingAddress string  `json:"receivingAddress"`
	Subscriptions    []*Sub  `json:"subscriptions,omitempty"`
	CreatedAt        *string `json:"createdAt,omitempty"`
}

type ProductModeUpdate struct {
	ProductID string `json:"ProductId"`
	Mode      Mode   `json:"mode"`
}

type ProductUpdate struct {
	Name             *string `json:"name,omitempty"`
	ReceivingAddress *string `json:"receivingAddress,omitempty"`
}

type Sub struct {
	Chain         int    `json:"chain"`
	Token         string `json:"token"`
	Amount        int    `json:"amount"`
	Active        bool   `json:"active"`
	Interval      string `json:"interval"`
	WalletAddress string `json:"walletAddress"`
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
