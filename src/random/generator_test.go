package random

import (
	"testing"
)

func TestGenerator_NextBool(t *testing.T) {
	r := NewGeneratorFromExtractable(i1)
	for i := 0; i < 16; i++ {
		want := true
		got := r.NextBool()
		if got != want {
			t.Error("Generator.NextBool == false, wanted true")
		}
	}
}

func TestGenerator_NextFloat64(t *testing.T) {
	r := NewGeneratorFromExtractable(i4)
	want := 9.36271761665039004007837775134E-2
	got := r.NextFloat64()
	if got != want {
		t.Errorf("Generator.NextFloat64 == %f, wanted %f", got, want)
	}
}

func TestGeneratorConfig(t *testing.T) {
	r := NewGenerator()
	if r == nil {
		t.Error("NewGenerator returned nil")
	}
	r = NewGeneratorFromConfig("prng")
	if r == nil {
		t.Error("NewGeneratorFromConfig(prng) returned nil")
	}
	r = NewGeneratorFromConfig("example")
	if r == nil {
		t.Error("NewGeneratorFromConfig(prng) returned nil")
	}
}
