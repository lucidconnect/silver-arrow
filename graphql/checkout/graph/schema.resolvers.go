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
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/gqlerror"
	"github.com/lucidconnect/silver-arrow/graphql/checkout/graph/generated"
	"github.com/lucidconnect/silver-arrow/graphql/checkout/graph/model"
	"github.com/lucidconnect/silver-arrow/service/erc4337"
	"github.com/lucidconnect/silver-arrow/service/wallet"
	"github.com/rs/zerolog/log"
)

// AddAccount is the resolver for the addAccount field.
func (r *mutationResolver) AddAccount(ctx context.Context, input model.Account) (string, error) {
	address := common.HexToAddress(input.Address)

	walletService := wallet.NewWalletService(r.Database, 0)
	// should check if the account is deployed
	// deploy if not deployed
	newAccount := wallet.Account{
		Email:   input.Email,
		Address: input.Address,
		Signer:  input.Signer,
	}

	err := walletService.AddAccount(newAccount)
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}

// CreatePaymentIntent is the resolver for the createPaymentIntent field.
func (r *mutationResolver) CreatePaymentIntent(ctx context.Context, input model.PaymentIntent) (string, error) {
	merchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return "", err
	}

	key, err := auth.KeyModeContext(ctx)
	if err != nil {
		log.Err(err).Msg("no key found in context")
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantAuthorisationFailed, err.Error(), ctx)
	}

	merchantId := merchant.ID
	signature, err := auth.SignatureContext(ctx, merchant.MerchantAccessKeys[0].PublicKey)
	if err != nil {
		log.Err(err).Msg("no signature in context")
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantAuthorisationFailed, err.Error(), ctx)
	}

	log.Info().Msgf("Authenticated Merchant: %v", merchantId)
	// validate signature
	// amount:token:interval:productId
	signatureCheck := fmt.Sprintf("%v", input.Amount) + ":" + input.Token + ":" + fmt.Sprintf("%v", input.Interval) + ":" + input.ProductID
	err = validateSignature(signatureCheck, signature, key.PublicKey)
	if err != nil {
		log.Debug().Err(err).Ctx(ctx).Send()
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantAuthorisationFailed, err.Error(), ctx)
	}

	// productId := parseUUID(input.ProductID)
	productId := uuid.MustParse(input.ProductID)
	product, err := r.Database.FetchProduct(productId)
	if err != nil {
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, "product not found", ctx)
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
			ProductID:      productId,
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
	merch, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return nil, err
	}
	_ = merch.ID

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

// FetchPayment is the resolver for the fetchPayment field.
func (r *queryResolver) FetchPayment(ctx context.Context, reference string) (*model.Payment, error) {
	panic(fmt.Errorf("not implemented: FetchPayment - fetchPayment"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
