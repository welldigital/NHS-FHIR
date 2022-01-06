package model

type IdentifierExtension struct {
	URL                  string       `json:"url"`
	ValueCodeableConcept Relationship `json:"valueCodeableConcept"`
}
