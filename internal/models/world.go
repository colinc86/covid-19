package models

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// World types contain a set of locations and world records.
type World struct {

	// The world's locations.
	Locations []*Location

	// The world's records.
	Records []*COVRecord
}

// MARK: Initializers

// NewWorld creates and returns a new world.
func NewWorld(locations []*Location, records []*COVRecord) *World {
	return &World{
		Locations: locations,
		Records:   records,
	}
}

// NewWorldFromPath creats and returns a new world object from the given
// CSV path.
func NewWorldFromPath(path string) (*World, error) {
	// Read records from the csv
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvReader := csv.NewReader(bufio.NewReader(file))
	csvReader.FieldsPerRecord = 6

	var world *Location
	var locations []*Location
	var records []*COVRecord
	location := ""

	for {
		record, readError := csvReader.Read()
		if readError != nil {
			if readError == io.EOF {
				break
			}

			return nil, err
		}

		covRecord, err := NewCOVRecord(record)
		if err != nil {
			if err.Error() == "header" {
				continue
			}

			return nil, err
		}

		if covRecord.Location != location && len(records) > 0 {
			if strings.ToLower(location) == "world" {
				world = NewLocation(location, records)
			} else {
				locations = append(locations, NewLocation(location, records))
			}

			location = covRecord.Location
			records = nil
		}

		records = append(records, covRecord)
	}

	if len(records) > 0 {
		if strings.ToLower(location) == "world" {
			world = NewLocation(location, records)
		} else {
			locations = append(locations, NewLocation(location, records))
		}
	}

	return &World{
		Locations: locations,
		Records:   world.Records,
	}, nil
}

// MARK: Exported methods

// TotalCases returns the total cases at the location.
func (w World) TotalCases() int {
	if len(w.Records) > 0 {
		return w.Records[len(w.Records)-1].TotalCases
	}
	return 0
}

// TotalDeaths returns the total cases at the location.
func (w World) TotalDeaths() int {
	if len(w.Records) > 0 {
		return w.Records[len(w.Records)-1].TotalCases
	}
	return 0
}

// ListData lists the world data.
func (w World) ListData() {
	fmt.Printf("%-32s %-32s %-12s %-12s %-12s %-12s", "Date", "Location", "New Cases", "New Deaths", "Total Cases", "Total Deaths")

	for _, r := range w.Records {
		fmt.Printf(r.String() + "\n")
	}
}

// ListLocationData lists the world data.
func (w World) ListLocationData() {
	fmt.Printf("%-32s %-32s %-12s %-12s %-12s %-12s", "Date", "Location", "New Cases", "New Deaths", "Total Cases", "Total Deaths")

	for _, l := range w.Locations {
		fmt.Printf(l.String() + "\n")
	}
}
