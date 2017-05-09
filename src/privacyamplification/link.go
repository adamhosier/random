package main

import "github.com/adamhosier/random/src/random"

type Message struct {
	name string
	value byte
}

// Perfect p2p link, no messages are lost
type PerfectLink struct {
	partner chan Message
}

func NewPerfectLink(partner chan Message) *PerfectLink {
	return &PerfectLink{partner}
}

func (pl *PerfectLink) Send(m Message) {
	pl.partner <- m
}

func (pl *PerfectLink) Receive() Message {
	return <-pl.partner
}

// Simulated lossy link, messges sent over these have a probability [p] of being lost
type LossyLink struct {
	partner chan Message
	p       float64
	rng     random.Generator
}

func NewLossyLink(partner chan Message) *LossyLink {
	return &LossyLink{partner, 0.2, random.NewGeneratorFromConfig("prng")}
}

func (ll *LossyLink) Send(m Message) {
	if ll.rng.NextNormalizedFloat() > ll.p {
		ll.partner <- m
	}
}

func (ll *LossyLink) Receive() Message {
	return <-ll.partner
}
