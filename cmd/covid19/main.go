package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/colinc86/covid-19/internal/commands"
	"github.com/urfave/cli/v2"
)

const (
	name     string = "covid19"
	version  string = "v0.1.0"
	helpName string = "covid19"
	usage    string = "Update and train from COVID-19 datasets."
)

func main() {
	// Seed rand
	rand.Seed(time.Now().Unix())

	// Setup the commands
	listHandler := commands.NewListCommandHandler()
	updateHandler := commands.NewUpdateCommandHandler()

	// Setup the application
	app := &cli.App{
		Name:     name,
		Version:  version,
		Compiled: time.Now(),
		HelpName: helpName,
		Usage:    usage,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "verbose",
				Usage:    "Verbose output.",
				Value:    false,
				Required: false,
			},
		},
		Commands: []*cli.Command{
			listHandler.Command(),
			updateHandler.Command(),
		},
		UseShortOptionHandling: true,
	}

	// Run the application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
