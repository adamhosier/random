package random

import "testing"

type inputTest struct{
  path       string
  want_error bool
}

func TestGetBits(t *testing.T) {
  cases := []inputTest{
	{ path: "inputs/webcam", want_error: false },
	{ path: "inputs/not_here", want_error: true },
  }
  for _, c := range cases {
	_, err := NewInput(c.path)
	if err == nil {
	  if c.want_error {
		t.Errorf("NewInput(\"%v\") expected an error to be thrown", c.path)
	  }
	} else {
	  if !c.want_error {
		t.Errorf("NewInput(\"%v\") threw an error which wasnt expected", c.path)
	  }
	}
  }
}