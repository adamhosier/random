package random

import (
	"testing"
)

var (
	i1 = &MockInput{
		MockGetBits: func(n int) *BitString {
		bs, _ := BitStringFromString("1000000000000000000000000000000000000000000000000000000000000000")
		return bs.Substring(0, n)
		},
	}
	i2 = &MockInput{
		MockGetBits: func(n int) *BitString {
			bs, _ := BitStringFromString("0000000000000000000000000000000000000000000000000000000000000000")
			return bs.Substring(0, n)
		},
	}
	i3 = &MockInput{
		MockGetBits: func(n int) *BitString {
			bs, _ := BitStringFromString("00000000")
			return bs.Substring(0, n)
		},
	}
)

type MockInput struct {
	MockGetBits func(int) *BitString
}

func (i *MockInput) GetBits(n int) *BitString {
	return i.MockGetBits(n)
}

func TestInnerProductExtractor(t *testing.T) {
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

func TestRandomWalkExtractor(t *testing.T) {
	extr := NewRandomWalkExtractor(i1, i3)
	want, _ := BitStringFromString("1000000000000000")
	got := extr.GetBits(16)
	if !got.Equals(want) {
		t.Errorf("RandomWalkExtractor.GetBits(16) == %q, expected %q", got, want)
	}
}
