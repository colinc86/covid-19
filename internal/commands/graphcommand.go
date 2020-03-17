// Package commands conatins the commands for the medina command line application.
package commands

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
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
				},
			},
		},
	}
}

// GraphDataSetAction graphs the full dataset.
func (h *GraphCommandHandler) GraphDataSetAction(c *cli.Context) error {
	// Read records from the csv
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}

	defer file.Close()

	csvReader := csv.NewReader(bufio.NewReader(file))
	csvReader.FieldsPerRecord = 6

	// Get the max value in the list
	for {
		record, readError := csvReader.Read()
		if readError != nil {
			if readError == io.EOF {
				break
			}

			return err
		}

		covRecord, err := models.NewCOVRecord(record)
		if err != nil {
			if err.Error() == "header" {
				continue
			}

			return err
		}

		if len(h.location) > 0 && strings.ToLower(covRecord.Location) == strings.ToLower(h.location) {
			fmt.Printf(covRecord.String() + "\n")
		} else if len(h.location) == 0 {
			fmt.Printf(covRecord.String() + "\n")
		}
	}

	file.Seek(0, 0)

	fmt.Printf("%-32s %-32s %-12s %-12s %-12s %-12s", "Date", "Location", "New Cases", "New Deaths", "Total Cases", "Total Deaths")

	for {
		record, readError := csvReader.Read()
		if readError != nil {
			if readError == io.EOF {
				break
			}

			return err
		}

		covRecord, err := models.NewCOVRecord(record)
		if err != nil {
			if err.Error() == "header" {
				continue
			}

			return err
		}

		if len(h.location) > 0 && strings.ToLower(covRecord.Location) == strings.ToLower(h.location) {
			fmt.Printf(covRecord.String() + "\n")
		} else if len(h.location) == 0 {
			fmt.Printf(covRecord.String() + "\n")
		}
	}

	file.Seek(0, 0)
	return nil
}
