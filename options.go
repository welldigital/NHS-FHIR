package client

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Options to configure the client with
type Options struct {
	*http.Client
	*AuthConfigOptions
	BaseURL   string
	UserAgent string
	*TracingOptions
}

type TracingOptions struct {
	// Enabled set to true to enable ALL tracing
	Enabled bool
	// TraceErrorsOnly set to true to limit tracing to errors only
	TraceErrorsOnly bool
	// Output allows you to configure where the log output is outputted to. Defaults to os.Stdout
	Output io.Writer
}

func newDefaultBaseURL() *url.URL {
	baseURL, _ := url.Parse(defaultBaseURL)

	// adds trailing slash
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}
	return baseURL
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
