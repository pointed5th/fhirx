package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/samply/golang-fhir-models/fhir-models/fhir"
)

// map of resourcetype name to resource
var USCoreProfileResources = map[string]fhir.ResourceType{
	// fhir.ResourceTypeAllergyIntolerance,
	"AllergyIntolerance": fhir.ResourceTypeAllergyIntolerance,
	// fhir.ResourceTypeCarePlan,
	"CarePlan": fhir.ResourceTypeCarePlan,
	// fhir.ResourceTypeCareTeam,
	"CareTeam": fhir.ResourceTypeCareTeam,
	// fhir.ResourceTypeCondition,
	"Condition": fhir.ResourceTypeCondition,
	// fhir.ResourceTypeDevice,
	"Device": fhir.ResourceTypeDevice,
	// fhir.ResourceTypeDiagnosticReport,
	"DiagnosticReport": fhir.ResourceTypeDiagnosticReport,
	// fhir.ResourceTypeDocumentReference,
	"DocumentReference": fhir.ResourceTypeDocumentReference,
	// fhir.ResourceTypeEncounter,
	"Encounter": fhir.ResourceTypeEncounter,
	// fhir.ResourceTypeGoal,
	"Goal": fhir.ResourceTypeGoal,
	// fhir.ResourceTypeImmunization,
	"Immunization": fhir.ResourceTypeImmunization,
	// fhir.ResourceTypeLocation,
	"Location": fhir.ResourceTypeLocation,
	// fhir.ResourceTypeMedication,
	"Medication": fhir.ResourceTypeMedication,
	// fhir.ResourceTypeMedicationRequest,
	"MedicationRequest": fhir.ResourceTypeMedicationRequest,
	// fhir.ResourceTypeObservation,
	"Observation": fhir.ResourceTypeObservation,
	// fhir.ResourceTypeOrganization,
	"Organization": fhir.ResourceTypeOrganization,
	// fhir.ResourceTypePatient,
	"Patient": fhir.ResourceTypePatient,
	// fhir.ResourceTypePractitioner,
	"Practitioner": fhir.ResourceTypePractitioner,
	// fhir.ResourceTypePractitionerRole,
	"PractitionerRole": fhir.ResourceTypePractitionerRole,
	// fhir.ResourceTypeProcedure,
	"Procedure": fhir.ResourceTypeProcedure,
	// fhir.ResourceTypeProvenance,
	"Provenance": fhir.ResourceTypeProvenance,
	// fhir.ResourceTypeQuestionnaireResponse,
	"QuestionnaireResponse": fhir.ResourceTypeQuestionnaireResponse,
	// fhir.ResourceTypeRelatedPerson,
	"RelatedPerson": fhir.ResourceTypeRelatedPerson,
	// fhir.ResourceTypeServiceRequest,
	"ServiceRequest": fhir.ResourceTypeServiceRequest,
}

// US Core Profile resources (https://www.hl7.org/fhir/us/core/#us-core-profiles)
// type USCoreProfile struct {
// 	ResourceTypeAllergyIntolerance    fhir.AllergyIntolerance
// 	ResourceTypeCarePlan              fhir.CarePlan
// 	ResourceTypeCareTeam              fhir.CareTeam
// 	ResourceTypeCondition             fhir.Condition
// 	ResourceTypeDevice                fhir.Device
// 	ResourceTypeDiagnosticReport      fhir.DiagnosticReport
// 	ResourceTypeDocumentReference     fhir.DocumentReference
// 	ResourceTypeEncounter             fhir.Encounter
// 	ResourceTypeGoal                  fhir.Goal
// 	ResourceTypeImmunization          fhir.Immunization
// 	ResourceTypeLocation              fhir.Location
// 	ResourceTypeMedication            fhir.Medication
// 	ResourceTypeMedicationRequest     fhir.MedicationRequest
// 	ResourceTypeObservation           fhir.Observation
// 	ResourceTypeOrganization          fhir.Organization
// 	ResourceTypePatient               fhir.Patient
// 	ResourceTypePractitioner          fhir.Practitioner
// 	ResourceTypePractitionerRole      fhir.PractitionerRole
// 	ResourceTypeProcedure             fhir.Procedure
// 	ResourceTypeProvenance            fhir.Provenance
// 	ResourceTypeQuestionnaireResponse fhir.QuestionnaireResponse
// 	ResourceTypeRelatedPerson         fhir.RelatedPerson
// 	ResourceTypeServiceRequest        fhir.ServiceRequest
// }

func GetCapabilityStatement() (*fhir.CapabilityStatement, error) {
	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	url := fmt.Sprintf("http://localhost:%s", port)
	title := "Capability Statement for FHIR Server"
	purpose := "Main EHR capability statement, published for contracting and operational support"
	name := "fhird"
	publisher := "fructose"
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
			// TODO: Add supported resources
			Resource:    restResource,
			SearchParam: nil,
			Operation:   nil,
			Compartment: nil,
		},
		{
			Id:                nil,
			Extension:         nil,
			ModifierExtension: nil,
			Mode:              fhir.RestfulCapabilityModeClient,
			Documentation:     nil,
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
			Name:    "FHIR Test Server",
			Version: &version,
		},
		Format: []string{FHIRXML.String(), FHIRJSON.String()},
		Rest:   rest,
	}, nil
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
