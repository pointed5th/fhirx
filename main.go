package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Name = "fhird"
	app.Description = "Experimental HL7 FHIR R4 Server"
	app.Version = "1.0.0"
	app.BashComplete = cli.DefaultAppComplete
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose output",
			Value: false,
		},
	}

	app.Action = func(c *cli.Context) error {
		var err error

		config := Config{
			Port:       os.Getenv("PORT"),
			Verbose:    c.Bool("verbose"),
			CGOEnabled: os.Getenv("CGO_ENABLED") == "false",
		}

		if config.Port == "" {
			config.Port = "9090"
		}

		server, err := NewFHIRD(config)

		if err != nil {
			return err
		}

		err = server.Serve()

		if err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
