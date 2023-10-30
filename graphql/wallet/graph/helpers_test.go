package graph

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lucidconnect/silver-arrow/auth"
	"github.com/stretchr/testify/assert"
)

func Test_validateSignature(t *testing.T) {
	pk, sk, err := auth.CreateAccessKey()
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	text := "hello world"
	signature := simpleSign(text, sk)
	type args struct {
		rawString string
		signature string
		pk        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"simple sign",
			args{
				text,
				signature,
				pk,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateSignature(tt.args.rawString, tt.args.signature, tt.args.pk); (err != nil) != tt.wantErr {
				t.Errorf("validateSignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func simpleSign(text, sk string) string {
	privateKey, err := crypto.HexToECDSA(sk[2:])
	if err !=nil {
		panic(err)
	} 

	raw := []byte(text)
	digest := crypto.Keccak256(raw)

	signature, _ := crypto.Sign(digest, privateKey)
	return hexutil.Encode(signature)
}
