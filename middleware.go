package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FHIRMIMEType int

const (
	FHIRJSON FHIRMIMEType = iota
	FHIRXML
)

func (f FHIRMIMEType) String() string {
	return []string{"application/fhir+json", "application/fhir+xml"}[f]
}

func (fserver *FHIRServer) RegisterMiddlewares() error {
	r := fserver.Base.Handler.(*chi.Mux)

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: &fserver.Logger}))
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(SetDefaultTimeZone)
	r.Use(middleware.AllowContentType(FHIRJSON.String(), FHIRXML.String()))

	return nil
}

func PostResourceHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var accept string

		if r.Header.Get("Accept") != "" {
			accept = r.Header.Get("Accept")
		} else {
			accept = "application/fhir+json"
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "Format", accept)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func SetDefaultTimeZone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Date", time.Now().Format(time.RFC1123))
		next.ServeHTTP(w, r)
	})
}

func SetWeakETag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `W/"weak"`)
		next.ServeHTTP(w, r)
	})
}
