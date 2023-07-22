package models

type Subscription struct {
	Amount        float64 `bson:"amount"`
	Active        bool    `bson:"active"`
	Interval      int     `bson:"interval"`
	Token         string  `bson:"token"` // really the token contract
	MerchantId    string  `bson:"merchant_id"`
	WalletAddress string  `bson:"wallet_address"`
}
