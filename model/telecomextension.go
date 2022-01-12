package model

// TelecomExtension extension for telecom
type TelecomExtension struct {
	URL         string   `json:"url"`
	ValueCoding Security `json:"valueCoding"`
}
