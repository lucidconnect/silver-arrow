package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/helicarrierstudio/silver-arrow/erc4337"
	"github.com/helicarrierstudio/silver-arrow/graph/generated"
	"github.com/helicarrierstudio/silver-arrow/graph/model"
	"github.com/helicarrierstudio/silver-arrow/wallet"
	"github.com/pkg/errors"
)

// AddAccount is the resolver for the addAccount field.
func (r *mutationResolver) AddAccount(ctx context.Context, input model.Account) (string, error) {
	address := common.HexToAddress(input.Address)
	walletService := wallet.NewWalletService(r.WalletRepository, r.Bundler)
	// should check if the account is deployed
	// deploy if not deployed
	err := walletService.AddAccount(input)
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}

// AddSubscription is the resolver for the addSubscription field.
func (r *mutationResolver) AddSubscription(ctx context.Context, input model.NewSubscription) (*model.ValidationData, error) {
	walletService := wallet.NewWalletService(r.WalletRepository, r.Bundler)
	validationData, userOp, err := walletService.AddSubscription(input)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("Userop hash", validationData.UserOpHash)
	err = r.Cache.Set(validationData.UserOpHash, userOp)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return validationData, nil
}

// ValidateSubscription is the resolver for the validateSubscription field.
func (r *mutationResolver) ValidateSubscription(ctx context.Context, input model.SubscriptionValidation) (*model.SubscriptionData, error) {
	time.Sleep(time.Second)
	opInterface, err := r.Cache.Get(input.UserOpHash)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	op, _ := opInterface.(map[string]any)
	sig, err := hexutil.Decode(erc4337.SUDO_MODE)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	partialSig, err := hexutil.Decode(input.SignedMessage)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	sig = append(sig, partialSig...)
	op["signature"] = hexutil.Encode(sig)
	walletService := wallet.NewWalletService(r.WalletRepository, r.Bundler)

	log.Println("User op", op)
	subData, key, err := walletService.ValidateSubscription(op)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// x := int64(subData.Amount)
	target := "0x1BB271879576fD79324156F539DD760756C9D061"

	err = walletService.ExecuteCharge(subData.WalletAddress, target, subData.MerchantID, subData.Token, key, int64(subData.Amount))
	if err != nil {
		err = errors.Wrap(err, "error occurred during first time charge execution - ")
		log.Println(err)
		return subData, err
	}
	return subData, nil
}

// CancelSubscription is the resolver for the cancelSubscription field.
func (r *mutationResolver) CancelSubscription(ctx context.Context, id string) (string, error) {
	panic(fmt.Errorf("not implemented: CancelSubscription - cancelSubscription"))
}

// FetchSubscriptions is the resolver for the fetchSubscriptions field.
func (r *queryResolver) FetchSubscriptions(ctx context.Context, account string) ([]*model.SubscriptionData, error) {
	_ = wallet.NewWalletService(r.WalletRepository, r.Bundler)
	panic(fmt.Errorf("not implemented: FetchSubscriptions - fetchSubscriptions"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
