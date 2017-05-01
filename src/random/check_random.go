/* Implements some of the NIST randomness statistical tests described in here
 * http://csrc.nist.gov/publications/nistpubs/800-22-rev1a/SP800-22rev1a.pdf
 *
 * Adam Hosier 2017
 */

package random

import (
	"fmt"
	"github.com/adamhosier/random/src/bitstring"
	"math"
)

const significance float64 = 0.01

type testResult struct {
	Name         string
	P            float64
	Significance float64
	Result       bool
}

// Runs all cbecks over [bs] to determine if it appears random or notm returing the result of every check
func CheckRandom(bs *bitstring.BitString) (bool, []testResult) {
	// Templates for templace matching checks
	t1, _ := bitstring.BitStringFromString("010101")
	t2, _ := bitstring.BitStringFromString("001100")
	t3, _ := bitstring.BitStringFromString("000111")

	// Run all tests
	results := []testResult{FrequencyCheck(bs), BlockFrequencyCheck(bs, bs.Length/10), RunsCheck(bs),
		LongestRunCheck(bs), NonOverlappingTemplateMatchingCheck(bs, t1), NonOverlappingTemplateMatchingCheck(bs, t2),
		NonOverlappingTemplateMatchingCheck(bs, t3), SerialCheck(bs), ApproximateEntropyCheck(bs),
		CumulativeSumsCheck(bs)}

	// Calculate overall result
	overallResult := true
	for _, res := range results {
		overallResult = overallResult && res.Result
	}

	return overallResult, results
}

// Tests that the proportion of ones and zeros are approximately equal
func FrequencyCheck(bs *bitstring.BitString) testResult {
	// Sum over each bit, where a zero is worth -1 and a one is worth 1
	var sum int
	for _, b := range bs.Data {
		if b {
			sum++
		} else {
			sum--
		}
	}
	// Calculate statistic
	s := math.Abs(float64(sum)) / math.Sqrt(float64(bs.Length))

	// Calculate P value
	p := math.Erfc(s / math.Sqrt2)

	return testResult{"Frequency", p, significance, p >= significance}
}

// Tests the proportion of ones in each block of size [m]
func BlockFrequencyCheck(bs *bitstring.BitString, M int) testResult {
	// Partition string into blocks of length m
	blocks := bs.Partition(M)

	// Calculate proportion of ones to zeros in each block
	numblocks := bs.Length / M
	pi := make([]float64, numblocks)
	for i, block := range blocks {
		pi[i] = block.Proportion()
	}

	// Calculate chisq statistic
	sum := 0.0
	for _, p := range pi {
		sum += math.Pow((p - 0.5), 2)
	}
	chi := 4.0 * float64(M) * sum

	// find p value using incomplete gamma function
	p := igamc(float64(numblocks)/2.0, chi/2.0)

	return testResult{fmt.Sprintf("Block Frequency (%d)", M), p, significance, p >= significance}
}

// Tests the amount of consecutive ones or zeros over the whole string
func RunsCheck(bs *bitstring.BitString) testResult {
	// Test the proportion of ones to zeros, as this must be valid for runs test to succeed
	pi := bs.Proportion()
	if math.Abs(pi-0.5) >= 2.0/math.Sqrt(10) {
		return testResult{"Consecutive Runs", 0, significance, false}
	}

	// Calculate runs test statistic
	sum := 0
	for i, b := range bs.Data {
		if i < bs.Length-1 && b != bs.Data[i+1] {
			sum++
		}
	}
	s := float64(sum + 1)

	// Calculate p value
	tmp := 2.0 * pi * (1.0 - pi)
	n := float64(bs.Length)
	p := math.Erfc(math.Abs(s-tmp*n) / (tmp * math.Sqrt(2*n)))

	return testResult{"Consecutive Runs", p, significance, p >= significance}
}

// Tests the longest run of ones in each block of size 8, 128 or 10^4
// Assumes length of [bs] must be either 128, 6272 or 750,000
func LongestRunCheck(bs *bitstring.BitString) testResult {
	// Find parameters for different lengths of [bs]
	M, k := 0, 0
	var pi []float64
	if bs.Length < 128 {
		return testResult{"Longest Runs", 0, significance, false}
	} else if bs.Length < 6272 {
		bs = bs.Substring(0, 128)
		M = 8
		k = 4
		pi = []float64{0.2148, 0.3672, 0.2305, 0.1875}
	} else if bs.Length < 750000 {
		bs = bs.Substring(0, 6272)
		M = 128
		k = 6
		pi = []float64{0.1174, 0.2430, 0.2493, 0.1752, 0.1027, 0.1124}
	} else {
		bs = bs.Substring(0, 750000)
		M = 100000
		k = 7
		pi = []float64{0.0882, 0.2092, 0.2483, 0.1933, 0.1208, 0.0675, 0.0727}
	}
	n := bs.Length / M

	// Find longest run in each block, recording the lowest and highest run sizes
	blocks := bs.Partition(M)
	longestRuns := make([]int, n)
	for i, block := range blocks {
		current := 0
		for _, b := range block.Data {
			if b {
				current++
			} else {
				current = 0
			}
			if current > longestRuns[i] {
				longestRuns[i] = current
			}
		}
	}

	// Categorise these run lengths, where each element contains the number of blocks with some run length
	// Also build the theoretical probabilities pi
	v := make([]int, k)
	for _, r := range longestRuns {
		if bs.Length == 128 {
			switch {
			case r <= 1:
				v[0]++
			case r >= 4:
				v[3]++
			default:
				v[r-1]++
			}
		} else if bs.Length == 6272 {
			switch {
			case r <= 4:
				v[0]++
			case r >= 9:
				v[5]++
			default:
				v[r-4]++
			}
		} else {
			switch {
			case r <= 10:
				v[0]++
			case r >= 16:
				v[6]++
			default:
				v[r-10]++
			}
		}
	}

	// Calculate chisq statistic
	chi := 0.0
	for i, p := range v {
		npi := pi[i] * float64(n)
		chi += math.Pow((float64(p)-npi), 2) / npi
	}

	// Calculate p value
	p := igamc(float64(k-1)/2.0, chi/2.0)

	return testResult{"Longest Runs", p, significance, p >= significance}
}

// Tests a [template] BitString against [bs] to find if the template occurs significantly often
func NonOverlappingTemplateMatchingCheck(bs, template *bitstring.BitString) testResult {
	// Set parameters, including theoretical mean and variance
	n := bs.Length
	m := template.Length
	N := 8     // number of blocks is fixed to 8
	M := n / N // block size
	mean := float64(M-m+1) / math.Pow(2.0, float64(m))
	variance := float64(M) * ((1 / math.Pow(2.0, float64(m))) - float64(2*m-1)/math.Pow(2.0, float64(2*m)))

	// Partition the string into blocks
	blocks := bs.Partition(M)
	w := make([]int, N)
	for i, block := range blocks {
		// Scan block for template matches
		for j := 0; j < M-m+1; j++ {
			if block.HasTemplateAt(template, j) {
				w[i]++
				j += m - 2 // skip the remainder of the pattern
			}
		}
	}

	// Compute chisq value
	chi := 0.0
	for _, v := range w {
		chi += math.Pow(float64(v)-mean, 2.0) / variance
	}

	// Compute p value
	p := igamc(float64(N)/2.0, chi/2.0)

	return testResult{fmt.Sprintf("Non Overlapping Templates (%q)", template), p, significance, p >= significance}
}

// Checks every binary block over a few sizes to ensure they don't occur too commonly
// Assumes bs.length >= 16
func SerialCheck(bs *bitstring.BitString) testResult {
	// Set parameters
	n := bs.Length
	m := int(math.Floor(math.Log2(float64(n)))) - 3

	// Extend string by adding m-1 bits from its start to its end
	ebs := bs.Extend(bs.First(m - 1))

	// Get all m, m-1 and m-2 bit blocks, and the number of times the occur in [bs]
	v1 := make([]int, int(math.Pow(2.0, float64(m))))
	for i, t := range bitstring.BitStringsOfLength(m) {
		for j := 0; j < n; j++ {
			if ebs.HasTemplateAt(t, j) {
				v1[i]++
			}
		}
	}
	v2 := make([]int, int(math.Pow(2.0, float64(m-1))))
	for i, t := range bitstring.BitStringsOfLength(m - 1) {
		for j := 0; j < n; j++ {
			if ebs.HasTemplateAt(t, j) {
				v2[i]++
			}
		}
	}
	v3 := make([]int, int(math.Pow(2.0, float64(m-2))))
	for i, t := range bitstring.BitStringsOfLength(m - 2) {
		for j := 0; j < n; j++ {
			if ebs.HasTemplateAt(t, j) {
				v3[i]++
			}
		}
	}

	// Compute psisq values
	sum := 0.0
	for _, v := range v1 {
		sum += math.Pow(float64(v), 2.0)
	}
	psi1 := (math.Pow(2.0, float64(m))/float64(n))*sum - float64(n)
	sum = 0.0
	for _, v := range v2 {
		sum += math.Pow(float64(v), 2.0)
	}
	psi2 := (math.Pow(2.0, float64(m-1))/float64(n))*sum - float64(n)
	sum = 0.0
	for _, v := range v3 {
		sum += math.Pow(float64(v), 2.0)
	}
	psi3 := (math.Pow(2.0, float64(m-2))/float64(n))*sum - float64(n)

	// Compute delta psysq
	dpsi1 := psi1 - psi2
	dpsi2 := psi1 - 2.0*psi2 + psi3

	// Compute p values
	p1 := igamc(math.Pow(2.0, float64(m-2)), dpsi1/2.0)
	p2 := igamc(math.Pow(2.0, float64(m-3)), dpsi2/2.0)

	if !(p1 >= significance) {
		return testResult{"Serial (p1)", p1, significance, false}
	}
	if !(p2 >= significance) {
		return testResult{"Serial (p2)", p2, significance, false}
	}
	return testResult{"Serial", (p1 + p2) / 2.0, significance, true}
}

// Compares all possible bit patterns of a certain length to see if any occur significantly often
func ApproximateEntropyCheck(bs *bitstring.BitString) testResult {
	// Set parameters
	n := bs.Length
	m := int(math.Floor(math.Log2(float64(n)))) - 5

	// Extend string by adding m bits from its start to its end
	ebs := bs.Extend(bs.First(m))

	// Count occurrences of each m bit block in bs
	numBlocks := int(math.Pow(2.0, float64(m)))
	c1 := make([]float64, numBlocks)
	for i := 0; i < n; i++ {
		c1[ebs.Substring(i, m).Int()]++
	}
	for i := 0; i < numBlocks; i++ {
		c1[i] /= float64(n)
	}

	// Compute phi for m bit blocks
	phi1 := 0.0
	for _, c := range c1 {
		if c != 0 {
			phi1 += c * math.Log(c)
		}
	}

	// Count occurrences of each m+1 bit block in bs
	numBlocks = int(math.Pow(2.0, float64(m+1)))
	c2 := make([]float64, numBlocks)
	for i := 0; i < n; i++ {
		c2[ebs.Substring(i, m+1).Int()]++
	}
	for i := 0; i < numBlocks; i++ {
		c2[i] /= float64(n)
	}

	// Compute phi for m+1 bit blocks
	phi2 := 0.0
	for _, c := range c2 {
		if c != 0 {
			phi2 += c * math.Log(c)
		}
	}

	// Compute chisq statistic
	chi := float64(2*n) * (math.Log(2.0) - (phi1 - phi2))

	// Compute p value
	p := igamc(math.Pow(2.0, float64(m-1)), chi/2.0)

	return testResult{"Approximate Entropy", p, significance, p >= significance}
}

// Find consecutive partial sums from the start, checking they don't get too high or low
func CumulativeSumsCheck(bs *bitstring.BitString) testResult {
	// Find partial sums, saving the max in z
	n := bs.Length
	s := make([]int, n)
	z := 0
	for i, b := range bs.Data {
		var v int
		if b {
			v = 1
		} else {
			v = -1
		}
		if i == 0 {
			s[0] = v
		} else {
			s[i] = s[i-1] + v
		}
		sa := int(math.Abs(float64(s[i])))
		if sa > z {
			z = sa
		}
	}

	// Find p value
	sum1, sum2 := 0.0, 0.0
	max := (n/z - 1) / 4
	for k := (-n/z + 1) / 4; k <= max; k++ {
		sqrtn := math.Sqrt(float64(n))
		sum1 += stdNormal(float64((4*k+1)*z)/sqrtn) - stdNormal(float64((4*k-1)*z)/sqrtn)
	}
	for k := (-n/z - 3) / 4; k <= max; k++ {
		sqrtn := math.Sqrt(float64(n))
		sum2 += stdNormal(float64((4*k+3)*z)/sqrtn) - stdNormal(float64((4*k+1)*z)/sqrtn)
	}
	p := 1.0 - sum1 + sum2

	return testResult{"Cumulative Sums", p, significance, p >= significance}
}
