package model

// Search Every matched patient in the result list includes a score to indicate how closely the patient
// matched the search parameters.
type Search struct {
	// Score = 1 if exact match, otherwise partial match.
	Score float64 `json:"score"`
}
