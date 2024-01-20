package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.36

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/conversions"
	"github.com/lucidconnect/silver-arrow/core/merchant"
	"github.com/lucidconnect/silver-arrow/gqlerror"
	"github.com/lucidconnect/silver-arrow/server/graphql/merchant/graph/generated"
	"github.com/lucidconnect/silver-arrow/server/graphql/merchant/graph/model"
	"github.com/rs/zerolog/log"
)

// AddProduct is the resolver for the addProduct field.
func (r *mutationResolver) AddProduct(ctx context.Context, input model.NewProduct) (*model.Product, error) {
	var product *merchant.Product
	activeMerchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return nil, err
	}
	if activeMerchant == nil {
		return nil, errors.New("merchant does not exist")
	}

	merchantService := merchant.NewMerchantService(r.Database, activeMerchant.ID)

	if input.PriceData == nil {
		return nil, errors.New("price information not supplied")
	}

	newProduct := merchant.ParseGraphqlInput(input)

	product, err = merchantService.CreateProduct(newProduct)
	if err != nil {
		return nil, err
	}
	// create price data
	amount := conversions.ParseFloatAmountToInt(input.PriceData.Token, input.PriceData.Amount)
	newPrice := &merchant.Price{
		Active:        true,
		Amount:        amount,
		Chain:         int64(input.PriceData.Chain),
		Token:         input.PriceData.Token,
		Interval:      merchant.RecuringInterval(input.PriceData.Interval),
		IntervalCount: int64(input.PriceData.IntervalCount),
		Type:          merchant.PriceType(input.PriceData.Type),
		TrialPeriod:   int64(*input.PriceData.TrialPeriod),
	}
	price, err := merchantService.CreatePrice(newPrice, product.ID.String())
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	productUpdate := map[string]interface{}{
		"default_price_id": price.ID,
	}
	if err = r.Database.UpdateProduct(product.ID, activeMerchant.ID, productUpdate); err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	// }

	result := &model.Product{
		Name:  product.Name,
		Owner: product.Owner,
		// Interval:         int(interval),
		// Amount:           input.Amount,
		ProductID:        product.ID.String(),
		ReceivingAddress: product.DepositAddress,
	}

	return result, nil
}

// UpdateProduct is the resolver for the updateProduct field.
func (r *mutationResolver) UpdateProduct(ctx context.Context, input model.ProductUpdate) (*model.Product, error) {
	panic(fmt.Errorf("not implemented: UpdateProduct - updateProduct"))
}

// CreateAccessKey is the resolver for the createAccessKey field.
func (r *mutationResolver) CreateAccessKey(ctx context.Context, input model.NewMerchantKey) (*model.MerchantAccessKey, error) {
	activeMerchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return nil, err
	}
	if activeMerchant == nil {
		return nil, errors.New("merchant does not exist")
	}

	merchantService := merchant.NewMerchantService(r.Database, activeMerchant.ID)

	accessKeys, err := merchantService.CreateAccessKeys(activeMerchant.OwnerAddress, input.Mode.String())
	if err != nil {
		return nil, err
	}
	return accessKeys, nil
}

// CreateMerchant is the resolver for the createMerchant field.
func (r *mutationResolver) CreateMerchant(ctx context.Context, input model.NewMerchant) (*model.Merchant, error) {
	merchantService := merchant.NewMerchantService(r.Database, uuid.Nil)
	result, err := merchantService.CreateMerchant(input)
	if err != nil {
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}
	return result, nil
}

// UpdateMerchantwebHookURL is the resolver for the updateMerchantwebHookUrl field.
func (r *mutationResolver) UpdateMerchantwebHookURL(ctx context.Context, webhookURL string) (*model.Merchant, error) {
	activeMerchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return nil, err
	}

	merchantService := merchant.NewMerchantService(r.Database, activeMerchant.ID)

	result, err := merchantService.UpdateMerchantWebhook(*activeMerchant, webhookURL)
	if err != nil {
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}
	return result, nil
}

// ToggleProductMode is the resolver for the toggleProductMode field.
func (r *mutationResolver) ToggleProductMode(ctx context.Context, input model.ProductModeUpdate) (model.Mode, error) {
	activeMerchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		return "nil", err
	}
	merchantService := merchant.NewMerchantService(r.Database, activeMerchant.ID)

	err = merchantService.UpdateProductMode(activeMerchant.ID, input.ProductID, input.Mode.String())
	if err != nil {
		return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}

	return input.Mode, nil
}

// CreatePaymentLink is the resolver for the createPaymentLink field.
func (r *mutationResolver) CreatePaymentLink(ctx context.Context, input model.NewPaymentLink) (string, error) {
	activeMerchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		log.Err(err).Send()
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantAuthorisationFailed, err.Error(), ctx)
	}
	merchantService := merchant.NewMerchantService(r.Database, activeMerchant.ID)

	productId, err := uuid.Parse(input.ProductID)
	if err != nil {
		log.Err(err).Send()
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, err.Error(), ctx)
	}

	product, err := r.Database.FetchProduct(productId)
	if err != nil {
		log.Err(err).Send()
		return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}

	if activeMerchant.ID != product.MerchantID {
		return "", gqlerror.ErrToGraphQLError(gqlerror.MerchantDataInvalid, err.Error(), ctx)
	}

	id, err := merchantService.CreatePaymentLink(input)
	if err != nil {
		return "", err
	}
	return id, nil
}

// DeletePaymentLink is the resolver for the deletePaymentLink field.
func (r *mutationResolver) DeletePaymentLink(ctx context.Context, id string) (string, error) {
	pid, err := uuid.Parse(id)
	if err != nil {
		log.Err(err).Caller().Msg("parsing uuid failed")
		return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}
	err = r.Database.DeletePaymentLink(pid)
	if err != nil {
		log.Err(err).Caller().Msgf("deleting payment link [%v] failed", id)
		return "", gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}
	return id, nil
}

// CreatePrice is the resolver for the createPrice field.
func (r *mutationResolver) CreatePrice(ctx context.Context, input model.NewPrice) (*model.PriceData, error) {
	activeMerchant, err := getAuthenticatedAndActiveMerchant(ctx)
	if err != nil {
		log.Err(err).Send()
		return nil, gqlerror.ErrToGraphQLError(gqlerror.MerchantAuthorisationFailed, err.Error(), ctx)
	}
	merchantService := merchant.NewMerchantService(r.Database, activeMerchant.ID)

	amount := conversions.ParseFloatAmountToInt(input.Token, input.Amount)
	newPrice := &merchant.Price{
		Active:        true,
		Amount:        amount,
		Chain:         int64(input.Chain),
		Token:         input.Token,
		Interval:      merchant.RecuringInterval(input.Interval),
		IntervalCount: int64(input.IntervalCount),
		Type:          merchant.PriceType(input.Type),
		TrialPeriod:   int64(*input.TrialPeriod),
	}
	price, err := merchantService.CreatePrice(newPrice, input.ProductID)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	priceData := &model.PriceData{
		ID:            price.ID.String(),
		Type:          model.PaymentType(price.Type),
		Active:        price.Active,
		Amount:        input.Amount,
		Token:         price.Token,
		Chain:         int(price.Chain),
		Interval:      model.IntervalType(price.Interval),
		IntervalCount: int(price.IntervalCount),
		ProductID:     price.ProductId,
		TrialPeriod:   int(price.TrialPeriod),
	}

	return priceData, nil
}

// UpdatePrice is the resolver for the updatePrice field.
func (r *mutationResolver) UpdatePrice(ctx context.Context, input *model.PriceUpdate) (*model.PriceData, error) {
	panic(fmt.Errorf("not implemented: UpdatePrice - updatePrice"))
}

// FetchOneProduct is the resolver for the fetchOneProduct field.
func (r *queryResolver) FetchOneProduct(ctx context.Context, id string, price *string) (*model.Product, error) {
	merchantService := merchant.NewMerchantService(r.Database, uuid.Nil)
	result, err := merchantService.FetchProduct(id)
	if err != nil {
		return nil, err
	}
	if price != nil {
		priceData, err := merchantService.RetrievePriceData(*price)
		if err != nil {
			return nil, err
		}
		result.PriceData = append(result.PriceData, priceData)
	} else {
		// return all prices
		priceData, err := merchantService.RetrieveProductPriceData(id)
		if err != nil {
			return nil, err
		}
		result.PriceData = priceData
	}
	return result, nil
}

// FetchProducts is the resolver for the fetchProducts field.
func (r *queryResolver) FetchProducts(ctx context.Context, owner string) ([]*model.Product, error) {
	merchantService := merchant.NewMerchantService(r.Database, uuid.Nil)
	result, err := merchantService.FetchProductsByOwner(owner)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FetchMerchantKey is the resolver for the fetchMerchantKey field.
func (r *queryResolver) FetchMerchantKey(ctx context.Context, input model.MerchantAccessKeyQuery) (string, error) {
	merchantService := merchant.NewMerchantService(r.Database, uuid.Nil)
	result, err := merchantService.FetchMerchantKey(input.MerchantAddress, input.Mode.String())
	if err != nil {
		return "", err
	}
	return result, nil
}

// FetchMerchantStats is the resolver for the fetchMerchantStats field.
func (r *queryResolver) FetchMerchantStats(ctx context.Context, owner string) (*model.MerchantStats, error) {
	merchantService := merchant.NewMerchantService(r.Database, uuid.Nil)
	stats, err := merchantService.SummarizeMerchant(owner)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// FetchMerchantInfo is the resolver for the fetchMerchantInfo field.
func (r *queryResolver) FetchMerchantInfo(ctx context.Context, owner string) (*model.Merchant, error) {
	merchant, err := r.Database.FetchMerchantByAddress(owner)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	merchantInfo := &model.Merchant{
		ID:         merchant.ID.String(),
		Name:       merchant.Name,
		Email:      merchant.Email,
		WebHookURL: merchant.WebhookUrl,
	}
	return merchantInfo, nil
}

// GetPaymentLink is the resolver for the getPaymentLink field.
func (r *queryResolver) GetPaymentLink(ctx context.Context, id string) (*model.PaymentLinkDetails, error) {
	merchantService := merchant.NewMerchantService(r.Database, uuid.Nil)

	paymentLinkQuery := merchant.PaymentLinkQueryParams{
		PaymentLinkId: &id,
	}
	paymentLinkDetails, err := merchantService.FetchPaymentLink(paymentLinkQuery)
	if err != nil {
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}

	return paymentLinkDetails, nil
}

// GetMerchantPaymentLinks is the resolver for the getMerchantPaymentLinks field.
func (r *queryResolver) GetMerchantPaymentLinks(ctx context.Context, merchantID string) ([]*model.PaymentLinkDetails, error) {
	var paymentLinkDetails []*model.PaymentLinkDetails
	mid, err := uuid.Parse(merchantID)
	if err != nil {
		log.Err(err).Caller().Msg("parsing uuid failed")
		return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, err.Error(), ctx)
	}
	paymentLinks, err := r.Database.FetchPaymentLinkByMerchant(mid)
	if err != nil {
		log.Err(err).Caller().Msg("parsing uuid failed")
	}

	for _, paymentLink := range paymentLinks {
		interval := conversions.ParseNanoSecondsToDay(paymentLink.Price.IntervalCount)
		amount := conversions.ParseTransferAmountFloat(paymentLink.Price.Token, paymentLink.Price.Amount)

		paymentLinkDetail := &model.PaymentLinkDetails{
			ID:            paymentLink.ID.String(),
			Mode:          paymentLink.Product.Mode,
			ProductID:     paymentLink.ProductID.String(),
			MerchantID:    paymentLink.MerchantID.String(),
			MerchantName:  paymentLink.MerchantName,
			ProductName:   paymentLink.Product.Name,
			IntervalCount: int(interval),
			Interval:      paymentLink.Price.Interval,
			CallbackURL:   paymentLink.CallbackURL,
			Amount:        amount,
			Token:         paymentLink.Price.Token,
			Chain:         int(paymentLink.Product.Chain),
		}
		paymentLinkDetails = append(paymentLinkDetails, paymentLinkDetail)
	}
	return paymentLinkDetails, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) UpdateMerchantDetails(ctx context.Context, input model.MerchantUpdate) (*model.Merchant, error) {

	return nil, gqlerror.ErrToGraphQLError(gqlerror.InternalError, "err.Error()", ctx)

}
