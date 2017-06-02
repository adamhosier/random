package main

import (
	"fmt"
	"github.com/adamhosier/random/src/random"
)

func main() {
	r := random.NewGeneratorFromConfig("prng")
	max := 100
	data := make([]int, max)
	for i := 0; i < 100000; i++ {
		n := r.NextIntBetween(0, max)
		data[n]++
	}
	for i := 0; i < max; i++ {
		fmt.Printf("%d: %d\n", i, data[i])
	}
}
