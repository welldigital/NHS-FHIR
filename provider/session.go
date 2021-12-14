package provider

import (
	"context"
	"errors"
	"net/url"

	"golang.org/x/oauth2"
)

// ISession provides the mechanism for authorizing a provider
// The session will be marshaled and persisted between requests to "tie"
// the start and the end of the authorization process with a
// 3rd party provider.
type ISession interface {
	// GetAuthURL returns the URL for the authentication end-point for the provider.
	GetAuthURL() (string, error)
	// Authorize should validate the data from the provider and return an access token
	// that can be stored for later access to the provider.
	Authorize(Provider, params url.Values) (string, error)
}

type Session struct {
	AuthURL string
	Token   *oauth2.Token
}

var errNonNilAuthUrl = errors.New("no auth url provided")

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errNonNilAuthUrl
	}
	return s.AuthURL, nil
}

// Authorize - uses the authURL to exchange for a token
func (s Session) Authorize(p Provider, params url.Values) (*oauth2.Token, error) {
	code := params.Get("code")

	state := params.Get("state")

	if state != p.state {
		// potential CSRF attack
		return nil, errors.New("redirect url state didnt match the providers state param")
	}

	// TODO: ctx we can get from another oauth2 method?
	ctx := context.Background()
	token, err := p.config.Exchange(ctx, code)

	if err != nil {
		return nil, err
	}

	if !token.Valid() {
		return nil, errors.New("invalid token received from provider")
	}

	s.Token = token

	return token, err
}
