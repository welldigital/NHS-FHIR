package model

type IdentifierElement struct {
	System    string                `json:"system"`   
	Value     string                `json:"value"`    
	Extension []IdentifierExtension `json:"extension"`
}
