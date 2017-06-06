package main

import (
	"fmt"
	"github.com/adamhosier/random/src/bitstring"
	"github.com/adamhosier/random/src/random"
)

const (
	VERBOSE = true
	SECRET_LEN = 512
	PARTITION_SIZE = 16
)

func Log(format string, args ...interface{}) {
	if VERBOSE {
		fmt.Printf(format, args...)
	}
}

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

// Recursively performs the error correction phase of the protocol, takes a list of partitioned BitStrings and will
// return the corrected string
func AliceParts(pl *PerfectLink, parts []*bitstring.BitString) *bitstring.BitString {
	// BitString to store the corrected result
	result := bitstring.NewBitString()

	for _, part := range parts {
		// If part.Length <= 1, we have found the error (base case)
		if part.Length > 1 {
			// Calculate and transmit parity information to Bob
			parity := part.Ones() % 2
			pl.Send(bitstring.BitStringFromInt(1, parity))

			// Receive information on whether Bob' parity matches
			success := pl.Receive().Int()
			if success == 0 {
				// On failure, recurse to localise error
				result = result.Extend(AliceParts(pl, part.PartitionExtra(part.Length/2)))
			} else {
				// On success, add the result to the array
				result = result.Extend(part.Substring(0, part.Length - 1))
			}
		}
	}
	return result
}

func BobParts(pl *PerfectLink, parts []*bitstring.BitString) *bitstring.BitString {
	// BitString to store the corrected result
	result := bitstring.NewBitString()

	for _, part := range parts {
		// If part.Length <= 1, we have found the error (base case)
		if part.Length > 1 {
			// Receive parity information for this block from Alice
			parity := pl.Receive().Int()

			if parity != part.Ones()%2 {
				// If parity differs, notify Alice of failure and recurse
				pl.Send0()
				result = result.Extend(BobParts(pl, part.PartitionExtra(part.Length/2)))
			} else {
				pl.Send1()
				// Else notify Alice of success and update our result
				result = result.Extend(part.Substring(0, part.Length - 1))
			}
		}
	}
	return result
}

func Alice(pl *PerfectLink, ll *LossyLink, done chan *bitstring.BitString) {
	// Randomly generate secret
	rng := random.NewGeneratorFromConfig("innerprod")
	secret := rng.GetBits(SECRET_LEN)

	// Send secret over the private, lossy channel
	Log("Alice: Sending the secret 0x%X over the private channel\n", secret.Bytes())
	ll.Send(secret)

	// Hash the secret, using values securely generated
	a := rng.GetBits(secret.Length)
	b := rng.GetBits(secret.Length)
	hash := Hash(secret, a, b)

	// Broadcast the randomly generated BitStrings, and the hash value
	Log("Alice: Sending the hash 0x%X with values 0x%X and 0x%X over the public channel\n",
		hash.Bytes(), a.Bytes(), b.Bytes())
	pl.Send(a)
	pl.Send(b)
	pl.Send(hash)

	status := pl.Receive().Int()
	if status == 1 {
		Log("Alice: Notified of success, secret agreed as 0x%X\n", secret.Bytes())
		done <- secret
		return
	}

	// Repeat error correction attempts until a secret is agreed
	secretAgreed := false
	round := 1
	for !secretAgreed {
		// Attempt to fix errors
		Shuffle(secret, a.Int())
		parts := secret.PartitionExtra(PARTITION_SIZE * round)
		secret = AliceParts(pl, parts)
		Log("Alice: Error correction attempt: 0x%X\n", secret.Bytes())

		// Recalculate hash to verify errors
		a = a.Substring(0, secret.Length)
		b = b.Substring(0, secret.Length)
		hash = Hash(secret, a, b)

		// Send this hash to Bob to validate it
		pl.Send(hash)
		secretAgreed = pl.Receive().Int() == 1
		round++
	}

	Log("Alice: Secret has been agreed to be 0x%X\n", secret.Bytes())

	done <- secret
}

func Bob(pl *PerfectLink, ll *LossyLink, done chan *bitstring.BitString) {
	// Wait for secret to be sent from Alice
	secret := ll.Receive()
	Log("Bob: Received the secret  0x%X\n", secret.Bytes())

	// Receive the hash information
	a := pl.Receive()
	b := pl.Receive()
	hashA := pl.Receive()
	hashB := Hash(secret, a, b)
	Log("Bob: Calculated hash as 0x%X\n", hashB.Bytes())

	// Compare the received hash, to our own calculated hash
	if hashA.Equals(hashB) {
		Log("Bob: Received the correct secret, the value is accepted")
		// Notify Alice of success
		pl.Send1()
		done <- secret
		return
	}

	Log("Bob: received an incorrect hash, trying single bit-twiddling\n")
	rng := random.NewGeneratorFromConfig("innerprod")
	// Try 100 random bit-twiddles
	for i := 0; i < 100; i++ {
		n := rng.NextIntBetween(0, secret.Length)
		secretCandidate := secret.Copy()
		secretCandidate.Invert(n)
		hashCandidate := Hash(secretCandidate, a, b)
		// Check if this modification fixes the error
		if hashCandidate.Equals(hashA) {
			Log("Bob: Bit twiddling corrected the secret to 0x%X\n", secretCandidate.Bytes())
			// Notify Alice of completion by sending a 1
			pl.Send1()
			done <- secretCandidate
			return
		}
	}

	// Notify Alice of failure by sending a 0
	pl.Send0()

	Log("Bob: Bit twiddling failed, continuing privacy amplification\n")

	// Repeat error correction attempts until secret is agreed
	secretAgreed := false
	round := 1
	for !secretAgreed {
		// Attempt to fix errors
		Shuffle(secret, a.Int())
		parts := secret.PartitionExtra(PARTITION_SIZE*round)
		secret = BobParts(pl, parts)
		Log("Bob: Error correction attempt:   0x%X\n", secret.Bytes())

		// Recalculate hash
		a = a.Substring(0, secret.Length)
		b = b.Substring(0, secret.Length)
		hashA = pl.Receive()
		hashB = Hash(secret, a, b)

		// Notify Alice if all errors are fixed
		secretAgreed = hashA.Equals(hashB)
		if !secretAgreed {
			pl.Send0()
		}
		round++
	}
	pl.Send1()
	Log("Bob: Secret has been agreed to be   0x%X (len %d)\n", secret.Bytes(), secret.Length)

	done <- secret
}

func main() {
	done := make(chan *bitstring.BitString)
	pl := NewPerfectLink()
	ll := NewLossyLink()
	go Alice(pl, ll, done)
	go Bob(pl, ll, done)

	// Wait for both clients to send finishing signal
	<-done
	<-done
	fmt.Println("DONE")
}