/* Implements some of the statistical tests described in here
 * http://csrc.nist.gov/publications/nistpubs/800-22-rev1a/SP800-22rev1a.pdf
 *
 * Adam Hosier 2017
 */

package random

import (
  "math"
)

const SIGNIFICANCE = 0.01

// Tests that the proportion of ones and zeros are approximately equal
func FrequencyCheck(bs *BitString) bool {
  // Sum over each bit, where a zero is worth -1 and a one is worth 1
  var sum int
  for _, b := range bs.data {
	if b { sum++ } else { sum-- }
  }
  // Calculate statistic
  s := math.Abs(float64(sum)) / math.Sqrt(float64(bs.length))

  // Calculate P value
  p := math.Erfc(s / math.Sqrt2)

  return p >= SIGNIFICANCE
}

// Tests the proportion of ones in each block of size [m]
func BlockFrequencyCheck(bs *BitString, m int) bool {
  numblocks := bs.length / m
  // Partition string into blocks of length m
  var pi = make([]float64, numblocks)
  for i := 0; i < numblocks; i++ {
	// Calculate pi value for this block
	pi[i] = proportion(bs.data[i*m:(i+1)*m], m)
  }

  // Calculate chisq statistic
  var sum float64
  for _, p := range pi { sum += math.Pow((p - 0.5), 2) }
  chi := 4.0 * float64(m) * sum

  // find p value using incomplete gamma function
  p := igameQ(float64(numblocks) / 2.0, chi / 2.0)

  return p >= SIGNIFICANCE
}

func RunsTest(bs *BitString) bool {
  // Test the proportion of ones to zeros, as this must be valid for runs test to succeed
  pi := proportion(bs.data, bs.length)
  if math.Abs(pi - 0.5) >= 2.0 / math.Sqrt(10) { return false }

  // Calculate runs test statistic
  var sum int
  for i, b := range bs.data {
	if i < bs.length - 1 && b != bs.data[i + 1] { sum++ }
  }
  s := float64(sum + 1)

  // Calculate p value
  tmp := 2.0 * pi * (1.0 - pi)
  n := float64(bs.length)
  p := math.Erfc(math.Abs(s - tmp * n) / (tmp * math.Sqrt(2 * n)))

  return p >= SIGNIFICANCE
}