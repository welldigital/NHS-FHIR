package model

// Period Business effective period when name was, is, or will be in use.
type Period struct {
	// Start yyyy-mm-dd
	Start string `json:"start"`
	// End yyyy-mm-dd
	End   string `json:"end"`
}
