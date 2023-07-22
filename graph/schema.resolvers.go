package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/helicarrierstudio/silver-arrow/graph/generated"
	"github.com/helicarrierstudio/silver-arrow/graph/model"
	"github.com/helicarrierstudio/silver-arrow/wallet"
)

// AddAccount is the resolver for the addAccount field.
func (r *mutationResolver) AddAccount(ctx context.Context, input model.Account) (string, error) {
	address := common.HexToAddress(input.Address)
	walletService := wallet.NewWalletService(r.WalletRepository)
	err := walletService.AddAccount(input)
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}

// AddSubscription is the resolver for the addSubscription field.
func (r *mutationResolver) AddSubscription(ctx context.Context, input model.NewSubscription) (*model.SubscriptionData, error) {
	walletService := wallet.NewWalletService(r.WalletRepository)
	subData, err := walletService.AddSubscription(input)
	if err != nil {
		return nil, err
	}
	return subData, nil
}

// CancelSubscription is the resolver for the cancelSubscription field.
func (r *mutationResolver) CancelSubscription(ctx context.Context, id string) (string, error) {
	panic(fmt.Errorf("not implemented: CancelSubscription - cancelSubscription"))
}

// FetchSubscriptions is the resolver for the fetchSubscriptions field.
func (r *queryResolver) FetchSubscriptions(ctx context.Context, account string) ([]*model.SubscriptionData, error) {
	_ = wallet.NewWalletService(r.WalletRepository)
	panic(fmt.Errorf("not implemented: FetchSubscriptions - fetchSubscriptions"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
