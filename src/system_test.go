package random

import (
	"testing"
	"fmt"
)

func TestSystem(t *testing.T) {
	webcamInput, _ := NewBufferedInput("../input_bin/webcam")
	webcamBits := webcamInput.GetBits(1000)
	webcamPassed, webcamResults := CheckRandom(webcamBits)
	fmt.Printf("WEBCAM INPUT\n---------------\nBitString: %q\n", webcamBits)
	if !webcamPassed {
		for _, res := range webcamResults {
			fmt.Printf("Test \"%s\"\n", res.name)
			if res.result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.p, res.significance)
			}
		}
	}

	fmt.Println()

	audioInput, _ := NewBufferedInput("../input_bin/audio")
	audioBits := audioInput.GetBits(1000)
	audioPassed, audioResults := CheckRandom(audioBits)
	fmt.Printf("AUDIO INPUT\n---------------\nBitString: %q\n", audioBits)
	if !audioPassed {
		for _, res := range audioResults {
			fmt.Printf("Test \"%s\"\n", res.name)
			if res.result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.p, res.significance)
			}
		}
	}

	fmt.Println()

	innerProdExtractor := NewInnerProductExtractor(webcamInput, audioInput)
	innerProdBits := innerProdExtractor.GetBits(1000)
	innerProdPassed, innerProdResults := CheckRandom(innerProdBits)
	fmt.Printf("INNER PRODUCT EXTRACTOR\n--------------------------\nBitString: %q\n", innerProdBits)
	if !innerProdPassed {
		for _, res := range innerProdResults {
			fmt.Printf("Test \"%s\"\n", res.name)
			if res.result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.p, res.significance)
			}
		}
	}
}