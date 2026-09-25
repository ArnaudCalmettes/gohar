package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// Run these with -benchmem. The core claims to allocate nothing on
// these paths, and a claim like that decays silently: an iterator whose
// closure starts escaping costs an allocation per call and nothing at
// the call site looks different.
//
// The operations themselves are single machine instructions, so the
// number that matters here is the allocation count, not the nanoseconds.

var (
	benchSet     harmony.PitchSet
	benchClass   harmony.PitchClass
	benchInt     int
	benchBool    bool
	benchSemi    harmony.Semitones
	benchPattern harmony.ScalePattern
)

func benchInput(b *testing.B) harmony.PitchSet {
	b.Helper()
	s, err := harmony.NewPitchSet(0, 2, 4, 5, 7, 9, 11)
	if err != nil {
		b.Fatal(err)
	}
	return s
}

func BenchmarkPitchSetTranspose(b *testing.B) {
	s := benchInput(b)
	for b.Loop() {
		benchSet = s.Transpose(7)
	}
}

func BenchmarkPitchSetLen(b *testing.B) {
	s := benchInput(b)
	for b.Loop() {
		benchInt = s.Len()
	}
}

func BenchmarkPitchSetContains(b *testing.B) {
	s := benchInput(b)
	for b.Loop() {
		benchBool = s.Contains(7)
	}
}

func BenchmarkPitchSetIsSubsetOf(b *testing.B) {
	s := benchInput(b)
	for b.Loop() {
		benchBool = s.IsSubsetOf(harmony.ChromaticPitchSet)
	}
}

// Classes and Rotations are the two range over func iterators on the
// hot path of recognition.
func BenchmarkPitchSetClasses(b *testing.B) {
	s := benchInput(b)
	for b.Loop() {
		for c := range s.Classes() {
			benchClass = c
		}
	}
}

func BenchmarkPitchSetRotations(b *testing.B) {
	s := benchInput(b)
	for b.Loop() {
		for n, rotated := range s.Rotations() {
			benchSemi = n
			benchSet = rotated
		}
	}
}

// Canonical exists to replace the twelve comparisons that Rotations
// implies. Comparing the two here is what tells whether it earns its
// complexity: if the gap is small, Rotations alone is enough and
// Canonical can go.
func BenchmarkPitchSetCanonical(b *testing.B) {
	s := benchInput(b)
	for b.Loop() {
		benchSet, benchSemi = s.Canonical()
	}
}

func BenchmarkChordPatternFold(b *testing.B) {
	for b.Loop() {
		benchSet = harmony.ChordDominantSeventh.Fold()
	}
}

func BenchmarkChordPatternOffsets(b *testing.B) {
	for b.Loop() {
		for n := range harmony.ChordDominantSeventh.Offsets() {
			benchSemi = n
		}
	}
}

func BenchmarkScalePatternMode(b *testing.B) {
	for b.Loop() {
		benchPattern, _ = harmony.ScaleMajor.Mode(6)
	}
}

// From is the generic method that replaced six near identical ones.
// Instantiating it on two different types is worth measuring separately:
// the PitchClass instantiation folds, the Pitch one does not.
func BenchmarkScalePatternFromPitchClass(b *testing.B) {
	for b.Loop() {
		for c := range harmony.ScaleMajor.From(harmony.PitchClass(0)) {
			benchClass = c
		}
	}
}

func BenchmarkScalePatternFromPitch(b *testing.B) {
	for b.Loop() {
		for p := range harmony.ScaleMajor.From(harmony.MiddleC) {
			benchClass = p.Class()
		}
	}
}
