package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/api"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// func (s *Server) CreateCheckoutSession() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		request := &api.NewCheckoutSession{}
// 		response := &httpResponse{}
// 		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
// 			log.Err(err).Caller().Msg("decoding request failed")
// 			response = &httpResponse{
// 				Status: http.StatusBadRequest,
// 			}
// 			writeJsonResponse(w, response)
// 			return
// 		}

// 		sessionId := uuid.New()
// 		productId, err := uuid.Parse(request.ProductId)
// 		if err != nil {
// 			log.Err(err).Caller().Msg("decoding request failed")
// 			response = &httpResponse{
// 				Status: http.StatusBadRequest,
// 				Error:  "product id is invalid",
// 			}
// 			writeJsonResponse(w, response)
// 			return
// 		}
// 		auth := strings.Split(r.Header.Get("Authorization"), " ")[1]
// 		merchant, err := s.database.FetchMerchantByPublicKey(auth)
// 		// merchantID, err := uuid.Parse("merchantId")
// 		if err != nil {
// 			log.Err(err).Msg("decoding request failed")
// 			response = &httpResponse{
// 				Status: http.StatusBadRequest,
// 				Error:  "merchant id is invalid",
// 			}
// 			writeJsonResponse(w, response)
// 			return
// 		}
// 		merchantID := merchant.ID

// 		paymentLink, err := s.database.FetchPaymentLinkByProduct(productId)
// 		if err != nil {
// 			if err == gorm.ErrRecordNotFound {
// 				response = &httpResponse{
// 					Status: http.StatusNotFound,
// 					Error:  "product does not exist",
// 				}
// 			} else {
// 				log.Err(err).Caller().Send()
// 				response = &httpResponse{
// 					Status: http.StatusInternalServerError,
// 				}
// 			}
// 			writeJsonResponse(w, response)
// 			return
// 		}
// 		newSession := &models.CheckoutSession{
// 			ID:            sessionId,
// 			Customer:      request.Customer,
// 			ProductID:     productId,
// 			MerchantID:    merchantID,
// 			PaymentLinkID: paymentLink.ID,
// 			PaymentLink:   *paymentLink,
// 		}

// 		if err = s.database.CreateCheckoutSession(newSession); err != nil {
// 			log.Err(err).Send()
// 			response = &httpResponse{
// 				Status: http.StatusInternalServerError,
// 				Error:  "failed to create checkout session, please contact support.",
// 			}
// 			writeJsonResponse(w, response)
// 			return
// 		}

// 		checkoutSession := &api.CheckoutSessionResponse{
// 			SessionId: sessionId.String(),
// 		}
// 		response = &httpResponse{
// 			Status: http.StatusOK,
// 			Data:   checkoutSession,
// 		}
// 		writeJsonResponse(w, response)
// 	}
// }

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
		productId, err := uuid.Parse(request.ProductId)
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
			log.Err(err).Msg("decoding request failed")
			response = &httpResponse{
				Status: http.StatusBadRequest,
				Error:  "merchant id is invalid",
			}
			writeJsonResponse(w, response)
			return
		}
		merchantID := merchant.ID

		paymentLink, err := s.database.FetchPaymentLinkByProduct(productId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// create a payment link
				paymentLink = &models.PaymentLink{
					ID:          uuid.New(),
					MerchantID:  merchant.ID,
					CallbackURL: request.CallbackUrl,
					ProductID:   productId,
					CreatedAt:   time.Now(),
				}
				err = s.database.CreatePaymentLink(paymentLink)
				if err != nil {
					log.Err(err).Msgf("could not create payment link for product [%v]", productId)
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
			ProductID:     productId,
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
