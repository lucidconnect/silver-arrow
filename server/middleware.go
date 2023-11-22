package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/rs/zerolog/log"
	"github.com/spruceid/siwe-go"
)

func (s *Server) JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header["Token"] == nil {
				authorizationValue := r.Header.Get("Authorization")
				if authorizationValue == "" {
					if strings.Contains(r.RequestURI, "graphiql") {
						next.ServeHTTP(w, r)
						return
					}
					goto jwtAuth
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

		jwtAuth:
			var secretKey = os.Getenv("JWT_SECRET")
			key, err := hex.DecodeString(secretKey)
			if err != nil {
				log.Err(err).Msg("decoding secret key failed")
			}

			log.Info().Msg(string(key))
			headerToken := r.Header.Get("Token")
			fmt.Println("token", headerToken)
			if headerToken != "" {
				token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
					_, ok := token.Method.(*jwt.SigningMethodHMAC)
					if !ok {
						log.Error().Msg("token method invalid")
						response := &httpResponse{Status: http.StatusBadRequest, Error: "invalid JWT"}
						writeJsonResponse(w, response)
						return nil, err
					}
					return key, nil
				})
				if err != nil {
					log.Err(err).Msg("parsing/validating token failed")
					response := &httpResponse{Status: http.StatusBadRequest, Error: "invalid JWT"}
					writeJsonResponse(w, response)
					return
				}
				if !token.Valid {
					response := &httpResponse{Status: http.StatusBadRequest, Error: "invalid JWT"}
					writeJsonResponse(w, response)
					return
				}

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					log.Error().Msg("parsing claims failed")
					response := &httpResponse{Status: http.StatusInternalServerError}
					writeJsonResponse(w, response)
					return
				}

				merchantAddress := claims["user"].(string)
				merchant, err := s.database.FetchMerchantByAddress(merchantAddress)
				if err != nil {
					log.Err(err).Send()
					next.ServeHTTP(w, r)
					return
				}
				merchantCtx := context.WithValue(r.Context(), auth.AuthMerchantCtxKey, merchant)
				r = r.WithContext(merchantCtx)

				next.ServeHTTP(w, r)
			}
		})
	}
}

func (s *Server) MerchantAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			app := os.Getenv("APP_ENV")

			if app != "production" {
				authorizationValue := r.Header.Get("Authorization")
				if authorizationValue == "" {
					if strings.Contains(r.RequestURI, "graphiql") {
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
			fmt.Println("siwe:", siweObj)

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
