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
