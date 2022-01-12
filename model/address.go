package model

// Address address details for a patient.
// only the home address is returned on a search.
// When a patient tagged as restricted or very restricted is retrieved, all addresses are removed from the response.
type Address struct {
	ID         string             `json:"id"`
	Period     Period             `json:"period"`
	Use        string             `json:"use"`
	Line       []string           `json:"line"`
	PostalCode string             `json:"postalCode"`
	Extension  []AddressExtension `json:"extension"`
	Text       *string            `json:"text,omitempty"`
}
