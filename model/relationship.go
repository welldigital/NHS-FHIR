package model

// Relationship The contact relationship wrapper object that holds the details of the relationship to the patient.
type Relationship struct {
	Coding []Security `json:"coding"`
}
