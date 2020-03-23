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

// NewCases returns the new cases at the location.
func (l Location) NewCases() int {
	if len(l.Records) > 0 {
		return l.Records[len(l.Records)-1].NewCases
	}
	return 0
}

// NewDeaths returns the new deaths at the location.
func (l Location) NewDeaths() int {
	if len(l.Records) > 0 {
		return l.Records[len(l.Records)-1].NewDeaths
	}
	return 0
}

// TotalCases returns the total cases at the location.
func (l Location) TotalCases() int {
	if len(l.Records) > 0 {
		return l.Records[len(l.Records)-1].TotalCases
	}
	return 0
}

// TotalDeaths returns the total deaths at the location.
func (l Location) TotalDeaths() int {
	if len(l.Records) > 0 {
		return l.Records[len(l.Records)-1].TotalDeaths
	}
	return 0
}

// TotalCasesSignal returns the location's records' total cases
// as a float slice.
func (l Location) TotalCasesSignal() []float64 {
	var signal []float64
	for _, r := range l.Records {
		signal = append(signal, float64(r.TotalCases))
	}
	return signal
}

// TotalDeathsSignal returns the location's records' total deaths
// as a float slice.
func (l Location) TotalDeathsSignal() []float64 {
	var signal []float64
	for _, r := range l.Records {
		signal = append(signal, float64(r.TotalDeaths))
	}
	return signal
}

// MARK: String interface methods

func (l Location) String() string {
	return fmt.Sprintf("%-32s %-12d %-12d %-12d %-12d", l.Name, l.NewCases(), l.NewDeaths(), l.TotalCases(), l.TotalDeaths())
}
