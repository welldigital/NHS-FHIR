package model

type Contact struct {
	ID           string           `json:"id"`
	Period       Period           `json:"period"`
	Relationship []Relationship   `json:"relationship"`
	Telecom      []ContactTelecom `json:"telecom"`
}
