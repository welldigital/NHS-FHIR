package model

// ResourceTelecom List of contact points for the patient; for example, phone numbers or email addresses.
// When a patient tagged as restricted or very restricted is retrieved,
// all contact points are removed from the response.
type ResourceTelecom struct {
	ID        string             `json:"id"`
	Period    Period             `json:"period"`
	System    string             `json:"system"`
	Value     string             `json:"value"`
	Use       string             `json:"use"`
	Extension []TelecomExtension `json:"extension,omitempty"`
}
