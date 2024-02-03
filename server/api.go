package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/conversions"
	"github.com/lucidconnect/silver-arrow/core"
	"github.com/lucidconnect/silver-arrow/core/merchant"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/api"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (s *Server) CreateCheckoutSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &api.NewCheckoutSession{}
		response := &httpResponse{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Err(err).Caller().Msg("decoding request failed")
			response = &httpResponse{
				Status: http.StatusBadRequest,
			}
			writeJsonResponse(w, response)
			return
		}

		sessionId := uuid.New()
		priceId, err := uuid.Parse(request.PriceId)
		if err != nil {
			log.Err(err).Caller().Msg("decoding request failed")
			response = &httpResponse{
				Status: http.StatusBadRequest,
				Error:  "product id is invalid",
			}
			writeJsonResponse(w, response)
			return
		}
		auth := strings.Split(r.Header.Get("Authorization"), " ")[1]
		merchant, err := s.database.FetchMerchantByPublicKey(auth)
		// merchantID, err := uuid.Parse("merchantId")
		if err != nil {
			log.Err(err).Caller().Send()
			response = &httpResponse{
				Status: http.StatusBadRequest,
				Error:  "merchant id is invalid",
			}
			writeJsonResponse(w, response)
			return
		}
		merchantID := merchant.ID
		price, err := s.database.FetchPrice(priceId)
		if err != nil {
			log.Err(err).Caller().Send()
			response = &httpResponse{
				Status: http.StatusBadRequest,
				Error:  "failed to load price object",
			}
			writeJsonResponse(w, response)
			return
		}

		paymentLink, err := s.database.FetchPaymentLinkByProduct(price.ProductID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// create a payment link
				paymentLink = &models.PaymentLink{
					ID:          uuid.New(),
					MerchantID:  merchant.ID,
					CallbackURL: request.CallbackUrl,
					ProductID:   price.ProductID,
					CreatedAt:   time.Now(),
				}
				err = s.database.CreatePaymentLink(paymentLink)
				if err != nil {
					log.Err(err).Msgf("could not create payment link for product [%v]", price.ProductID)
					response = &httpResponse{
						Status: http.StatusInternalServerError,
					}
					writeJsonResponse(w, response)
					return
				}
			} else {
				log.Err(err).Caller().Send()
				response = &httpResponse{
					Status: http.StatusInternalServerError,
				}
				writeJsonResponse(w, response)
				return
			}
		}
		newSession := &models.CheckoutSession{
			ID:            sessionId,
			Customer:      request.Customer,
			ProductID:     price.ProductID,
			MerchantID:    merchantID,
			PaymentLinkID: paymentLink.ID,
			CallbackURL:   request.CallbackUrl,
		}

		if err = s.database.CreateCheckoutSession(newSession); err != nil {
			log.Err(err).Send()
			response = &httpResponse{
				Status: http.StatusInternalServerError,
				Error:  "failed to create checkout session, please contact support.",
			}
			writeJsonResponse(w, response)
			return
		}
		var url string
		environment := os.Getenv("APP_ENV")
		switch environment {
		case "staging":
			url = fmt.Sprintf("https://pay.staging.lucidconnect.xyz/c/%v", sessionId.String())
		case "production":
			url = fmt.Sprintf("https://pay.lucidconnect.xyz/c/%v", sessionId.String())
		}
		checkoutSession := &api.CheckoutSessionResponse{
			SessionId: sessionId.String(),
			Url:       url,
		}
		response = &httpResponse{
			Status: http.StatusOK,
			Data:   checkoutSession,
		}
		writeJsonResponse(w, response)
	}
}

func (s *Server) CreateNewProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &api.NewProduct{}
		response := &httpResponse{}

		auth := strings.Split(r.Header.Get("Authorization"), " ")[1]
		activeMerchant, err := s.database.FetchMerchantByPublicKey(auth)
		if err != nil {
			log.Err(err).Msg("decoding request failed")
			response = &httpResponse{
				Status: http.StatusUnauthorized,
				Error:  "merchant is invalid",
			}
			writeJsonResponse(w, response)
			return
		}

		merchantService := merchant.NewMerchantService(s.database, activeMerchant.ID)

		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Err(err).Caller().Msg("decoding request failed")
			response = &httpResponse{
				Status: http.StatusBadRequest,
			}
			writeJsonResponse(w, response)
			return
		}
		request.Owner = activeMerchant.OwnerAddress
		// productID := uuid.New()

		newProduct := merchant.ParseNewApiProduct(*request)
		product, err := merchantService.CreateProduct(newProduct)
		if err != nil {
			log.Err(err).Send()
			response = &httpResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}
			writeJsonResponse(w, response)
			return
		}

		productResponse := &api.ProductResponse{
			ID:             product.ID.String(),
			Name:           request.Name,
			FirstChargeNow: true,
		}
		// create price
		priceId := uuid.New()
		amount := conversions.ParseFloatAmountToIntDenomination(request.PriceData.Token, request.PriceData.Amount)
		var trialPeriod int64
		if request.PriceData.TrialPeriod != 0 {
			trialPeriod = int64(request.PriceData.TrialPeriod)
		}
		newPrice := &merchant.Price{
			Active:       true,
			Amount:       amount,
			Chain:        request.PriceData.Chain,
			Token:        request.PriceData.Token,
			IntervalUnit: core.RecuringInterval(request.PriceData.Interval),
			Interval:     int64(request.PriceData.IntervalCount),
			Type:         request.PriceData.Type,
			TrialPeriod:  trialPeriod,
		}
		price, err := merchantService.CreatePrice(newPrice, product.ID.String())
		if err != nil {
			log.Err(err).Send()
			response = &httpResponse{
				Status: http.StatusInternalServerError,
				Error:  err.Error(),
			}
			writeJsonResponse(w, response)
			return
		}

		productUpdate := map[string]interface{}{
			"default_price_id": price.ID,
		}

		if err = s.database.UpdateProduct(product.ID, activeMerchant.ID, productUpdate); err != nil {
			log.Err(err).Caller().Send()
		}

		productResponse.DefaultPriceData = api.PriceDataResponse{
			ID:            priceId.String(),
			Active:        true,
			Amount:        amount,
			Token:         request.PriceData.Token,
			Chain:         request.PriceData.Chain,
			Type:          request.PriceData.Type,
			Interval:      request.PriceData.Interval,
			IntervalCount: request.PriceData.IntervalCount,
			ProductID:     product.ID.String(),
		}
		response = &httpResponse{
			Status: http.StatusInternalServerError,
			Data:   product,
		}
		writeJsonResponse(w, response)
	}
}
