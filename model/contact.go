package model

// Contact A list of patient contacts. Only emergency contacts are returned and only emergency contacts should be added.
// Any other contacts should be added to the patients Related Person.
type Contact struct {
	ID           string           `json:"id"`
	Period       Period           `json:"period"`
	Relationship []Relationship   `json:"relationship"`
	Telecom      []ContactTelecom `json:"telecom"`
}
