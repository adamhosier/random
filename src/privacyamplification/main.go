package main

import (
	"fmt"
	"github.com/adamhosier/random/src/bitstring"
	"github.com/adamhosier/random/src/random"
)

func Hash(msg, a, b *bitstring.BitString) *bitstring.BitString {
	return msg.BinaryMul(a).BinaryAdd(b)
}

// Shuffles the BitString [bs] given a seed to a prng [seed]
func Shuffle(bs *bitstring.BitString, seed int) {
	rng := random.NewGeneratorFromExtractable(random.NewPseudoRandomExtractor(seed))
	for i := 0; i < bs.Length-2; i++ {
		j := rng.NextIntBetween(i, bs.Length)
		bs.Data[i], bs.Data[j] = bs.Data[j], bs.Data[i]
	}
}

func AliceParts(pl *PerfectLink, parts []*bitstring.BitString) bool {
	for _, part := range parts {
		parity := part.Ones() % 2
		pl.Send(bitstring.BitStringFromInt(1, parity))
		success := pl.Receive()
		if success == 0 {
			// try subblocks
		}
	}
}

func BobParts(pl *PerfectLink, parts []*bitstring.BitString) bool {
	for _, part := range parts {
		parity := pl.Receive().Int()
		// If parity differs, notify Alice
		if parity != part.Ones()%2 {
			pl.Send0()
			// call subblocks
		} else {
			pl.Send1()
		}
	}
}

func Alice(pl *PerfectLink, ll *LossyLink, done chan bool) {
	// Randomly generate secret
	rng := random.NewGeneratorFromConfig("innerprod")
	secret := rng.GetBits(512)

	// Send secret over the private, lossy channel
	fmt.Printf("Alice: Sending the secret 0x%X over the private channel\n", secret.Bytes())
	ll.Send(secret)

	// Hash the secret, using values securely generated
	a := rng.GetBits(secret.Length)
	b := rng.GetBits(secret.Length)
	hash := Hash(secret, a, b)

	// Broadcast the randomly generated BitStrings, and the hash value
	fmt.Printf("Alice: Sending the hash 0x%X with values 0x%X and 0x%X over the public channel\n",
		hash.Bytes(), a.Bytes(), b.Bytes())
	pl.Send(a)
	pl.Send(b)
	pl.Send(hash)

	status := pl.Receive().Int()
	if status == 1 {
		fmt.Printf("Alice: Notified of success, secret agreed as 0x%X\n", secret.Bytes())
		done <- true
		return
	}

	fmt.Println("Alice: Notified of bit-twiddling failure")
	// We expect there to be one error every 8 bits, so partition our secret into blocks of length 8
	Shuffle(secret, a.Int())
	parts := secret.Partition(32)
	AliceParts(pl, parts)
	done <- true
}

func Bob(pl *PerfectLink, ll *LossyLink, done chan bool) {
	// Wait for secret to be sent from Alice
	secret := ll.Receive()
	fmt.Printf("Bob: Received the secret 0x%X\n", secret.Bytes())

	// Receive the hash information
	a := pl.Receive()
	b := pl.Receive()
	otherHash := pl.Receive()
	hash := Hash(secret, a, b)
	fmt.Printf("Bob: Calculated hash as 0x%X\n", hash.Bytes())

	// Compare the received hash, to our own calculated hash
	if hash.Equals(otherHash) {
		fmt.Println("Bob: Received the correct secret, the value is accepted")
		// Notify Alice of success
		pl.Send1()
		done <- true
		return
	}

	fmt.Println("Bob: received an incorrect hash, trying single bit-twiddling")
	rng := random.NewGeneratorFromConfig("innerprod")
	// Try 100 random bit-twiddles
	for i := 0; i < 100; i++ {
		n := rng.NextIntBetween(0, secret.Length)
		secretCandidate := secret.Copy()
		secretCandidate.Invert(n)
		hashCandidate := Hash(secretCandidate, a, b)
		// Check if this modification fixes the error
		if hashCandidate.Equals(otherHash) {
			fmt.Printf("Bob: Bit twiddling corrected the secret to 0x%X\n", secretCandidate.Bytes())
			// Notify Alice of completion by sending a 1
			pl.Send1()
			done <- true
			return
		}
	}

	// Notify Alice of failure by sending a 0
	pl.Send0()

	fmt.Println("Bob: Bit twiddling failed, continuing privacy amplification")
	Shuffle(secret, a.Int())
	parts := secret.Partition(32)
	BobParts(pl, parts)

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
