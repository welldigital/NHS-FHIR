package model

// IdentifierExtension generic identifier extension for FHIR
type IdentifierExtension struct {
	URL                  string       `json:"url"`
	ValueCodeableConcept Relationship `json:"valueCodeableConcept"`
}
