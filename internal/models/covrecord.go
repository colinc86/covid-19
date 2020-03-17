package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// COVRecord represents a COVID-19 record for the specified location at the
// given time.
type COVRecord struct {
	// The 24 hour observation period of this record.
	Date time.Time

	// The location of this record.
	Location string

	// The number of new cases.
	NewCases int

	// The number of new deaths.
	NewDeaths int

	// The total number of cases for this location.
	TotalCases int

	// The total number of deaths for this location.
	TotalDeaths int
}

// MARK: Initializers

// NewCOVRecord creates and returns a new COVID-19 record with the given row
// from a CSV dataset.
func NewCOVRecord(record []string) (*COVRecord, error) {
	if len(record) != 6 {
		return nil, fmt.Errorf("expected 6 columns but found %d", len(record))
	}

	if record[0] == "date" {
		return nil, errors.New("header")
	}

	date, err := time.Parse("2006-01-02", record[0])
	if err != nil {
		return nil, err
	}

	newCases, err := strconv.Atoi(record[2])
	if err != nil && len(record[2]) > 0 {
		return nil, err
	}

	newDeaths, err := strconv.Atoi(record[3])
	if err != nil && len(record[3]) > 0 {
		return nil, err
	}

	totalCases, err := strconv.Atoi(record[4])
	if err != nil && len(record[4]) > 0 {
		return nil, err
	}

	totalDeaths, err := strconv.Atoi(record[5])
	if err != nil && len(record[5]) > 0 {
		return nil, err
	}

	return &COVRecord{
		Date:        date,
		Location:    record[1],
		NewCases:    newCases,
		NewDeaths:   newDeaths,
		TotalCases:  totalCases,
		TotalDeaths: totalDeaths,
	}, nil
}

// MARK: String interface methods

func (c COVRecord) String() string {
	return fmt.Sprintf("%-32v %-32s %-12d %-12d %-12d %-12d", c.Date, c.Location, c.NewCases, c.NewDeaths, c.TotalCases, c.TotalDeaths)
}
