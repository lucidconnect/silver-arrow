package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/lucidconnect/silver-arrow/server/api"
	"github.com/rs/zerolog/log"
)

func (s *Server) CreateCheckoutSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := &api.NewCheckoutSession{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Err(err).Msg("decoding request failed")
			w.WriteHeader(http.StatusBadRequest)
		}

		sessionId := uuid.New()
		productId, err := uuid.Parse(request.ProductId)
		if err != nil {
			log.Err(err).Msg("decoding request failed")
			w.WriteHeader(http.StatusBadRequest)
		}
		merchantID, err := uuid.Parse("merchantId")
		if err != nil {
			log.Err(err).Msg("decoding request failed")
			w.WriteHeader(http.StatusBadRequest)
		}
		newSession := &models.CheckoutSession{
			ID:          sessionId,
			Chain:       request.Chain,
			Token:       request.Token,
			Amount:      request.Amount,
			Customer:    request.Customer,
			Interval:    request.Interval,
			ProductID:   productId,
			MerchantID:  merchantID,
			PaymentType: request.PaymentMode.String(),
			ChargeLater: request.ChargeLater,
		}

		if err = s.database.CreateCheckoutSession(newSession); err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
		}

		response := &api.CheckoutSessiontResponse{
			Id: sessionId.String(),
		}

		json.NewEncoder(w).Encode(response)
	}
}
