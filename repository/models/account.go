package models

type Wallet struct {
	Email          string `bson:"email,omitempty"`
	SignerAddress  string `bson:"signer_address,omitempty"`
	AccountAddress string `bson:"account_address,omitempty"`
	Subscriptions  []Subscription
}
