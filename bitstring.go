package random

import "errors"

type BitString struct {
  length     int
  data   []bool
}

// Builds a BitString from a string of zeros and ones e.g. "01101001"
func BitStringFromString(s string) (*BitString, error) {
  n := len(s)
  var bs = make([]bool, n)
  for i, c := range s {
	if c == '1' {
	  bs[i] = true
	} else if c == '0' {
	  bs[i] = false
	} else {
	  return nil, errors.New("bitstring: Invalid character in BitStringFromString")
	}
  }
  return &BitString{n, bs}, nil
}

// Gets the value of the bit at position [i] in [bs]
func (bs *BitString) At(i int) bool {
  return bs.data[i]
}

// Partitions [bs] into blocks of length [len] discarding extra bits at the end
func (bs *BitString) Partition(len int) []*BitString {
  bss := make([]*BitString, bs.length / len)
  for i := 0; i < bs.length / len; i++ {
  	bss[i] = &BitString{len, bs.data[i*len:(i+1)*len]}
  }
  return bss
}

// Count number of ones in [bs]
func (bs *BitString) Ones() int {
  var sum int
  for _, b := range bs.data { if b { sum++ } }
  return sum
}

// Calculates the proportion of ones to zeros in [bs]
func (bs *BitString) Proportion() float64 {
  return float64(bs.Ones()) / float64(bs.length)
}

// Tests if [block] matches [template] at position i
func (bs *BitString) HasTemplateAt(template *BitString, i int) bool {
  for j, b := range template.data {
	if bs.At(i + j) != b { return false }
  }
  return true
}

// Converts [bs] to a string of zeros and ones
func (bs *BitString) String() string {
  var s string
  for _, b := range bs.data { if b { s += "1"} else { s += "0" } }
  return s
}