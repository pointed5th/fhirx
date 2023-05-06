package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var SupportedContentTypes = []string{
	"application/x-www-form-urlencoded",
	"application/fhir+json",
	"application/fhir+xml",
	// the next two are for backwards compatibility with DSTU2 and STU3
	"application/json+fhir",
	"application/xml+fhir",
}

func (fserver *FHIRServer) RegisterMiddlewares() error {
	r := fserver.Base.Handler.(*chi.Mux)

	// keep logger first
	// r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.StripSlashes)
	r.Use(SetDefaultTimeZone)
	r.Use(middleware.AllowContentType(SupportedContentTypes...))
	//r.Use(ClientContentTypePrefer)
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
