package server

import (
	"context"
	"encoding/hex"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/rs/zerolog/log"
)

func (s *Server) JWTMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header["Token"] == nil {
				authorizationValue := r.Header.Get("Authorization")
				if authorizationValue == "" {

					// if strings.Contains(r.RequestURI, "graphiql") {
					next.ServeHTTP(w, r)
					return
					// }
					// goto jwtAuth
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

			// jwtAuth:
			var secretKey = os.Getenv("JWT_SECRET")
			key, err := hex.DecodeString(secretKey)
			if err != nil {
				log.Err(err).Msg("decoding secret key failed")
			}

			headerToken := r.Header.Get("Token")
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

				log.Info().Msgf("token claims %v", token.Claims)

				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					log.Error().Msg("parsing claims failed")
					response := &httpResponse{Status: http.StatusInternalServerError}
					writeJsonResponse(w, response)
					return
				}
				// siwe, err := parseSiweClaim(claims)
				// if err != nil {
				// 	log.Err(err).Msg("parsing siwe claims failed")
				// 	response := &httpResponse{Status: http.StatusInternalServerError}
				// 	writeJsonResponse(w, response)
				// 	return
				// }
				merchantAddress := claims["address"].(string)
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

// func parseSiweClaim(claim interface{}) (*siwe.Message, error) {
// 	fmt.Println(claim)
// 	claimStr, ok := claim.(string)
// 	if !ok {
// 		err := errors.New("parsing claim failed")
// 		return nil, err
// 	}

// 	siweClaim, err := siwe.ParseMessage(claimStr)
// 	if err != nil {
// 		log.Err(err).Send()
// 		return nil, err
// 	}
// 	return siweClaim, nil
// }

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
			// use the public key to fetch the key,
			// determine it's mode (test or live?)
			// fetch the merchant with the merchantID,
			// attach the mode to the context
			// attach the merchant to the context
			// attach the public key to the context

			merchant, err := s.database.FetchMerchantByPublicKey(authorizationValue)
			if err != nil {
				log.Err(err).Send()
				next.ServeHTTP(w, r)
				return
			}

			key := merchant.MerchantAccessKeys[0]
			merchantCtx := context.WithValue(r.Context(), auth.MerchantCtxKey, merchant)
			r = r.WithContext(merchantCtx)

			signatureCtx := context.WithValue(r.Context(), auth.AuthSignatureCtxKey, signature)
			r = r.WithContext(signatureCtx)

			modeCtx := context.WithValue(r.Context(), auth.ModeCtxKey, &key)
			r = r.WithContext(modeCtx)
			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) PaymentLinkMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header["Token"] == nil {
				next.ServeHTTP(w, r)
				return
			}

			headerToken := r.Header.Get("Token")
			claims, err := parseJwt(headerToken)
			if err != nil {
				log.Info().Msgf("token: %v", headerToken)
				log.Err(err).Msg("error parsing jwt")
				response := &httpResponse{Status: http.StatusUnauthorized, Error: "invalid jwt"}
				writeJsonResponse(w, response)
				return
			}

			merchantId, productId, err := parsePaymentLinkClaims(claims)
			if err != nil {
				log.Err(err).Send()
				response := &httpResponse{Status: http.StatusUnauthorized, Error: "invalid merchant/product id"}
				writeJsonResponse(w, response)
				return
			}

			// merchantAddress := claims["address"].(string)
			merchant, err := s.database.FetchMerchantById(merchantId)
			if err != nil {
				log.Err(err).Send()
				response := &httpResponse{Status: http.StatusUnauthorized, Error: "merchant not found"}
				writeJsonResponse(w, response)
				return
			}
			merchantCtx := context.WithValue(r.Context(), auth.MerchantCtxKey, merchant)
			r = r.WithContext(merchantCtx)

			product, err := s.database.FetchProduct(productId)
			if err != nil {
				log.Err(err).Send()
				response := &httpResponse{Status: http.StatusUnauthorized, Error: "product not found"}
				writeJsonResponse(w, response)
				return
			}
			productCtx := context.WithValue(r.Context(), auth.ProductCtxKey, product)
			r = r.WithContext(productCtx)

			next.ServeHTTP(w, r)
		})
	}
}
