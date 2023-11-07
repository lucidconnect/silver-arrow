package model

type NewSubscription struct {
	Chain         int
	Email         string
	Token         string
	Amount        int
	Interval      int
	ProductID     string
	WalletAddress string
	OwnerAddress  string
}
