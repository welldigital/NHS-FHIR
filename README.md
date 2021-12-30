# NHS-FHIR

[![Go Reference](https://pkg.go.dev/badge/github.com/welldigital/nhs-fhir.svg)](https://pkg.go.dev/github.com/welldigital/nhs-fhir)
[![Build](https://github.com/welldigital/NHS-FHIR/actions/workflows/test.yml/badge.svg)](https://github.com/welldigital/NHS-FHIR/actions/workflows/test.yml)
[![semantic-release: angular](https://img.shields.io/badge/semantic--release-angular-e10079?logo=semantic-release)](https://github.com/semantic-release/semantic-release)
[![report](https://goreportcard.com/badge/github.com/welldigital/nhs-fhir)](https://goreportcard.com/report/github.com/welldigital/nhs-fhir)


A Go client for [NHS FHIR API](https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#top).

This client allows you to access the PDS (Patient Demographic Service) which is the national database of NHS patient details.

You can retrieve a patients name, date of birth, address, registered GP and much more.


- [NHS-FHIR](#nhs-fhir)
	- [Installing](#installing)
	- [Getting started](#getting-started)
		- [Authentication](#authentication)
		- [Authentication with JWT](#authentication-with-jwt)
			- [Authentication with AWS KMS](#authentication-with-aws-kms)
	- [Services](#services)
		- [Patient Service](#patient-service)
	- [Roadmap](#roadmap)
	- [Contributing](#contributing)
	- [Testing](#testing)
	- [Release](#release)
		- [Pre-releases](#pre-releases)

## Installing

To install this library use: 

```
go get github.com/welldigital/nhs-fhir
```

## Getting started

You use this library by creating a new client and calling methods on the client.


```go
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

### Authentication
The easiest and recommended way to do this is using the oauth2 library, but you can always use any other library that provides a http.Client. If you have an OAuth2 access token you can use it like so:

```go
import (
	"golang.org/x/oauth2"
	client "github.com/welldigital/nhs-fhir"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "... your access token ..."},
	)
	tc := oauth2.NewClient(ctx, ts)

	c := client.NewClient(tc)

	p, resp, err := cli.Patient.Get(ctx, "9000000009")
	// ...
}
```

### Authentication with JWT

This example shows how you can gain access to the restricted API's (integration, prod) using your RSA private key.
Make sure you follow the instructions outlined on the NHS [website](https://digital.nhs.uk/developer/guides-and-documentation/security-and-authorisation/application-restricted-restful-apis-signed-jwt-authentication#step-1-create-an-application) first.

```go

import (
	"log"
	client "github.com/welldigital/nhs-fhir"
)

func main() {
	opts := &client.Options{
			AuthConfigOptions: &client.AuthConfigOptions{
				ClientID:          "your-nhs-app-id",
				Kid:               "test-1",
				BaseURL:           "https://int.api.service.nhs.uk",
				PrivateKeyPemFile: "path/to/private/key/key.pem",
			},
			BaseURL: "https://int.api.service.nhs.uk",
	}
	cli, err := client.NewClientWithOptions(opts)

	if err != nil {
		log.Fatalf("error init client: %v", err)
	}

	p, resp, err := cli.Patient.Get(ctx, "9449304424")
}

```

#### Authentication with AWS KMS

This example shows you can authenticate with the AWS Key management service (KMS) by using the KMS to sign your jwt token.
Don't forget to pass in your aws credentials!

```go

import(
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/golang-jwt/jwt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/welldigital/jwt-go-aws-kms/jwtkms"
	client "github.com/welldigital/nhs-fhir"
)

func Signer(config *jwtkms.Config) func(token *jwt.Token, key interface{}) (string, error) {
	return func(token *jwt.Token, key interface{}) (string, error) {
		return token.SignedString(config)
	}
}

func main() {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("REGION"))
		
	if err != nil {
		panic(err)
	}

	kmsClient := kms.NewFromConfig(awsCfg)

	ctx := context.Background()

	kmsConfig := jwtkms.NewKMSConfig(kmsClient, "your private key id from kms", false)

	signerFunc := Signer(kmsConfig.WithContext(ctx))

	opts := &client.Options{
		AuthConfigOptions: &client.AuthConfigOptions{
			ClientID:          "your-nhs-app-id",
			Kid:               "test-1",
			BaseURL:           "https://int.api.service.nhs.uk",
			Signer:            signerFunc,
			SigningMethod:     jwtkms.SigningMethodRS512,
		},
		BaseURL: "https://int.api.service.nhs.uk",
	}
	cli, err := client.NewClientWithOptions(opts)

	if err != nil {
		log.Fatalf("error init client: %v", err)
	}
	p, resp, err := cli.Patient.Get(ctx, "9449304424")

}

```

## Services

The client contains services which can be used to get the data you require.

### Patient Service

The patient service contains methods for getting a patient from the PDS either using their NHS number or the `PatientSearchOptions`.


## Roadmap

The following pieces of work still need to be done: 

- [Updating patient details](https://digital.nhs.uk/developer/api-catalogue/personal-demographics-service-fhir#api-Default-update-patient-partial)
- Rate limiting
- Better Error handling

## Contributing

If you wish to contribute to the project then open a Pull Request outlining what you want to do and why. 

We can then discuss how the feature might be done and then you can create a new branch from which you can develop this feature. Please add tests where appropriate and add documentation where necessary.


## Testing

Tests are written preferably in a [table driven manner](https://mj-go.in/golang/table-driven-tests-in-go) where it makes sense. 

Consider using interfaces as it makes the process of testing easier because we can control external parts of the system and only test the parts we are interested in. Read [this](https://nathanleclaire.com/blog/2015/10/10/interfaces-and-composition-for-effective-unit-testing-in-golang/) for more info.

To assist in testing we use a tool called [moq](https://github.com/matryer/moq) which generates a struct from any interface. This then allows us to mock an interface in test code.

## Release

Releases are handled automatically by [semantic-release](https://github.com/semantic-release/semantic-release) which is run whenever a commit is pushed to the branch named 'main'. This is done by the github-action found in `.github/workflows/release.yml`.

Semantic release requires a [github personal access token](https://github.com/settings/tokens) to make releases to protected branches.
Your commit messages need to be in the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) format e.g.

```
feat(feature-name): your message
```

or

```
fix: your message
```

- `feat` stands for feature. Other acceptable values are `feat`, `fix`, `chore`, `docs` and `BREAKING CHANGE`. This determines what version release is made.
- `feature-name` (optional) is the name of the feature, must be lower case with no gaps. This will be included as items in the change log.
- `your message` is the changes you've made in the commit. This will make up most of the auto generated change log.

### Pre-releases

This repo supports the use of pre-releases. Any work which will have a lot of breaking changes should be done on either an `alpha` or `beta` branch which are both marked as prerelease branches as shown in `.releaserc`. This is to avoid creating a lot of un-necessary versions.

For more information on how pre-release branches work see the [documentation](https://semantic-release.gitbook.io/semantic-release/usage/workflow-configuration#pre-release-branches).