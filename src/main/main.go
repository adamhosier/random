package main

import (
	"github.com/adamhosier/random/src/random"
)

func main() {
	r := random.NewGeneratorFromConfig("prng")
	max := 10
	data := make([]int, max)
	for i := 0; i < 10000; i++ {
		n := r.NextIntBetween(0, max)
		data[n]++
	}
	for i := 0; i < max; i++ {
		println(data[i])
	}
}

/*
func main() {
	fmt.Println("\nWEBCAM INPUT")
	webcamInput := random.NewInput("input_bin/audio")
	test(webcamInput, 1000)

	fmt.Println("\nAUDIO INPUT")
	audioInput := random.NewInput("input_bin/audio")
	test(audioInput, 1000)

	fmt.Println("\nINNER PRODUCT EXTRACTOR")
	test(random.NewInnerProductExtractor(webcamInput, audioInput), 1000)

	fmt.Println("\nRANDOM WALK EXTRACTOR")
	randomOrgInput := random.NewInput("input_bin/randomorg")
	test(random.NewRandomWalkExtractor(audioInput, randomOrgInput), 500)

	fmt.Println("\nPSEUDO-RANDOM EXTRACTOR")
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
*/
