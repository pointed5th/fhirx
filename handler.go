package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var version = 1
var base = fmt.Sprintf("/api/v%d", version)

type IndexPageContext struct {
	Title  string `json:"title"`
	Header string `json:"header"`
}

func (f *FHIRD) RegisterHandlers() {
	r := f.Base.Handler.(*chi.Mux)

	r.Route(base, func(r chi.Router) {
		r.Get("/metadata", MetadataHandler)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("Not Allowed"))
	})

	f.USCoreProfileResourcesHandlers()
}

func (f *FHIRD) USCoreProfileResourcesHandlers() {
	r := f.Base.Handler.(*chi.Mux)
	// l := f.Logger.Sugar()

	for k, _ := range USCoreProfileResources {
		rr := k
		r.Route(fmt.Sprintf("%s/%s", base, rr), func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				// generalParams := r.Context().Value(ParamsCtxKey).(Paramateres)
				// l.Info().Str("method", r.Method).Str("resource", k).Interface("params", generalParams).Str("path", r.URL.Path).Msg("GET")

				w.WriteHeader(200)
				w.Write([]byte(fmt.Sprintf("GET %s", r.URL.Path)))
			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				// generalParams := r.Context().Value(ParamsCtxKey).(Paramateres)

				// l.Info().Str("method", r.Method).Str("resource", k).Interface("params", generalParams).Str("path", r.URL.Path).Msg("POST")

				w.WriteHeader(200)
				w.Write([]byte(fmt.Sprintf("POST %s", r.URL.Path)))
			})
		})
	}
}

func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	stmt, err := GetCapabilityStatement()

	if err != nil {
		fmt.Printf("error getting capability statement: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}

	cap, err := json.Marshal(stmt)

	if err != nil {
		fmt.Printf("error marshalling capability statement: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}

	w.Header().Set("Content-Type", "application/fhir+json")
	w.WriteHeader(http.StatusOK)
	w.Write(cap)
}
