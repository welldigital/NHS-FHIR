package model

type Security struct {
	System  string  `json:"system"`           
	Code    string  `json:"code"`             
	Display string  `json:"display"`          
	Version *string `json:"version,omitempty"`
}
