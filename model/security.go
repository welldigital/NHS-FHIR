package model

// Security The level of security on the patients record, which affects which fields are populated on retrieval
type Security struct {
	// System URI of the value set specification
	System  string  `json:"system"`
	// Code can be the following: 'U' = unrestricted,
	//'R' = restricted (any sensitive data e.g. address, gp, telecom is removed)
	// 'V' = very restricted (all patient data removed except id, identifier, meta)
	Code    string  `json:"code"`
	// Display 'unrestricted', 'restricted', 'very restricted', 'redacted'
	Display string  `json:"display"`
	Version *string `json:"version,omitempty"`
}
