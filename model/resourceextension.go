package model

type ResourceExtension struct {
	URL            string            `json:"url"`
	Extension      []FluffyExtension `json:"extension"`
	ValueReference *ValueReference   `json:"valueReference,omitempty"`
	ValueAddress   *ValueAddress     `json:"valueAddress,omitempty"`
}
