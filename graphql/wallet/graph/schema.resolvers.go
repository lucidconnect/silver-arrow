package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/graphql/wallet/graph/generated"
	"github.com/lucidconnect/silver-arrow/graphql/wallet/graph/model"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/lucidconnect/silver-arrow/service/merchant"
	"github.com/lucidconnect/silver-arrow/service/wallet"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// AddAccount is the resolver for the addAccount field.
func (r *mutationResolver) AddAccount(ctx context.Context, input model.Account) (string, error) {
	address := common.HexToAddress(input.Address)

	walletService := wallet.NewWalletService(r.Database, r.TurnkeyService)
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
	merchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return nil, err
	}
	merchantId := merchant.ID
	signature, err := auth.SignatureContext(ctx, merchant.PublicKey)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("Authenticated Merchant: %v", merchantId)
	// validate signature
	// amount:token:interval:productId
	signatureCheck := fmt.Sprintf("%v", input.Amount) + ":" + input.Token + ":" + fmt.Sprintf("%v", input.Interval) + ":" + input.ProductID
	err = validateSignature(signatureCheck, signature, merchant.PublicKey)
	if err != nil {
		log.Err(err).Ctx(ctx).Send()
		err = errors.New("request signature is invalid")
		return nil, err
	}
	
	walletService := wallet.NewWalletService(r.Database, r.TurnkeyService)
	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}
	validationData, userOp, err := walletService.AddSubscription(merchantId, input, usePaymaster, common.Big0, int64(input.Chain))
	if err != nil {
		return nil, err
	}
	fmt.Println("Userop hash", validationData.UserOpHash)
	err = r.Cache.Set(validationData.UserOpHash, userOp)
	if err != nil {
		return nil, err
	}
	return validationData, nil
}

// ValidateSubscription is the resolver for the validateSubscription field.
func (r *mutationResolver) ValidateSubscription(ctx context.Context, input model.SubscriptionValidation) (*model.SubscriptionData, error) {
	merch, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return nil, err
	}
	_ = merch.ID

	walletService := wallet.NewWalletService(r.Database, r.TurnkeyService)
	merchantService := merchant.NewMerchantService(r.Database)

	time.Sleep(time.Second)
	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}

	opInterface, err := r.Cache.Get(input.UserOpHash)
	if err != nil {
		return nil, err
	}
	op, _ := opInterface.(map[string]any)
	sig, err := hexutil.Decode(erc4337.SUDO_MODE)
	if err != nil {
		return nil, err
	}
	partialSig, err := hexutil.Decode(input.SignedMessage)
	if err != nil {
		return nil, err
	}
	fmt.Println("partial signature - ", partialSig)
	sig = append(sig, partialSig...)
	op["signature"] = hexutil.Encode(sig)

	chain := int64(input.Chain)
	subData, key, err := walletService.ValidateSubscription(op, chain)
	if err != nil {
		return nil, err
	}
	product, err := merchantService.FetchProduct(subData.ProductID)
	if err != nil {
		return nil, err
	}
	target := product.ReceivingAddress
	// x := int64(subData.Amount)
	// Delay for a few seconds to allow the changes to be propagated onchain
	time.Sleep(15 * time.Second)
	err = walletService.ExecuteCharge(subData.WalletAddress, target, subData.Token, key, int64(subData.Amount), chain, usePaymaster)
	if err != nil {
		err = errors.Wrap(err, "ExecuteCharge() - error occurred during first time charge execution - ")
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
	merchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return nil, err
	}
	_ = merchant.ID

	_ = wallet.NewWalletService(r.Database, r.TurnkeyService)
	panic(fmt.Errorf("not implemented: FetchSubscriptions - fetchSubscriptions"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
