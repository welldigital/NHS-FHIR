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
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

// Client manages communication with the NHS FHIR API.
type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	withAuth   bool
	httpClient *http.Client

	Patient *PatientService
}

//go:generate moq -out client_moq.go . IClient
// IClient interface for Client
type IClient interface {
	newRequest(method, path string, body interface{}) (*http.Request, error)
	do(ctx context.Context, req *http.Request, v interface{}) (*Response, error)
}

var errNonNilContext = errors.New("context must be non-nil")

const (
	defaultBaseURL = "https://sandbox.api.service.nhs.uk/personal-demographics/FHIR/R4/"
)

// NewClientWithOptions takes in some options to create the client with.
// If no options are given then its treated the same as NewClient with nil passed as the http client.
func NewClientWithOptions(opts ...ClientOptions) *Client {
	c := &Client{}

	if len(opts) == 0 {
		c = NewClient(nil)
	} else {
		for _, opt := range opts {
			// applies the options on the client
			opt(c)
		}
	}

	WithServices()(c)

	return c
}

// NewClient returns a new FHIR client. If a nil httpClient is provided then a new http.client will be used.
// To use API methods requiring auth then provide a http.Client which will perform the authentication for you e.g. oauth2
func NewClient(httpClient *http.Client) *Client {
	c := &Client{}
	if httpClient == nil {
		WithDefaultHttpClient()(c)
	}

	baseURL := newDefaultBaseURL()

	WithBaseUrl(baseURL)(c)
	WithServices()(c)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL().ResolveReference(rel)
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

	// TODO: how do we identify that the request should be made this way?
	if c.withAuth {
		// TODO: add token in header
		bearerToken := ""
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

// baseURL provides a way to guarantee retrieving a baseURL
func (c *Client) baseURL() *url.URL {
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
