package model

type PurpleExtension struct {
	URL         string       `json:"url"`
	ValueCoding *ValueCoding `json:"valueCoding,omitempty"`
	ValueString *string      `json:"valueString,omitempty"`
}
