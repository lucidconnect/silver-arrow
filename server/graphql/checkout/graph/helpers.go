package graph

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/lucidconnect/silver-arrow/repository/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func getAuthenticatedAndActiveMerchant(ctx context.Context) (*models.Merchant, error) {
	useAuthStr := os.Getenv("USE_AUTH")
	useAuth, _ := strconv.ParseBool(useAuthStr)
	if !useAuth {
		return &models.Merchant{}, nil
	}

	merchant, err := auth.ForContext(ctx)
	if err != nil {
		err = errors.Wrapf(err, "merchant authorization failed %v", ctx)
		log.Err(err).Send()
		return nil, err
	}

	return merchant, nil
}

func validateSignature(rawString, signature, pk string) error {
	raw := []byte(rawString)
	hash := crypto.Keccak256(raw)
	sigBytes := (hexutil.MustDecode(signature))
	pbk, _ := crypto.Ecrecover(hash, sigBytes)

	pub, _ := crypto.UnmarshalPubkey(pbk)

	fmt.Println("pk -", pk)

	recoveredAddress := crypto.PubkeyToAddress(*pub)
	fmt.Println(recoveredAddress)

	if recoveredAddress.Hex() != pk {
		return errors.New("signature invalid")
	}
	return nil
}

// func parseUUID(mid string) uuid.UUID {
// 	b, _ := base64.RawURLEncoding.DecodeString(mid)
// 	id, _ := uuid.FromBytes(b)
// 	return id
// }

func getActiveMerchant(ctx context.Context) (*models.Merchant, error) {
	merchant, err := auth.ForContext(ctx)
	if err != nil {
		err = errors.Wrapf(err, "merchant authorization failed %v", ctx)
		log.Err(err).Send()
		return nil, err
	}

	return merchant, nil
}

func getActiveProduct(ctx context.Context) (*models.Product, error) {
	product, err := auth.ProductContext(ctx)
	if err != nil {
		err = errors.Wrapf(err, "no product found in context %v", ctx)
		log.Err(err).Send()
		return nil, err
	}
	return product, nil
}
