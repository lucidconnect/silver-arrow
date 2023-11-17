package erc20

import (
	"testing"
)

func TestGetChainExplorer(t *testing.T) {
	type args struct {
		chain int64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"fetch Polygon Mumbai explorer",
			args{
				chain: 80001,
			},
			"https://mumbai.polygonscan.com/",
			false,
		},
		{
			"fail",
			args{
				chain: 1,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadSupportedChains("../tokens/chains.json")
			got, err := GetChainExplorer(tt.args.chain)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChainExplorer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetChainExplorer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadSupportedChains(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"load chains.json",
			args{
				path: "../tokens/chains.json",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadSupportedChains(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("LoadSupportedChains() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
