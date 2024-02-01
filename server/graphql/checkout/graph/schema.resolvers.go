package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/core"
	"github.com/lucidconnect/silver-arrow/core/gateway"
	"github.com/lucidconnect/silver-arrow/core/service/erc4337"
	"github.com/lucidconnect/silver-arrow/core/wallet"
	"github.com/lucidconnect/silver-arrow/gqlerror"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/graphql/checkout/graph/generated"
	"github.com/lucidconnect/silver-arrow/server/graphql/checkout/graph/model"
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
	var sessionId uuid.UUID
	merchant, err := getActiveMerchant(ctx)
	if err != nil {
		return "", err
	}
	product, err := getActiveProduct(ctx)
	if err != nil {
		return "", err
	}

	paymentLink, err := r.Database.FetchPaymentLinkByProduct(product.ID)
	if err != nil {
		log.Err(err).Caller().Send()
		return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, "", ctx)
	}
	merchantId := merchant.ID

	log.Info().Msgf("Authenticated Merchant: %v", merchantId)

	if input.ProductID != product.ID.String() {
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, "invalid product id supplied", ctx)
	}
	if merchantId != product.MerchantID {
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, "product not found", ctx)
	}

	if input.CheckoutSessionID != nil {
		sessionId, err = uuid.Parse(*input.CheckoutSessionID)
		if err != nil {
			log.Err(err).Caller().Send()
			return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, "malformed session id", ctx)
		}
	} else {
		// Create a new session
		sessionId = uuid.New()
		newSession := &models.CheckoutSession{
			ID:            sessionId,
			Customer:      input.WalletAddress,
			ProductID:     product.ID,
			MerchantID:    merchant.ID,
			PaymentLinkID: paymentLink.ID,
		}
		if err = r.Database.CreateCheckoutSession(newSession); err != nil {
			log.Err(err).Send()
			return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, "", ctx)
		}
	}
	pg := gateway.NewPaymentGateway(r.Database, int64(input.Chain))

	paymentIntent := core.PaymentIntent{
		Type:              core.PaymentType(input.Type),
		ProductId:         input.ProductID,
		PriceId:           input.PriceID,
		WalletAddress:     input.WalletAddress,
		FirstChargeNow:    input.FirstChargeNow,
		OwnerAddress:      input.OwnerAddress,
		Email:             *input.Email,
		Source:            input.WalletAddress,
		CheckoutSessionId: sessionId,
	}
	userop, useropHash, err := pg.CreatePaymentIntent(paymentIntent)
	if err != nil {
		return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, "creating payment intent failed", ctx)
	}
	err = r.Cache.Set(useropHash, userop)
	if err != nil {
		log.Err(err).Send()
		return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, "", ctx)
	}

	return useropHash, nil
}

// ValidatePaymentIntent is the resolver for the validatePaymentIntent field.
func (r *mutationResolver) ValidatePaymentIntent(ctx context.Context, input model.RequestValidation) (*model.TransactionData, error) {
	gatewayService := gateway.NewPaymentGateway(r.Database, int64(input.Chain))
	// merchantService := merchant.NewMerchantService(r.Database)
	var result *model.TransactionData
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
	data, err := gatewayService.ValidatePaymentIntent(op, chain, input.Type.String())
	if err != nil {
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, "payment intent validation failed", ctx)
	}

	result = &model.TransactionData{
		ID:     data.ID,
		Token:  data.Token,
		Amount: data.Amount,
		// Interval:            subData.Interval,
		ProductID:           data.ProductID,
		WalletAddress:       data.WalletAddress,
		CreatedAt:           data.CreatedAt,
		TransactionExplorer: data.TransactionExplorer,
	}

	return result, nil
}

// FetchPayment is the resolver for the fetchPayment field.
func (r *queryResolver) FetchPayment(ctx context.Context, reference string) (*model.Payment, error) {
	panic(fmt.Errorf("not implemented: FetchPayment - fetchPayment"))
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
