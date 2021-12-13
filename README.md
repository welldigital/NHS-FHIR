# NHS-FHIR
A Go client for [NHS FHIR API](https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#top).

This client allows you to access the PDS (Patient Demographic Service) which is the national database of NHS patient details.

You can retrieve a patients name, date of birth, address, registered GP and much more.

Table of contents
------------------
- [Installing](installing)
- [Getting started](#getting-started)
- [Services](#services)
  - [Patient Service](#patient-service)
- [Roadmap](#roadmap)
- [Contributing](#contributing)

## Installing

To install this library use: 

```
go get github.com/welldigital/nhs-fhir
```

## Getting started

You use this library by creating a new client and calling methods on the client.


```
package main

import (
	"context"
	"fmt"

	client "github.com/welldigital/nhs-fhir"
)

func main() {
	cli := client.NewClient(nil)
	ctx := context.Background()
	p, resp, err := cli.Patient.Get(ctx, "9000000009")

	if err != nil {
		panic(err)
	}

	fmt.Println(p)

	fmt.Println(resp)
}

```

## Services
### Patient Service

The patient service contains methods for getting a patient from the PDS either using their NHS number or the `PatientSearchOptions`.


## Roadmap

The following pieces of work still need to be done: 

- [Authentication](https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#api-description__security-and-authorisation) 
- [Updating patient details](https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#api-Default-update-patient-partial)
- Rate limiting
- Better Error handling
- Handling refresh tokens

## Contributing

If you wish to contribute to the project then open a Pull Request outlining what you want to do and why. 

We can then discuss how the feature might be done and then you can create a new branch from which you can develop this feature. Please add tests where appropriate and add documentation where necessary.

