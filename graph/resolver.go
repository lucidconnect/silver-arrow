package graph

import (
	"github.com/helicarrierstudio/silver-arrow/erc4337"
	"github.com/helicarrierstudio/silver-arrow/repository"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	WalletRepository repository.WalletRepository
	Bundler          *erc4337.ERCBundler
	Cache            repository.CacheWizard
}
