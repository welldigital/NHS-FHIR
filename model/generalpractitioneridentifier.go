package model

// GeneralPractitionerIdentifier information identifying the GP
type GeneralPractitionerIdentifier struct {
	System string `json:"system"`
	Value  string `json:"value"`
	Period Period `json:"period"`
}
