package random

import (
	"testing"
	"fmt"
)

func TestSystem(t *testing.T) {
	webcamInput, _ := NewInput("../input_bin/webcam")
	webcamBits := webcamInput.GetBits(1000)
	_, results := CheckRandom(webcamBits)
	fmt.Printf("BitString: %q\n", webcamBits)
	for _, res := range results {
		fmt.Printf("Test \"%s\"\n", res.name)
		if res.result {
			fmt.Println("\tPASSED")
		} else {
			fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.p, res.significance)
		}
	}
}