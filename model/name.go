package model

type Name struct {
	ID     string   `json:"id"`
	Use    string   `json:"use"`
	Period Period   `json:"period"`
	Given  []string `json:"given"`
	Family string   `json:"family"`
	Prefix []string `json:"prefix"`
	Suffix []string `json:"suffix"`
}
