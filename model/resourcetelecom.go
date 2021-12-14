package model

type ResourceTelecom struct {
	ID        string             `json:"id"`
	Period    Period             `json:"period"`
	System    string             `json:"system"`
	Value     string             `json:"value"`
	Use       string             `json:"use"`
	Extension []TelecomExtension `json:"extension,omitempty"`
}
