package model

// GeneralPractitioner General Practice (not practitioner) with which the patient is, or was, registered.
// Always contains zero or one general practitioner object.
// When a patient tagged as restricted or very restricted is retrieved,
// the General Practice is removed from the response.
type GeneralPractitioner struct {
	ID         string                        `json:"id"`
	Type       string                        `json:"type"`
	Identifier GeneralPractitionerIdentifier `json:"identifier"`
}
