package client

import "net/http"

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	return &Response{
		Response: r,
	}
}
