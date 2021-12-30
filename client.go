/*

Package Client is a Go client for the NHS FHIR API.

For more information about the API visit this link here:
https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#api-description__overview

Usage:

You use this library by creating a new client and calling methods on the client.

Example:

	package main

	import (
		client "github.com/welldigital/nhs-fhir"
		"fmt"
		"context"
	)

	func main() {
		c, err := client.NewClient()
		if err != nil {
			panic(err)
		}

		ctx := context.Background()
		patient, _, err := c.Patient.Get(ctx, "9000000009")
	}

*/

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

// Client manages communication with the NHS FHIR API.
type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	withAuth   bool
	httpClient *http.Client

	Patient *PatientService

	accessToken AccessTokenResponse
	jwt         string
	authConfig  *AuthConfigOptions
}

//go:generate moq -out client_moq.go . IClient
// IClient interface for Client
type IClient interface {
	newRequest(method, path string, body interface{}) (*http.Request, error)
	do(ctx context.Context, req *http.Request, v interface{}) (*Response, error)
	postForm(ctx context.Context, url string, data url.Values, v interface{}) (*Response, error)
	baseURLGetter() *url.URL
}

var errNonNilContext = errors.New("context must be non-nil")

const (
	sandboxURL     = "https://sandbox.api.service.nhs.uk/"
	defaultBaseURL = sandboxURL
)

// NewClientWithOptions takes in some options to create the client with.
// If no options are given then its treated the same as NewClient(nil)
func NewClientWithOptions(opts *Options) (*Client, error) {
	c := &Client{}

	if opts == nil {
		return NewClient(nil), nil
	}

	if opts.Client != nil {
		c.httpClient = opts.Client
	} else {
		c.httpClient = newDefaultHttpClient()
	}
	c.UserAgent = opts.UserAgent

	if opts.AuthConfigOptions != nil {
		c.withAuth = true
		c.authConfig = opts.AuthConfigOptions
	}

	if opts.BaseURL != "" {
		baseURL, err := url.Parse(opts.BaseURL)
		if err != nil {
			return nil, err
		}
		c.BaseURL = baseURL
	} else {
		c.BaseURL = newDefaultBaseURL()
	}

	patientService := PatientService{client: c}
	c.Patient = &patientService

	return c, nil
}

// NewClient returns a new FHIR client. If a nil httpClient is provided then a new http.client will be used.
// To use API methods requiring auth then provide a http.Client which will perform the authentication for you e.g. oauth2
func NewClient(httpClient *http.Client) *Client {
	c := &Client{}
	if httpClient == nil {
		c.httpClient = newDefaultHttpClient()
	}

	baseURL := newDefaultBaseURL()

	c.BaseURL = baseURL

	patientService := PatientService{client: c}
	c.Patient = &patientService

	return c
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURLGetter().ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	// Every request to NHS API should contain a unique id otherwise we receive a 429
	req.Header.Set("X-Request-ID", uuid.New().String())

	// sandbox doesnt have auth
	if c.withAuth && c.baseURLGetter().String() != sandboxURL {
		bearerToken, err := c.getAccessToken(context.Background())
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClientGetter().Do(req)

	// use the error stored in context as likely to be more informative
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, &RateLimitError{}
	}

	err = json.NewDecoder(resp.Body).Decode(v)

	r := newResponse(resp)
	r.RequestID = req.Header.Get("X-Request-ID")
	return r, err
}

// getAccessToken return a valid access token
// we check if the access token is valid and if not then we generate a new one
// note: the nhs fhir api does not currently provide us with a way to get a refresh token
// instead it's advised to grab a new access token
// https://digital.nhs.uk/developer/guides-and-documentation/security-and-authorisation/application-restricted-restful-apis-signed-jwt-authentication#step-8-refresh-token
func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	if c.accessToken.AccessToken != "" && !c.accessToken.HasExpired() {
		return c.accessToken.AccessToken, nil
	}
	if c.jwt == "" {
		jwt, err := generateSecret(*c.authConfig)
		if err != nil {
			return c.accessToken.AccessToken, err
		}
		c.jwt = *jwt
	}
	token, _, err := c.generateAccessToken(ctx, c.jwt)
	if err != nil {
		return "", err
	}
	return token.AccessToken, err
}

func (c *Client) postForm(ctx context.Context, url string, data url.Values, v interface{}) (*Response, error) {
	resp, err := c.httpClientGetter().PostForm(url, data)

	// use the error stored in context as likely to be more informative
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, &RateLimitError{}
	}

	err = json.NewDecoder(resp.Body).Decode(v)

	return newResponse(resp), err

}

// GenerateToken gets the access token using a signed token
func (c *Client) generateAccessToken(ctx context.Context, jwt string) (*AccessTokenResponse, *Response, error) {

	path := "/oauth2/token"

	opts := AccessTokenRequest{
		GrantType:           "client_credentials",
		ClientAssertionType: "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
		JWT:                 jwt,
	}

	data, err := query.Values(opts)
	if err != nil {
		return nil, nil, err
	}
	tokenRes := &AccessTokenResponse{}

	res, err := c.postForm(ctx, c.authConfig.BaseURL+path, data, tokenRes)

	if err != nil {
		return nil, res, fmt.Errorf("error generating access token: %v", err)
	}
	c.accessToken = *tokenRes
	return tokenRes, res, err
}

// baseURL retrieves a baseURL, if not set then we return the default url
func (c *Client) baseURLGetter() *url.URL {
	if c.BaseURL == nil {
		return newDefaultBaseURL()
	}
	return c.BaseURL
}

// httpClientGetter provides a way to get the underlying http client
// if the client was initialized using a struct then this guarantees that the behaviour will be normal
func (c *Client) httpClientGetter() *http.Client {
	if c.httpClient == nil {
		return newDefaultHttpClient()
	}
	return c.httpClient
}
