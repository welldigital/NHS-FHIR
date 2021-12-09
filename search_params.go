package client

import (
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
	GE Prefix = "ge"
	LE Prefix = "le"
)

// String returns the prefix as a string.
func (p Prefix) String() string {
	return string(p)
}

type DateParam struct {
	Prefix Prefix
	Value  time.Time
}

// String converts DateParam to a string. Useful for logging or as parameters to other funcs
func (d *DateParam) String() string {
	formattedTime := d.Value.Format("2006-01-02")
	return d.Prefix.String() + formattedTime
}
