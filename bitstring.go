package random

import "errors"

type BitString struct {
  length     int
  data   []bool
}

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

func (bs *BitString) String() string {
  var s string
  for _, b := range bs.data { if b { s += "1"} else { s += "0" } }
  return s
}