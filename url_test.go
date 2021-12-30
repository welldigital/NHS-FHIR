package client

import (
	"errors"
	"testing"
)

func TestIsAbsoluteUrl(t *testing.T) {

	tests := []struct {
		name string
		url  string
		err  error
	}{
		{
			name: "missing scheme",
			url:  "google.com",
			err:  ErrUrlSchemeMissing,
		},
		{
			name: "missing host",
			url:  "http://",
			err:  ErrUrlHostMissing,
		},
		{
			name: "invalid url",
			url:  ":::/not.valid/a//a??a?b=&&c#hi",
			err:  errors.New("parse \":::/not.valid/a//a??a?b=&&c\": missing protocol scheme"),
		},
		{
			name: "valid absolute url",
			url:  "https://https://int.api.service.nhs.uk",
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAbsoluteUrl(tt.url)
			if got == nil || tt.err == nil {
				if got != tt.err {
					t.Errorf("got = %v, want = %v", got, tt.err)
				}
			} else {

				if got.Error() != tt.err.Error() {
					t.Errorf("IsAbsoluteUrl() = %v, want %v", got, tt.err)
				}
			}

		})
	}
}
