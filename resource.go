package fhird

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// US Core Profile resources (https://www.hl7.org/fhir/us/core/#us-core-profiles)
var USCoreProfileResources = map[string]fhir.ResourceType{
	"AllergyIntolerance":    fhir.ResourceTypeAllergyIntolerance,
	"CarePlan":              fhir.ResourceTypeCarePlan,
	"CareTeam":              fhir.ResourceTypeCareTeam,
	"Condition":             fhir.ResourceTypeCondition,
	"Device":                fhir.ResourceTypeDevice,
	"DiagnosticReport":      fhir.ResourceTypeDiagnosticReport,
	"DocumentReference":     fhir.ResourceTypeDocumentReference,
	"Encounter":             fhir.ResourceTypeEncounter,
	"Goal":                  fhir.ResourceTypeGoal,
	"Immunization":          fhir.ResourceTypeImmunization,
	"Location":              fhir.ResourceTypeLocation,
	"Medication":            fhir.ResourceTypeMedication,
	"MedicationDispense":    fhir.ResourceTypeMedicationDispense,
	"MedicationRequest":     fhir.ResourceTypeMedicationRequest,
	"Observation":           fhir.ResourceTypeObservation,
	"Organization":          fhir.ResourceTypeOrganization,
	"Patient":               fhir.ResourceTypePatient,
	"Practitioner":          fhir.ResourceTypePractitioner,
	"PractitionerRole":      fhir.ResourceTypePractitionerRole,
	"Procedure":             fhir.ResourceTypeProcedure,
	"Provenance":            fhir.ResourceTypeProvenance,
	"QuestionnaireResponse": fhir.ResourceTypeQuestionnaireResponse,
	"RelatedPerson":         fhir.ResourceTypeRelatedPerson,
	"ServiceRequest":        fhir.ResourceTypeServiceRequest,
	"Specimen":              fhir.ResourceTypeSpecimen,
}

func CapabilityStatement() *fhir.CapabilityStatement {
	port := os.Getenv("PORT")

	if port == "" {
		port = "9090"
	}

	url := fmt.Sprintf("http://localhost:%s", port)
	title := "Capability Statement for the FHIR Server"
	purpose := "Experimental HL7 FHIR R4 Server"
	name := "fhird"
	publisher := "fhird admin"
	copyright := "Copyright (c) 2023"
	experimental := true
	version := "1.0.0"

	var restResource []fhir.CapabilityStatementRestResource

	for _, resource := range USCoreProfileResources {
		restResource = append(restResource, fhir.CapabilityStatementRestResource{
			Type: resource,
			Interaction: []fhir.CapabilityStatementRestResourceInteraction{
				{
					Code: fhir.TypeRestfulInteractionRead,
				},
				{
					Code: fhir.TypeRestfulInteractionVread,
				},
				{
					Code: fhir.TypeRestfulInteractionUpdate,
				},
				{
					Code: fhir.TypeRestfulInteractionDelete,
				},
			},
		})
	}

	rest := []fhir.CapabilityStatementRest{
		{
			Id:                nil,
			Extension:         nil,
			ModifierExtension: nil,
			Mode:              fhir.RestfulCapabilityModeServer,
			Documentation:     nil,
			Security:          nil,
			Resource:          restResource,
			SearchParam:       nil,
			Operation:         nil,
			Compartment:       nil,
		},
	}

	return &fhir.CapabilityStatement{
		Name:         &name,
		Id:           nil,
		Url:          &url,
		Purpose:      &purpose,
		Title:        &title,
		FhirVersion:  fhir.FHIRVersion4_0_1,
		Experimental: &experimental,
		Publisher:    &publisher,
		Copyright:    &copyright,
		Kind:         fhir.CapabilityStatementKindCapability,
		Status:       fhir.PublicationStatusDraft,
		Date:         time.Now().Format(time.RFC3339),
		Software: &fhir.CapabilityStatementSoftware{
			Name:    "FHIRD",
			Version: &version,
		},
		Format: []string{FHIRXML.String(), FHIRJSON.String()},
		Rest:   rest,
	}
}

func NewPatient(r *http.Request) (fhir.Patient, error) {
	var patient fhir.Patient

	err := json.NewDecoder(r.Body).Decode(&patient)

	if err != nil {
		return patient, err
	}

	id := "22"

	patient.Id = &id

	sys := "https://example.com/"
	identifier := fhir.Identifier{
		System: &sys,
		Value:  &id,
	}

	var identifiers []fhir.Identifier

	patient.Identifier = append(identifiers, identifier)

	return patient, nil
}
