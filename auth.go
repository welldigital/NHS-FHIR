package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	*jwt.StandardClaims
}

type AccessTokenResponse struct {
	// AccessToken used to call NHS restricted API's
	AccessToken string `json:"access_token"`
	// ExpiresIn the time in seconds that the token will expire in.
	ExpiresIn int64 `json:"expires_in,string"`
	// TokenType = "bearer"
	TokenType string `json:"token_type"`
	// timestamp of when the token was issued in milliseconds
	IssuedAt int64 `json:"issued_at,string"`
}

func (a AccessTokenResponse) ExpiryTime() time.Time {

	return time.Unix(0, (a.IssuedAt+a.ExpiresIn)*int64(time.Millisecond))
}

func (a AccessTokenResponse) HasExpired() bool {
	now := time.Now()
	return now.After(a.ExpiryTime()) || now.Equal(a.ExpiryTime())
}

type AccessTokenRequest struct {
	GrantType           string `url:"grant_type"`
	ClientAssertionType string `url:"client_assertion_type"`
	JWT                 string `url:"client_assertion"`
}

// SigningFunc used to sign your own JWT tokens
// returns a signed token
type SigningFunc = func(token *jwt.Token, key interface{}) (string, error)

// AuthConfigOptions the options used for JWT Auth
type AuthConfigOptions struct {
	// BaseURL the url for auth
	BaseURL string

	// ClientID the api key of your nhs application
	ClientID string
	// Kid is a header used as a key identifier to identify which key to look up to sign a particular token
	// When used with a JWK, the kid value is used to match a JWK kid parameter value.
	Kid string
	// PrivateKeyPemFile file location to your private RSA key
	PrivateKeyPemFile string

	// PrivateKey the value of your private key
	PrivateKey []byte

	// Signer is a function you can use to sign your own tokens
	Signer SigningFunc

	// SigningMethod to be used when signing/verifing tokens, must be RSA
	SigningMethod jwt.SigningMethod
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

var ErrBaseUrlMissing = errors.New("auth base url is missing but required")
var ErrKidMissing = errors.New("kid is missing but required")
var ErrClientIdMissing = errors.New("client id is missing but required")
var ErrKeyMissing = errors.New("private key or private key file must be specified")
var ErrInvalidSigningMethodAlg = errors.New("signing method must be RSA")

func (c AuthConfigOptions) Validate() error {
	if c.BaseURL == "" {
		return ErrBaseUrlMissing
	}
	if err := IsAbsoluteUrl(c.BaseURL); err != nil {
		return err
	}

	if c.Kid == "" {
		return ErrKidMissing
	}
	if c.ClientID == "" {
		return ErrClientIdMissing
	}

	if len(c.PrivateKey) == 0 && c.PrivateKeyPemFile == "" && c.Signer == nil {
		return ErrKeyMissing
	}

	if !isNil(c.SigningMethod) && !strings.Contains(c.SigningMethod.Alg(), "RS") {
		return ErrInvalidSigningMethodAlg
	}

	return nil
}

func generateSecret(config AuthConfigOptions) (*string, error) {

	err := config.Validate()

	if err != nil {
		return nil, err
	}

	claims := jwt.StandardClaims{
		Audience:  config.BaseURL + "/oauth2/token",
		Id:        uuid.NewString(),
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Issuer:    config.ClientID,
		Subject:   config.ClientID,
	}

	var jwtToken *jwt.Token

	if config.SigningMethod != nil {
		jwtToken = jwt.NewWithClaims(config.SigningMethod, claims)
	} else {
		jwtToken = jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	}
	jwtToken.Header["kid"] = config.Kid

	secretKey := config.PrivateKey

	if config.PrivateKeyPemFile != "" {
		secretKey, err = ioutil.ReadFile(config.PrivateKeyPemFile)

		if err != nil {
			return nil, err
		}
	}

	var tokenSigned string

	if config.Signer == nil {
		key, err := jwt.ParseRSAPrivateKeyFromPEM(secretKey)
		if err != nil {
			return nil, fmt.Errorf("error parsing RSA private key: %v", err)
		}

		tokenSigned, err = jwtToken.SignedString(key)
		if err != nil {
			return nil, fmt.Errorf("error signing jwt with key: %v", err)
		}

	} else {
		// use custom signer to get signed jwt token
		tokenSigned, err = config.Signer(jwtToken, nil)
		if err != nil {
			return nil, fmt.Errorf("error signing jwt with key using custom signer: %v", err)
		}
	}

	if err != nil {
		return nil, err
	}
	return &tokenSigned, nil
}
