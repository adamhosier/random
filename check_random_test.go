package random

import (
  "testing"
  "os"
)

type test struct{
  in   *BitString;
  want bool;
}

var (
  // Test case strings
  testStrings = []string{
	"1100100100001111110110101010001000100001011010001100001000110100110001001100011001100010100010111000",
	"0010100110101110101011101101001010101110100101101100101110100101011100001111011010110100101001001001",
	"0000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111",
	"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	"0000000011111111000000000000000000001111110000000011111100000000001111111000000000000000000111111111",
	"0011111100000000000000001111111111111111111111000000000000000000000000000000000000111111111111111000" +
		"0000000001111111111111111000",
	"1100110000010101011011000100110011100000000000100100110101010001000100111101011010000000110101111100" +
		"1100111001101101100010110010",
  }
  // Will hold respective BitStrings for string cases above
  testBitStrings = make([]*BitString, len(testStrings))
)

// Ensures setup code is run before tests start
func TestMain(m *testing.M) {
  // Convert each testString to BitString
  for i, s := range testStrings {
	if bs, err := BitStringFromString(s); err != nil {
	  os.Exit(1)
	} else {
	  testBitStrings[i] = bs
	}
  }
  // Run tests
  code := m.Run()
  os.Exit(code)
}

// TEST CASES START

func TestFrequencyCheck(t *testing.T) {
  cases := []test{
	{ testBitStrings[0], true },
	{ testBitStrings[1], true },
	{ testBitStrings[2], true },
	{ testBitStrings[3], false },
	{ testBitStrings[4], false },
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
	{ testBitStrings[0], true },
	{ testBitStrings[1], true },
	{ testBitStrings[2], false },
	{ testBitStrings[3], false },
	{ testBitStrings[4], false },
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
	{ testBitStrings[0], true },
	{ testBitStrings[1], false },
	{ testBitStrings[2], false },
	{ testBitStrings[3], false },
	{ testBitStrings[4], false },
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
	{ testBitStrings[5], false },
	{ testBitStrings[6], true },
  }
  for _, c := range cases {
	got := LongestRunCheck(c.in)
	if got != c.want {
	  t.Errorf("LongestRunCheck(%q) == %t, want %t", c.in, got, c.want)
	}
  }
}


