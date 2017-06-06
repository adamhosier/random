package random

import (
	"github.com/adamhosier/random/src/bitstring"
	"math"
)

type Extractable interface {
	GetBits(int) *bitstring.BitString
}

const defaultBlockSize int = 8

type InnerProductExtractor struct {
	input1    Extractable
	input2    Extractable
	blockSize int         // Number of blocks to compute the inner product over
}

// Creates a new inner product extractor combining [i1] and [i2]
func NewInnerProductExtractor(i1, i2 Extractable) *InnerProductExtractor {
	return &InnerProductExtractor{i1, i2, defaultBlockSize}
}

// Gets a BitString of length [n] containing the inner product over GF(2) of two inputs
func (e *InnerProductExtractor) GetBits(n int) *bitstring.BitString {
	bs := bitstring.BitStringOfLength(n)

	// Get a list of blocks containing [n] bits
	input1 := e.input1.GetBits(e.blockSize*n).Partition(n)
	input2 := e.input2.GetBits(e.blockSize*n).Partition(n)

	// Compute the inner product over these bits
	for i := 0; i < e.blockSize; i++ {
		bs = bs.BinaryAdd(input1[i].BinaryMul(input2[i]))
	}
	return bs
}

// Random walk extractor
type RandomWalkExtractor struct {
	input1 Extractable // Fast, weak random input
	input2 Extractable // Slow, strong random input
	d      int         // Number of neighbours in each GraphNode, must be in [2^n | n],
}

type randomGraphNode struct {
	value      *bitstring.BitString
	neighbours []*randomGraphNode
}

func NewRandomWalkExtractor(i1, i2 Extractable) *RandomWalkExtractor {
	return &RandomWalkExtractor{i1, i2, 8}
}

// This function used to generate the entire random graph, then randomly traverse it. Below is an implementation of
// an optimised function, that lazily generates graph nodes as they are needed for a huge performance improvement.
func (e *RandomWalkExtractor) GetBits(n int) *bitstring.BitString {
	// Cache of already constructed nodes, initially empty
	gns := make(map[string]randomGraphNode)
	rng := NewPseudoRandomExtractor(e.input1.GetBits(64).Int())

	// Calculate number of steps to reach a random point
	steps := 2 * int(math.Log2(float64(n)))

	// Get start node from weak input
	startBits := e.input1.GetBits(n)
	gns[startBits.Hash()] = randomGraphNode{startBits, make([]*randomGraphNode, e.d)}

	// Lazily traverse tree
	current := gns[startBits.Hash()]
	for i := 0; i < steps; i++ {
		// Lazily construct a (weak) random graph around the selected node if it hasn't been visited yet
		if current.neighbours[0] == nil {
			for j := 0; j < e.d; j++ {
				// Get bits from weak input
				bits := rng.GetBits(n)

				// If this node hasn't been visited, add it to the graph then set it as a neighbour
				if _, exists := gns[bits.Hash()]; !exists {
					gns[bits.Hash()] = randomGraphNode{bits, make([]*randomGraphNode, e.d)}
				}
				node := gns[bits.Hash()]
				current.neighbours[j] = &node
			}
		}

		// Select one of these neighbours using the strong generator
		next := e.input2.GetBits(int(math.Log2(float64(e.d)))).Int()
		current = *current.neighbours[next]
	}

	return current.value
}

// Pseudo-random extractor (used for PRNG)
type PseudoRandomExtractor struct {
	seed int
}

func NewPseudoRandomExtractor(seed int) *PseudoRandomExtractor {
	return &PseudoRandomExtractor{(seed ^ 0x5DEECE66D) & (1<<48 - 1)}
}

func (e *PseudoRandomExtractor) GetBits(n int) *bitstring.BitString {
	// The linear congruential prng gets up to 32 bits at a time
	result := bitstring.NewBitString()
	for i := 0; i < n; i += 32 {
		numBits := int(math.Min(float64(32), float64(n-i)))

		// update seed
		e.seed = (e.seed*0x5DEECE66D + 0xB) & ((1 << 48) - 1)

		// return next value in the sequence
		b := e.seed >> uint(48-numBits)

		// Add these bits to our result
		result = result.Extend(bitstring.BitStringFromInt(numBits, b))
	}

	return result
}
