package random

import "testing"

type inputTest struct {
	path      string
	wantError bool
}

func TestNewInput(t *testing.T) {
	cases := []inputTest{
		{path: "../input_bin/webcam", wantError: false},
		{path: "../input_bin/not_here", wantError: true},
	}
	for _, c := range cases {
		_, err := NewInput(c.path)
		if err == nil {
			if c.wantError {
				t.Errorf("NewInput(\"%v\") expected an error to be thrown", c.path)
			}
		} else {
			if !c.wantError {
				t.Errorf("NewInput(\"%v\") threw an error which wasnt expected", c.path)
			}
		}
	}
}

func TestGetBits(t *testing.T) {
	cases := []inputTest{
		{path: "../input_bin/webcam", wantError: false},
	}

	for _, c := range cases {
		i, _ := NewInput(c.path)
		bs := i.GetBits(100)
		if bs.Length != 100 {
			t.Error("GetBits() retreived the wrong amount of bits")
		}
	}
}
