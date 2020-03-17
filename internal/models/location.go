package models

import "fmt"

// Location types contain a set of records of a location.
type Location struct {

	// The location's name.
	Name string

	// The location's records.
	Records []*COVRecord
}

// MARK: Initializers

// NewLocation creates and returns a new location.
func NewLocation(name string, records []*COVRecord) *Location {
	return &Location{
		Name:    name,
		Records: records,
	}
}

// MARK: Exported methods

// TotalCases returns the total cases at the location.
func (l Location) TotalCases() int {
	if len(l.Records) > 0 {
		return l.Records[len(l.Records)-1].TotalCases
	}
	return 0
}

// TotalDeaths returns the total cases at the location.
func (l Location) TotalDeaths() int {
	if len(l.Records) > 0 {
		return l.Records[len(l.Records)-1].TotalCases
	}
	return 0
}

// MARK: String interface methods

func (l Location) String() string {
	return fmt.Sprintf("%-32s %-12d %-12d", l.Name, l.TotalCases(), l.TotalDeaths())
}
