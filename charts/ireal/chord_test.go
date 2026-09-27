package ireal

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

func TestReadChord(t *testing.T) {
	for in, want := range map[string]struct {
		root    naming.SpelledNote
		pattern harmony.ChordPattern
	}{
		"C^7":   {naming.SpelledNote{Letter: naming.LetterC}, harmony.ChordMajorSeventh},
		"Bb-7":  {naming.SpelledNote{Letter: naming.LetterB, Accidental: naming.FlatSign}, harmony.ChordMinorSeventh},
		"F#h7":  {naming.SpelledNote{Letter: naming.LetterF, Accidental: naming.SharpSign}, harmony.ChordHalfDiminished},
		"Eo7":   {naming.SpelledNote{Letter: naming.LetterE}, harmony.ChordDiminishedSeventh},
		"G7sus": {naming.SpelledNote{Letter: naming.LetterG}, harmony.ChordDominantSeventhSus4},
		"A-^7":  {naming.SpelledNote{Letter: naming.LetterA}, harmony.ChordMinorMajorSeventh},
		"D6":    {naming.SpelledNote{Letter: naming.LetterD}, harmony.ChordMajorSixth},
	} {
		sym, _ := lexChord(in)
		r, err := sym.Read()
		if err != nil || r.Root != want.root || r.Pattern != want.pattern {
			t.Errorf("%s: %+v, %v", in, r, err)
		}
	}
}

func TestReadChordBassAndErrors(t *testing.T) {
	sym, _ := lexChord("C-7/Bb")
	r, err := sym.Read()
	if err != nil || !r.HasBass || r.Bass != (naming.SpelledNote{Letter: naming.LetterB, Accidental: naming.FlatSign}) {
		t.Errorf("%+v, %v", r, err)
	}
	if _, err := (ChordSymbol{Root: "W"}).Read(); err != ErrInvisibleRoot {
		t.Errorf("W: %v", err)
	}
	if _, err := (ChordSymbol{Root: "G", Quality: "7us", Custom: true}).Read(); err == nil {
		t.Errorf("a typo read")
	}
}

// Every reading is already in the form harmony normalises to: the
// table and the rules of construction never disagree, so that a chord
// read from a chart and the same chord played compare equal.
func TestQualitiesAreNormalised(t *testing.T) {
	all := map[string][]int{}
	for q, o := range qualityOffsets {
		all[q] = o
	}
	for q, o := range handQualities {
		all["hand "+q] = o
	}
	for q, o := range all {
		st := make([]harmony.Semitones, len(o))
		for i, n := range o {
			st[i] = harmony.Semitones(n)
		}
		p, err := harmony.NewChordPattern(st...)
		if err != nil {
			t.Errorf("%q: %v", q, err)
			continue
		}
		if n := p.Normalize(); n != p {
			t.Errorf("%q: %v normalises to %v", q, p, n)
		}
	}
}
