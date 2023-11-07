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
	"github.com/lucidconnect/silver-arrow/api/graphql/wallet/graph/generated"
	"github.com/lucidconnect/silver-arrow/api/graphql/wallet/graph/model"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/erc20"
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

// CreatePaymentIntent is the resolver for the createPaymentIntent field.
func (r *mutationResolver) CreatePaymentIntent(ctx context.Context, input model.PaymentIntent) (string, error) {
	merchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return "", err
	}
	merchantId := merchant.ID
	signature, err := auth.SignatureContext(ctx, merchant.PublicKey)
	if err != nil {
		return "", err
	}
	log.Info().Msgf("Authenticated Merchant: %v", merchantId)
	// validate signature
	// amount:token:interval:productId
	signatureCheck := fmt.Sprintf("%v", input.Amount) + ":" + input.Token + ":" + fmt.Sprintf("%v", input.Interval) + ":" + input.ProductID
	err = validateSignature(signatureCheck, signature, merchant.PublicKey)
	if err != nil {
		log.Err(err).Ctx(ctx).Send()
		err = errors.New("request signature is invalid")
		return "", err
	}

	walletService := wallet.NewWalletService(r.Database, r.TurnkeyService)
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
		newSubscription := model.NewSubscription{
			Chain:         input.Chain,
			Token:         input.Token,
			Email:         *input.Email,
			Amount:        input.Amount,
			Interval:      input.Interval,
			ProductID:     input.ProductID,
			WalletAddress: input.WalletAddress,
			OwnerAddress:  input.OwnerAddress,
		}
		validationData, userOp, err := walletService.AddSubscription(merchantId, newSubscription, usePaymaster, common.Big0, int64(input.Chain))
		if err != nil {
			return "", err
		}
		fmt.Println("Userop hash", validationData.UserOpHash)
		err = r.Cache.Set(validationData.UserOpHash, userOp)
		if err != nil {
			return "", err
		}
		useropHash = validationData.UserOpHash
	default:
		return "", errors.New("unsupported payment type")
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
	product, err := merchantService.FetchProduct(*subData.ProductID)
	if err != nil {
		return nil, err
	}
	target := product.ReceivingAddress
	// x := int64(subData.Amount)
	// Delay for a few seconds to allow the changes to be propagated onchain
	time.Sleep(15 * time.Second)
	transactionHash, err := walletService.ExecuteCharge(subData.WalletAddress, target, subData.Token, key, int64(subData.Amount), chain, usePaymaster)
	if err != nil {
		err = errors.Wrap(err, "ExecuteCharge() - error occurred during first time charge execution - ")
		return subData, err
	}
	explorer, err := erc20.GetChainExplorer(chain)
	if err != nil {
		log.Err(err).Send()
	}
	transactionDetails := fmt.Sprintf("%v/tx/%v", explorer, transactionHash)

	subData.TransactionHash = &transactionHash
	subData.TransactionExplorer = &transactionDetails

	return subData, nil
}

// ModifySubscriptionState is the resolver for the modifySubscriptionState field.
func (r *mutationResolver) ModifySubscriptionState(ctx context.Context, input model.SubscriptionMod) (string, error) {
	var err error
	var result string
	walletService := wallet.NewWalletService(r.Database, nil)

	switch input.Toggle {
	case model.StatusToggleCancel:
		// cancel subscription
		result, err = walletService.CancelSubscription(input.SubscriptionID)
		if err != nil {
			log.Err(err).Send()
			return "", err
		}
	case model.StatusToggleDisable:
		// temporary disbale
		result, err = walletService.DisableSubscription(input.SubscriptionID)
		if err != nil {
			log.Err(err).Send()
			return "", err
		}
	case model.StatusToggleEnable:
		// reenable subscription
		result, err = walletService.EnableSubscription(input.SubscriptionID)
		if err != nil {
			log.Err(err).Send()
			return "", err
		}
	}
	return result, err
}

// InitiateTransferRequest is the resolver for the initiateTransferRequest field.
func (r *mutationResolver) InitiateTransferRequest(ctx context.Context, input model.NewTransferRequest) (string, error) {
	var sponsored bool
	switch os.Getenv("USE_PAYMASTER") {
	case "TRUE":
		sponsored = true
	default:
		sponsored = false
	}

	walletService := wallet.NewWalletService(r.Database, nil)
	validationData, userop, err := walletService.InitiateTransfer(input.Sender, input.Target, input.Token, input.Amount, int64(input.Chain), sponsored)
	if err != nil {
		log.Err(err).Send()
		return "", errors.New("internal server error")
	}

	err = r.Cache.Set(validationData.UserOpHash, userop)
	if err != nil {
		return "", err
	}

	return validationData.UserOpHash, nil
}

// ValidateTransferRequest is the resolver for the validateTransferRequest field.
func (r *mutationResolver) ValidateTransferRequest(ctx context.Context, input model.RequestValidation) (*model.TransactionData, error) {
	opInterface, err := r.Cache.Get(input.UserOpHash)
	if err != nil {
		log.Err(err).Send()
		return nil, errors.New("internal server error")
	}
	op, _ := opInterface.(map[string]any)
	sig, err := hexutil.Decode(erc4337.SUDO_MODE)
	if err != nil {
		log.Err(err).Send()
		return nil, errors.New("internal server error")
	}
	partialSig, err := hexutil.Decode(input.SignedMessage)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, errors.New("invalid signature format")
	}
	fmt.Println("partial signature - ", partialSig)
	sig = append(sig, partialSig...)
	op["signature"] = hexutil.Encode(sig)

	chain := int64(input.Chain)

	walletService := wallet.NewWalletService(r.Database, nil)

	td, err := walletService.ValidateTransfer(op, chain)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, errors.New("internal server error")
	}

	return td, nil
}

// FetchSubscriptions is the resolver for the fetchSubscriptions field.
func (r *queryResolver) FetchSubscriptions(ctx context.Context, account string) ([]*model.SubscriptionData, error) {
	ws := wallet.NewWalletService(r.Database, r.TurnkeyService)
	subs, err := ws.FetchSubscriptions(account)
	if err != nil {
		err = errors.New("failed to fetch subscriptions")
		return nil, err
	}
	return subs, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
