package model

// IdentifierElement Identifier and system of identification used for this Patient.
type IdentifierElement struct {
	System    string                `json:"system"`
	Value     string                `json:"value"`
	Extension []IdentifierExtension `json:"extension"`
}
