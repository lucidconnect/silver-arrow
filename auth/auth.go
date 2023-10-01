package auth

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateAccessKey() (publicKey, privateKey string, err error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return
	}

	publicKey = hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey))
	privateKey = hexutil.EncodeBig(pk.D)
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
