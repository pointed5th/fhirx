package fhird

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
	"github.com/google/fhir/go/fhirversion"
	"github.com/google/fhir/go/jsonformat"
)

type Config struct {
	Port        string
	FHIRVersion fhirversion.Version
	InContainer bool
	CGOEnabled  bool
	HTTPTimeout time.Duration
	DataDir     string
	Timezone    string
}

type Server struct {
	Srv     *http.Server
	Logger  *Logger
	Config  Config
	Profile *USCoreProfile
	CapStmt []byte // unmarshall once
}

func DefaultConfig() Config {
	return Config{
		Port:        "9292",
		DataDir:     "./data",
		InContainer: false,
		CGOEnabled:  false,
		FHIRVersion: fhirversion.R4,
		HTTPTimeout: 6 * time.Second,
		Timezone:    "UTC",
	}
}

func NewServer(config Config) (*Server, error) {
	cap, err := DefaultCapability()

	if err != nil {
		return nil, err
	}

	marshaller, err := jsonformat.NewMarshaller(true, "", " ", config.FHIRVersion)

	if err != nil {
		return nil, err
	}

	capjson, err := marshaller.MarshalResource(cap)

	if err != nil {
		return nil, err
	}

	return &Server{
		Srv: &http.Server{
			Handler: chi.NewRouter(),
			Addr:    ":" + config.Port,
		},
		Config:  config,
		Logger:  DefaultLogger(),
		Profile: NewUSCoreProfile(config.FHIRVersion),
		CapStmt: capjson,
	}, nil
}

func (f *Server) Start() error {
	f.MountMiddlewares()
	f.MountHandlers()
	f.PrintRoutes()

	err := f.Srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (f *Server) GenerateDocs() error {
	r := f.Srv.Handler.(*chi.Mux)

	rjson, err := os.OpenFile("routes.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	defer rjson.Close()

	_, err = rjson.WriteString("\n")

	if err != nil {
		return err
	}

	_, err = rjson.WriteString(docgen.JSONRoutesDoc(r))

	if err != nil {
		return err
	}

	return nil
}

func (f *Server) PrintRoutes() error {
	router := f.Srv.Handler.(*chi.Mux)

	fmt.Println()

	theader := fmt.Sprintf("%-6s | %-6s\n", "METHOD", "ROUTE")
	tsep := fmt.Sprintf("%-6s + %-6s\n", "------", "------")
	fmt.Print(theader)
	fmt.Print(tsep)

	fmt.Printf("%-6s | %-6s\n", "GET", "/ping")
	fmt.Printf("%-6s---%-6s\n", "------", "------")
	walker := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route != "/" && route[len(route)-1:] == "/" {
			route = route[:len(route)-1]
		}

		fmt.Printf("%-6s | %-6s\n", method, route)

		// skip the last route
		fmt.Printf("%-6s---%-6s\n", "------", "------")
		return nil
	}

	if err := chi.Walk(router, walker); err != nil {
		return err
	}
	// fmt.Print(tsep + "\n")

	return nil
}
