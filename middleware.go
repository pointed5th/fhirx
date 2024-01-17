package fhird

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
type CtxKey string

const (
	FHIRJSON FHIRMIMEType = iota
	FHIRXML
	JSON
	XML

	SummaryTrue SummaryParamValue = iota
	SummaryText
	SummaryData
	SummaryCount
	SummaryFalse

	ParamsKey CtxKey = "params"
)

type Paramaters struct {
	Format   string `json:"_format,omitempty"`
	Pretty   bool   `json:"_pretty,omitempty"`
	Summary  string `json:"_summary,omitempty"`
	Elements string `json:"_elements,omitempty"`
}

type RequestCtx struct {
	Headers     http.Header `json:"headers"`
	Host        string      `json:"host"`
	Method      string      `json:"method"`
	RemoteAddr  string      `json:"remote_address"`
	RequestURI  string      `json:"request_uri"`
	URL         string      `json:"url"`
	UserAgent   string      `json:"user_agent"`
	Paramateres Paramaters  `json:"parameters"`
	ContentBody string      `json:"content_body"`
}

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
	case JSON:
		return "application/json"
	case XML:
		return "application/xml"
	default:
		return ""
	}
}

func (s Paramaters) String() string {
	return fmt.Sprintf("format=%s pretty=%t summary=%s elements=%s", s.Format, s.Pretty, s.Summary, s.Elements)
}

func (f *Server) MountMiddlewares() {
	r := f.Srv.Handler.(*chi.Mux)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(f.Config.HTTPTimeout))
	r.Use(middleware.AllowContentType("application/fhir+json", "application/json"))
	r.Use(middleware.SetHeader("Content-Type", "application/fhir+json"))
	r.Use(SetTimeZone)
	r.Use(middleware.Heartbeat("/ping"))
}

func (f *Server) ParseURLParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params Paramaters

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

		ctx := context.WithValue(r.Context(), ParamsKey, params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetTimeZone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Date", time.Now().UTC().Format(time.RFC3339))
		next.ServeHTTP(w, r)
	})
}
