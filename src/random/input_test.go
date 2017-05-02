package random

import "testing"

type inputTest struct {
	path      string
	wantError bool
}

func TestGetBits(t *testing.T) {
	cases := []inputTest{
		{path: "../../input_bin/webcam", wantError: false},
	}

	for _, c := range cases {
		i := NewInput(c.path)
		bs := i.GetBits(100)
		if bs.Length != 100 {
			t.Error("GetBits() retreived the wrong amount of bits")
		}
	}
}
