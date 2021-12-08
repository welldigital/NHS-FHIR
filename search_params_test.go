package client

import (
	"testing"
	"time"

	"github.com/google/go-querystring/query"
)

func TestDateParam_EncodeValues(t *testing.T) {

	type SearchParams struct {
		Date *DateParam `url:"birth-date,omitempty"`
	}
	type fields struct {
		Prefix Prefix
		Value  time.Time
	}

	date := time.Date(2020, 07, 01, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		fields        fields
		expQueryParam string
		wantErr       bool
	}{
		{
			name: "encode date param as url values",
			fields: fields{
				Prefix: EQ,
				Value:  date,
			},
			expQueryParam: "birth-date=eq2020-07-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DateParam{
				Prefix: tt.fields.Prefix,
				Value:  tt.fields.Value,
			}

			params := SearchParams{Date: d}

			v, err := query.Values(params)
			if err != nil {
				t.Errorf("expected nil err but got %v", err)
			}
			if got := v.Encode(); got != tt.expQueryParam {
				t.Errorf("query params dont match! got = %v, does not equal want: %v", got, tt.expQueryParam)
			}
		})
	}
}

func TestDateParam_String(t *testing.T) {
	type fields struct {
		Prefix Prefix
		Value  time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "to string",
			fields: fields{
				Prefix: GE,
				Value:  time.Date(2005, time.August, 05, 0, 0, 0, 0, time.UTC),
			},
			want: "ge2005-08-05",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DateParam{
				Prefix: tt.fields.Prefix,
				Value:  tt.fields.Value,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("DateParam.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
