package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var version = 1
var base = fmt.Sprintf("/api/v%d", version)

func (fserver *FHIRD) RegisterHandlers() {
	r := fserver.Base.Handler.(*chi.Mux)

	r.Route(base, func(r chi.Router) {
		r.Get("/ping", PingHandler)
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
}

func (fserver *FHIRD) USCoreProfileResourcesHandlers() {
	r := fserver.Base.Handler.(*chi.Mux)

	l := fserver.Logger

	for k, _ := range USCoreProfileResources {
		r.Route(fmt.Sprintf("%s/%s", base, k), func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				generalParams := r.Context().Value(ParamsCtxKey).(Paramateres)

				l.Info().Str("method", r.Method).Str("resource", k).Interface("params", generalParams).Str("path", r.URL.Path).Msg("GET")

				w.WriteHeader(200)
				w.Write([]byte(fmt.Sprintf("GET %s", r.URL.Path)))
			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				generalParams := r.Context().Value(ParamsCtxKey).(Paramateres)

				l.Info().Str("method", r.Method).Str("resource", k).Interface("params", generalParams).Str("path", r.URL.Path).Msg("POST")

				w.WriteHeader(200)
				w.Write([]byte(fmt.Sprintf("POST %s", r.URL.Path)))
			})
		})
	}
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	pingEndpoint := "/ping"

	if r.Method == "GET" || r.Method == "HEAD" && r.URL.Path == pingEndpoint {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("."))
		return
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
