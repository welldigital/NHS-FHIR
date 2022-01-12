package model

// AddressExtension Postal Address File (PAF) key associated with this address formatted as a FHIR extension.
// Empty if no PAF key for the address is known, or an object specifying the code system of the address key and the value of the address key.
type AddressExtension struct {
	URL       string                `json:"url"`
	Extension []AddressKeyExtension `json:"extension"`
}
