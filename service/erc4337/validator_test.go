package erc4337

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/holiman/uint256"
)

func Test_parseUint48(t *testing.T) {
	type args struct {
		value uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"0",
			args{
				0,
			},
			[]byte{0, 0, 0, 0, 0, 0},
			false,
		},
		{
			"0xFFFFFFFFFFFF",
			args{
				281474976710655,
			},
			[]byte{255, 255, 255, 255, 255, 255},
			false,
		},
		{
			"1000000000000",
			args{
				281474976710656,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseUint48(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUint48() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseUint48() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uint256(t *testing.T) {
	l := big.NewInt(32)
	lBytes := l.Bytes()
	fmt.Println(lBytes)
	fmt.Println("length ", l.BitLen())
	uintl, _ := uint256.FromBig(l)

	padded := uintl.PaddedBytes(32)
	fmt.Println(padded)
	fmt.Println(len(padded))

	fmt.Println(uint256.MustFromBig(new(big.Int).SetBytes(padded)))
	t.Fail()
}

func Test_EnableSignature(t *testing.T) {

}
