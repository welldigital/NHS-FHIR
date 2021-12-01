/*

Package Client is a Go client for the NHS FHIR API.

For more information about the API visit this link here:
https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#api-description__overview

Usage:

You use this library by creating a new client and calling methods on the client.

Example:

	package main

	import (
		"github.com/welldigital/NHS-FHIR"
		"fmt"
	)

	func main() {
		c, err := client.NewClient()
		if err != nil {
			panic(err)
		}

		url := "http://foo.com/bar"
		body := nil

		req, err := c.NewRequest("GET", url, body)

		if err != nil {
			fmt.Println(err)
		}

		// do something with req

		resp, err := c.Do(req)

		// handle resp etc...
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
	"strings"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

var errNonNilContext = errors.New("context must be non-nil")

const (
	defaultBaseURL = "https://sandbox.api.service.nhs.uk/personal-demographics/FHIR/R4/"
)

// NewClient returns a new FHIR client. If a nil httpClient is provided then a new http.client will be used.
// To use API methods requiring auth then provide a http.Client which will perform the authentication for you e.g. oauth2
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	// adds trailing slash
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	return &Client{
		BaseURL:    baseURL,
		httpClient: httpClient,
	}
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
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
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)

	// use the error stored in context as likely to be more informative
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}
