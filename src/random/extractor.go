package random

import (
	"math"
	"github.com/adamhosier/random/src/bitstring"
)

type Extractable interface {
	GetBits(int) *bitstring.BitString
}

const defaultBlockSize int = 8

type InnerProductExtractor struct {
	input1    Extractable
	input2    Extractable
	blockSize int // Size of subsequences to take the inner product over to produce one bit
}

// Creates a new inner product extractor combining [i1] and [i2]
func NewInnerProductExtractor(i1, i2 Extractable) *InnerProductExtractor {
	return &InnerProductExtractor{i1, i2, defaultBlockSize}
}

// Gets a BitString of length [n] containing the inner product over GF(2) of two inputs
func (e *InnerProductExtractor) GetBits(n int) *bitstring.BitString {
	bs := &bitstring.BitString{}
	for i := 0; i < n; i++ {
		b1 := e.input1.GetBits(e.blockSize)
		b2 := e.input2.GetBits(e.blockSize)
		bs.Add(b1.InnerProduct(b2)%2 == 1)
	}
	return bs
}



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
    return &RandomWalkExtractor{i1, i2, 32}
}

// This function used to generate the entire random graph, then randomly traverse it. Below is an implementation of
// an optimised function, that lazily generates graph nodes as they are needed for a huge performance improvement.
func (e *RandomWalkExtractor) GetBits(n int) *bitstring.BitString {
	// Cache of already constructed nodes, initially empty
	gns := make(map[string]randomGraphNode)

	// Calculate number of steps to reach a random point
	steps := 10 * int(math.Log2(float64(n)))

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
				bits := e.input1.GetBits(n)

				// If this node hasn't been visited, add it to the graph then set it as a neighbour
				if _, exists := gns[bits.Hash()]; !exists {
					gns[bits.Hash()] = randomGraphNode{startBits, make([]*randomGraphNode, e.d)}
				}
				node := gns[bits.Hash()]
				current.neighbours[j] = &node
			}
		}

		// Select one of these neighbours using the strong generator
		n = e.input2.GetBits(int(math.Log2(float64(e.d)))).Int()
		current = *current.neighbours[n]
	}

	return current.value
}