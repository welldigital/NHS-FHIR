package model

type TelecomExtension struct {
	URL         string   `json:"url"`        
	ValueCoding Security `json:"valueCoding"`
}
