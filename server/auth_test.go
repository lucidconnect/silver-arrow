package server

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_parseJwt(t *testing.T) {
	godotenv.Load()
	type args struct {
		jwToken string
	}
	tests := []struct {
		name string
		args args
		// want    jwt.MapClaims
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "jwt test",
			args: args{
				jwToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJtZXJjaGFudC1pZCI6IjY5MTg1NTE0LTA2NjQtNDdhZi1iOTJjLTVlMzg1NmE2Y2ExMyIsInByb2R1Y3QtaWQiOiJhZmFiYmIxOC1kNzAzLTQ2MTAtYjFmYi04ZjRhMThlNzY4YmUiLCJpYXQiOjE3MDMzMjk4Nzh9.s3xNOhQALZ8Cx0t4bkpDlxv2dymVEJ5J430w5xHVxeU",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseJwt(tt.args.jwToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("parseJwt() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestPaymentlinkJwt(t *testing.T) {
	godotenv.Load()
	jwt, err := generatePaymentLinkJwt("afabbb18-d703-4610-b1fb-8f4a18e768be", "69185514-0664-47af-b92c-5e3856a6ca13")
	assert.NoError(t, err)
	fmt.Println(jwt)
	_, err = parseJwt(jwt)
	assert.NoError(t, err)
}
