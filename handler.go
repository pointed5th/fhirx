package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Page struct {
	Title string
}

func (fserver *FServer) RegisterHandlers() {
	r := fserver.Base.Handler.(*chi.Mux)

	apiVersion := "v1"
	base := fmt.Sprintf("/api/%s", apiVersion)
	pingEndpoint := "/ping"

	r.Route(base, func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "HEAD" && r.URL.Path == pingEndpoint {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("."))
				return
			}
		})

		r.Get("/metadata", func(w http.ResponseWriter, r *http.Request) {
			cap, err := GetCapabilityStatement()

			if err != nil {
				fmt.Errorf("error getting capability statement: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}

			capJson, err := json.Marshal(cap)

			if err != nil {
				fmt.Errorf("error marshalling capability statement: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("server error"))
			}

			w.Header().Set("Content-Type", "application/fhir+json")

			w.WriteHeader(http.StatusOK)
			w.Write(capJson)
		})

		// US Core Profile resource routes
		for _, resource := range USCoreProfileResources {
			r.Get(fmt.Sprintf("/%s", resource), func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/fhir+json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("US Core Profile Resource"))
			})
		}
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, err := w.Write([]byte("Not Found"))
		if err != nil {
			fmt.Errorf("error writing response: %s", err.Error())
			log.Fatal(err)
		}
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		_, err := w.Write([]byte("Not Allowed"))
		if err != nil {
			log.Fatal(err)
		}
	})
}

func PostPatientHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("Patient Created"))
	if err != nil {
		log.Fatal(err)
	}
}
func EscapeHTML(s string) template.HTML {
	return template.HTML(html.UnescapeString(s))
}
