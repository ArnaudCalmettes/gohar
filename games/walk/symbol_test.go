package main

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
)

// The runs of a symbol, the raised one between brackets.
func TestSymbolOf(t *testing.T) {
	for in, want := range map[string]string{
		"D7(b13)":    "D7[(♭13)]",
		"Am7b5":      "Am7[♭5]",
		"Bbmaj7":     "B♭maj7",
		"Cmaj7#5":    "Cmaj7[♯5]",
		"C7b5(b9)":   "C7[♭5(♭9)]",
		"Gm6":        "Gm6",
		"Cm(maj7,9)": "Cm(maj7,9)",
		"C13":        "C13",
		"C7alt":      "C7alt",
		"D7(b9)/F#":  "D7[(♭9)]/F♯",
		"N.C.":       "N.C.",
	} {
		c, err := chordpro.ReadChord(in)
		if err != nil {
			t.Fatal(err)
		}
		s := symbolOf(c)
		got := s.line
		if s.raised != "" {
			got += "[" + s.raised + "]"
		}
		got += s.bass
		if got != want {
			t.Errorf("%s: drawn %s, want %s", in, got, want)
		}
	}
}
