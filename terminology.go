package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

type TerminologyService interface {
	Ping() error
	Capability() (*fhir.Bundle, error)
	// ValueSet(name string) (interface{}, error)
	// CodeSystem(name string) (interface{}, error)
	// ConceptMap(name string) (interface{}, error)
}

type LOINCTerminologyService struct {
	Name    string
	Version string
	BaseUrl string
	Client  *http.Client
	Auth    map[string]string
}

func NewLOINCTerminologyService(username, pass string) (*LOINCTerminologyService, error) {
	if username == "" || pass == "" {
		return nil, errors.New("username or password not set: check .env file")
	}

	l := &LOINCTerminologyService{
		Name:    "LOINC",
		Version: "2.74",
		BaseUrl: "https://fhir.loinc.org",
		Client:  &http.Client{},
		Auth: map[string]string{
			"username": username,
			"password": pass,
		},
	}

	return l, nil
}

func (l *LOINCTerminologyService) Ping() error {
	var err error

	req, err := http.NewRequest("GET", l.BaseUrl+"/CodeSystem/?url=http://loinc.org", nil)

	if err != nil {
		fmt.Println("error creating request")
		return err
	}

	req.SetBasicAuth(l.Auth["username"], l.Auth["password"])

	resp, err := l.Client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("LOINC terminology service is not available")
	}

	return nil
}

func (l *LOINCTerminologyService) Capability() (*fhir.Bundle, error) {
	var err error

	req, err := http.NewRequest("GET", l.BaseUrl+"/CodeSystem/?url=http://loinc.org", nil)

	if err != nil {
		fmt.Println("error creating request")
		return nil, err
	}

	req.SetBasicAuth(l.Auth["username"], l.Auth["password"])

	resp, err := l.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("LOINC terminology service is not available")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var bundle fhir.Bundle

	err = json.Unmarshal(body, &bundle)

	if err != nil {
		return nil, err
	}
	return &bundle, nil
}
