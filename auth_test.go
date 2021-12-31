package client

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestAccessTokenResponse_ExpiryTime(t *testing.T) {
	type fields struct {
		ExpiresIn int64
		IssuedAt  int64
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "get correct expiry time",
			fields: fields{
				ExpiresIn: 599,
				IssuedAt:  1640774690275,
			},
			want: time.Date(2021, 12, 29, 10, 54, 49, 275*1000000, time.Local),
		},
		{
			name: "test 2",
			fields: fields{
				ExpiresIn: 3600,          // 1hr in seconds
				IssuedAt:  1640945411000, // milliseconds
			},
			want: time.Date(2021, 12, 31, 11, 10, 11, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccessTokenResponse{
				ExpiresIn: tt.fields.ExpiresIn,
				IssuedAt:  tt.fields.IssuedAt,
			}
			if got := a.ExpiryTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccessTokenResponse.ExpiryTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthConfigOptions_Validate(t *testing.T) {
	type fields struct {
		BaseURL           string
		ClientID          string
		Kid               string
		PrivateKeyPemFile string
		PrivateKey        []byte
		Signer            SigningFunc
		SigningMethod     *jwt.SigningMethodRSA
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name: "invalid kid field",
			fields: fields{
				BaseURL:           "https://https://int.api.service.nhs.uk",
				ClientID:          "123",
				PrivateKeyPemFile: "abc.pem",
				PrivateKey:        []byte{},
			},
			err: ErrKidMissing,
		},
		{
			name: "missing client id",
			fields: fields{
				BaseURL:           "https://https://int.api.service.nhs.uk",
				Kid:               "test",
				PrivateKeyPemFile: "abc.pem",
				PrivateKey:        []byte{},
			},
			err: ErrClientIdMissing,
		},
		{
			name: "missing key/file with no signing func",
			fields: fields{
				BaseURL:  "https://https://int.api.service.nhs.uk",
				ClientID: "123",
				Kid:      "test",
			},
			err: ErrKeyMissing,
		},
		{
			name: "base url cant be empty",
			fields: fields{
				BaseURL:           "",
				ClientID:          "234",
				Kid:               "test",
				PrivateKeyPemFile: "file.pem",
			},
			err: ErrBaseUrlMissing,
		},
		{
			name: "base url must be a valid url",
			fields: fields{
				BaseURL: "http://",
			},
			err: ErrUrlHostMissing,
		},
		{
			name: "invalid signing method algorithm",
			fields: fields{
				BaseURL:  "https://google.com",
				ClientID: "1234",
				Kid:      "test",
				Signer: func(token *jwt.Token, key interface{}) (string, error) {
					return "", nil
				},
				SigningMethod: (*jwt.SigningMethodRSA)(jwt.SigningMethodHS256),
			},
			err: ErrInvalidSigningMethodAlg,
		},
		{
			name: "valid signing method",
			fields: fields{
				BaseURL:  "https://https://int.api.service.nhs.uk",
				ClientID: "1",
				Kid:      "1",
				Signer: func(token *jwt.Token, key interface{}) (string, error) {
					return "", nil
				},
				SigningMethod: jwt.SigningMethodRS256,
			},
			err: nil,
		},
		{
			name: "valid auth config",
			fields: fields{
				BaseURL:           "https://https://int.api.service.nhs.uk",
				ClientID:          "123",
				Kid:               "test",
				PrivateKeyPemFile: "file.pem",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := AuthConfigOptions{
				BaseURL:           tt.fields.BaseURL,
				ClientID:          tt.fields.ClientID,
				Kid:               tt.fields.Kid,
				PrivateKeyPemFile: tt.fields.PrivateKeyPemFile,
				PrivateKey:        tt.fields.PrivateKey,
				Signer:            tt.fields.Signer,
				SigningMethod:     tt.fields.SigningMethod,
			}
			if err := c.Validate(); !errors.Is(err, tt.err) {
				t.Errorf("AuthConfigOptions.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
