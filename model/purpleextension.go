package model

// AddressKeyExtension Specification of address key system and address key value.
// Contains exactly two items: one describing the code system the Address Key uses,
// and the other specifying the value of the Address Key.
type AddressKeyExtension struct {
	URL         string       `json:"url"`
	ValueCoding *ValueCoding `json:"valueCoding,omitempty"`
	ValueString *string      `json:"valueString,omitempty"`
}
