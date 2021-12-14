package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/welldigital/nhs-fhir/provider"
)

// ClientOptions options to configure the client with
type ClientOptions func(*Client)

// WithOauth returns a ClientOption that sets up a client to use with Auth
func WithOauth(clientKey, secret, callbackURL string) ClientOptions {
	return func(c *Client) {
		p := provider.New(clientKey, secret, callbackURL)
		state := p.CreateState()
		sess, err := p.BeginAuth(state)

		if err != nil {
			// TODO: handle err - dont panic
			panic(err)
		}

		// we cant do this bit as it's done from the users browser / they sign in.
		//  not needed - wrong method of auth. We want JWT.
		var params url.Values
		token, err := sess.Authorize(*p, params)

		if err != nil {
			// TODO: handle err - dont panic
			panic(err)
		}

		client := p.Client(context.Background(), token)

		c.httpClient = client
	}
}

// WithHTTPClient returns a ClientOption that specifies the http client
// used to perform requests with. Fallbacks to defaultHttpClient if none given
func WithHTTPClient(hc *http.Client) ClientOptions {
	return func(c *Client) {
		if hc == nil {
			c.httpClient = newDefaultHttpClient()
		} else {
			c.httpClient = hc
		}
	}
}

func WithJWTAuth(clientId string) ClientOptions {
	return func(c *Client) {
		c.withAuth = true

		claims := jwt.StandardClaims{
			Audience:  "https://dev.api.service.nhs.uk/oauth2/token",
			Id:        uuid.NewString(),
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			Issuer:    clientId,
			Subject:   clientId,
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

		secretKey, err := ioutil.ReadFile("secret/nhs-well-dev.key")

		if err != nil {
			log.Fatal(err)
		}

		tokenSigned, err := jwtToken.SignedString(secretKey)

		if err != nil {
			log.Fatal(err)
		}

		service := AuthService{
			client: c,
		}

		res, _, err := service.GenerateAccessToken(context.Background(), tokenSigned)

		// what to do with token?

		fmt.Println(res)

	}
}

// WithUserAgent returns a ClientOption that specifies the user agent
// used to tell the server what device you are identifying as
func WithUserAgent(u string) ClientOptions {
	return func(c *Client) {
		c.UserAgent = u
	}
}

// WithServices attaches the services to the client
func WithServices() ClientOptions {
	return func(c *Client) {
		patientService := PatientService{client: c}
		c.Patient = &patientService
	}
}

func newDefaultBaseURL() *url.URL {
	baseURL, _ := url.Parse(defaultBaseURL)

	// adds trailing slash
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}
	return baseURL
}

// WithBaseUrl returns a ClientOption that specifies the url for which
// to send requests to. Defaults to the defaultBaseUrl if nil
func WithBaseUrl(url *url.URL) ClientOptions {
	return func(c *Client) {
		if url == nil {
			c.BaseURL = newDefaultBaseURL()
		} else {
			c.BaseURL = url
		}
	}
}

func newDefaultHttpClient() *http.Client {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

func WithDefaultHttpClient() ClientOptions {
	return func(c *Client) {
		c.httpClient = newDefaultHttpClient()
	}
}
