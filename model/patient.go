package model

// Patient patient resource
type Patient struct {
	ResourceType         string                `json:"resourceType"`
	ID                   string                `json:"id"`
	Identifier           []IdentifierElement   `json:"identifier"`
	Meta                 Meta                  `json:"meta"`
	Name                 []Name                `json:"name"`
	Gender               string                `json:"gender"`
	BirthDate            string                `json:"birthDate"`
	MultipleBirthInteger int64                 `json:"multipleBirthInteger"`
	DeceasedDateTime     string                `json:"deceasedDateTime"`
	Address              []Address             `json:"address"`
	Telecom              []ResourceTelecom     `json:"telecom"`
	Contact              []Contact             `json:"contact"`
	GeneralPractitioner  []GeneralPractitioner `json:"generalPractitioner"`
	Extension            []ResourceExtension   `json:"extension"`
}
