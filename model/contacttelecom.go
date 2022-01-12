package model

// ContactTelecom List of contact points for the patient; for example, phone numbers or email addresses.
// When a patient tagged as restricted or very restricted is retrieved,
// all contact points are removed from the response
type ContactTelecom struct {
	System string `json:"system"`
	Value  string `json:"value"`
}
