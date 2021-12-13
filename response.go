package client

import "net/http"

// Response is for all API responses, it contains the http response.
type Response struct {
	*http.Response
	// RequestID contains a string which is used to uniquely identify the request
	// Used for debugging or support
	RequestID string
}

func newResponse(r *http.Response) *Response {
	return &Response{
		Response: r,
	}
}
