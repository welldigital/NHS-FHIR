package client

import (
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// addParamsToURL is a helper func which adds parameters to a url.
// The parameters are added in alphabetical order, see example below.
/* Example:
type PeopleOpts struct {
	Foo string `url:"foo"`
	Bar string `url:"bar"`
}
opts := PeopleOpts{"abc", "kazoo"}
out, _ := addParamsToURL("/people", opts)

fmt.Println(out) // "/people?bar=kazoo&foo=abc"

*/
func addParamsToURL(urlString string, params interface{}) (string, error) {
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return urlString, nil
	}
	u, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}
	vs, err := query.Values(params)
	if err != nil {
		return urlString, err
	}
	u.RawQuery = vs.Encode()
	return u.String(), nil
}
