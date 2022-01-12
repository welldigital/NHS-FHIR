// To parse and unparse this JSON data, add this code to your project and do:
//
//    result, err := UnmarshalResult(bytes)
//    bytes, err = result.Marshal()

package model

import "encoding/json"

// UnmarshalResult unmarshals the data into Result object
func UnmarshalResult(data []byte) (Result, error) {
	var r Result
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal marshals a Result object into json
func (r *Result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Result wrapper for results that come back from NHS FHIR API.
type Result struct {
	ResourceType string  `json:"resourceType"`
	Type         string  `json:"type"`
	Timestamp    string  `json:"timestamp"`
	Total        int64   `json:"total"`
	Entry        []Entry `json:"entry"`
}
