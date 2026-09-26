package harmony

import "fmt"

// A Tonality is a tonal context: a tonic together with a heptatonic
// pattern.
//
// It exists only where invoking a key has meaning, which is where one
// can speak of degrees, of functions and of cadences. A whole tone or
// octatonic pattern has no tonality, and that is not a degraded case to
// repair: it is an answer.
//
// # Why the fields are unexported
//
// Every other value type in this package exposes its representation,
// because the invariant is cheap to state and the ceremony is not worth
// it. Tonality is the exception: its whole reason to exist is that the
// pattern has seven notes, and an exported field would let a caller
// build a five note Tonality with a composite literal, defeating the
// point. Construction goes through [NewTonality].
//
// Comparison with == still works, and the zero value is still
// meaningful as "no tonality".
//
// # Where it is used
//
// Two packages consume it, which is why it lives in the core rather
// than in naming. The naming package spells by degree, and analysis
// reasons in functions and cadences. Neither owns it.
type Tonality struct {
	tonic   PitchClass
	pattern ScalePattern
}

// NewTonality builds a tonal context.
//
// Fails when `p` is not heptatonic. That failure is the type doing its
// job, not an obstacle: a caller holding a pentatonic pattern has
// learned that no tonal reading applies, which is information the
// analysis engine should propagate rather than swallow.
func NewTonality(tonic PitchClass, p ScalePattern) (Tonality, error) {
	if !tonic.IsValid() {
		return Tonality{}, fmt.Errorf("gohar: tonic %d is out of range", tonic)
	}
	if !p.IsHeptatonic() {
		return Tonality{}, fmt.Errorf(
			"gohar: a %d note pattern has no tonality", p.Len())
	}
	if !p.Contains(0) {
		return Tonality{}, fmt.Errorf("gohar: pattern does not contain its own tonic")
	}
	return Tonality{tonic: tonic, pattern: p}, nil
}

// IsZero reports whether `t` is the zero value, meaning no tonal context.
//
// The zero value is not C major and must never be treated as one. An
// engine that has inferred nothing yet, and a player working outside
// any key, both produce a zero Tonality, and an interface that renders
// it as a named key lies to the player. Spelling in the absence of a
// tonality is a convention that belongs to naming, not a default value
// smuggled in here.
func (t Tonality) IsZero() bool {
	return t.pattern == 0
}

// Tonic returns the tonic of `t`.
func (t Tonality) Tonic() PitchClass {
	return t.tonic
}

// Pattern returns the pattern of `t`, always heptatonic.
func (t Tonality) Pattern() ScalePattern {
	return t.pattern
}

// Scale widens `t` to a [Scale], losing the heptatonic guarantee.
//
// One way only. Narrowing back goes through [NewTonality], which
// checks.
func (t Tonality) Scale() Scale {
	return Scale{Tonic: t.tonic, Pattern: t.pattern}
}

// Set returns the classes belonging to `t`.
func (t Tonality) Set() PitchSet {
	if t.IsZero() {
		return EmptyPitchSet
	}
	return t.pattern.At(t.tonic)
}

// Contains reports whether `c` belongs to `t`.
func (t Tonality) Contains(c PitchClass) bool {
	return t.Set().Contains(c)
}

// Degree returns the class sitting at degree `d`, and whether `d` is in
// range. For a Tonality the range is always 1 through 7.
func (t Tonality) Degree(d Degree) (PitchClass, bool) {
	if t.IsZero() {
		return 0, false
	}
	return t.Scale().Degree(d)
}

// DegreeOf returns the degree that `c` occupies in `t`, and whether `c`
// belongs to `t` at all.
//
// This is the operation both consumers are built on. Spelling by degree
// needs it to assign one letter per rank, and functional analysis needs
// it to recognise that a chord is built on the fifth degree. A class
// outside `t` has no degree, and the false return is the honest answer:
// an accidental note is not a degree of the key, it is a note foreign
// to it, and the caller decides what that means.
func (t Tonality) DegreeOf(c PitchClass) (Degree, bool) {
	if t.IsZero() {
		return 0, false
	}
	for d, n := range t.pattern.Offsets() {
		if t.tonic.Transpose(n) == c {
			return d, true
		}
	}
	return 0, false
}

// Transpose moves `t` to the tonic `n` semitones away, keeping its pattern.
//
// The pattern is unchanged, so the result is heptatonic and cannot
// fail. This makes Tonality satisfy [Transposable].
func (t Tonality) Transpose(n Semitones) Tonality {
	if t.IsZero() {
		return t
	}
	return Tonality{tonic: t.tonic.Transpose(n), pattern: t.pattern}
}

// String returns the tonic and pattern of `t` numerically, for debugging
// only. It never returns a key name.
func (t Tonality) String() string {
	if t.IsZero() {
		return "no tonality"
	}
	return fmt.Sprintf("%d/%012b", uint8(t.tonic), uint16(t.pattern))
}
