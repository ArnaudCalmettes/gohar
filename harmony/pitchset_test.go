package harmony_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sink keeps the compiler from eliminating the loops in the allocation
// test below.
var sink harmony.PitchClass

func mustSet(t *testing.T, classes ...harmony.PitchClass) harmony.PitchSet {
	t.Helper()
	s, err := harmony.NewPitchSet(classes...)
	require.NoError(t, err)
	return s
}

// The high bits being clear is what makes == the right comparison for
// this type. Every operation has to preserve it, so every operation is
// checked here rather than in its own test.
func TestPitchSetHighBitsStayClear(t *testing.T) {
	s := mustSet(t, 0, 4, 7)
	other := mustSet(t, 1, 5, 8, 11)

	ops := map[string]harmony.PitchSet{
		"union":                   s.Union(other),
		"intersect":               s.Intersect(other),
		"difference":              s.Difference(other),
		"with":                    s.With(11),
		"without":                 s.Without(0),
		"transpose up":            s.Transpose(5),
		"transpose down":          s.Transpose(-5),
		"transpose past an octave": s.Transpose(25),
		"chromatic":               harmony.ChromaticPitchSet,
	}

	for name, got := range ops {
		t.Run(name+" leaves bits twelve and above clear", func(t *testing.T) {
			assert.Zero(t, uint16(got)&0xf000)
		})
	}
}

func TestPitchSetEqualityFollowsMembership(t *testing.T) {
	t.Run("order of construction does not matter", func(t *testing.T) {
		assert.Equal(t, mustSet(t, 0, 4, 7), mustSet(t, 7, 0, 4))
	})

	t.Run("duplicates do not matter", func(t *testing.T) {
		assert.Equal(t, mustSet(t, 0, 4, 7), mustSet(t, 0, 4, 4, 7, 7))
	})

	t.Run("the zero value is the empty set", func(t *testing.T) {
		var s harmony.PitchSet
		assert.Equal(t, harmony.EmptyPitchSet, s)
		assert.True(t, s.IsEmpty())
		assert.Zero(t, s.Len())
	})
}

func TestNewPitchSetRejectsOutOfRangeClasses(t *testing.T) {
	_, err := harmony.NewPitchSet(0, 4, 12)
	assert.Error(t, err,
		"class twelve is out of range and folding it would hide a caller "+
			"who confused a distance with a class",
	)
}

// Transpose is a rotation over twelve bits, not a shift. A shift would
// drop members silently, so the note count is the property that catches
// it.
func TestPitchSetTransposeIsARotation(t *testing.T) {
	s := mustSet(t, 0, 4, 7, 11)

	t.Run("a full octave is the identity", func(t *testing.T) {
		assert.Equal(t, s, s.Transpose(12))
	})

	t.Run("moving up then back down returns the original", func(t *testing.T) {
		assert.Equal(t, s, s.Transpose(5).Transpose(-5))
	})

	t.Run("the note count survives every distance", func(t *testing.T) {
		for n := harmony.Semitones(-12); n <= 12; n++ {
			assert.Equal(t, s.Len(), s.Transpose(n).Len(),
				"no member may fall off the end at distance %d", n)
		}
	})

	t.Run("a member that passes the top re-enters at the bottom", func(t *testing.T) {
		got := mustSet(t, 11).Transpose(1)
		assert.Equal(t, mustSet(t, 0), got)
	})

	t.Run("the chromatic set is its own transposition", func(t *testing.T) {
		for n := harmony.Semitones(0); n < 12; n++ {
			assert.Equal(t, harmony.ChromaticPitchSet, harmony.ChromaticPitchSet.Transpose(n))
		}
	})
}

func TestPitchSetRotations(t *testing.T) {
	s := mustSet(t, 0, 4, 7)

	t.Run("there are twelve rotations and the first is the identity", func(t *testing.T) {
		var distances []harmony.Semitones
		first := true
		for n, got := range s.Rotations() {
			if first {
				assert.Equal(t, harmony.Semitones(0), n)
				assert.Equal(t, s, got)
				first = false
			}
			assert.Equal(t, s.Transpose(n), got,
				"each pair holds the distance that produced it")
			distances = append(distances, n)
		}
		assert.Len(t, distances, 12)
	})
}

// Canonical turns twelve comparisons into one lookup. The property that
// makes it usable as a key is that all twelve rotations agree on it.
func TestPitchSetCanonical(t *testing.T) {
	s := mustSet(t, 0, 4, 7)
	want, _ := s.Canonical()

	t.Run("every rotation shares one canonical form", func(t *testing.T) {
		for _, rotated := range s.Rotations() {
			got, _ := rotated.Canonical()
			assert.Equal(t, want, got)
		}
	})

	t.Run("the returned distance reaches the canonical form", func(t *testing.T) {
		for _, rotated := range s.Rotations() {
			got, n := rotated.Canonical()
			assert.Equal(t, got, rotated.Transpose(n))
		}
	})

	t.Run("the canonical form is its own canonical form", func(t *testing.T) {
		got, n := want.Canonical()
		assert.Equal(t, want, got)
		assert.Equal(t, harmony.Semitones(0), n)
	})

	t.Run("the empty set is canonical", func(t *testing.T) {
		got, n := harmony.EmptyPitchSet.Canonical()
		assert.Equal(t, harmony.EmptyPitchSet, got)
		assert.Equal(t, harmony.Semitones(0), n)
	})
}

func TestPitchSetClasses(t *testing.T) {
	t.Run("members come out in ascending numeric order", func(t *testing.T) {
		var got []harmony.PitchClass
		for c := range mustSet(t, 7, 0, 4, 11).Classes() {
			got = append(got, c)
		}
		assert.Equal(t, []harmony.PitchClass{0, 4, 7, 11}, got)
	})

	t.Run("the count matches Len", func(t *testing.T) {
		s := mustSet(t, 1, 3, 6, 8, 10)
		n := 0
		for range s.Classes() {
			n++
		}
		assert.Equal(t, s.Len(), n)
	})

	t.Run("the empty set yields nothing", func(t *testing.T) {
		for range harmony.EmptyPitchSet.Classes() {
			t.Fatal("the empty set yielded a member")
		}
	})
}

func TestPitchSetSubset(t *testing.T) {
	triad := mustSet(t, 0, 4, 7)
	seventh := mustSet(t, 0, 4, 7, 10)

	t.Run("a triad is a subset of the seventh chord built on it", func(t *testing.T) {
		assert.True(t, triad.IsSubsetOf(seventh))
	})

	t.Run("the seventh chord is not a subset of its triad", func(t *testing.T) {
		assert.False(t, seventh.IsSubsetOf(triad))
	})

	t.Run("every set is a subset of itself", func(t *testing.T) {
		assert.True(t, triad.IsSubsetOf(triad))
	})

	t.Run("the empty set is a subset of everything", func(t *testing.T) {
		assert.True(t, harmony.EmptyPitchSet.IsSubsetOf(triad))
	})
}

// The core documents that it does not allocate on its hot paths. That
// is a property to verify, not an intention to write in a comment, and
// iterators are the place where an escaping closure would break it
// without anything looking wrong at the call site.
func TestPitchSetDoesNotAllocate(t *testing.T) {
	s := mustSet(t, 0, 2, 4, 5, 7, 9, 11)

	t.Run("iterating the members allocates nothing", func(t *testing.T) {
		avg := testing.AllocsPerRun(100, func() {
			for c := range s.Classes() {
				sink = c
			}
		})
		assert.Zero(t, avg)
	})

	t.Run("iterating the rotations allocates nothing", func(t *testing.T) {
		avg := testing.AllocsPerRun(100, func() {
			for _, rotated := range s.Rotations() {
				sink = harmony.PitchClass(rotated.Len())
			}
		})
		assert.Zero(t, avg)
	})
}
