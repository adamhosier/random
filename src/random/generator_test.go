package random

import (
	"testing"
)

func TestGenerator_NextBool(t *testing.T) {
	r := NewGenerator(i1)
	for i := 0; i < 16; i++ {
		want := true
		got := r.NextBool()
		if got != want {
			t.Error("Generator.NextBool == false, wanted true")
		}
	}
}

func TestGenerator_NextFloat64(t *testing.T) {
	r := NewGenerator(i4)
	want := 9.36271761665039004007837775134E-2
	got := r.NextFloat64()
	if got != want {
		t.Errorf("Generator.NextFloat64 == %f, wanted %f", got, want)
	}
}