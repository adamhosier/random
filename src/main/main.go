package main

import (
	"fmt"
	"github.com/adamhosier/random/src/random"
	"log"
)

func main() {
	webcamInput, err := random.NewInput("input_bin/webcam")
	if err != nil {
		log.Fatal("Error initialising the webcam input")
	}
	webcamBits := webcamInput.GetBits(1000)
	webcamPassed, webcamResults := random.CheckRandom(webcamBits)
	fmt.Printf("WEBCAM INPUT\n---------------\nBitString: %q\n", webcamBits)
	if !webcamPassed {
		for _, res := range webcamResults {
			fmt.Printf("Test \"%s\"\n", res.Name)
			if res.Result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.P, res.Significance)
			}
		}
	}

	fmt.Println()

	audioInput, err := random.NewInput("input_bin/audio")
	if err != nil {
		log.Fatal("Error initialising the audio input")
	}
	audioBits := audioInput.GetBits(1000)
	audioPassed, audioResults := random.CheckRandom(audioBits)
	fmt.Printf("AUDIO INPUT\n---------------\nBitString: %q\n", audioBits)
	if !audioPassed {
		for _, res := range audioResults {
			fmt.Printf("Test \"%s\"\n", res.Name)
			if res.Result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.P, res.Significance)
			}
		}
	}

	fmt.Println()

	innerProdExtractor := random.NewInnerProductExtractor(webcamInput, audioInput)
	innerProdBits := innerProdExtractor.GetBits(500)
	innerProdPassed, innerProdResults := random.CheckRandom(innerProdBits)
	fmt.Printf("INNER PRODUCT EXTRACTOR\n--------------------------\nBitString: %q\n", innerProdBits)
	if !innerProdPassed {
		for _, res := range innerProdResults {
			fmt.Printf("Test \"%s\"\n", res.Name)
			if res.Result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.P, res.Significance)
			}
		}
	}

	fmt.Println()

	randomOrgInput, err := random.NewInput("input_bin/randomorg")
	randomWalkExtractor := random.NewRandomWalkExtractor(audioInput, randomOrgInput)
	randomWalkBits := randomWalkExtractor.GetBits(128)
	randomWalkPassed, randomWalkResults := random.CheckRandom(randomWalkBits)
	fmt.Printf("RANDOM WALK EXTRACTOR\n--------------------------\nBitString: %q\n", randomWalkBits)
	if !randomWalkPassed {
		for _, res := range randomWalkResults {
			fmt.Printf("Test \"%s\"\n", res.Name)
			if res.Result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.P, res.Significance)
			}
		}
	}

	fmt.Println()

	prng := random.NewPseudoRandomExtractor(randomOrgInput.GetBits(64).Int())
	prngBits := prng.GetBits(1000)
	prngPassed, prngResults := random.CheckRandom(prngBits)
	fmt.Printf("PSEUDO-RANDOM EXTRACTOR\n--------------------------\nBitString: %q\n", prngBits)
	if !prngPassed {
		for _, res := range prngResults {
			fmt.Printf("Test \"%s\"\n", res.Name)
			if res.Result {
				fmt.Println("\tPASSED")
			} else {
				fmt.Printf("\tFAILED (p = %g at %f significance)\n", res.P, res.Significance)
			}
		}
	}

}
