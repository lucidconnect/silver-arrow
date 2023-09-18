package models

type Wallet struct {
	ID                   uint   `gorm:"primaryKey"`
	Email                string `gorm:"index"`
	SignerAddress        string `gorm:"index"`
	WalletAddress        string `gorm:"index"`
	TurnkeySubOrgID      string
	TurnkeySubOrgName    string
	TurnkeyPrivateKeyTag string
	Subscriptions        []Subscription
}
