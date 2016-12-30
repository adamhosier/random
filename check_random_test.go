package random

import "testing"

type test struct{
  in   string;
  want bool;
}

func TestFrequencyCheck(t *testing.T) {
  cases := []test{
	{ "00110011",                       true },
	{ "01101100",                       true },
	{ "000000000000000111111111111111", true },
	{ "000000100000100",                false },
  }
  for _, c := range cases {
	if bs, err := BitStringFromString(c.in); err != nil {
	  t.Error(err)
	} else {
	  got := FrequencyCheck(bs)
	  if got != c.want {
		t.Errorf("frequency_check(%q) == %t, want %t", c.in, got, c.want)
	  }
	}
  }
}
