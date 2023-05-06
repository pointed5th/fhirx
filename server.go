package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Config struct {
	Port        string
	Verbose     bool
	InContainer bool
	CGOEnabled  bool
}

type FServer struct {
	Base   *http.Server
	Config Config
	Logger zerolog.Logger
}

func NewServer(c Config) *FServer {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	server := &FServer{
		Base: &http.Server{
			Handler: chi.NewRouter(),
			Addr:    ":" + c.Port,
		},
		Config: c,
		Logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
	return server
}

func (fserver *FServer) Serve() error {
	l := fserver.Logger

	err := fserver.RegisterMiddlewares()

	if err != nil {
		return err
	}

	fserver.RegisterHandlers()

	router := fserver.Base.Handler.(*chi.Mux)

	if fserver.Config.Verbose {
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

	l.Info().Msgf("Listening on port %s", fserver.Config.Port)

	err = fserver.Base.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
