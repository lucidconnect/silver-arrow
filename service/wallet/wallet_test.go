package wallet

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_daysToNanoSeconds(t *testing.T) {
	duration := daysToNanoSeconds(28)
	fmt.Println(duration.Nanoseconds())
	assert.Equal(t, time.Duration(2419200000000000), duration)

	days := nanoSecondsToDay(2419200000000000)
	fmt.Println(days)
	assert.Equal(t, int64(28), days)
}

func Test_amountToWei(t *testing.T) {
	w, _ := new(big.Int).SetString("50000000000000000000", 10)
	f, _ := new(big.Int).SetString("34424200000000000000", 10)
	type args struct {
		amount any
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"whole number",
			args{
				amount: big.NewInt(50),
			},
			w,
			false,
		},
		{
			"floating point",
			args{
				amount: big.NewFloat(34.4242),
			},
			f,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := amountToWei(tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("amountToWei() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("amountToWei() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_amountToWei(t *testing.T) {
// 	type args struct {
// 		amount any
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *big.Int
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := amountToWei(tt.args.amount)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("amountToWei() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("amountToWei() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_parseTransferAmount(t *testing.T) {
	type args struct {
		token  string
		amount float64
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		// TODO: Add test cases.
		{
			"USDC",
			args{
				"USDC",
				149.99,
			},
			big.NewInt(149990000),
		},
		{
			"ETH",
			args{
				"ETH",
				0.8239,
			},
			big.NewInt(823900000000000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTransferAmount(tt.args.token, tt.args.amount); got != tt.want {
				t.Errorf("parseTransferAmount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTransferAmountFloat(t *testing.T) {
	type args struct {
		token  string
		amount int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
		{
			"USDC",
			args{
				"USDC",
				149990000,
			},
			149.99,
		},
		{
			"ETH",
			args{
				"ETH",
				823900000000000000,
			},
			0.8239,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTransferAmountFloat(tt.args.token, tt.args.amount); got != tt.want {
				t.Errorf("parseTransferAmountFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_amountToMwei(t *testing.T) {
	type args struct {
		amount int64
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		// TODO: Add test cases.
		{
			"non nil result",
			args{
				amount: 1000000,
			},
			big.NewInt(1000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := amountToMwei(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("amountToMwei() = %v, want %v", got, tt.want)
			}
		})
	}
}
