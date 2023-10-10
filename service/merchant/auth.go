package merchant

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/helicarrierstudio/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
)

var merchantCtxKey = &contextKey{"merchant"}

// Authenticated merchant => PublicKey:PrivateKey pair
// var authMerchantCtxKey = &contextKey{"authMerchant"}

// var authSignatureCtxKey = &contextKey{"authSignature"}

type contextKey struct {
	name string
}

func (m *MerchantService) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationValue := r.Header.Get("Authorization")
			log.Info().Msgf("merchant public key - %v", authorizationValue)
			// privateKeyValue := r.Header.Get("Private-Key")
			signature := r.Header.Get("X-Lucid-Request-Signature")
			requestHash := r.Header.Get("Lucid-Request-Hash")
			hash, err := hexutil.Decode(fmt.Sprintf("0x%v", requestHash))
			if err != nil {
				log.Err(err).Send()
				next.ServeHTTP(w, r)
				return
			}

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
			merchant, err := m.repository.FetchMerchantByPublicKey(authorizationValue)
			if err != nil {
				log.Err(err).Send()
				next.ServeHTTP(w, r)
				return
			}

			signatureBytes, err := hexutil.Decode(signature)
			if err != nil {
				log.Err(err).Send()
				next.ServeHTTP(w, r)
				return
			}

			pub, err := crypto.SigToPub(hash, signatureBytes)
			if err != nil {
				log.Err(err).Send()
				next.ServeHTTP(w, r)
				return
			}

			if strings.Compare(authorizationValue, hexutil.Encode(crypto.CompressPubkey(pub))) != 0 {
				next.ServeHTTP(w, r)
				return
			}

			baseCtx := context.WithValue(r.Context(), merchantCtxKey, merchant)
			r = r.WithContext(baseCtx)

			// signatureCtx := context.WithValue(r.Context(), authSignatureCtxKey, signature)
			// r = r.WithContext(signatureCtx)

			// authMerchant, err := engine.FetchMerchantByAuthKey(authorizationValue, privateKeyValue)
			// if err != nil {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }

			// ctx := context.WithValue(r.Context(), authMerchantCtxKey, authMerchant)
			// r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) (*models.Merchant, error) {
	raw, _ := ctx.Value(merchantCtxKey).(*models.Merchant)
	if raw == nil {
		return nil, errors.New("Token is Invalid")
	}
	return raw, nil
}
