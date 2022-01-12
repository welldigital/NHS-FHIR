package model

// Meta meta data for this resource
type Meta struct {
	// Version The NHS Digital assigned version of the patient resource.
	VersionID string     `json:"versionId"`
	Security  []Security `json:"security"`
}
