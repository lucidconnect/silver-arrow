package auth

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lucidconnect/silver-arrow/repository"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/rs/zerolog/log"
)

var (
	merchantCtxKey      = &contextKey{"merchant"}
	AuthMerchantCtxKey  = &contextKey{"authMerchant"}
	authSignatureCtxKey = &contextKey{"authSignature"}
)

type contextKey struct {
	name string
}

func Middleware(db repository.DB) func(http.Handler) http.Handler {
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
			merchant, err := db.FetchMerchantByPublicKey(authorizationValue)
			if err != nil {
				log.Err(err).Send()
				next.ServeHTTP(w, r)
				return
			}

			merchantCtx := context.WithValue(r.Context(), merchantCtxKey, merchant)
			r = r.WithContext(merchantCtx)

			signatureCtx := context.WithValue(r.Context(), authSignatureCtxKey, signature)
			r = r.WithContext(signatureCtx)

			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) (*models.Merchant, error) {
	raw, _ := ctx.Value(merchantCtxKey).(*models.Merchant)
	if raw == nil {
		return nil, errors.New("invalid token")
	}
	return raw, nil
}

func SignatureContext(ctx context.Context, pk string) (string, error) {
	raw, _ := ctx.Value(authSignatureCtxKey).(string)
	if raw == "" {
		return "", errors.New("invalid signature")
	}
	return raw, nil
}

func AuthMerchantContext(ctx context.Context) (*models.Merchant, error) {
	raw, _ := ctx.Value(AuthMerchantCtxKey).(*models.Merchant)
	if raw == nil {
		return nil, errors.New("invalid merchant context")
	}
	return raw, nil
}

func CreateAccessKey() (publicKey, privateKey string, err error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return
	}

	privateKey = hexutil.EncodeBig(pk.D)

	raw := []byte("text")
	digest := crypto.Keccak256(raw)

	sig, _ := simpleSign(digest, privateKey)

	publicKey, err = simpleRecover(digest, sig)
	if err != nil {
		return
	}
	// publicKey = crypto.PubkeyToAddress(pk.PublicKey).Hex()
	// publicKey = hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey))
	return
}

// should return a byte array consisting of the publicKey, merchantid
// the public key is unique for each subscription hence can be used to identify the subscription
func CreateaWhitelistData(merchantId uint32, key []byte) ([]byte, error) {
	whitelistData := []byte{}

	if len(key) < 20 {
		return nil, errors.New("INVALID KEY")
	}

	if merchantId == 0 {
		return nil, errors.New("merchantId can not be 0")
	}

	tmp := make([]byte, 4)
	binary.LittleEndian.PutUint32(tmp, merchantId)
	whitelistData = append(whitelistData, key...)
	whitelistData = append(whitelistData, tmp...)

	return whitelistData, nil
}

func simpleSign(digest []byte, key string) (string, error) {
	privateKey, err := crypto.HexToECDSA(key[2:])
	if err != nil {
		return "", err
	}

	// raw := []byte(text)
	// digest := crypto.Keccak256(raw)

	signature, err := crypto.Sign(digest, privateKey)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(signature), nil
}

func simpleRecover(digest []byte, signature string) (string, error) {
	signatureBytes := (hexutil.MustDecode(signature))

	pbk, err := crypto.Ecrecover(digest, signatureBytes)
	if err != nil {
		return "", err
	}
	pub, err := crypto.UnmarshalPubkey(pbk)
	if err != nil {
		return "", err
	}

	publicKey := hexutil.Encode(crypto.CompressPubkey(pub))
	fmt.Println("compressed public key - ", publicKey)

	recoveredAddress := crypto.PubkeyToAddress(*pub)

	return recoveredAddress.Hex(), nil
}
