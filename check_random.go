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

  return p > SIGNIFICANCE
}

// Tests the proportion of ones in each block of size [m]
func BlockFrequencyCheck(bs *BitString, m int) bool {
  numblocks := bs.length / m
  // Partition string into blocks of length m
  var pi = make([]float64, numblocks)
  for i := 0; i < numblocks; i++ {
	// Calculate pi value for this block
	var sum int
	for _, b := range bs.data[i*m:(i+1)*m] {
	  if b { sum++ }
	}
	pi[i] = float64(sum) / float64(m)
  }

  // Calculate chisq statistic
  var sum float64
  for _, p := range pi { sum += math.Pow((p - 0.5), 2) }
  chi := 4.0 * float64(m) * sum

  // find p value using incomplete gamma function
  p := igameQ(float64(numblocks) / 2.0, chi / 2.0)

  return p > SIGNIFICANCE
}

