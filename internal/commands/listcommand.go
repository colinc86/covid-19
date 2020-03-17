// Package commands conatins the commands for the medina command line application.
package commands

import (
	"fmt"
	"strings"

	"github.com/colinc86/covid-19/internal/models"
	"github.com/urfave/cli/v2"
)

// ListCommandHandler handles list commands.
type ListCommandHandler struct {
	Name        string
	Aliases     []string
	Usage       string
	Description string

	// MARK: Private properties
	location string
	world    bool
}

// MARK: Initializers

// NewListCommandHandler creates and returns a new list command handler.
func NewListCommandHandler() *ListCommandHandler {
	return &ListCommandHandler{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "Lists the dataset.",
		Description: `List the COVID-19 world data set.
		
		Examples:
			# List a summary of location data
			covid19 list data
			
			# List the data for a specific location
			covid19 list data -l [location]
			
			# List the world data
			covid19 list data -w`,
	}
}

// MARK: Public methods

// Command creates and returns the handler's command.
func (h *ListCommandHandler) Command() *cli.Command {
	return &cli.Command{
		Name:        h.Name,
		Aliases:     h.Aliases,
		Usage:       h.Usage,
		Description: h.Description,
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:    "data",
				Aliases: []string{"d"},
				Action:  h.ListDataSetAction,
				Usage:   "The COVID-19 dataset.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "location",
						Aliases:     []string{"l"},
						Usage:       "Filter by location.",
						Required:    false,
						Destination: &h.location,
					},
					&cli.BoolFlag{
						Name:        "world",
						Aliases:     []string{"w"},
						Usage:       "List world data.",
						Required:    false,
						Destination: &h.world,
					},
				},
			},
		},
	}
}

// ListDataSetAction lists the full dataset.
func (h *ListCommandHandler) ListDataSetAction(c *cli.Context) error {
	// Get the world locations
	world, err := models.NewWorldFromPath(localPath)
	if err != nil {
		return err
	}

	fmt.Printf("%-32s %-32s %-12s %-12s %-12s %-12s", "Date", "Location", "New Cases", "New Deaths", "Total Cases", "Total Deaths")

	for _, l := range world.Locations {
		if len(h.location) > 0 && strings.ToLower(l.Name) != strings.ToLower(h.location) {
			continue
		}

		for _, r := range l.Records {
			fmt.Printf(r.String() + "\n")
		}
	}
	return nil
}
