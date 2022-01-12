package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.httpClient == c2.httpClient {
		t.Error("NewClient returned same http.Clients, but they should differ")
	}

	if c.Patient == nil {
		t.Errorf("NewClient() returned nil for patient service")
	}
}

func TestDo(t *testing.T) {

	ctx := context.Background()
	reqWithNoBody := httptest.NewRequest("GET", "/foo", nil)
	// not allowed to set RequestURI
	reqWithNoBody.RequestURI = ""

	type Result struct {
		Foo string
	}

	expected := Result{Foo: "foo"}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		b, err := json.Marshal(expected)
		if err != nil {
			t.Errorf("expected err to be nil got %v", err)
		}
		w.Write(b)
	}))
	defer svr.Close()

	c := NewClient(svr.Client())

	url, err := url.Parse(svr.URL)

	c.BaseURL = url

	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	reqWithNoBody.URL = url

	reqWithNoBody.Header.Set("Content-Type", "application/json")
	reqWithNoBody.Header.Set("Accept", "application/json")
	reqWithNoBody.Header.Set("X-Request-ID", "1")

	result := new(Result)
	resp, err := c.do(ctx, reqWithNoBody, result)

	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	if *result != expected {
		t.Errorf("expected res to be %s got %s", expected, result)
	}

	// response should have the request ID
	if len(resp.RequestID) == 0 {
		t.Errorf("expected res to contain a request ID but got none")
	}

}

// Test that an error caused by the internal http client's do() function
// does not leak the client secret.
func TestDo_sanitizeURL(t *testing.T) {

	type TestClient struct {
		ClientID     string
		ClientSecret string
		*http.Client
	}

	tp := &TestClient{
		ClientID:     "id",
		ClientSecret: "secret",
	}
	unauthedClient := NewClient(tp.Client)
	unauthedClient.BaseURL = &url.URL{Scheme: "http", Host: "127.0.0.1:0", Path: "/"} // Use port 0 on purpose to trigger a dial TCP error, expect to get "dial tcp 127.0.0.1:0: connect: can't assign requested address".
	req, err := unauthedClient.newRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("newRequest returned unexpected error: %v", err)
	}
	ctx := context.Background()
	_, err = unauthedClient.do(ctx, req, nil)
	if err == nil {
		t.Fatal("Expected error to be returned.")
	}
	if strings.Contains(err.Error(), "client_secret=secret") {
		t.Errorf("do error contains secret, should be redacted:\n%q", err)
	}
}

func TestNewRequest(t *testing.T) {

	type TestBody struct {
		Foo string
	}

	c := NewClient(nil)

	inURL, outURL := "foo", defaultBaseURL+"foo"
	inBody, outBody := &TestBody{Foo: "bar"}, `{"Foo":"bar"}`+"\n"
	req, _ := c.newRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("newRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("newRequest(%q) Body is %v, want %v", inBody, got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("newRequest() User-Agent is %v, want %v", got, want)
	}

	// test that each request contains a unique guid
	req2, _ := c.newRequest("GET", inURL, inBody)

	if id1, id2 := req.Header.Get("X-Request-ID"), req2.Header.Get("X-Request-ID"); id1 == id2 {
		t.Errorf("NewRequest() X-Request-ID ")
	}
}

func TestNewClientWithOptions(t *testing.T) {

	urlStr := "http://hello-world.com"
	u, _ := url.Parse(urlStr)

	type args struct {
		opts *Options
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "no options is same as NewClient()",
			args: args{
				opts: nil,
			},
			want:    NewClient(nil),
			wantErr: false,
		},
		{
			name: "sets the default base url to sandbox if not given",
			want: &Client{
				BaseURL:    newDefaultBaseURL(),
				httpClient: newDefaultHTTPClient(),
				withAuth:   false,
			},
			wantErr: false,
		},
		{
			name: "sets the base url",
			args: args{
				opts: &Options{
					BaseURL: urlStr,
				},
			},
			want: &Client{
				BaseURL: u,
			},
		},
		{
			name: "sets the auth config",
			args: args{
				opts: &Options{
					AuthConfigOptions: &AuthConfigOptions{
						BaseURL:           urlStr,
						ClientID:          "123",
						Kid:               "test",
						PrivateKeyPemFile: "file.pem",
					},
				},
			},
			want: &Client{
				withAuth: true,
				authConfig: &AuthConfigOptions{
					BaseURL:           urlStr,
					ClientID:          "123",
					Kid:               "test",
					PrivateKeyPemFile: "file.pem",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientWithOptions(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && !isClientEqual(got, tt.want) {
				t.Errorf("NewClientWithOptions() = %v, want %v", got, tt.want)
			}

		})
	}

}

func isClientEqual(got *Client, want *Client) bool {
	if got.accessToken != want.accessToken {
		fmt.Println("access token not equal")
		return false
	}

	if got.jwt != want.jwt {
		fmt.Println("jwt not equal")
		return false
	}

	if got.withAuth != want.withAuth {
		fmt.Println("withAuth not equal")
		return false
	}
	if !reflect.DeepEqual(&got.authConfig, &want.authConfig) {
		fmt.Println("auth config not equal")
		return false
	}

	if !isAuthConfigEqual(got.authConfig, want.authConfig) {
		return false
	}

	return true
}

func isAuthConfigEqual(got *AuthConfigOptions, want *AuthConfigOptions) bool {

	if got == nil && want == nil {
		return true
	}

	if got == nil && want != nil || got != nil && want == nil {
		return false
	}

	if got.BaseURL != want.BaseURL {
		return false
	}

	if got.ClientID != want.ClientID {
		return false
	}

	if got.Kid != want.Kid {
		return false
	}

	if !bytes.Equal(got.PrivateKey, want.PrivateKey) {
		return false
	}

	if got.PrivateKeyPemFile != want.PrivateKeyPemFile {
		return false
	}

	return true

}

func Test_redactFieldFromHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/foo", nil)
	field := "Auth"
	req.Header.Set(field, "my-secret")

	redactFieldFromHeader(&req.Header, field)

	res := req.Header.Get(field)
	assert.Equal(t, redactedString, res)

}

func TestClient_dumpHTTP(t *testing.T) {

	type args struct {
		req  *http.Request
		resp *http.Response
	}
	tests := []struct {
		name    string
		opts    Options
		args    args
		wantErr bool
	}{
		{
			name: "trace",
			opts: Options{
				TracingOptions: &TracingOptions{
					Enabled:         true,
					TraceErrorsOnly: false,
				},
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://google.com", nil),
				resp: &http.Response{
					StatusCode: http.StatusBadGateway,
				},
			},
			wantErr: false,
		},
		{
			name: "strips out Authorization header",
			opts: Options{
				TracingOptions: &TracingOptions{
					Enabled:         true,
					TraceErrorsOnly: false,
				},
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "https://google.com", nil),
				resp: &http.Response{
					StatusCode: http.StatusBadRequest,
				},
			},
			wantErr: false,
		},
	}
	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if index == 1 {
				tt.args.req.Header.Add("Authorization", "my-secret")
			}

			client, err := NewClientWithOptions(&tt.opts)
			var b bytes.Buffer
			client.tracingConfig.Output = &b
			if err != nil {
				t.Errorf("couldnt init client: %v", err)
			}
			if err := client.dumpHTTP(tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Client.dumpHTTP() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := b.String()

			_, found := tt.args.req.Header["Authorization"]

			if found && !strings.Contains(got, redactedString) {
				t.Errorf("auth header found but not redacted")
			}

			if !strings.Contains(got, fmt.Sprint(tt.args.resp.StatusCode)) {
				t.Errorf("trace didnt output response: %v", got)
			}
			if !strings.Contains(got, tt.args.req.Host) {
				t.Errorf("trace didnt output request: %v", got)
			}
		})
	}
}
