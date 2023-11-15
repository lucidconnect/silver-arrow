package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spruceid/siwe-go"
)

// var sessionStore = sessions.NewCookieStore([]byte("siwe-quickstart-secret"))
// cookies:
//
// siwe, nonce

var sessionName = "xyz.lucidconnect.auth"

func (s *Server) GetNonce() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.sessionStore.Get(r, sessionName)
		session.ID = uuid.NewString()
		nonce := siwe.GenerateNonce()
		session.Values["nonce"] = nonce
		if err := json.NewEncoder(w).Encode(session.Values["nonce"]); err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusBadRequest)
		}
		session.Save(r, w)

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(session.Values["nonce"].(string)))
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) VerifyMerchant() http.HandlerFunc {
	type requestBody struct {
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}
	type responseBody struct {
		Valid bool   `json:"status"`
		Error string `json:"string,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := &requestBody{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		session, _ := s.sessionStore.Get(r, sessionName)
		message := request.Message
		signature := request.Signature
		nonce := session.Values["nonce"].(string)
		siweObj, err := siwe.ParseMessage(message)
		if err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pkey, err := siweObj.Verify(signature, nil, &nonce, nil)
		if err != nil {
			log.Err(err).Send()
			response := &responseBody{Valid: false, Error: err.Error()}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Err(err).Send()
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		response := &responseBody{Valid: true}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Err(err).Send()
			w.WriteHeader(http.StatusInternalServerError)
		}

		address := crypto.PubkeyToAddress(*pkey)
		session.Values["siwe"] = address
		session.Options.MaxAge = int(24 * time.Hour.Seconds())

		session.Save(r, w)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}
