package model

// ValueCoding URL of specification of wrapper FHIR extension
type ValueCoding struct {
	System string `json:"system"`
	Code   string `json:"code"`
}
