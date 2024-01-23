package fhird

import (
	"os"
	"time"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
)

const (
	DEFAULT_TIMEZONE = "UTC"
	DEFAULT_PORT     = "9292"
)

type Config struct {
	Port         string
	FHIRVersion  fhirversion.Version
	USCDIVersion USCDIVersion
	InContainer  bool
	CGOEnabled   bool
	HTTPTimeout  time.Duration
	DataDir      string
	Timezone     string
}

func DefaultConfig() Config {
	return Config{
		Port:         DEFAULT_PORT,
		DataDir:      "./data",
		InContainer:  false,
		CGOEnabled:   false,
		FHIRVersion:  fhirversion.R4,
		USCDIVersion: USCDIv4,
		HTTPTimeout:  6 * time.Second,
		Timezone:     DEFAULT_TIMEZONE,
	}
}

func DefaultCapabilityStatement() (*CapabilityStatement, error) {
	file, err := os.ReadFile("./capability_statement.json")

	if err != nil {
		return nil, err
	}

	um, err := jsonformat.NewUnmarshaller("UTC", fhirversion.R4)

	if err != nil {
		return nil, err
	}

	unmarshalled, err := um.Unmarshal(file)

	if err != nil {
		return nil, err
	}

	resource := unmarshalled.(*r4pb.ContainedResource)

	return &CapabilityStatement{
		Resource: resource.GetCapabilityStatement(),
		Raw:      file,
	}, nil
}
