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
	"github.com/lucidconnect/silver-arrow/gqlerror"
	"github.com/lucidconnect/silver-arrow/server/graphql/paymentLink/graph/generated"
	"github.com/lucidconnect/silver-arrow/server/graphql/paymentLink/graph/model"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/lucidconnect/silver-arrow/service/merchant"
	"github.com/lucidconnect/silver-arrow/service/wallet"
	"github.com/rs/zerolog/log"
)

// CreatePaymentIntent is the resolver for the createPaymentIntent field.
func (r *mutationResolver) CreatePaymentIntent(ctx context.Context, input model.PaymentIntent) (string, error) {
	merchant, err := getActiveMerchant(ctx)
	if err != nil {
		return "", err
	}
	product, err := getActiveProduct(ctx)
	if err != nil {
		return "", err
	}

	merchantId := merchant.ID

	log.Info().Msgf("Authenticated Merchant: %v", merchantId)

	if input.ProductID != product.ID.String() {
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, "invalid product id supplied", ctx)
	}
	if merchantId != product.MerchantID {
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, "product not found", ctx)
	}

	walletService := wallet.NewWalletService(r.Database, int64(input.Chain))
	var usePaymaster bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		usePaymaster = true
	default:
		usePaymaster = false
	}
	var useropHash string

	switch input.Type {
	case model.PaymentTypeRecurring:
		var nextCharge time.Time

		if input.FirstChargeNow {
			nextCharge = time.Now()
		}

		var email string
		if input.Email != nil {
			email = *input.Email
		}
		newSubscription := wallet.NewSubscription{
			Chain:          input.Chain,
			Token:          input.Token,
			Email:          email,
			Amount:         input.Amount,
			Interval:       input.Interval,
			ProductID:      product.ID,
			ProductName:    product.Name,
			OwnerAddress:   input.OwnerAddress,
			WalletAddress:  input.WalletAddress,
			DepositAddress: product.DepositAddress,
			NextChargeDate: &nextCharge,
		}

		validationData, userOp, err := walletService.AddSubscription(merchantId, newSubscription, usePaymaster, common.Big0, int64(input.Chain))
		if err != nil {
			return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, "Couldn't add subscription to user's wallet", ctx)
		}
		fmt.Println("Userop hash", validationData.UserOpHash)
		err = r.Cache.Set(validationData.UserOpHash, userOp)
		if err != nil {
			log.Err(err).Send()
			return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, "Couldn't add subscription to user's wallet", ctx)
		}
		useropHash = validationData.UserOpHash
	default:
		return "", gqlerror.ErrToGraphQLError(gqlerror.NilError, "unsupported payment type", ctx)
	}

	return useropHash, nil
}

// ValidatePaymentIntent is the resolver for the validatePaymentIntent field.
func (r *mutationResolver) ValidatePaymentIntent(ctx context.Context, input model.RequestValidation) (*model.TransactionData, error) {
	walletService := wallet.NewWalletService(r.Database, int64(input.Chain))
	// merchantService := merchant.NewMerchantService(r.Database)

	// time.Sleep(time.Second)

	opInterface, err := r.Cache.Get(input.UserOpHash)
	if err != nil {
		log.Err(err).Send()
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, "Subscription validation failed", ctx)
	}
	op, _ := opInterface.(map[string]any)
	sig, err := hexutil.Decode(erc4337.SUDO_MODE)
	if err != nil {
		log.Err(err).Send()
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, "Subscription validation failed", ctx)
	}
	partialSig, err := hexutil.Decode(input.SignedMessage)
	if err != nil {
		return nil, err
	}
	fmt.Println("partial signature - ", partialSig)
	sig = append(sig, partialSig...)
	op["signature"] = hexutil.Encode(sig)

	chain := int64(input.Chain)
	subData, err := walletService.ValidateSubscription(op, chain)
	if err != nil {
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, "Subscription validation failed", ctx)
	}

	result := &model.TransactionData{
		ID:                  subData.ID,
		Token:               subData.Token,
		Amount:              subData.Amount,
		Interval:            subData.Interval,
		ProductID:           subData.ProductID,
		WalletAddress:       subData.WalletAddress,
		CreatedAt:           subData.CreatedAt,
		TransactionExplorer: subData.TransactionExplorer,
	}
	return result, nil
}

// GetPaymentLink is the resolver for the getPaymentLink field.
func (r *queryResolver) GetPaymentLink(ctx context.Context, id string) (*model.PaymentLinkDetails, error) {
	merchantService := merchant.NewMerchantService(r.Database)

	paymentLinkQuery := merchant.PaymentLinkQueryParams{
		PaymentLinkId: &id,
	}
	pd, err := merchantService.FetchPaymentLink(paymentLinkQuery)
	if err != nil {
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}
	paymentLinkDetails := &model.PaymentLinkDetails{
		ID: pd.ID,
		// Mode: pd.Mode,
		ProductID:    pd.ProductID,
		MerchantID:   pd.MerchantID,
		Amount:       pd.Amount,
		Token:        pd.Token,
		Chain:        pd.Chain,
		ProductName:  pd.ProductName,
		MerchantName: pd.MerchantName,
		Interval:     pd.Interval,
		CallbackURL:  pd.CallbackURL,
	}
	return paymentLinkDetails, nil
}

// GetBillingHistory is the resolver for the getBillingHistory field.
func (r *queryResolver) GetBillingHistory(ctx context.Context, walletAddress string, productID string) ([]*model.BillingHistory, error) {
	var billingHistory []*model.BillingHistory

	walletService := wallet.NewWalletService(r.Database, 0)
	bh, err := walletService.FetchUserBillingHistory(walletAddress, productID)
	if err != nil {
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}

	for _, b := range bh {

		billing := &model.BillingHistory{
			Date:        b.Date,
			Amount:      b.Amount,
			ExplorerURL: b.ExplorerURL,
		}

		billingHistory = append(billingHistory, billing)

	}
	return billingHistory, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
