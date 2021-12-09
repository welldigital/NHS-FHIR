package model

type GeneralPractitioner struct {
	ID         string                        `json:"id"`        
	Type       string                        `json:"type"`      
	Identifier GeneralPractitionerIdentifier `json:"identifier"`
}
