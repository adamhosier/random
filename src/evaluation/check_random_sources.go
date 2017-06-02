package main

import (
	"fmt"
	"github.com/adamhosier/random/src/random"
)

func main() {
	fmt.Println("\nWEBCAM INPUT")
	webcamInput := random.NewInput("input_bin/webcam")
	test(webcamInput, 1000)

	fmt.Println("\nAUDIO INPUT")
	audioInput := random.NewInput("input_bin/audio")
	test(audioInput, 1000)

	fmt.Println("\nINNER PRODUCT EXTRACTOR")
	test(random.NewInnerProductExtractor(webcamInput, audioInput), 1000)

	fmt.Println("\nRANDOM WALK EXTRACTOR")
	sysTimeInput := random.NewInput("input_bin/time")
	test(random.NewRandomWalkExtractor(webcamInput, sysTimeInput), 500)

	fmt.Println("\nPSEUDO-RANDOM EXTRACTOR")
	randomOrgInput := random.NewInput("input_bin/randomorg")
	test(random.NewPseudoRandomExtractor(randomOrgInput.GetBits(64).Int()), 1000)
}

func test(e random.Extractable, numBits int) {
	bits := e.GetBits(numBits)
	passed, results := random.CheckRandom(bits)
	if !passed {
		for _, res := range results {
			fmt.Printf("Test \"%s\"\n", res.Name)
			if res.Result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.P, res.Significance)
			}
		}
	} else {
		fmt.Println("ALL PASSED")
	}
}
