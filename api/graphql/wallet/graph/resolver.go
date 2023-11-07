package graph

import (
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/service/turnkey"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Cache          repository.CacheWizard
	Database       repository.Database
	TurnkeyService *turnkey.TurnkeyService
}
