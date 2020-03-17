// Package commands conatins the commands for the medina command line application.
package commands

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/colinc86/covid-19/internal/models"
	"github.com/urfave/cli/v2"
)

// GraphCommandHandler handles list commands.
type GraphCommandHandler struct {
	Name        string
	Aliases     []string
	Usage       string
	Description string

	// MARK: Private properties
	location string
	graph    string
}

// MARK: Initializers

// NewGraphCommandHandler creates and returns a new graph command handler.
func NewGraphCommandHandler() *GraphCommandHandler {
	return &GraphCommandHandler{
		Name:    "graph",
		Aliases: []string{"g"},
		Usage:   "graphs the dataset.",
		Description: `Graph the COVID-19 world data set.
		
		Examples:
			# Graph all data
			covid19 graph data
			
			# Graph the data for a specific location
			covid19 graph data -l [location]`,
	}
}

// MARK: Public methods

// Command creates and returns the handler's command.
func (h *GraphCommandHandler) Command() *cli.Command {
	return &cli.Command{
		Name:        h.Name,
		Aliases:     h.Aliases,
		Usage:       h.Usage,
		Description: h.Description,
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:    "data",
				Aliases: []string{"d"},
				Action:  h.GraphDataSetAction,
				Usage:   "The COVID-19 dataset.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "location",
						Aliases:     []string{"l"},
						Usage:       "Filter by location.",
						Required:    false,
						Destination: &h.location,
					},
					&cli.StringFlag{
						Name:        "value",
						Aliases:     []string{"v"},
						Usage:       "Value by newCases, newDeaths, totalCases or totalDeaths.",
						Required:    false,
						Destination: &h.graph,
					},
				},
			},
		},
	}
}

// GraphDataSetAction graphs the full dataset.
func (h *GraphCommandHandler) GraphDataSetAction(c *cli.Context) error {
	if os.Getenv("UPDATE_DATA") == "true" {
		// Update our data set
		err := updateDataset(localPath, dataSetURL)
		if err != nil {
			return err
		}
	}

	// Get the world locations
	world, err := models.NewWorldFromPath(localPath)
	if err != nil {
		return err
	}

	if len(h.graph) == 0 {
		h.graph = "totalCases"
	}

	// Get the total value for the location in question
	total := 0
	if len(h.location) == 0 || strings.ToLower(h.location) == "world" {
		switch strings.ToLower(h.graph) {
		case "newcases", "newdeaths":
			total = 0
		case "totalcases":
			total = world.TotalCases()
		case "totaldeaths":
			total = world.TotalDeaths()
		}
	} else if len(h.location) > 0 {
		for _, l := range world.Locations {
			if strings.ToLower(l.Name) != strings.ToLower(h.location) {
				continue
			}

			switch strings.ToLower(h.graph) {
			case "newcases":
				total = l.NewCases()
			case "newdeaths":
				total = l.NewDeaths()
			case "totalcases":
				total = l.TotalCases()
			case "totaldeaths":
				total = l.TotalDeaths()
			}

			break
		}
	}

	// Get the header value name
	name := ""
	switch strings.ToLower(h.graph) {
	case "newcases":
		name = "New Cases"
	case "newdeaths":
		name = "New Deaths"
	case "totalcases":
		name = "Total Cases"
	case "totaldeaths":
		name = "Total Deaths"
	}

	fmt.Printf("%-32s %-12s\n", "Date", name)

	// Draw the graphs
	if len(h.location) == 0 || strings.ToLower(h.location) == "world" {
		for _, r := range world.Records {
			value := 0
			switch strings.ToLower(h.graph) {
			case "newcases":
				value = r.NewCases
			case "newdeaths":
				value = r.NewDeaths
			case "totalcases":
				value = r.TotalCases
			case "totaldeaths":
				value = r.TotalDeaths
			}

			bar := ""
			ticks := int(math.Ceil(float64(value) / (float64(total) / 40.0)))
			for i := 0; i < ticks; i++ {
				bar += "#"
			}

			fmt.Printf("%-32v %-12d %s\n", r.Date, value, bar)
		}
	} else if len(h.location) > 0 {
		for _, l := range world.Locations {
			if strings.ToLower(l.Name) != strings.ToLower(h.location) {
				continue
			}

			for _, r := range l.Records {
				value := 0
				switch strings.ToLower(h.graph) {
				case "newcases":
					value = r.NewCases
				case "newdeaths":
					value = r.NewDeaths
				case "totalcases":
					value = r.TotalCases
				case "totaldeaths":
					value = r.TotalDeaths
				}

				bar := ""
				ticks := int(math.Ceil(float64(value) / (float64(total) / 40.0)))
				for i := 0; i < ticks; i++ {
					bar += "#"
				}

				fmt.Printf("%-32v %-12d %s\n", r.Date, value, bar)
			}

			break
		}
	}

	return nil
}
