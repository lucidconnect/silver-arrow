package erc20

import (
	"testing"
)

func TestLoadSupportedTokens(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"load tokens.json",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadSupportedTokens(); (err != nil) != tt.wantErr {
				t.Errorf("LoadSupportedTokens() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTokenAddres(t *testing.T) {
	type args struct {
		token string
		chain int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"get usdc address on mumbai",
			args{
				"USDC",
				80001,
			},
			"0x0FA8781a83E46826621b3BC094Ea2A0212e71B23",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = LoadSupportedTokens()
			if got := GetTokenAddress(tt.args.token, tt.args.chain); got != tt.want {
				t.Errorf("GetTokenAddres() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initFromJsonString(t *testing.T) {
	type args struct {
		jsonString string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"test string",
			args{
				`[
				{
					"name": "USDC",
					"chain": 80001,
					"address": "0x0FA8781a83E46826621b3BC094Ea2A0212e71B23"
				},
				{
					"name": "USDC",
					"chain": 10,
					"address": "0x0FA8781a83E46826621b3BC094Ea2A0212e71B23"
				}
			]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initFromJsonString(tt.args.jsonString)
		})
	}
}

// func Test_updateTokenCache(t *testing.T) {
// 	type args struct {
// 		cache   map[string]map[int64]string
// 		token   string
// 		address string
// 		chain   int64
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			updateTokenCache(tt.args.cache, tt.args.token, tt.args.address, tt.args.chain)
// 		})
// 	}
// }
