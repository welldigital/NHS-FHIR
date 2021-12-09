package model

type Entry struct {
	FullURL string `json:"fullUrl"`
	Search  Search `json:"search"`
	// In reality resource is type interface{} and would require unmarshalling the type based on ResourceType
	// But we can simply for now as NHS API only provides patient details
	Resource Patient `json:"resource"`
}
