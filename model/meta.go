package model

type Meta struct {
	VersionID string     `json:"versionId"`
	Security  []Security `json:"security"` 
}
