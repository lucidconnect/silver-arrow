package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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
		merchantID, err := uuid.Parse("merchantId")
		if err != nil {
			log.Err(err).Msg("decoding request failed")
			response = &httpResponse{
				Status: http.StatusBadRequest,
				Error:  "merchant id is invalid",
			}
			writeJsonResponse(w, response)
			return
		}

		paymentLink, err := s.database.FetchPaymentLinkByProduct(productId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				response = &httpResponse{
					Status: http.StatusNotFound,
					Error:  "product does not exist",
				}
			} else {
				log.Err(err).Caller().Send()
				response = &httpResponse{
					Status: http.StatusInternalServerError,
				}
			}
			writeJsonResponse(w, response)
			return
		}
		newSession := &models.CheckoutSession{
			ID:            sessionId,
			Customer:      request.Customer,
			ProductID:     productId,
			MerchantID:    merchantID,
			PaymentLinkID: paymentLink.ID,
			PaymentLink:   *paymentLink,
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

		checkoutSession := &api.CheckoutSessionResponse{
			SessionId: sessionId.String(),
		}
		response = &httpResponse{
			Status: http.StatusOK,
			Data:   checkoutSession,
		}
		writeJsonResponse(w, response)
	}
}
