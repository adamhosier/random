package random

import (
	"github.com/adamhosier/random/src/bitstring"
	"testing"
)

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
	want, _ := bitstring.BitStringFromString("1000000000000000000000000000000000000000000000000000000000000000")
	got := extr.GetBits(64)
	if !got.Equals(want) {
		t.Errorf("RandomWalkExtractor.GetBits(16) == %q, expected %q", got, want)
	}
}

func TestPseudoRandomExtractor(t *testing.T) {
	extr := NewPseudoRandomExtractor(0)
	bs := extr.GetBits(32)
	if bs.Length != 32 {
		t.Errorf("PseudoRandomExtractor.GetBits(32) contained %d bits", bs.Length)
	}
}
