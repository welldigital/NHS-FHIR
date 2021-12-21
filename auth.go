package client

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	ExpiresIn string `json:"expires_in"`
	// TokenType = "bearer"
	TokenType string `json:"token_type"`
	IssuedAt  string `json:"issued_at"`
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

	// SigningMethod to be used when signing/verifing tokens
	SigningMethod jwt.SigningMethod
}

func (c AuthConfigOptions) Validate() error {
	if c.Kid == "" {
		return errors.New("kid is missing but required")
	}
	if c.ClientID == "" {
		return errors.New("client id is missing but required")
	}

	if len(c.PrivateKey) == 0 && c.PrivateKeyPemFile == "" {
		return errors.New("private key or private key file must be specified")
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

func parseTokenFromSignedTokenString(tokenString string) (*jwt.Token, error) {
	publicKey, err := ioutil.ReadFile("/Users/joshdando/Documents/projects/well-digital/nhs-fhir/secret/nhs-well-dev.key.pub")
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %v\n", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("error parsing RSA public key: %v\n", err)
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	return parsedToken, nil
}
