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
	graphHandler := commands.NewGraphCommandHandler()
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
				Name:     "update",
				Usage:    "Updates the dataset before executing the given command.",
				Aliases:  []string{"u"},
				Value:    false,
				Required: false,
			},
		},
		Commands: []*cli.Command{
			graphHandler.Command(),
			listHandler.Command(),
			updateHandler.Command(),
		},
		Before:                 before,
		UseShortOptionHandling: true,
	}

	// Run the application
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// MARK: Application callbacks

func before(c *cli.Context) error {
	if c.Bool("update") {
		os.Setenv("UPDATE_DATA", "true")
	} else {
		os.Setenv("UPDATE_DATA", "false")
	}

	return nil
}
