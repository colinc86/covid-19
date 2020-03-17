// Package commands conatins the commands for the medina command line application.
package commands

import (
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

const dataSetURL = "https://covid.ourworldindata.org/data/full_data.csv"
const localPath = "/usr/local/var/covid_full_data.csv"

// UpdateCommandHandler handles update commands.
type UpdateCommandHandler struct {
	Name        string
	Aliases     []string
	Usage       string
	Description string
}

// MARK: Initializers

// NewUpdateCommandHandler creates and returns a new update command handler.
func NewUpdateCommandHandler() *UpdateCommandHandler {
	return &UpdateCommandHandler{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Update the dataset.",
		Description: `Update the COVID-19 world data set.
		
		Examples:
			# Update all data
			covid19 update data`,
	}
}

// MARK: Public methods

// Command creates and returns the handler's command.
func (h *UpdateCommandHandler) Command() *cli.Command {
	return &cli.Command{
		Name:        h.Name,
		Aliases:     h.Aliases,
		Usage:       h.Usage,
		Description: h.Description,
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:    "data",
				Aliases: []string{"d"},
				Action:  h.UpdateDataSetAction,
				Usage:   "The COVID-19 dataset.",
			},
		},
	}
}

// UpdateDataSetAction updates the full dataset.
func (h *UpdateCommandHandler) UpdateDataSetAction(c *cli.Context) error {
	// Update our data set
	err := updateDataset(localPath, dataSetURL)
	if err != nil {
		return err
	}

	return nil
}

// MARK: Unexported methods

// updateDataset updates the dataset at the given url and saves it to filepath.
func updateDataset(filepath string, url string) error {
	s := NewSpinnerWithTitle("Updating dataset...")
	s.Start()
	defer s.Stop()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
