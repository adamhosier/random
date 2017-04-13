package random

type Extractor interface {
	GetBits(int) *BitString
}

const defaultBlockSize int = 32

type InnerProductExtractor struct {
	input1    Input
	input2    Input
	blockSize int   // Size of subsequences to take the inner product over to produce one bit
}

// Creates a new inner product extractor combining [i1] and [i2]
func NewInnerProductExtractor(i1, i2 Input) *InnerProductExtractor {
	return &InnerProductExtractor{i1, i2, defaultBlockSize}
}

// Gets a BitString of length [n] containing the inner product over GF(2) of two inputs
func (e *InnerProductExtractor) GetBits(n int) *BitString {
	bs := &BitString{}
	for i := 0; i < n; i++ {
		b1 := e.input1.GetBits(e.blockSize)
		b2 := e.input2.GetBits(e.blockSize)
		bs.Add(b1.InnerProduct(b2) % 2 == 1)
	}
	return bs
}