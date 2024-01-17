package fhird

import (
	"os"

	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
	r4pb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/bundle_and_contained_resource_go_proto"
	cspb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/capability_statement_go_proto"
)

func DefaultCapability() (*cspb.CapabilityStatement, error) {
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

	return resource.GetCapabilityStatement(), nil
}
