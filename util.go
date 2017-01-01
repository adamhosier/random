package random

import (
  "math"
)

const (
  NUM_ITERATIONS = 100
  EPSILON = 3.0e-7
)

func AlmostEqual(val, diff float64) bool {
  return math.Abs(diff) <= math.Abs(val * EPSILON)
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
	d = an * d + b
	if math.Abs(d) < minValue { d = minValue }

	c = an / c + b
	if math.Abs(c) < minValue { c = minValue }

	d = 1 / d
	diff := d*c
	fracEst *= diff

	if AlmostEqual(1, diff - 1) {
	  lg , _ := math.Lgamma(a)
	  return fracEst * math.Exp(-x + a * math.Log(x) - lg)
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
	  return sum * math.Exp(-x + a * math.Log(x) - lg)
	}
  }

  // Should never get here
  return 0.0
}

// Incomplete Gamma function (P)
func igamcP(a, x float64) float64 {
  if x < 0 || a <= 0 { return 0.0 }

  if x < a + 1 {
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
