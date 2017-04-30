package random

import "math"

type Generator struct {
	e Extractable
}

func NewGenerator(e Extractable) *Generator {
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
			fraction += math.Pow(2, float64(-(i+1)))
		}
	}
	exponent := math.Pow(2, float64(g.next(11) - 1023))
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
