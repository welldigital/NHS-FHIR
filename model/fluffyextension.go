package model

// FluffyExtension FHIR extensions for resources in wrapper array
type FluffyExtension struct {
	URL                  string        `json:"url"`
	ValueCodeableConcept *Relationship `json:"valueCodeableConcept,omitempty"`
	ValueDateTime        *string       `json:"valueDateTime,omitempty"`
	ValueBoolean         *bool         `json:"valueBoolean,omitempty"`
	ValueString          *string       `json:"valueString,omitempty"`
}
