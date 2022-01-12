package model

// ResourceExtension Wrapper array for the patient's pharmacies, death notification status, communication details,
// contact preferences and place of birth; these are all FHIR extensions.
type ResourceExtension struct {
	URL            string            `json:"url"`
	Extension      []FluffyExtension `json:"extension"`
	ValueReference *ValueReference   `json:"valueReference,omitempty"`
	ValueAddress   *ValueAddress     `json:"valueAddress,omitempty"`
}
