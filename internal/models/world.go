package models

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
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

		if covRecord.Location != location && len(records) > 0 && len(location) > 0 {
			if strings.ToLower(location) == "world" {
				world = NewLocation(location, records)
			} else {
				locations = append(locations, NewLocation(location, records))
			}

			location = covRecord.Location
			records = nil
		} else if covRecord.Location != location {
			location = covRecord.Location
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

// Sort sorts the world data by the given descriptor and order.
func (w *World) Sort(descriptor string, order string) {
	sort.Slice(w.Locations, func(i, j int) bool {
		iLocation := w.Locations[i]
		jLocation := w.Locations[j]

		switch strings.ToLower(descriptor) {
		case "name":
			if order == "desc" {
				return strings.Compare(iLocation.Name, jLocation.Name) < 0
			}
			return strings.Compare(iLocation.Name, jLocation.Name) > 0
		case "newcases":
			if order == "desc" {
				return iLocation.NewCases() < jLocation.NewCases()
			}
			return iLocation.NewCases() > jLocation.NewCases()
		case "newdeaths":
			if order == "desc" {
				return iLocation.NewDeaths() < jLocation.NewDeaths()
			}
			return iLocation.NewDeaths() > jLocation.NewDeaths()
		case "totalcases":
			if order == "desc" {
				return iLocation.TotalCases() < jLocation.TotalCases()
			}
			return iLocation.TotalCases() > jLocation.TotalCases()
		case "totaldeaths":
			if order == "desc" {
				return iLocation.TotalDeaths() < jLocation.TotalDeaths()
			}
			return iLocation.TotalDeaths() > jLocation.TotalDeaths()
		default:
			if order == "desc" {
				return strings.Compare(iLocation.Name, jLocation.Name) < 0
			}
			return strings.Compare(iLocation.Name, jLocation.Name) > 0
		}
	})
}

// TotalCasesSignal returns the world's records' total cases
// as a float slice.
func (w World) TotalCasesSignal() []float64 {
	var signal []float64
	for _, r := range w.Records {
		signal = append(signal, float64(r.TotalCases))
	}
	return signal
}

// TotalDeathsSignal returns the world's records' total deaths
// as a float slice.
func (w World) TotalDeathsSignal() []float64 {
	var signal []float64
	for _, r := range w.Records {
		signal = append(signal, float64(r.TotalDeaths))
	}
	return signal
}

// TotalCasesSignalForLocation returns the location's records' total
// cases as a float slice.
func (w World) TotalCasesSignalForLocation(location string) []float64 {
	for _, l := range w.Locations {
		if strings.ToLower(location) == strings.ToLower(l.Name) {
			return l.TotalCasesSignal()
		}
	}
	return nil
}

// TotalDeathsSignalForLocation returns the location's records' total
// deaths as a float slice.
func (w World) TotalDeathsSignalForLocation(location string) []float64 {
	for _, l := range w.Locations {
		if strings.ToLower(location) == strings.ToLower(l.Name) {
			return l.TotalDeathsSignal()
		}
	}
	return nil
}
