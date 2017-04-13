package random

import (
	"testing"
)

type MockInput struct {
	MockGetBits func(int) *BitString
}

func (i *MockInput) GetBits(n int) *BitString {
	return i.MockGetBits(n)
}

func TestInnerProductExtractor(t *testing.T) {
	i1 := &MockInput{
		MockGetBits: func(n int) *BitString {
			bs, _ := BitStringFromString("1000000000000000000000000000000000000000000000000000000000000000")
			return bs.Substring(0, n)
		},
	}
	i2 := &MockInput{
		MockGetBits: func(n int) *BitString {
			bs, _ := BitStringFromString("0000000000000000000000000000000000000000000000000000000000000000")
			return bs.Substring(0, n)
		},
	}

	extr := NewInnerProductExtractor(i1, i2)
	if extr.GetBits(1).At(0) {
		t.Error("InnerProductExtractor.GetBits(1) == 1, expected 0")
	}

	extr = NewInnerProductExtractor(i1, i1)
	if !extr.GetBits(1).At(0) {
		t.Error("InnerProductExtractor.GetBits(0) == 0, expected 1")
	}

	extr = NewInnerProductExtractor(i2, i2)
	if extr.GetBits(1).At(0) {
		t.Error("InnerProductExtractor.GetBits(1) == 1, expected 0")
	}
}
