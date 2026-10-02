package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

func chordPattern(t *testing.T, offsets ...harmony.Semitones) harmony.ChordPattern {
	t.Helper()
	p, err := harmony.NewChordPattern(offsets...)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// The qualities the app writes, as a lead sheet shows them by default.
func TestChordSymbol(t *testing.T) {
	for _, tc := range []struct {
		offsets []harmony.Semitones
		want    string
	}{
		{[]harmony.Semitones{4, 7}, ""},
		{[]harmony.Semitones{3, 7}, "m"},
		{[]harmony.Semitones{3, 6}, "dim"},
		{[]harmony.Semitones{4, 8}, "+"},
		{[]harmony.Semitones{7}, "5"},
		{[]harmony.Semitones{2, 7}, "sus2"},
		{[]harmony.Semitones{5, 7}, "sus4"},
		{[]harmony.Semitones{4, 7, 14}, "add9"},
		{[]harmony.Semitones{3, 7, 14}, "madd9"},
		{[]harmony.Semitones{4, 7, 9}, "6"},
		{[]harmony.Semitones{4, 7, 9, 14}, "6/9"},
		{[]harmony.Semitones{3, 7, 9}, "m6"},
		{[]harmony.Semitones{3, 7, 9, 14}, "m6/9"},
		{[]harmony.Semitones{3, 7, 8}, "m♭6"},
		{[]harmony.Semitones{3, 8}, "m♯5"},
		{[]harmony.Semitones{4, 7, 11}, "maj7"},
		{[]harmony.Semitones{4, 7, 11, 14}, "maj9"},
		{[]harmony.Semitones{4, 7, 11, 14, 21}, "maj13"},
		{[]harmony.Semitones{4, 7, 11, 18}, "maj7(♯11)"},
		{[]harmony.Semitones{4, 7, 11, 14, 18}, "maj9(♯11)"},
		{[]harmony.Semitones{4, 8, 11}, "maj7♯5"},
		{[]harmony.Semitones{3, 7, 10}, "m7"},
		{[]harmony.Semitones{3, 7, 10, 14}, "m9"},
		{[]harmony.Semitones{3, 7, 10, 14, 17}, "m11"},
		{[]harmony.Semitones{3, 7, 10, 14, 17, 21}, "m13"},
		{[]harmony.Semitones{3, 7, 10, 14, 21}, "m9(13)"},
		{[]harmony.Semitones{3, 7, 11}, "m(maj7)"},
		{[]harmony.Semitones{3, 7, 11, 14}, "m(maj7,9)"},
		{[]harmony.Semitones{3, 7, 11, 14, 18}, "m(maj7,9,♯11)"},
		{[]harmony.Semitones{3, 6, 10}, "m7♭5"},
		{[]harmony.Semitones{3, 6, 10, 14}, "m9♭5"},
		{[]harmony.Semitones{3, 6, 9}, "dim7"},
		{[]harmony.Semitones{3, 6, 9, 23}, "dim7(♮14)"},
		{[]harmony.Semitones{4, 7, 10}, "7"},
		{[]harmony.Semitones{4, 7, 10, 14}, "9"},
		{[]harmony.Semitones{4, 7, 10, 14, 21}, "13"},
		{[]harmony.Semitones{4, 6, 10}, "7♭5"},
		{[]harmony.Semitones{4, 8, 10}, "7♯5"},
		{[]harmony.Semitones{4, 7, 10, 13}, "7(♭9)"},
		{[]harmony.Semitones{4, 7, 10, 15}, "7(♯9)"},
		{[]harmony.Semitones{4, 7, 10, 18}, "7(♯11)"},
		{[]harmony.Semitones{4, 7, 10, 20}, "7(♭13)"},
		{[]harmony.Semitones{4, 7, 10, 14, 18, 21}, "13(♯11)"},
		{[]harmony.Semitones{4, 7, 10, 14, 17, 21}, "13(11)"},
		{[]harmony.Semitones{4, 7, 10, 13, 21}, "7(♭9,13)"},
		{[]harmony.Semitones{4, 6, 10, 14}, "9♭5"},
		{[]harmony.Semitones{4, 8, 10, 15}, "7♯5(♯9)"},
		{[]harmony.Semitones{4, 7, 10, 13, 18}, "7(♭9,♯11)"},
		{[]harmony.Semitones{4, 6, 10, 13, 15, 20}, "7alt"},
		{[]harmony.Semitones{4, 7, 10, 17}, "7(11)"},
		{[]harmony.Semitones{5, 7, 10}, "7sus4"},
		{[]harmony.Semitones{5, 7, 10, 14}, "9sus4"},
		{[]harmony.Semitones{5, 7, 10, 14, 21}, "13sus4"},
		{[]harmony.Semitones{5, 7, 10, 13}, "7sus4(♭9)"},
	} {
		p := chordPattern(t, tc.offsets...)
		got, ok := naming.ChordStyle{}.Symbol(p)
		if !ok || got != tc.want {
			t.Errorf("C%v: got %q (%v), want %q", tc.offsets, got, ok, tc.want)
		}
	}
}

// Each sign the style offers, against its default.
func TestChordStyle(t *testing.T) {
	maj9 := chordPattern(t, 4, 7, 11, 14)
	minMaj7 := chordPattern(t, 3, 7, 11)
	min7 := chordPattern(t, 3, 7, 10)
	half9 := chordPattern(t, 3, 6, 10, 14)
	half := chordPattern(t, 3, 6, 10)
	dim := chordPattern(t, 3, 6)
	dim7 := chordPattern(t, 3, 6, 9)
	for _, tc := range []struct {
		style naming.ChordStyle
		p     harmony.ChordPattern
		want  string
	}{
		{naming.ChordStyle{MajorSeventh: naming.DeltaSeventh}, maj9, "Δ9"},
		{naming.ChordStyle{MajorSeventh: naming.NaturalSeventh}, maj9, "♮9"},
		{naming.ChordStyle{MajorSeventh: naming.DeltaSeventh, Minus: true}, minMaj7, "-(Δ7)"},
		{naming.ChordStyle{MajorSeventh: naming.NaturalSeventh}, minMaj7, "m(♮7)"},
		{naming.ChordStyle{Minus: true}, min7, "-7"},
		{naming.ChordStyle{HalfDiminishedSign: true}, half, "ø"},
		{naming.ChordStyle{HalfDiminishedSign: true}, half9, "ø9"},
		{naming.ChordStyle{DiminishedSign: true}, dim, "°"},
		{naming.ChordStyle{DiminishedSign: true}, dim7, "°7"},
	} {
		if got, _ := tc.style.Symbol(tc.p); got != tc.want {
			t.Errorf("%+v: got %q, want %q", tc.style, got, tc.want)
		}
	}
}

// A first octave that is no chord is not named.
func TestChordSymbolUnknown(t *testing.T) {
	if got, ok := (naming.ChordStyle{}).Symbol(chordPattern(t, 1, 2)); ok {
		t.Errorf("C(♭2, 2): got %q, want none", got)
	}
}
