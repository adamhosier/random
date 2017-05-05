package random

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"runtime"
	"path"
)

// Generator config structure
type GeneratorConfig struct {
	Extractor *ExtractableConfig `json:"extractor"`
}

type ExtractableConfig struct {
	Type          string             `json:"type"`
	Seed          int                `json:"seed"`
	SeedGenerator *ExtractableConfig `json:"seedGenerator"`
	Path          string             `json:"path"`
	Input1        *ExtractableConfig `json:"input1"`
	Input2        *ExtractableConfig `json:"input2"`
}

type Generator struct {
	e Extractable
}

// Creates a random number genertor using the configuration defined in 'default.json'
func NewGenerator() *Generator {
	return NewGeneratorFromConfig("default")
}

// Creates a random number generator using the configuration defined at [path]
func NewGeneratorFromConfig(name string) *Generator {
	// Allow config to be read cross-package
	_, thisfile, _, _ := runtime.Caller(0)
	filepath := fmt.Sprintf("%s/config/%s.json", path.Dir(thisfile), name)

	// Read config file
	file, err := os.Open(filepath)
	if err != nil {
		panic("NewGenerator: Could not load config file")
	}
	var config GeneratorConfig
	json.NewDecoder(file).Decode(&config)

	// Build the generator from the config
	return NewGeneratorFromExtractable(configureExtractable(*config.Extractor))
}

func configureExtractable(config ExtractableConfig) Extractable {
	switch config.Type {
	case "pseudorandom":
		if config.SeedGenerator != nil {
			return NewPseudoRandomExtractor(configureExtractable(*config.SeedGenerator).GetBits(64).Int())
		} else {
			return NewPseudoRandomExtractor(config.Seed)
		}
	case "input":
		_, thisfile, _, _ := runtime.Caller(0)
		return NewInput(fmt.Sprintf("%s/%s", path.Dir(path.Dir(path.Dir(thisfile))), config.Path))
	case "innerproduct":
		return NewInnerProductExtractor(configureExtractable(*config.Input1), configureExtractable(*config.Input2))
	case "randomwalk":
		return NewRandomWalkExtractor(configureExtractable(*config.Input1), configureExtractable(*config.Input2))
	default:
		panic("NewGenerator: Invalid generator config (extractable object)")
	}
}

// Creates a random number generator using a user defined configuration
func NewGeneratorFromExtractable(e Extractable) *Generator {
	return &Generator{e}
}

// Gets a bool from the extractable
func (g *Generator) NextBool() bool {
	return g.e.GetBits(1).At(0)
}

// Gets a 64 bit float following IEEE 754 double-precision binary format
func (g *Generator) NextFloat64() float64 {
	sign := g.next(1)
	fraction := 1.0
	for i, b := range g.e.GetBits(52).Data {
		if b {
			fraction += math.Pow(2, float64(-(i + 1)))
		}
	}
	exponent := math.Pow(2, float64(g.next(11)-1023))
	return math.Pow(-1, float64(sign)) * fraction * exponent
}

// Gets a 64 bit integer
func (g *Generator) NextInt() int {
	return g.next(64)
}

// Gets an integer consisting of n bits of randomness, with n < 64
func (g *Generator) next(n int) int {
	return g.e.GetBits(n).Int()
}
