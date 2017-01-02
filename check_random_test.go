package random

import (
  "testing"
  "os"
)

type test struct{
  in       *BitString;
  template *BitString;
  want     bool;
}

var (
  // Test case strings
  testStrings = []string{
	"1100100100001111110110101010001000100001011010001100001000110100110001001100011001100010100010111000",
	"0010100110101110101011101101001010101110100101101100101110100101011100001111011010110100101001001001",
	"0000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111",
	"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	"0000000011111111000000000000000000001111110000000011111100000000001111111000000000000000000111111111",
	"1010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
	"0011111100000000000000001111111111111111111111000000000000000000000000000000000000111111111111111000" +
		"0000000001111111111111111000",
	"1100110000010101011011000100110011100000000000100100110101010001000100111101011010000000110101111100" +
		"1100111001101101100010110010",
  }
  templateStrings = []string{
	"100001", "010101", "110011", "111111",
  }
  // Will hold respective BitStrings for string cases above
  testBitStrings = make([]*BitString, len(testStrings))
  templateBitStrings = make([]*BitString, len(templateStrings))
)

// Ensures setup code is run before tests start
func TestMain(m *testing.M) {
  // Convert each testString and templateString to BitString
  for i, s := range testStrings {
	if bs, err := BitStringFromString(s); err != nil {
	  os.Exit(1)
	} else {
	  testBitStrings[i] = bs
	}
  }
  for i, s := range templateStrings {
	if bs, err := BitStringFromString(s); err != nil {
	  os.Exit(1)
	} else {
	  templateBitStrings[i] = bs
	}
  }
  // Run tests
  code := m.Run()
  os.Exit(code)
}

// TEST CASES START

func TestFrequencyCheck(t *testing.T) {
  cases := []test{
	{ in: testBitStrings[0], want: true },
	{ in: testBitStrings[1], want: true },
	{ in: testBitStrings[2], want: true },
	{ in: testBitStrings[3], want: false },
	{ in: testBitStrings[4], want: false },
	{ in: testBitStrings[5], want: true },
	{ in: testBitStrings[6], want: true },
	{ in: testBitStrings[7], want: true },
  }
  for _, c := range cases {
	got := FrequencyCheck(c.in)
	if got != c.want {
	  t.Errorf("FrequencyCheck(%q) == %t, want %t", c.in, got, c.want)
	}
  }
}

func TestBlockFrequencyCheck(t *testing.T) {
  cases := []test{
	{ in: testBitStrings[0], want: true },
	{ in: testBitStrings[1], want: true },
	{ in: testBitStrings[2], want: false },
	{ in: testBitStrings[3], want: false },
	{ in: testBitStrings[4], want: false },
	{ in: testBitStrings[5], want: true },
	{ in: testBitStrings[6], want: false },
	{ in: testBitStrings[7], want: true },
  }
  for _, c := range cases {
	got := BlockFrequencyCheck(c.in, 10)
	if got != c.want {
	  t.Errorf("BlockFrequencyCheck(%q) == %t, want %t", c.in, got, c.want)
	}

  }
}

func TestRunsCheck(t *testing.T) {
  cases := []test{
	{ in: testBitStrings[0], want: true },
	{ in: testBitStrings[1], want: false },
	{ in: testBitStrings[2], want: false },
	{ in: testBitStrings[3], want: false },
	{ in: testBitStrings[4], want: false },
	{ in: testBitStrings[5], want: false },
	{ in: testBitStrings[6], want: false },
	{ in: testBitStrings[7], want: true },
  }
  for _, c := range cases {
	got := RunsCheck(c.in)
	if got != c.want {
	  t.Errorf("RunsCheck(%q) == %t, want %t", c.in, got, c.want)
	}
  }
}

func TestLongestRunCheck(t *testing.T) {
  cases := []test{
	{ in: testBitStrings[6], want: false },
	{ in: testBitStrings[7], want: true },
  }
  for _, c := range cases {
	got := LongestRunCheck(c.in)
	if got != c.want {
	  t.Errorf("LongestRunCheck(%q) == %t, want %t", c.in, got, c.want)
	}
  }
}

func TestNonOverlappingTemplateMatchingCheck(t *testing.T) {
  cases := []test{
	{ testBitStrings[0], templateBitStrings[0], true },
	{ testBitStrings[0], templateBitStrings[1], true },
	{ testBitStrings[0], templateBitStrings[2], true },
	{ testBitStrings[0], templateBitStrings[3], true },
	{ testBitStrings[1], templateBitStrings[0], true },
	{ testBitStrings[1], templateBitStrings[1], true },
	{ testBitStrings[1], templateBitStrings[2], true },
	{ testBitStrings[1], templateBitStrings[3], true },
	{ testBitStrings[2], templateBitStrings[0], true },
	{ testBitStrings[2], templateBitStrings[1], true },
	{ testBitStrings[2], templateBitStrings[2], true },
	{ testBitStrings[2], templateBitStrings[3], false },
	{ testBitStrings[5], templateBitStrings[0], true },
	{ testBitStrings[5], templateBitStrings[1], false },
	{ testBitStrings[5], templateBitStrings[2], true },
	{ testBitStrings[5], templateBitStrings[3], true },
  }
  for _, c := range cases {
	got := NonOverlappingTemplateMatchingCheck(c.in, c.template)
	if got != c.want {
	  t.Errorf("NonOverlappingTemplateMatchingCheck(%q, %q) == %t, want %t", c.in, c.template, got, c.want)
	}
  }
}

func TestSerialCheck(t *testing.T) {
  cases := []test{
	{ in: testBitStrings[0], want: true },
	{ in: testBitStrings[1], want: false },
	{ in: testBitStrings[2], want: false },
	{ in: testBitStrings[3], want: false },
	{ in: testBitStrings[4], want: false },
	{ in: testBitStrings[5], want: false },
	{ in: testBitStrings[6], want: false },
	{ in: testBitStrings[7], want: true },
  }
  for _, c := range cases {
	got := SerialCheck(c.in)
	if got != c.want {
	  t.Errorf("SerialCheck(%q) == %t, want %t", c.in, got, c.want)
	}
  }
}

func TestApproximateEntropyCheck(t *testing.T) {
  cases := []test{
	{ in: testBitStrings[0], want: true },
	{ in: testBitStrings[1], want: false },
	{ in: testBitStrings[2], want: false },
	{ in: testBitStrings[3], want: false },
	{ in: testBitStrings[4], want: false },
	{ in: testBitStrings[5], want: false },
	{ in: testBitStrings[6], want: false },
	{ in: testBitStrings[7], want: true },
  }
  for _, c := range cases {
	got := ApproximateEntropyCheck(c.in)
	if got != c.want {
	  t.Errorf("ApproximateEntropyCheck(%q) == %t, want %t", c.in, got, c.want)
	}
  }
}

func TestCumulativeSumsCheck(t *testing.T) {
  cases := []test{
	{ in: testBitStrings[0], want: true },
	{ in: testBitStrings[1], want: true },
	{ in: testBitStrings[2], want: false },
	{ in: testBitStrings[3], want: false },
	{ in: testBitStrings[4], want: false },
	{ in: testBitStrings[5], want: true },
	{ in: testBitStrings[6], want: true },
	{ in: testBitStrings[7], want: true },
  }
  for _, c := range cases {
	got := CumulativeSumsCheck(c.in)
	if got != c.want {
	  t.Errorf("CumulativeSumsCheck(%q) == %t, want %t", c.in, got, c.want)
	}
  }
}



