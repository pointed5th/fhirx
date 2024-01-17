package fhird

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/fhir/go/fhirversion"
	uscore "github.com/google/fhir/go/proto/google/fhir/proto/r4/uscore_go_proto"
)

var _ = uscore.USCoreAllergyIntolerance{}

type USCoreProfile struct {
	Version            fhirversion.Version
	SupportedResources map[string]string
	Router             *chi.Router
}

func NewUSCoreProfile(version fhirversion.Version) *USCoreProfile {
	return &USCoreProfile{
		Version: version,
		SupportedResources: map[string]string{
			"Patient": "Patient",
		},
	}
}
