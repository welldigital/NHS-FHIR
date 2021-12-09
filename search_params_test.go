package client

import (
	"testing"
	"time"
)

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
