package fhird

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/docgen"
)

type Server struct {
	*http.Server
	Logger        *Logger
	USCoreProfile *USCoreProfile
	Auth          *Auth
	Config        Config
}

func NewServer(config Config) (*Server, error) {
	uscore, err := NewUSCoreProfile(config)

	if err != nil {
		return nil, err
	}

	return &Server{
		Server: &http.Server{
			Handler: chi.NewRouter(),
			Addr:    ":" + config.Port,
		},
		Logger:        DefaultLogger(),
		USCoreProfile: uscore,
		Config:        config,
		Auth:          NewAuth(),
	}, nil
}

func (s *Server) Serve() error {
	r := s.Handler.(*chi.Mux)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Timeout(s.Config.HTTPTimeout))
	r.Use(middleware.AllowContentType("application/fhir+json", "application/json"))
	r.Use(middleware.SetHeader("Content-Type", "application/fhir+json"))
	r.Use(SetTimeZone)
	r.Use(ParseURLParams)
	r.Use(middleware.Heartbeat("/ping"))

	authm := s.Auth.Middleware()
	authRoutes, avaRoutes := s.Auth.Handlers()

	r.With(authm.Auth).Route("/protected", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("protected\n"))
		})
	})

	r.Mount("/auth", authRoutes)  // add auth handlers
	r.Mount("/avatar", avaRoutes) // add avatar handler

	r.With(s.USCoreProfile.ValidateRequestedResource).Route("/", func(r chi.Router) {
		r.Route("/CapabilityStatement", s.USCoreProfile.CapabilityStatement)
		r.Route("/AllergyIntolerance", s.USCoreProfile.AllergyIntolerance)
		r.Route("/CareTeam", s.USCoreProfile.CareTeam)
		r.Route("/DocumentReference", s.USCoreProfile.DocumentReference)
		r.Route("/Observation", s.USCoreProfile.Observation)
		r.Route("/Patient", s.USCoreProfile.Patient)
		r.Route("/Encounter", s.USCoreProfile.Encounter)
		r.Route("/Location", s.USCoreProfile.Location)
		r.Route("/Goal", s.USCoreProfile.Goal)
		r.Route("/Coverage", s.USCoreProfile.Coverage)
		r.Route("/Immunization", s.USCoreProfile.Immunization)
		r.Route("/Device", s.USCoreProfile.Device)
		r.Route("/Medication", s.USCoreProfile.Medication)
	})

	s.PrintRoutes()

	err := s.GenerateDocs()

	if err != nil {
		return err
	}

	err = s.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (f *Server) GenerateDocs() error {
	r := f.Handler.(*chi.Mux)

	rjson, err := os.OpenFile("server_info.json", os.O_WRONLY|os.O_CREATE, 0644)

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
	router := f.Handler.(*chi.Mux)

	fmt.Println()

	theader := fmt.Sprintf("%-6s | %-6s\n", "METHOD", "ROUTE")
	tsep := fmt.Sprintf("%-6s + %-6s\n", "------", "------")
	fmt.Print(theader)
	fmt.Print(tsep)

	// fmt.Printf("%-6s | %-6s\n", "GET", "/ping")
	// fmt.Printf("%-6s---%-6s\n", "------", "------")

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
