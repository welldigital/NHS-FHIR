package model

type Address struct {
	ID         string             `json:"id"`        
	Period     Period             `json:"period"`    
	Use        string             `json:"use"`       
	Line       []string           `json:"line"`      
	PostalCode string             `json:"postalCode"`
	Extension  []AddressExtension `json:"extension"` 
}
