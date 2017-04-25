package random

import "math"

type Extractable interface {
	GetBits(int) *BitString
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
func (e *InnerProductExtractor) GetBits(n int) *BitString {
	bs := &BitString{}
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
	value      *BitString
	neighbours []*randomGraphNode
}

func NewRandomWalkExtractor(i1, i2 Extractable) *RandomWalkExtractor {
    return &RandomWalkExtractor{i1, i2, 32}
}

func (e *RandomWalkExtractor) GetBits(n int) *BitString {
	// Get all possible outputs
	bss := BitStringsOfLength(n)

	// Build a random graph
	gns := make([]randomGraphNode, len(bss))
	for i, bs := range bss {
		gns[i] = randomGraphNode{bs, make([]*randomGraphNode, e.d)}
	}

	// Connect graph using weak input
	for _, gn := range gns {
		for j := 0; j < e.d; j++ {
			gn.neighbours[j] = &gns[e.input1.GetBits(int(math.Log2(float64(len(gns))))).Int()]
		}
	}

	// Find start node from weak input
	start := &gns[e.input1.GetBits(n).Int()]

	// Calculate number of steps to reach a random point
	steps := 10 * int(math.Log2(float64(n)))

	// Randomly walk based on the strong input
	var current *randomGraphNode = start
	for i := 0; i < steps; i++ {
		current = current.neighbours[e.input2.GetBits(int(math.Log2(float64(e.d)))).Int()]
	}

	return current.value
}