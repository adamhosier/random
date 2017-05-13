package main

import (
	"fmt"
	"github.com/adamhosier/random/src/bitstring"
	"github.com/adamhosier/random/src/random"
)

func Hash(msg, a, b *bitstring.BitString) *bitstring.BitString {
	return msg.BinaryMul(a).BinaryAdd(b)
}

func Alice(pl *PerfectLink, ll *LossyLink, done chan bool) {
	// Define secret
	secret, err := bitstring.BitStringFromBytes(&[]byte{0xFF, 0xEE, 0xDD, 0xCC})
	if err != nil {
		panic(err)
		done <- true
	}

	// Send secret over the private, lossy channel
	ll.Send(secret)

	// Hash the secret, using values securely generated
	rng := random.NewGeneratorFromConfig("innerprod")
	a := rng.GetBits(secret.Length)
	b := rng.GetBits(secret.Length)
	hash := Hash(secret, a, b)

	// Broadcast the randomly generated BitStrings, and the hash value
	pl.Send(a)
	pl.Send(b)
	pl.Send(hash)

	fmt.Printf("Alice: 0x%X\n", secret.Bytes())

	done <- true
}

func Bob(pl *PerfectLink, ll *LossyLink, done chan bool) {
	// Wait for secret to be sent from Alice
	secret := ll.Receive()

	// Receive the hash information
	a := pl.Receive()
	b := pl.Receive()
	otherHash := pl.Receive()
	hash := Hash(secret, a, b)

	fmt.Printf("Bob: 0x%X\n", secret.Bytes())

	if hash.Equals(otherHash) {
		fmt.Println("Bob has received the correct secret, the value is accepted")
	} else {
		fmt.Println("Bob received an incorrect hash")
	}

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

