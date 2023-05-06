package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FHIRMIMEType int
type SummaryParamValue int
type ParamsCtxKeyType string

type Paramateres struct {
	Format   string `json:"_format,omitempty"`
	Pretty   bool   `json:"_pretty,omitempty"`
	Summary  string `json:"_summary,omitempty"`
	Elements string `json:"_elements,omitempty"`
}

const (
	FHIRJSON FHIRMIMEType = iota
	FHIRXML

	SummaryTrue SummaryParamValue = iota
	SummaryText
	SummaryData
	SummaryCount
	SummaryFalse

	ParamsCtxKey ParamsCtxKeyType = "params"
)

func (s SummaryParamValue) String() string {
	switch s {
	case SummaryTrue:
		return "true"
	case SummaryText:
		return "text"
	case SummaryData:
		return "data"
	case SummaryCount:
		return "count"
	case SummaryFalse:
		return "false"
	default:
		return ""
	}
}

func (f FHIRMIMEType) String() string {
	switch f {
	case FHIRJSON:
		return "application/fhir+json"
	case FHIRXML:
		return "application/fhir+xml"
	default:
		return ""
	}
}

func (s Paramateres) String() string {
	return fmt.Sprintf("format=%s pretty=%t summary=%s elements=%s", s.Format, s.Pretty, s.Summary, s.Elements)
}

func (fserver *FHIRServer) RegisterMiddlewares() error {
	r := fserver.Base.Handler.(*chi.Mux)

	r.Use(middleware.RealIP)
	// r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: &fserver.Logger}))
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(SetDefaultTimeZone)
	r.Use(middleware.AllowContentType(FHIRJSON.String(), FHIRXML.String()))
	r.Use(SetGeneralParameters)
	return nil
}

// SetGeneralParameters sets any provided parameters in the request context
func SetGeneralParameters(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params Paramateres

		if r.URL.Query().Get("_format") != "" {
			params.Format = r.URL.Query().Get("_format")
		}

		if r.URL.Query().Get("_pretty") != "" {
			params.Pretty = true
		}

		if r.URL.Query().Get("_summary") != "" {
			switch r.URL.Query().Get("_summary") {
			case "true":
				params.Summary = SummaryTrue.String()

			case "text":
				params.Summary = SummaryText.String()

			case "data":
				params.Summary = SummaryData.String()

			case "count":
				params.Summary = SummaryCount.String()

			case "false":
				params.Summary = SummaryFalse.String()

			default:
				params.Summary = SummaryFalse.String()
			}
		}

		if r.URL.Query().Get("_elements") != "" {
			params.Elements = r.URL.Query().Get("_elements")
			// TODO: parse elements and validate given elements based on the resource
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ParamsCtxKey, params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetDefaultTimeZone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Date", time.Now().Format(time.RFC1123))
		next.ServeHTTP(w, r)
	})
}

// TODO: implement weak etag based on the resource version change
func SetWeakETag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `W/"weak"`)
		next.ServeHTTP(w, r)
	})
}
