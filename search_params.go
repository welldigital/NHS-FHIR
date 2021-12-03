package client

import (
	"net/url"
	"strings"
	"time"
)

type Gender string

const (
	Male    Gender = "male"
	Female  Gender = "female"
	Other   Gender = "other"
	Unknown Gender = "unknown"
)

// String returns gender as a string
func (g Gender) String() string {
	return string(g)
}

// Prefix is an enum representing FHIR parameter prefixes.  The following
// description is from the FHIR DSTU2 specification:
//
// For the ordered parameter types number, date, and quantity, a prefix
// to the parameter value may be used to control the nature of the matching.
type Prefix string

// Constant values for the Prefix enum.
const (
	EQ Prefix = "eq"
	NE Prefix = "ne"
	GT Prefix = "gt"
	LT Prefix = "lt"
	GE Prefix = "ge"
	LE Prefix = "le"
	SA Prefix = "sa"
	EB Prefix = "eb"
	AP Prefix = "ap"
)

// String returns the prefix as a string.
func (p Prefix) String() string {
	return string(p)
}

type DateParam struct {
	Prefix Prefix
	Value  time.Time
}

// EncodeValues
/*
	Example 1
	// where Params is a struct containing parameters for a request
	type p Params struct {
		Date DateParam `url:death-date`
	}

	d := DateParam{ "eq", "2010-10-22"}
	p := Params{Date: d}
	v, _ := query.Values(p)
	fmt.Print(v.Encode()) // will output: "death-date=eq2010-10-22"

*/
func (d *DateParam) EncodeValues(key string, v *url.Values) error {
	formattedTime := d.Value.Format("2006-01-02")
	val := d.Prefix.String() + formattedTime
	v.Set(strings.ToLower(key), val)
	return nil
}
