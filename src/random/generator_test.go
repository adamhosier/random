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

func TestGenerator_NextIntBetween(t *testing.T) {
	r := NewGenerator()
	f := func(min, max int) {
		for i := 0; i < 10 * (max - min); i++ {
			n := r.NextIntBetween(min, max)
			if n < min || n >= max {
				t.Errorf("Generator.NextIntBetween(%d, %d) produced value outside the range %d", min, max, n)
			}
		}
	}
	f(0, 100)
	f(-10, 10)
	f(100, 200)
}

func TestGenerator_NextNormalizedFloat(t *testing.T) {
	r := NewGenerator()
	for i := 0; i < 10000; i++ {
		n := r.NextNormalizedFloat()
		if n < 0 || n > 1 {
			t.Errorf("Generator.NextNormalizedFloat produced value outside the range [0,1]: %f", n)
		}
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
