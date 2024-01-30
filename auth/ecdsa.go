package auth

import (
	"crypto/sha256"
	"encoding/json"
)

func RecoverPublicKeyFromSignature(signature string, request any) (string, error) {
	digest, err := serializeAndHashRequest(request)
	if err != nil {
		return "", err
	}

	pk, err := simpleRecover(digest, signature)
	if err != nil {
		return "", err
	}
	return pk, nil
}

func serializeAndHashRequest(request any) ([]byte, error) {
	serializedRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	checksum := sha256.New().Sum(serializedRequest)
	return checksum, nil
}
