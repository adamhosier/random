/* Implements some of the statistical tests described in here
 * http://csrc.nist.gov/publications/nistpubs/800-22-rev1a/SP800-22rev1a.pdf
 *
 * Adam Hosier 2017
 */

package random

import (
  "math"
)

const SIGNIFICANCE = 0.01;

// Tests that the proportion of ones and zeros in [s] are approximately equal
func FrequencyCheck(bs *BitString) bool {
  // Sum over each bit, where a zero is worth -1 and a one is worth 1
  var sum int
  for _, b := range bs.data {
	if b {
	  sum += 1
	} else {
	  sum -= 1
	}
  }
  // Calculate statistic
  s := math.Abs(float64(sum)) / math.Sqrt(float64(bs.length));

  // Calculate P value
  p := math.Erfc(s / math.Sqrt2);

  // if P value < SIGNIFICANCE, the input is non random
  return p > SIGNIFICANCE;
}

