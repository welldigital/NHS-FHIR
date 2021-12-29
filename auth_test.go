package client

import (
	"reflect"
	"testing"
	"time"
)

func TestAccessTokenResponse_ExpiryTime(t *testing.T) {
	type fields struct {
		ExpiresIn int64
		IssuedAt  int64
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "get correct expiry time",
			fields: fields{
				ExpiresIn: 599,
				IssuedAt:  1640774690275,
			},
			want: time.Date(2021, 12, 29, 10, 44, 50, 874000000, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AccessTokenResponse{
				ExpiresIn: tt.fields.ExpiresIn,
				IssuedAt:  tt.fields.IssuedAt,
			}
			if got := a.ExpiryTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccessTokenResponse.ExpiryTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
