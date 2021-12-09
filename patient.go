package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Joshswooft/nhs/cmd/validation"
	"github.com/welldigital/NHS-FHIR/model"
)

type PatientService = service

const path = "Patient"

// Get gets a patient from the PDS using the patients NHS number as the id.
// id = The patient's NHS number. The primary identifier of a patient, unique within NHS England and Wales. Always 10 digits and must be a valid NHS number.
func (p *PatientService) Get(ctx context.Context, id string) (*model.Patient, error) {
	err := validation.NhsNumberValidator(id)
	if err != nil {
		return nil, err
	}
	req, err := p.client.NewRequest(http.MethodGet, fmt.Sprintf(path+"/%v", id), nil)

	if err != nil {
		return nil, err
	}

	patient := &model.Patient{}
	_, err = p.client.Do(ctx, req, patient)

	if err != nil {
		return nil, err
	}

	return patient, nil
}

// PatientSearchOptions is the options we pass into the request for searching a patient
type PatientSearchOptions struct {
	// A fuzzy search is performed, including checks for homophones, transposed names and historic information.
	// You cant use wildcards with fuzzy search
	FuzzyMatch *bool `url:"_fuzzy-match,omitempty"`
	// The search only returns results where the score field is 1.0. Use this with care - it is unlikely to work with fuzzy search or wildcards.
	ExactMatch *bool `url:"_exact-match,omitempty"`
	// The search looks for matches in historic information such as previous names and addresses.
	// This parameter has no effect for a fuzzy search, which always includes historic information.
	History *bool `url:"_history,omitempty"`
	// For application-restricted access, this must be 1
	MaxResults int `url:"_max-results"`
	// if used with wildcards, fuzzy match must be false. Wildcards must contain at least two characters, this matches Smith, Smythe. Not case-sensitive.
	Family *string `url:"family,omitempty"`
	// The patients given name, can be used with wildcards. E.g. Jane Anne Smith
	// Use * as a wildcard but not in the first two characters and not in fuzzy search mode
	Given  *[]string `url:"given,omitempty"`
	Gender *Gender   `url:"gender,omitempty"`
	// Format: <eq|ge|le>yyyy-mm-dd e.g. eq2021-08-01
	BirthDate []*string `url:"birthdate,omitempty"`
	// For a fuzzy search, this is ignored for matching but included in the score calculation.
	// Format: <eq|ge|le>yyyy-mm-dd e.g. eq2021-08-01
	DeathDate *[]string `url:"death-date,omitempty"`
	// Not case sensitive. Spaces are ignored, for example LS16AE and LS1 6AE both match LS1 6AE
	Postcode *string `url:"address-postcode,omitempty"`
	// The Organisation Data Service (ODS) code of the patient's registered GP practice.
	// Not case sensitive. For a fuzzy search, this is ignored for matching but included in the score calculation.
	// Example: Y12345
	GeneralPractioner *string `url:"general-practitioner,omitempty"`
}

// Search searches for a patient in the PDS
// The behaviour of this endpoint depends on your access mode:
//https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#api-Default-search-patient
func (p *PatientService) Search(ctx context.Context, opts PatientSearchOptions) ([]*model.Patient, error) {
	url, err := addParamsToUrl(path, opts)

	if err != nil {
		return nil, err
	}

	req, err := p.client.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	result := &model.Result{}

	_, err = p.client.Do(ctx, req, result)

	if err != nil {
		return nil, err
	}

	if len(result.Entry) == 0 {
		return nil, errors.New("user not found")
	}

	patients := make([]*model.Patient, len(result.Entry))

	for i, entry := range result.Entry {
		patients[i] = &entry.Resource
	}

	return patients, nil

}
