package client

import (
	"errors"
	"net/url"
)

var ErrUrlSchemeMissing error = errors.New("url scheme is empty")

var ErrUrlHostMissing error = errors.New("url host is empty")

func IsAbsoluteUrl(str string) error {
	// url.Parse allows relative urls and also accepts empty "http://"
	u, err := url.Parse(str)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		return ErrUrlSchemeMissing
	}
	if u.Host == "" {
		return ErrUrlHostMissing
	}
	return nil
}
