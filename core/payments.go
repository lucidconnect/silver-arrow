package core

import (
	"time"

	"github.com/google/uuid"
)

type PaymentType string
type PaymentState string

const (
	PaymentTypeSingle    PaymentType = "single"
	PaymentTypeRecurring PaymentType = "recurring"

	// PaymentStatusInitiated PaymentStatus = "initiated"
	PaymentStatusPending PaymentState = "pending"
	PaymentStatusSuccess PaymentState = "success"
	PaymentStatusFailed  PaymentState = "failed"
)

/*** PaymentIntent is the entry point resource to initiate a payment on a user's erc4337 wallet
 * PaymentIntent is created by merchant and sent to Lucid
 * if payment type is recurring, a subscription is first created before a payment can be made
 * an erc4337 user operation is created, and the hash is returned to the merchant requiring the user's signature.
 */
type PaymentIntent struct {
	Type  PaymentType
	Email string
	// Chain             int64
	// Token             string
	Source string
	// Amount            float64
	Interval          int
	ProductId         string
	PriceId           string
	OwnerAddress      string
	WalletAddress     string
	FirstChargeNow    bool
	CheckoutSessionId uuid.UUID
}

type Payment struct {
	Id                uuid.UUID
	Type              PaymentType
	State             PaymentState
	Chain             int64
	Token             string
	TokenAddress      string
	Amount            int64
	Source            string
	Merchant          string // name of merchant
	MerchantId        uuid.UUID
	ProductId         uuid.UUID
	SubscriptionId    uuid.UUID
	CheckoutSessionID uuid.UUID
	Sponsored         bool
	Reference         uuid.UUID
	UserOpHash        string
	TransactionHash   string
	BlockExplorerTx   string
	CreatedAt         int64
}

type NewSubscription struct {
	Chain             int64
	Email             string
	Token             string
	Amount            int64
	Interval          int64
	IntervalUnit      string
	ProductID         uuid.UUID
	CheckoutSessionID uuid.UUID
	ProductName       string
	OwnerAddress      string
	WalletAddress     string
	DepositAddress    string
	NextChargeDate    time.Time
}

type PaymentReceipt struct {
	Type            PaymentType
	Chain           int64
	Status          PaymentState
	TransactionHash string
	BlockExplorerTx string
}

type BillingHistory struct {
	Date        time.Time
	Amount      float64
	ExplorerURL string
}
