package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	*jwt.StandardClaims
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AccessTokenRequest struct {
	GrantType           string `url:"grant_type"`
	ClientAssertionType string `url:"client_assertion_type"`
	JWT                 string `url:"client_assertion"`
}

type AuthService = service

// GenerateToken gets the access token
func (a *AuthService) GenerateAccessToken(ctx context.Context, jwt string) (*AccessTokenResponse, *Response, error) {

	// TODO: path should only be /oauth2/token
	path := "https://dev.api.service.nhs.uk/oauth2/token"

	opts := AccessTokenRequest{
		GrantType:           "client_credentials",
		ClientAssertionType: "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
		JWT:                 jwt,
	}

	url, err := addParamsToURL(path, opts)

	if err != nil {
		return nil, nil, err
	}

	req, err := a.client.newRequest(http.MethodPost, url, nil)

	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var tokenRes AccessTokenResponse

	resp, err := a.client.do(ctx, req, tokenRes)

	if err != nil {
		return nil, resp, err
	}

	fmt.Println(resp)

	return &tokenRes, resp, err
}

// TODO: handle token refresh
