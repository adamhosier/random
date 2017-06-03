package main

import (
	"fmt"
	"github.com/adamhosier/random/src/random"
)

func main() {
	//fmt.Print("(")

	rng := random.NewGeneratorFromConfig("prng")

	rng.NextInt()
	for i := 0; i < 1000; i++ {
		//fmt.Printf("%f, ", rng.NextNormalizedFloat())
		fmt.Printf("%d, ", i)
	}

	//fmt.Print(")\n")
}
