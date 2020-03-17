package commands

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/superhawk610/bar"
)

// NewBarWithTitle creates a new bar with the given title and number of ticks.
func NewBarWithTitle(title string, n int) *bar.Bar {
	return bar.NewWithOpts(
		bar.WithDimensions(n, 30),
		bar.WithFormat(fmt.Sprintf("%s :eta :bar :percent", title)),
		bar.WithDisplay(
			"",
			"◼︎",
			"◼︎",
			" ",
			"",
		),
	)
}

// NewSpinnerWithTitle creates a new spinner with the given title.
func NewSpinnerWithTitle(title string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 50*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s", title)
	return s
}
