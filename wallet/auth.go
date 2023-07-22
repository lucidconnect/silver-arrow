package wallet

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

func CreateAccessKey() ([]byte, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	signer := crypto.PubkeyToAddress(privateKey.PublicKey)

	// encrypt and store the private key (not ideal)
	// I plan to use hashicorp vault to manage the private keys, for now a database will do.
	crypto.FromECDSA(privateKey)

	return signer.Bytes(), nil
}

// should return a byte array consisting of the publicKey, merchantid
// the public key is unique for each subscription hence can be used to identify the subscription
func CreateaWhitelistData(merchantId uint32, key []byte) ([]byte, error) {
	var whitelistData []byte

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