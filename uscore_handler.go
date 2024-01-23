package fhird

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/fhir/go/fhirversion"
	cspb "github.com/google/fhir/go/proto/google/fhir/proto/r4/core/resources/capability_statement_go_proto"
	uscore "github.com/google/fhir/go/proto/google/fhir/proto/r4/uscore_go_proto"
)

// United States Core Data for Interoperability
type USCDIVersion int

const (
	USCDIv1 USCDIVersion = iota + 1
	USCDIv2
	USCDIv3
	USCDIv4
)

type CapabilityStatement struct {
	Resource *cspb.CapabilityStatement
	Raw      []byte
}

// https://build.fhir.org/ig/HL7/US-Core/index.html#us-core-profiles
type USCoreProfile struct {
	USCDIVersion USCDIVersion
	FHIRVersion  fhirversion.Version
	Profiles     map[string]interface{}
	CapStmt      CapabilityStatement
}

func NewUSCoreProfile(config Config) (*USCoreProfile, error) {
	cap, err := DefaultCapabilityStatement()

	if err != nil {
		return nil, err
	}

	return &USCoreProfile{
		USCDIVersion: config.USCDIVersion,
		FHIRVersion:  config.FHIRVersion,
		CapStmt:      *cap,
		Profiles: map[string]interface{}{
			"Patient":             uscore.USCorePatientProfile{},
			"CapabilityStatement": cspb.CapabilityStatement{},
		},
	}, nil
}

func (u *USCoreProfile) ValidateRequestedResource(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource := chi.URLParam(r, "*")

		if _, ok := u.Profiles[resource]; !ok {
			w.WriteHeader(404)
			render.JSON(w, r, fmt.Sprintf("Resource %s not found", resource))
			return
		}

		next.ServeHTTP(w, r)
	})
}

type HTTPErrResponse struct {
	Error string `json:"error"`
}

// CapabilityStatement https://build.fhir.org/ig/HL7/US-Core/CapabilityStatement-us-core-server.html
func (u *USCoreProfile) CapabilityStatement(r chi.Router) {
	// reject if its not a GET request
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if len(u.CapStmt.Raw) == 0 {
			w.WriteHeader(500)
			render.JSON(w, r, HTTPErrResponse{Error: "CapabilityStatement not found"})
			return
		}
		w.WriteHeader(200)
		w.Write(u.CapStmt.Raw)
	})
}

// AllergyIntolerance https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-allergyintolerance.html
func (u *USCoreProfile) AllergyIntolerance(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("AllergyIntolerance Search"))
	})
}

// CareTeam https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-careteam.html
func (u *USCoreProfile) CareTeam(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("CareTeam Search"))
	})
}

// DocumentReference https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-documentreference.html
func (u *USCoreProfile) DocumentReference(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("DocumentReference Search"))
	})
}

// Observation https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-observation-clinical-result.html
func (u *USCoreProfile) Observation(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Observation Search"))
	})
}

// Encounter https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-encounter.html
func (u *USCoreProfile) Encounter(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Encounter Search"))
	})
}

// Patient https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-patient.html
func (u *USCoreProfile) Patient(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Patient Search"))
	})
}

// Loction https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-location.html
func (u *USCoreProfile) Location(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Location Search"))
	})
}

// Goal https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-goal.html
func (u *USCoreProfile) Goal(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Goal Search"))
	})
}

// Coverage https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-coverage.html
func (u *USCoreProfile) Coverage(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Coverage Search"))
	})
}

// Immunization  https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-immunization.html
func (u *USCoreProfile) Immunization(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Immunization Search"))
	})
}

// Device https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-implantable-device.html
func (u *USCoreProfile) Device(r chi.Router) {
	// dev := uscore.USCoreImplantableDeviceProfile{}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Device Search"))
	})
}

// Medication https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-medication.html
func (u *USCoreProfile) Medication(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Medication Search"))
	})
}

// MedicationRequest https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-medicationrequest.html
func (u *USCoreProfile) MedicationRequest(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("MedicationRequest"))
	})
}

// MedicationDispense https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-medicationdispense.html
func (u *USCoreProfile) MedicationDispense(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("MedicationDispense"))
	})
}
