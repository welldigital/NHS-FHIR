package client

import "testing"

func Test_addParamsToUrl(t *testing.T) {
	type args struct {
		urlString string
		params    interface{}
	}

	type SearchParams struct {
		Count      int    `url:"count"`
		Page       int    `url:"page"`
		Name       string `url:"name"`
		ExactMatch bool   `url:"exact_match,omitempty"`
	}

	type PeopleOpts struct {
		Foo string `url:"foo"`
		Bar string `url:"bar"`
	}

	type ArrayExample struct {
		People []string `url:"friends"`
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "example test",
			args: args{
				urlString: "/people",
				params:    PeopleOpts{"abc", "kazoo"},
			},
			want: "/people?bar=kazoo&foo=abc",
		},
		{
			name: "params with some omitted fields",
			args: args{
				urlString: "/foo",
				params: SearchParams{
					Count: 10,
					Page:  0,
					Name:  "Bob",
				},
			},
			want: "/foo?count=10&name=Bob&page=0",
		},
		{
			name: "nil interface",
			args: args{
				urlString: "/foo",
				params:    nil,
			},
			want: "/foo",
		},
		{
			name: "bad url",
			args: args{
				urlString: "<% 00 7F #>",
				params:    nil,
			},
			wantErr: true,
			want:    "<% 00 7F #>",
		},
		{
			name: "handles arrays",
			args: args{
				urlString: "/foo",
				params:    ArrayExample{People: []string{"Mohammed", "Luke", "Jennifer"}},
			},
			want: "/foo?friends=Mohammed&friends=Luke&friends=Jennifer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addParamsToUrl(tt.args.urlString, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("addParamsToUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("addParamsToUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
