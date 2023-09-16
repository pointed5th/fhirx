package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
	"github.com/rs/zerolog"
)

type Config struct {
	Port        string
	Verbose     bool
	InContainer bool
	CGOEnabled  bool
}

type FHIRD struct {
	Base   *http.Server
	Config Config
	Logger *Logger
}

func NewFHIRD(c Config) (*FHIRD, error) {
	l, err := NewConsoleLogger()

	if err != nil {
		return nil, err
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err != nil {
		return nil, err
	}

	return &FHIRD{
		Base: &http.Server{
			Handler: chi.NewRouter(),
			Addr:    ":" + c.Port,
		},
		Config: c,
		Logger: l,
	}, nil
}

func (f *FHIRD) Serve() error {
	l := f.Logger.Sugar()

	var err error

	err = f.RegisterMiddlewares()

	if err != nil {
		return err
	}

	f.RegisterHandlers()

	if err != nil {
		return err
	}

	l.Info("registered handlers successfully")

	if err != nil {
		return errors.New("failed to create logger")
	}

	l.Infof("listening on port %s", f.Config.Port)

	err = f.Base.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (f *FHIRD) GenerateDocs() error {
	r := f.Base.Handler.(*chi.Mux)

	doc := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
		ProjectPath: "https://github.com/hawyar/fhird",
		Intro:       "fructose API docs",
	})

	readme, err := os.OpenFile("README.md", os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer readme.Close()

	_, err = readme.WriteString("\n")

	if err != nil {
		return err
	}

	_, err = readme.WriteString(doc)

	if err != nil {
		return err
	}

	l := f.Logger.Sugar()

	l.Info("Generated docs successfully")

	return nil
}
func (f *FHIRD) PrintRoutes() error {
	router := f.Base.Handler.(*chi.Mux)

	if f.Config.Verbose {
		fmt.Println()
		theader := fmt.Sprintf("%-6s | %-6s\n", "METHOD", "ROUTE")
		tsep := fmt.Sprintf("%-6s + %-6s\n", "------", "------")
		fmt.Print(theader)
		fmt.Print(tsep)

		walker := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			if route != "/" && route[len(route)-1:] == "/" {
				route = route[:len(route)-1]
			}

			fmt.Printf("%-6s | %-6s\n", method, route)
			return nil
		}

		if err := chi.Walk(router, walker); err != nil {
			return err
		}

		fmt.Print(tsep + "\n")
	}

	return nil
}
