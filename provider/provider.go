package provider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"golang.org/x/oauth2"
)

// IProvider interface for a resource provider for authentication
type IProvider interface {
	BeginAuth(state string) (Session, error)
	RefreshToken(refreshToken string) (*oauth2.Token, error) //Get new access token based on the refresh token
}

// TODO: fill these out properly
const (
	authURL  string = "authorize"
	tokenURL string = "token"
)

type Provider struct {
	ClientKey    string
	ClientSecret string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	state        string
}

func newConfig(p *Provider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     p.ClientKey,
		ClientSecret: p.ClientSecret,
		RedirectURL:  p.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}
}

func New(clientKey, secret, callbackURL string) *Provider {
	p := &Provider{
		CallbackURL:  callbackURL,
		ClientKey:    clientKey,
		ClientSecret: secret,
	}
	p.config = newConfig(p)

	return p
}

func (p *Provider) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return p.config.Client(ctx, token)
}

// CreateState creates a state token used to prevent CSRF attacks
func (p *Provider) CreateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

//BeginAuth calls the underlying oauth2 AuthCodeURL()
// State is a token to protect the user from CSRF attacks. You must always provide a non-empty string and validate that
// it matches the the state query parameter on your redirect callback.
// See http://tools.ietf.org/html/rfc6749#section-10.12 for more info.
func (p *Provider) BeginAuth(state string) (*Session, error) {
	return &Session{
		AuthURL: p.config.AuthCodeURL(state),
	}, nil
}

//RefreshToken get new access token based on the refresh token
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	token := &oauth2.Token{RefreshToken: refreshToken}
	// TODO: ctx
	ctx := context.Background()
	ts := p.config.TokenSource(ctx, token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}
	return newToken, err
}
