package main

import (
	"fmt"
	"github.com/adamhosier/random/src/bitstring"
)

func Alice(pl *PerfectLink, ll *LossyLink, done chan bool) {
	// Define secret
	secret, err := bitstring.BitStringFromBytes(&[]byte{0xFF, 0xEE, 0xDD, 0xCC})

	if err != nil {
		panic(err)
		done <- true
	}

	// Send secret over the private, lossy channel
	ll.Send(secret)

	done <- true
}

func Bob(pl *PerfectLink, ll *LossyLink, done chan bool) {
	// Wait for secret to be sent from Alice
	secret := ll.Receive()

	fmt.Printf("%X\n", secret.Bytes())
	done <- true
}

func main() {
	done := make(chan bool)
	pl := NewPerfectLink()
	ll := NewLossyLink()
	go Alice(pl, ll, done)
	go Bob(pl, ll, done)

	// Wait for both clients to send finishing signal
	<-done
	<-done
	fmt.Println("DONE")
}

