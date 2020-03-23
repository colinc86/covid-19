// Package commands conatins the commands for the medina command line application.
package commands

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/colinc86/covid-19/internal/models"
	"github.com/colinc86/go-genetics"
	"github.com/urfave/cli/v2"
)

// PredictCommandHandler handles list commands.
type PredictCommandHandler struct {
	Name        string
	Aliases     []string
	Usage       string
	Description string

	// MARK: Private properties
	location string
	days     uint
	signal   []float64
}

// MARK: Initializers

// NewPredictCommandHandler creates and returns a new list command handler.
func NewPredictCommandHandler() *PredictCommandHandler {
	return &PredictCommandHandler{
		Name:    "predict",
		Aliases: []string{"p"},
		Usage:   "Predicts future values in the data set.",
		Description: `Predicts future values in the dataset by extrapolating
		a Sigmoid curve.
		
		Examples:
			# Predict world data for the next day
			covid19 predict data
			
			# Predict data for a specific location for the next day
			covid19 predict data -l [location]
			
			# Predict data number days out
			covid19 predict data -d [number]`,
	}
}

// MARK: Public methods

// Command creates and returns the handler's command.
func (h *PredictCommandHandler) Command() *cli.Command {
	return &cli.Command{
		Name:        h.Name,
		Aliases:     h.Aliases,
		Usage:       h.Usage,
		Description: h.Description,
		Subcommands: []*cli.Command{
			&cli.Command{
				Name:    "data",
				Aliases: []string{"d"},
				Action:  h.PredictDataSetAction,
				Usage:   "The COVID-19 dataset.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "location",
						Aliases:     []string{"l"},
						Usage:       "Filter by location.",
						Required:    false,
						Destination: &h.location,
					},
					&cli.UintFlag{
						Name:        "days",
						Aliases:     []string{"d"},
						Usage:       "Days of prediction.",
						Required:    false,
						Value:       1,
						Destination: &h.days,
					},
				},
			},
		},
	}
}

// PredictDataSetAction predicts a value from the dataset.
func (h *PredictCommandHandler) PredictDataSetAction(c *cli.Context) error {
	// Validate days
	if h.days < 1 {
		h.days = 1
	}

	// Get the data set
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

	// Get the current series
	var totalCases []float64
	if len(h.location) > 0 {
		totalCases = world.TotalCasesSignalForLocation(h.location)
		// totalDeaths = world.TotalDeathsSignalForLocation(h.location)
	} else {
		totalCases = world.TotalCasesSignal()
		// totalDeaths = world.TotalDeathsSignal()
	}

	// Get sigmoid function coefficients and solve
	h.signal = totalCases
	casesCoefficients := h.analyzeSignal("cases", totalCases)

	// Print the current bars and predicted past bars
	for i, actualValue := range totalCases {
		predictedValue := h.sigmoidFunction(casesCoefficients, i)

		bar := ""
		ticks := int(math.Ceil(float64(actualValue) / (totalCases[len(totalCases)-1] / 40.0)))
		for i := 0; i < ticks; i++ {
			bar += "#"
		}

		pbar := ""
		pticks := int(math.Ceil(float64(predictedValue) / (totalCases[len(totalCases)-1] / 40.0)))
		for i := 0; i < pticks; i++ {
			pbar += "+"
		}

		fmt.Printf("          %-12d %s\n", int(actualValue), bar)
		fmt.Printf("%-9d %-12d %s\n", -1*len(totalCases)+i+1, int(predictedValue), pbar)
	}

	// Print the predicted future bars
	for i := 0; i < int(h.days); i++ {
		predictedValue := h.sigmoidFunction(casesCoefficients, len(totalCases)+i)

		pbar := ""
		pticks := int(math.Ceil(float64(predictedValue) / (totalCases[len(totalCases)-1] / 40.0)))
		for i := 0; i < pticks; i++ {
			pbar += "+"
		}

		fmt.Printf("\n%-9d %-12d %s\n", i+1, int(predictedValue), pbar)
	}

	fmt.Printf("coeff: %v\n", casesCoefficients)

	// h.signal = totalDeaths
	// deathsCoefficients := h.analyzeSignal("deaths", totalDeaths)
	// deathsPrediction := h.sigmoidFunction(deathsCoefficients, len(totalDeaths)+h.days)

	return nil
}

// MARK: Unexported methods

// analyzeCases analyzes the signal.
func (h *PredictCommandHandler) analyzeSignal(name string, signal []float64) []float64 {
	s := NewSpinnerWithTitle(fmt.Sprintf("Analyzing %s...", name))
	s.Start()
	defer s.Stop()

	// Create our evolver configuration
	config := genetics.NewEvolverConfiguration(
		genetics.NewSelectionMethod(genetics.SelectionMethodTypeTournament),
		genetics.NewCrossoverMethod(genetics.CrossoverMethodTypePoint, 1),
		1,
		0.5,
		0.2,
	)

	// Generate a population of chromosomes
	population := genetics.GeneratePopulation(5, 3, func(i, j int) float64 {
		if j == 0 {
			return 1000000 * rand.Float64()
		} else if j == 1 {
			return rand.Float64()
		}
		return 365.0 * rand.Float64()
	})

	// Create our evolver
	evolver := genetics.NewEvolver(config, h.fitnessFunction, h.mutationFunction)

	// Evolve until we meet our desired error rate
	var fittestChromosome *genetics.Chromosome
	count := 0
	evolver.Evolve(population, func(c *genetics.EvolverConfiguration, pop genetics.Population) bool {
		fittestChromosome = pop[len(pop)-1]
		count++
		return count < 2000000
	})

	return fittestChromosome.Genes
}

// fitnessFunction checks for the amount of error in the s-curve defined
// by the coefficients given by the chromosome.
func (h *PredictCommandHandler) fitnessFunction(chromosome *genetics.Chromosome) float64 {
	if len(h.signal) == 0 {
		return 0.0
	}

	return 1.0 / h.totalError(chromosome.Genes)
}

// totalError calculates and returns the total error of the genes against the
// signal.
func (h *PredictCommandHandler) totalError(genes []float64) float64 {
	if len(h.signal) == 0 {
		return 0.0
	}

	totalErr := 0.0
	for i, v := range h.signal {
		totalErr += math.Abs(h.sigmoidFunction(genes, i) - v)
	}
	return totalErr / float64(len(h.signal))
}

// mutationFunction mutates the given gene in the chromosome.
func (h *PredictCommandHandler) mutationFunction(chromosome *genetics.Chromosome, gene int) float64 {
	step := 1.0
	if gene == 0 {
		step = 10.0
	} else if gene == 1 {
		step = 0.0001
	} else if gene == 2 {
		step = 1.0
	}

	sign := rand.Intn(2)

	currentValue := chromosome.Genes[gene]
	if sign == 0 {
		if currentValue-step < 0.0 {
			currentValue = 0.0
		} else {
			currentValue -= step
		}
	} else {
		value := 0.0
		if gene == 0 {
			value = 1000000
		} else if gene == 1 {
			value = 1.0
		} else if gene == 2 {
			value = 365.0
		}

		if currentValue+step > value {
			currentValue = value
		} else {
			currentValue += step
		}
	}

	return currentValue
}

// sigmoidFunction calculates the value of the sigmoid function with the
// given coefficients at x.
func (h *PredictCommandHandler) sigmoidFunction(c []float64, x int) float64 {
	return c[0] / (1 + math.Exp(-1.0*c[1]*float64(x+1)+c[2]))
}
