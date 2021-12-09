// To parse and unparse this JSON data, add this code to your project and do:
//
//    result, err := UnmarshalResult(bytes)
//    bytes, err = result.Marshal()

package model

import "encoding/json"

func UnmarshalResult(data []byte) (Result, error) {
	var r Result
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Result struct {
	ResourceType string  `json:"resourceType"`
	Type         string  `json:"type"`
	Timestamp    string  `json:"timestamp"`
	Total        int64   `json:"total"`
	Entry        []Entry `json:"entry"`
}
