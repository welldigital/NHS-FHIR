package client

import (
	"errors"
	"net/url"
)

// ErrURLSchemeMissing error for when url scheme is missing e.g. https://
var ErrURLSchemeMissing = errors.New("url scheme is empty")
// ErrURLHostMissing error for when url host is missing e.g. google.com
var ErrURLHostMissing = errors.New("url host is empty")

// IsAbsoluteURL returns an error if the given string is not a valid absolute url
func IsAbsoluteURL(str string) error {
	// url.Parse allows relative urls and also accepts empty "http://"
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		return ErrURLSchemeMissing
	}
	if u.Host == "" {
		return ErrURLHostMissing
	}
	return nil
}
