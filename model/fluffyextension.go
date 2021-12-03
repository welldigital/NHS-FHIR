package model

type FluffyExtension struct {
	URL                  string        `json:"url"`                           
	ValueCodeableConcept *Relationship `json:"valueCodeableConcept,omitempty"`
	ValueDateTime        *string       `json:"valueDateTime,omitempty"`       
}
