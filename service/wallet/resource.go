package wallet

import (
	"time"

	"github.com/google/uuid"
)

/*** PaymentIntent is the entry point resource to initiate a payment on a user's erc4337 wallet
 * PaymentIntent is created by merchant and sent to Lucid
 * if payment type is recurring, a subscription is first created before a payment can be made
 * an erc4337 user operation is created, and the hash is returned to the merchant requiring the user's signature.
 */
type PaymentIntent struct {
	Type          PaymentType
	Email         string
	Chain         int64
	Token         string
	Source        string
	Amount        string
	Interval      int
	ProductId     string
	OwnerAddress  string
	WalletAddress string
}

type Account struct {
	Email   *string `json:"email,omitempty"`
	Address string  `json:"address"`
	Signer  *string `json:"signer,omitempty"`
}

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
type PaymentRequestValidation struct {
	Chain     int64
	Hash      string
	Signature string
}

type PaymentReceipt struct {
	Type            PaymentType
	Chain           int64
	Status          PaymentStatus
	TransactionHash string
	BlockExplorerTx string
}

type PaymentType string
type PaymentStatus string

// const (
// 	PaymentTypeSingle    PaymentType = "single"
// 	PaymentTypeRecurring PaymentType = "recurring"

// 	// PaymentStatusInitiated PaymentStatus = "initiated"
// 	PaymentStatusPending   PaymentStatus = "pending"
// 	PaymentStatusSuccess   PaymentStatus = "success"
// 	PaymentStatusFailed    PaymentStatus = "failed"
// )

// func (pt PaymentType) IsValid() bool {
// 	switch pt {
// 	case PaymentTypeSingle, PaymentTypeRecurring:
// 		return true
// 	}
// 	return false
// }

func (pt PaymentType) String() string {
	return string(pt)
}

func (pt PaymentStatus) String() string {
	return string(pt)
}
