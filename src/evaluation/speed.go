package main

import (
	"github.com/adamhosier/random/src/random"
	"fmt"
	crand "crypto/rand"
	"time"
	"math/rand"
)

const (
	NUM_BITS = 100000
)

func main() {
	// crypto.rand
	start := time.Now()
	crand.Reader.Read(make([]byte, NUM_BITS))
	crpyto := time.Since(start).Nanoseconds() / 1000
	fmt.Printf("Crpyto: %dms\n", crpyto)

	// go prng
	start = time.Now()
	rand.Read(make([]byte, NUM_BITS))
	goPrng := time.Since(start).Nanoseconds() / 1000
	fmt.Printf("Go PRNG: %dms\n", goPrng)


	// prng
	prng := time_gen(random.NewGeneratorFromConfig("prng"))
	fmt.Printf("Our PRNG: %dms\n", prng)

	// webcam
	camE := random.NewInput("input_bin/webcam")
	camE.GetBits(1)
	camera := time_gen(random.NewGeneratorFromExtractable(camE))
	fmt.Printf("Camera: %dms\n", camera)

	// audio
	micE := random.NewInput("input_bin/audio")
	micE.GetBits(1)
	audio := time_gen(random.NewGeneratorFromExtractable(micE))
	fmt.Printf("Audio: %dms\n", audio)

	// innerprod
	innerprodE := random.NewInnerProductExtractor(camE, camE)
	innerprod := time_gen(random.NewGeneratorFromExtractable(innerprodE))
	fmt.Printf("Innerprod: %dms\n", innerprod)

	// random walk
	rwalk := time_gen(random.NewGeneratorFromExtractable(random.NewRandomWalkExtractor(camE, innerprodE)))
	fmt.Printf("Rwalk: %dms\n", rwalk)

}

func time_gen(g *random.Generator) int64 {
	start := time.Now()
	g.GetBits(NUM_BITS)
	elapsed := time.Since(start)

	return elapsed.Nanoseconds() / 1000
}

