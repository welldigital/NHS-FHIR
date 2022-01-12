package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/welldigital/nhs-fhir/model"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func createString(s string) *string {
	return &s
}

func createBool(b bool) *bool {
	return &b
}

func TestPatientService_Get(t *testing.T) {

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		p       *PatientService
		args    args
		want    *model.Patient
		wantErr bool
	}{
		{
			wantErr: true,
			name:    "adds authentication headers",
			args: args{
				ctx: context.Background(),
				id:  "9000000009",
			},
			p: &service{
				client: &Client{
					withAuth: true,
					BaseURL: &url.URL{
						Scheme: "https",
						Host:   "test.com",
					},
					accessToken: AccessTokenResponse{
						AccessToken: "token",
						ExpiresIn:   600,
						TokenType:   "bearer",
						IssuedAt:    time.Now().UnixNano() / int64(time.Millisecond),
					},
					jwt: "jwt",
					httpClient: &http.Client{
						Timeout: 1 * time.Millisecond,
						Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
							if r.URL.Path == path+"/9000000009" {
								_, keyInMap := r.Header["Authorization"]
								if !keyInMap {
									t.Error("expected to find Authorization header in request but didnt: ", r.Header)
								}
								return &http.Response{StatusCode: 401}, nil
							}
							return &http.Response{}, nil
						}),
					},
					authConfig: &AuthConfigOptions{
						BaseURL:  "https://test.com",
						ClientID: "123",
						Kid:      "1",
						Signer: func(token *jwt.Token, key interface{}) (string, error) {
							return "", nil
						},
					},
				},
			},
		},

		{
			name: "invalid nhs number",
			p: &service{
				client: &IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						return &Response{}, nil
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						return &http.Request{}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "123",
			},
			wantErr: true,
		},
		{
			name: "bad request",
			p: &service{
				client: &IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						return &Response{}, nil
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						return &http.Request{}, errors.New("bang")
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "9000000009",
			},
			wantErr: true,
		},
		{
			name: "bad response",
			p: &service{
				&IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						return &Response{}, errors.New("fail")
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						return &http.Request{}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "9000000009",
			},
			wantErr: true,
		},
		{
			name: "gets a dummy patient from sandbox",
			p: &service{
				client: &IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						patient := `{
							"resourceType": "Patient",
							"id": "9000000009",
							"identifier": [
								{
									"system": "https://fhir.nhs.uk/Id/nhs-number",
									"value": "9000000009",
									"extension": [
										{
											"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NHSNumberVerificationStatus",
											"valueCodeableConcept": {
												"coding": [
													{
														"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-NHSNumberVerificationStatus",
														"version": "1.0.0",
														"code": "01",
														"display": "Number present and verified"
													}
												]
											}
										}
									]
								}
							],
							"meta": {
								"versionId": "2",
								"security": [
									{
										"system": "http://terminology.hl7.org/CodeSystem/v3-Confidentiality",
										"code": "U",
										"display": "unrestricted"
									}
								]
							},
							"name": [
								{
									"id": "123",
									"use": "usual",
									"period": {
										"start": "2020-01-01",
										"end": "2021-12-31"
									},
									"given": [
										"Jane"
									],
									"family": "Smith",
									"prefix": [
										"Mrs"
									],
									"suffix": [
										"MBE"
									]
								}
							],
							"gender": "female",
							"birthDate": "2010-10-22",
							"multipleBirthInteger": 1,
							"deceasedDateTime": "2010-10-22T00:00:00+00:00",
							"generalPractitioner": [
								{
									"id": "254406A3",
									"type": "Organization",
									"identifier": {
										"system": "https://fhir.nhs.uk/Id/ods-organization-code",
										"value": "Y12345",
										"period": {
											"start": "2020-01-01",
											"end": "2021-12-31"
										}
									}
								}
							],
							"extension": [
								{
									"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NominatedPharmacy",
									"valueReference": {
										"identifier": {
											"system": "https://fhir.nhs.uk/Id/ods-organization-code",
											"value": "Y12345"
										}
									}
								},
								{
									"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-PreferredDispenserOrganization",
									"valueReference": {
										"identifier": {
											"system": "https://fhir.nhs.uk/Id/ods-organization-code",
											"value": "Y23456"
										}
									}
								},
								{
									"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-MedicalApplianceSupplier",
									"valueReference": {
										"identifier": {
											"system": "https://fhir.nhs.uk/Id/ods-organization-code",
											"value": "Y34567"
										}
									}
								},
								{
									"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-DeathNotificationStatus",
									"extension": [
										{
											"url": "deathNotificationStatus",
											"valueCodeableConcept": {
												"coding": [
													{
														"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-DeathNotificationStatus",
														"version": "1.0.0",
														"code": "2",
														"display": "Formal - death notice received from Registrar of Deaths"
													}
												]
											}
										},
										{
											"url": "systemEffectiveDate",
											"valueDateTime": "2010-10-22T00:00:00+00:00"
										}
									]
								},
								{
									"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NHSCommunication",
									"extension": [
										{
											"url": "language",
											"valueCodeableConcept": {
												"coding": [
													{
														"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-HumanLanguage",
														"version": "1.0.0",
														"code": "fr",
														"display": "French"
													}
												]
											}
										},
										{
											"url": "interpreterRequired",
											"valueBoolean": true
										}
									]
								},
								{
									"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-ContactPreference",
									"extension": [
										{
											"url": "PreferredWrittenCommunicationFormat",
											"valueCodeableConcept": {
												"coding": [
													{
														"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-PreferredWrittenCommunicationFormat",
														"code": "12",
														"display": "Braille"
													}
												]
											}
										},
										{
											"url": "PreferredContactMethod",
											"valueCodeableConcept": {
												"coding": [
													{
														"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-PreferredContactMethod",
														"code": "1",
														"display": "Letter"
													}
												]
											}
										},
										{
											"url": "PreferredContactTimes",
											"valueString": "Not after 7pm"
										}
									]
								},
								{
									"url": "http://hl7.org/fhir/StructureDefinition/patient-birthPlace",
									"valueAddress": {
										"city": "Manchester",
										"district": "Greater Manchester",
										"country": "GBR"
									}
								}
							],
							"telecom": [
								{
									"id": "789",
									"period": {
										"start": "2020-01-01",
										"end": "2021-12-31"
									},
									"system": "phone",
									"value": "01632960587",
									"use": "home"
								},
								{
									"id": "OC789",
									"period": {
										"start": "2020-01-01",
										"end": "2021-12-31"
									},
									"system": "other",
									"value": "01632960587",
									"use": "home",
									"extension": [
										{
											"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-OtherContactSystem",
											"valueCoding": {
												"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-OtherContactSystem",
												"code": "textphone",
												"display": "Minicom (Textphone)"
											}
										}
									]
								}
							],
							"contact": [
								{
									"id": "C123",
									"period": {
										"start": "2020-01-01",
										"end": "2021-12-31"
									},
									"relationship": [
										{
											"coding": [
												{
													"system": "http://terminology.hl7.org/CodeSystem/v2-0131",
													"code": "C",
													"display": "Emergency Contact"
												}
											]
										}
									],
									"telecom": [
										{
											"system": "phone",
											"value": "01632960587"
										}
									]
								}
							],
							"address": [
								{
									"id": "456",
									"period": {
										"start": "2020-01-01",
										"end": "2021-12-31"
									},
									"use": "home",
									"line": [
										"1 Trevelyan Square",
										"Boar Lane",
										"City Centre",
										"Leeds",
										"West Yorkshire"
									],
									"postalCode": "LS1 6AE",
									"extension": [
										{
											"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-AddressKey",
											"extension": [
												{
													"url": "type",
													"valueCoding": {
														"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-AddressKeyType",
														"code": "PAF"
													}
												},
												{
													"url": "value",
													"valueString": "12345678"
												}
											]
										}
									]
								},
								{
									"id": "T456",
									"period": {
										"start": "2020-01-01",
										"end": "2021-12-31"
									},
									"use": "temp",
									"text": "Student Accommodation",
									"line": [
										"1 Trevelyan Square",
										"Boar Lane",
										"City Centre",
										"Leeds",
										"West Yorkshire"
									],
									"postalCode": "LS1 6AE",
									"extension": [
										{
											"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-AddressKey",
											"extension": [
												{
													"url": "type",
													"valueCoding": {
														"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-AddressKeyType",
														"code": "PAF"
													}
												},
												{
													"url": "value",
													"valueString": "12345678"
												}
											]
										}
									]
								}
							]
						}`
						r := ioutil.NopCloser(bytes.NewReader([]byte(patient)))
						err := json.NewDecoder(r).Decode(v)
						return newResponse(&http.Response{Status: "200", Body: r}), err
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						assert.Equal(t, path, "personal-demographics/FHIR/R4/Patient/2983396339")
						url, err := url.Parse(path)
						if err != nil {
							t.Errorf("parsing url caused an unexpected err: %v", err)
						}
						return &http.Request{Method: http.MethodGet, URL: url}, err
					},
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "2983396339",
			},
			want: &model.Patient{
				ResourceType: "Patient",
				ID:           "9000000009",
				Identifier: []model.IdentifierElement{
					{
						System: "https://fhir.nhs.uk/Id/nhs-number",
						Value:  "9000000009",
						Extension: []model.IdentifierExtension{
							{
								URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NHSNumberVerificationStatus",
								ValueCodeableConcept: model.Relationship{
									Coding: []model.Security{
										{
											System:  "https://fhir.hl7.org.uk/CodeSystem/UKCore-NHSNumberVerificationStatus",
											Code:    "01",
											Display: "Number present and verified",
											Version: createString("1.0.0"),
										},
									},
								},
							},
						},
					},
				},
				Meta: model.Meta{
					VersionID: "2",
					Security: []model.Security{
						{
							System:  "http://terminology.hl7.org/CodeSystem/v3-Confidentiality",
							Code:    "U",
							Display: "unrestricted",
						},
					},
				},
				Name: []model.Name{
					{
						ID:  "123",
						Use: "usual",
						Period: model.Period{
							Start: "2020-01-01",
							End:   "2021-12-31",
						},
						Given:  []string{"Jane"},
						Family: "Smith",
						Prefix: []string{"Mrs"},
						Suffix: []string{"MBE"},
					},
				},
				Gender:               "female",
				BirthDate:            "2010-10-22",
				MultipleBirthInteger: 1,
				DeceasedDateTime:     "2010-10-22T00:00:00+00:00",
				GeneralPractitioner: []model.GeneralPractitioner{
					{
						ID:   "254406A3",
						Type: "Organization",
						Identifier: model.GeneralPractitionerIdentifier{
							Period: model.Period{
								Start: "2020-01-01",
								End:   "2021-12-31",
							},
							System: "https://fhir.nhs.uk/Id/ods-organization-code",
							Value:  "Y12345",
						},
					},
				},
				Extension: []model.ResourceExtension{
					{
						URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NominatedPharmacy",
						ValueReference: &model.ValueReference{
							Identifier: model.IdentifierElement{
								System: "https://fhir.nhs.uk/Id/ods-organization-code",
								Value:  "Y12345",
							},
						},
					},
					{
						URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-PreferredDispenserOrganization",
						ValueReference: &model.ValueReference{
							Identifier: model.IdentifierElement{
								System: "https://fhir.nhs.uk/Id/ods-organization-code",
								Value:  "Y23456",
							},
						},
					},
					{
						URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-MedicalApplianceSupplier",
						ValueReference: &model.ValueReference{
							Identifier: model.IdentifierElement{
								System: "https://fhir.nhs.uk/Id/ods-organization-code",
								Value:  "Y34567",
							},
						},
					},
					{
						URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-DeathNotificationStatus",
						Extension: []model.FluffyExtension{
							{
								URL: "deathNotificationStatus",
								ValueCodeableConcept: &model.Relationship{
									Coding: []model.Security{
										{
											System:  "https://fhir.hl7.org.uk/CodeSystem/UKCore-DeathNotificationStatus",
											Code:    "2",
											Display: "Formal - death notice received from Registrar of Deaths",
											Version: createString("1.0.0"),
										},
									},
								},
							},
							{
								URL:           "systemEffectiveDate",
								ValueDateTime: createString("2010-10-22T00:00:00+00:00"),
							},
						},
					},
					{
						URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NHSCommunication",
						Extension: []model.FluffyExtension{
							{
								URL: "language",
								ValueCodeableConcept: &model.Relationship{
									Coding: []model.Security{
										{
											System:  "https://fhir.hl7.org.uk/CodeSystem/UKCore-HumanLanguage",
											Code:    "fr",
											Display: "French",
											Version: createString("1.0.0"),
										},
									},
								},
							},
							{
								URL:          "interpreterRequired",
								ValueBoolean: createBool(true),
							},
						},
					},
					{
						URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-ContactPreference",
						Extension: []model.FluffyExtension{
							{
								URL: "PreferredWrittenCommunicationFormat",
								ValueCodeableConcept: &model.Relationship{
									Coding: []model.Security{
										{
											System:  "https://fhir.hl7.org.uk/CodeSystem/UKCore-PreferredWrittenCommunicationFormat",
											Code:    "12",
											Display: "Braille",
										},
									},
								},
							},
							{
								URL: "PreferredContactMethod",
								ValueCodeableConcept: &model.Relationship{
									Coding: []model.Security{
										{
											System:  "https://fhir.hl7.org.uk/CodeSystem/UKCore-PreferredContactMethod",
											Code:    "1",
											Display: "Letter",
										},
									},
								},
							},
							{
								URL:         "PreferredContactTimes",
								ValueString: createString("Not after 7pm"),
							},
						},
					},
					{
						URL: "http://hl7.org/fhir/StructureDefinition/patient-birthPlace",
						ValueAddress: &model.ValueAddress{
							City:     "Manchester",
							District: "Greater Manchester",
							Country:  "GBR",
						},
					},
				},
				Address: []model.Address{
					{
						ID: "456",
						Period: model.Period{
							Start: "2020-01-01",
							End:   "2021-12-31",
						},
						Use: "home",
						Line: []string{"1 Trevelyan Square",
							"Boar Lane",
							"City Centre",
							"Leeds",
							"West Yorkshire"},
						PostalCode: "LS1 6AE",
						Extension: []model.AddressExtension{
							{
								URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-AddressKey",
								Extension: []model.AddressKeyExtension{
									{
										URL: "type",
										ValueCoding: &model.ValueCoding{
											System: "https://fhir.hl7.org.uk/CodeSystem/UKCore-AddressKeyType",
											Code:   "PAF",
										},
									},
									{
										URL:         "value",
										ValueString: createString("12345678"),
									},
								},
							},
						},
					},
					{
						ID:   "T456",
						Text: createString("Student Accommodation"),
						Period: model.Period{
							Start: "2020-01-01",
							End:   "2021-12-31",
						},
						Use: "temp",
						Line: []string{"1 Trevelyan Square",
							"Boar Lane",
							"City Centre",
							"Leeds",
							"West Yorkshire"},
						PostalCode: "LS1 6AE",
						Extension: []model.AddressExtension{
							{
								URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-AddressKey",
								Extension: []model.AddressKeyExtension{
									{
										URL: "type",
										ValueCoding: &model.ValueCoding{
											System: "https://fhir.hl7.org.uk/CodeSystem/UKCore-AddressKeyType",
											Code:   "PAF",
										},
									},
									{
										URL:         "value",
										ValueString: createString("12345678"),
									},
								},
							},
						},
					},
				},
				Telecom: []model.ResourceTelecom{
					{
						ID: "789",
						Period: model.Period{
							Start: "2020-01-01",
							End:   "2021-12-31",
						},
						System: "phone",
						Value:  "01632960587",
						Use:    "home",
					},
					{
						ID: "OC789",
						Period: model.Period{
							Start: "2020-01-01",
							End:   "2021-12-31",
						},
						System: "other",
						Value:  "01632960587",
						Use:    "home",
						Extension: []model.TelecomExtension{
							{
								URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-OtherContactSystem",
								ValueCoding: model.Security{
									System:  "https://fhir.hl7.org.uk/CodeSystem/UKCore-OtherContactSystem",
									Code:    "textphone",
									Display: "Minicom (Textphone)",
								},
							},
						},
					},
				},
				Contact: []model.Contact{
					{
						ID: "C123",
						Period: model.Period{
							Start: "2020-01-01",
							End:   "2021-12-31",
						},
						Relationship: []model.Relationship{
							{
								Coding: []model.Security{
									{
										System:  "http://terminology.hl7.org/CodeSystem/v2-0131",
										Code:    "C",
										Display: "Emergency Contact",
									},
								},
							},
						},
						Telecom: []model.ContactTelecom{
							{
								System: "phone",
								Value:  "01632960587",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, _, err := tt.p.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatientService.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatientService.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatientService_Search(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts PatientSearchOptions
	}
	tests := []struct {
		name    string
		p       *PatientService
		args    args
		want    []*model.Patient
		wantErr bool
	}{
		{
			name: "user not found",
			p: &service{
				&IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						return &Response{}, nil
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						return &http.Request{}, nil
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				opts: PatientSearchOptions{},
			},
			want:    []*model.Patient{},
			wantErr: false,
		},
		{
			name: "bad request",
			p: &service{
				&IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						return &Response{}, nil
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						return &http.Request{}, errors.New("bad request")
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				opts: PatientSearchOptions{},
			},
			wantErr: true,
		},
		{
			name: "bad response",
			p: &service{
				&IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						return &Response{}, errors.New("bad response")
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						return &http.Request{}, nil
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				opts: PatientSearchOptions{},
			},
			wantErr: true,
		},
		{
			name: "finds a patient",
			p: &service{
				&IClientMock{
					doFunc: func(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
						results := `{
							"resourceType": "Bundle",
							"type": "searchset",
							"timestamp": "2021-12-08T14:29:56+00:00",
							"total": 1,
							"entry": [
								{
									"fullUrl": "https://api.service.nhs.uk/personal-demographics/FHIR/R4/Patient/9000000017",
									"search": {
										"score": 0.8976
									},
									"resource": {
										"resourceType": "Patient",
										"id": "9000000017",
										"identifier": [
											{
												"system": "https://fhir.nhs.uk/Id/nhs-number",
												"value": "9000000017",
												"extension": [
													{
														"url": "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NHSNumberVerificationStatus",
														"valueCodeableConcept": {
															"coding": [
																{
																	"system": "https://fhir.hl7.org.uk/CodeSystem/UKCore-NHSNumberVerificationStatus",
																	"version": "1.0.0",
																	"code": "01",
																	"display": "Number present and verified"
																}
															]
														}
													}
												]
											}
										],
										"meta": {
											"versionId": "2",
											"security": [
												{
													"system": "http://terminology.hl7.org/CodeSystem/v3-Confidentiality",
													"code": "U",
													"display": "unrestricted"
												}
											]
										}
									}
										
								}
							]
						}`
						r := ioutil.NopCloser(bytes.NewReader([]byte(results)))
						err := json.NewDecoder(r).Decode(v)
						return newResponse(&http.Response{Status: "200", Body: r}), err
					},
					newRequestFunc: func(method, path string, body interface{}) (*http.Request, error) {
						assert.Equal(t, http.MethodGet, method)
						assert.Equal(t, "personal-demographics/FHIR/R4/Patient?_fuzzy-match=true&_max-results=1&address-postcode=M123&birthdate=lt2021-01-01&birthdate=ge2020-10-02&given=Smith", path)
						return &http.Request{}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				opts: PatientSearchOptions{
					MaxResults: 1,
					FuzzyMatch: createBool(true),
					Given:      &[]string{"Smith"},
					BirthDate: []*string{
						createString("lt2021-01-01"),
						createString("ge2020-10-02"),
					},
					Postcode: createString("M123"),
				},
			},
			wantErr: false,
			want: []*model.Patient{
				{
					ResourceType: "Patient",
					ID:           "9000000017",
					Identifier: []model.IdentifierElement{
						{
							System: "https://fhir.nhs.uk/Id/nhs-number",
							Value:  "9000000017",
							Extension: []model.IdentifierExtension{
								{
									URL: "https://fhir.hl7.org.uk/StructureDefinition/Extension-UKCore-NHSNumberVerificationStatus",
									ValueCodeableConcept: model.Relationship{
										Coding: []model.Security{
											{
												System:  "https://fhir.hl7.org.uk/CodeSystem/UKCore-NHSNumberVerificationStatus",
												Code:    "01",
												Display: "Number present and verified",
												Version: createString("1.0.0"),
											},
										},
									},
								},
							},
						},
					},
					Meta: model.Meta{
						VersionID: "2",
						Security: []model.Security{
							{
								System:  "http://terminology.hl7.org/CodeSystem/v3-Confidentiality",
								Code:    "U",
								Display: "unrestricted",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := tt.p.Search(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatientService.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatientService.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
