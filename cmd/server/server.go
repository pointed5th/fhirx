package main

import (
	"log"
	"os"

	"github.com/hawyar/fhird"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "fhird",
		Usage:   "Experimental US Core Server",
		Version: "1.0.0",
		Action: func(c *cli.Context) error {
			server, err := fhird.NewServer(fhird.DefaultConfig())

			if err != nil {
				return err
			}

			return server.Serve()
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
