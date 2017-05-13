package main

import (
	"github.com/adamhosier/random/src/random"
	"github.com/adamhosier/random/src/bitstring"
)

const (
	defaultBitCorruptionRate = 0.1 // Rate at which bits are flipped when sending over a lossy link [0..1]
)
// Perfect p2p link, no messages are lost
type PerfectLink struct {
	ch chan *bitstring.BitString
}

func NewPerfectLink() *PerfectLink {
	return &PerfectLink{make(chan *bitstring.BitString)}
}

func (pl *PerfectLink) Send(bs *bitstring.BitString) {
	pl.ch <- bs
}

func (pl *PerfectLink) Receive() *bitstring.BitString {
	return <-pl.ch
}

// Simulated lossy link, messges sent over these have a probability [p] of being lost
type LossyLink struct {
	ch  chan *bitstring.BitString
	p   float64
	rng *random.Generator
}

func NewLossyLink() *LossyLink {
	return &LossyLink{make(chan *bitstring.BitString), defaultBitCorruptionRate, random.NewGeneratorFromConfig("prng")}
}

func (ll *LossyLink) Send(bs *bitstring.BitString) {
	// Simulate corruption of message
	cp := bs.Copy()
	for i := 0; i < bs.Length; i++ {
		if ll.rng.NextNormalizedFloat() < ll.p {
			cp.Invert(i)
		}
	}
	ll.ch <- cp
}

func (ll *LossyLink) Receive() *bitstring.BitString {
	return <-ll.ch
}
