package model

type ResourceExtension struct {
	URL       string            `json:"url"`      
	Extension []FluffyExtension `json:"extension"`
}
