package graph

import "github.com/helicarrierstudio/silver-arrow/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	WalletRepository   repository.WalletRepository
	MerchantRepository repository.MerchantRepository
}
