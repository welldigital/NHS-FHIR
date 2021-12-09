package model

type AddressExtension struct {
	URL       string            `json:"url"`      
	Extension []PurpleExtension `json:"extension"`
}
