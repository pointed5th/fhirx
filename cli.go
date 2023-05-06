package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

const (
	defaultPort = "9090"
)

func Run() error {
	app := cli.NewApp()

	app.Name = "fhird"
	app.Description = "HL7 FHIR server implementation for fructose.dev"
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
			config.Port = defaultPort
		}

		server := NewServer(config)

		err = server.Serve()

		if err != nil {
			return err
		}

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		return err
	}

	return nil
}
