package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/rs/zerolog/log"
	"github.com/spruceid/siwe-go"
)

func (s *Server) MerchantAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app := os.Getenv("APP_ENV")

			if app != "production" {
				authorizationValue := r.Header.Get("Authorization")
				if authorizationValue == "" {
					if strings.Contains(r.RequestURI, "graphiql"){
						next.ServeHTTP(w, r)
						return
					}
					goto siweAuth
				}

				if strings.HasPrefix(authorizationValue, "Bearer ") {
					authorizationValue = authorizationValue[7:]
				} else {
					next.ServeHTTP(w, r)
					return
				}

				merchant, err := s.database.FetchMerchantByPublicKey(authorizationValue)
				if err != nil {
					log.Err(err).Send()
					next.ServeHTTP(w, r)
					return
				}

				merchantCtx := context.WithValue(r.Context(), auth.AuthMerchantCtxKey, merchant)
				r = r.WithContext(merchantCtx)

				next.ServeHTTP(w, r)
				return
			}

		siweAuth:
			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				log.Err(err).Send()
			}
			siweObj := session.Values["siwe"]

			if siweObj == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			siweMsg, ok := siweObj.(*siwe.Message)
			if !ok {
				err := errors.New("parsing siwe object failed")
				log.Err(err).Send()
				w.WriteHeader(http.StatusForbidden)
				return
			}

			merchantAddress := siweMsg.GetAddress().Hex()
			merchant, err := s.database.FetchMerchantByAddress(merchantAddress)
			if err != nil {
				log.Err(err).Send()
				w.WriteHeader(http.StatusForbidden)
				return
			}
			merchantCtx := context.WithValue(r.Context(), auth.AuthMerchantCtxKey, merchant)
			r = r.WithContext(merchantCtx)

			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) CheckoutMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authorizationValue := r.Header.Get("Authorization")
			// log.Info().Msgf("merchant public key - %v", authorizationValue)
			signature := r.Header.Get("X-Lucid-Request-Signature")
			if authorizationValue == "" {
				next.ServeHTTP(w, r)
				return
			}
			if strings.HasPrefix(authorizationValue, "Bearer ") {
				authorizationValue = authorizationValue[7:]
			} else {
				next.ServeHTTP(w, r)
				return
			}
			merchant, err := s.database.FetchMerchantByPublicKey(authorizationValue)
			if err != nil {
				log.Err(err).Send()
				next.ServeHTTP(w, r)
				return
			}

			merchantCtx := context.WithValue(r.Context(), auth.MerchantCtxKey, merchant)
			r = r.WithContext(merchantCtx)

			signatureCtx := context.WithValue(r.Context(), auth.AuthSignatureCtxKey, signature)
			r = r.WithContext(signatureCtx)

			next.ServeHTTP(w, r)
		})
	}
}
