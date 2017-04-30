package random

import (
	"math"
	"github.com/adamhosier/random/src/bitstring"
)

const (
	NUM_ITERATIONS = 100
	EPSILON        = 3.0e-7
)

func AlmostEqual(val, diff float64) bool {
	return math.Abs(diff) <= math.Abs(val*EPSILON)
}

// Uses the modified Lentz method to compute the incomplete Gamma function using a continued fraction
func igamcContinuedFraction(a, x float64) float64 {
	minValue := math.SmallestNonzeroFloat64

	// Start computing terms at n = 1.
	c := math.MaxFloat64
	b := x + 1 - a
	d := 1 / b

	fracEst := d
	for n := 1; n <= NUM_ITERATIONS; n++ {
		an := -float64(n) * (float64(n) - a)
		b += 2.0
		d = an*d + b
		if math.Abs(d) < minValue {
			d = minValue
		}

		c = an/c + b
		if math.Abs(c) < minValue {
			c = minValue
		}

		d = 1 / d
		diff := d * c
		fracEst *= diff

		if AlmostEqual(1, diff-1) {
			lg, _ := math.Lgamma(a)
			return fracEst * math.Exp(-x+a*math.Log(x)-lg)
		}
	}

	// Should never get here
	return 0.0
}

// Computes incomplete Gamma by summing a series
func igamcSeries(a, x float64) float64 {
	an := a
	termVal := 1 / an
	sum := 1 / an

	for n := 0; n < NUM_ITERATIONS; n++ {
		an++
		termVal *= x / float64(an)
		sum += termVal
		if AlmostEqual(sum, termVal) {
			// Note that x^a = exp(a * log(x)).
			lg, _ := math.Lgamma(a)
			return sum * math.Exp(-x+a*math.Log(x)-lg)
		}
	}

	// Should never get here
	return 0.0
}

// Incomplete Gamma function (P)
func igamcP(a, x float64) float64 {
	if x < 0 || a <= 0 {
		return 0.0
	}

	if x < a+1 {
		return igamcSeries(a, x)
	} else {
		// Continued fraction converges more quickly in this range
		return 1 - igamcContinuedFraction(a, x)
	}
}

// Incomplete Gamma function (Q)
func igamc(a, x float64) float64 {
	return 1 - igamcP(a, x)
}

// Standard normal CDF
func stdNormal(x float64) float64 {
	return 0.5 * (math.Erfc(-x / math.Sqrt2))
}


// Mock input structure for testing
type MockInput struct {
	MockGetBits func(int) *bitstring.BitString
}

func (i *MockInput) GetBits(n int) *bitstring.BitString {
	return i.MockGetBits(n)
}

var (
	i1 = &MockInput{
		MockGetBits: func(n int) *bitstring.BitString {
			bs, _ := bitstring.BitStringFromString("1000000000000000000000000000000000000000000000000000000000000000")
			return bs.Substring(0, n)
		},
	}
	i2 = &MockInput{
		MockGetBits: func(n int) *bitstring.BitString {
			bs, _ := bitstring.BitStringFromString("0000000000000000000000000000000000000000000000000000000000000000")
			return bs.Substring(0, n)
		},
	}
	i3 = &MockInput{
		MockGetBits: func(n int) *bitstring.BitString {
			bs, _ := bitstring.BitStringFromString("00000000")
			return bs.Substring(0, n)
		},
	}
	i4 = &MockInput{
		MockGetBits: func(n int) *bitstring.BitString {
			bs, _ := bitstring.BitStringFromString("0111111101111111001101011011101001101110011100101100")
			return bs.Substring(0, n)
		},
	}
)

