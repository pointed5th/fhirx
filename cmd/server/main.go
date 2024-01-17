package main

import (
	"log"
	"os"

	"github.com/hawyar/fhird"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "fhird",
		Usage:       "Experimental HL7 FHIR R4 Server",
		Description: "fhird is an experimental HL7 FHIR R4",
		Version:     "1.0.0",
		Commands: []*cli.Command{
			{
				Name:        "start",
				Usage:       "start the server",
				Description: "command to start the server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "port",
						Usage:       "port to listen on",
						Value:       fhird.DEFAULT_PORT,
						DefaultText: fhird.DEFAULT_PORT,
					},
				},
				Action: func(c *cli.Context) error {
					server, err := fhird.NewServer(fhird.DefaultConfig())

					if err != nil {
						return err
					}

					return server.Start()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
