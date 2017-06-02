package main

import (
	"crypto/rand"
	"fmt"
	"github.com/adamhosier/random/src/random"
	"math/big"
)

func main() {

	NUM_PTS := 100000
	MAX := 50
	M := MAX * MAX

	// crypto.rand
	crypto := make([]int, M)
	for i := 0; i < NUM_PTS; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(MAX)))
		y, _ := rand.Int(rand.Reader, big.NewInt(int64(MAX)))
		crypto[int(x.Int64())*MAX+int(y.Int64())] += 1
	}
	pr(crypto)

	// prng
	prng := make([]int, M)
	prngen := random.NewGeneratorFromConfig("prng")
	for i := 0; i < NUM_PTS; i++ {
		x := prngen.NextIntBetween(0, MAX)
		y := prngen.NextIntBetween(0, MAX)
		prng[x*MAX+y] += 1
	}
	pr(prng)

	// webcam
	camera := make([]int, M)
	cameragen := random.NewGeneratorFromExtractable(random.NewInput("input_bin/webcam"))
	for i := 0; i < NUM_PTS; i++ {
		x := cameragen.NextIntBetween(0, MAX)
		y := cameragen.NextIntBetween(0, MAX)
		camera[x*MAX+y] += 1
	}
	pr(camera)

	// innerprod
	innerprod := make([]int, M)
	innerprodgen := random.NewGeneratorFromConfig("innerprod")
	for i := 0; i < NUM_PTS; i++ {
		x := innerprodgen.NextIntBetween(0, MAX)
		y := innerprodgen.NextIntBetween(0, MAX)
		innerprod[x*MAX+y] += 1
	}
	pr(innerprod)

	// audio
	audio := make([]int, M)
	audiogen := random.NewGeneratorFromExtractable(random.NewInput("input_bin/audio"))
	for i := 0; i < NUM_PTS; i++ {
		x := audiogen.NextIntBetween(0, MAX)
		y := audiogen.NextIntBetween(0, MAX)
		audio[x*MAX+y] += 1
	}
	pr(audio)

	// random walk
	rwalk := make([]int, M)
	rwalkgen := random.NewGeneratorFromConfig("randomwalk")
	for i := 0; i < NUM_PTS; i++ {
		n := rwalkgen.NextIntBetween(0, M)
		rwalk[n] += 1
		if i%1000 == 0 {
			fmt.Printf("%d%% ", i/1000)
		}
	}
	fmt.Println()
	pr(rwalk)
}

func pr(data []int) {
	fmt.Print("(")
	for _, n := range data {
		fmt.Printf("%v, ", n)
	}
	fmt.Print(")\n")
}
