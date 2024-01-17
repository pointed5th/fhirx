package fhird

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) MountHandlers() {
	r := s.Srv.Handler.(*chi.Mux)

	r.Get("/metadata", s.GetCapabilityStmt)

	r.Route("/", func(r chi.Router) {
		r.Get("/CapabilityStatement", s.GetCapabilityStmt)
	})

	// Patient
	// https://build.fhir.org/ig/HL7/US-Core/StructureDefinition-us-core-patient.html
	r.Route("/Patient", func(r chi.Router) {
		r.Get("/", s.PatientSearch)
	})
}

func (s *Server) PatientSearch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Patient Search"))
}

func (s *Server) GetCapabilityStmt(w http.ResponseWriter, r *http.Request) {
	if len(s.CapStmt) == 0 {
		w.WriteHeader(500)
		w.Write([]byte("Capability Statement is empty"))
		return
	}
	w.WriteHeader(200)
	w.Write(s.CapStmt)
}
