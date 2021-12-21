// +build integration

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	client "github.com/welldigital/nhs-fhir"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func TestPatientGet(t *testing.T) {

	tests := []struct {
		name          string
		opts          *client.Options
		id            string
		expStatusCode int
		expGender     string
		expError      error
	}{
		{
			name:          "default sandbox client with nil opts",
			opts:          nil,
			id:            "9000000009",
			expStatusCode: 200,
			expGender:     client.Female.String(),
		},
		{
			name: "integration client with no auth",
			opts: &client.Options{
				BaseURL: "https://int.api.service.nhs.uk",
			},
			id:            "9449304424",
			expStatusCode: 401,
		},
		{
			name: "custom http client",
			opts: &client.Options{
				Client: &http.Client{
					Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
						assert.Equal(t, "https://sandbox.api.service.nhs.uk/personal-demographics/FHIR/R4/Patient/9000000009", r.URL.String())
						return http.DefaultTransport.RoundTrip(r)
					}),
				},
			},
			id:            "9000000009",
			expStatusCode: 200,
			expGender:     client.Female.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c, err := client.NewClientWithOptions(tt.opts)
			if err != nil {
				t.Errorf("client should not error")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

			p, res, err := c.Patient.Get(ctx, tt.id)

			assert.Equal(t, tt.expError, err)
			assert.Equal(t, tt.expStatusCode, res.StatusCode)

			if err != nil {
				if p == nil {
					t.Errorf("NewClient returned nil for patient GET")
				} else {
					fmt.Println("p: ", p)

					assert.Equal(t, tt.expGender, p.Gender)
				}
			}

			defer cancel()

		})
	}

}
