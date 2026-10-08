package grids

import (
	"errors"
	"fmt"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
	"github.com/ArnaudCalmettes/gohar/harmony/naming"
)

// InKey transposes `t` so that the tonality heard in it lands on
// `tonic`: the blues in F, in C for the lessons. The tonality is the
// first the analysis hears, never the key the chart declares.
func InKey(t Tune, tonic naming.SpelledNote) (Tune, error) {
	heard := analysis.Hear(t.Grid, nil).Tune
	if len(heard) == 0 {
		return Tune{}, fmt.Errorf("walk: %s: no tonality heard", t.Title)
	}
	from := naming.TonicSpelling(heard[0])
	steps := (int(tonic.Letter) - int(from.Letter) + naming.LetterCount) % naming.LetterCount
	return Transpose(t, steps, from.Class().Up(tonic.Class()))
}

// errSpelling is a note no accidental can write once transposed, past
// a double sharp or a double flat.
var errSpelling = errors.New("cannot spell")

// Transpose moves every chord of `t` by an interval: `steps` letters
// and `by` semitones up, a fifth being 4 and 7. The written chords are
// moved letter by letter, as a musician transposes a chart: from F to
// D♭, B♭7 becomes G♭7, never F♯7. The rhythm stays as it is.
func Transpose(t Tune, steps int, by harmony.Semitones) (Tune, error) {
	out := Tune{Title: t.Title, Tempo: t.Tempo, Grid: t.Grid, Written: make([]chordpro.Chord, len(t.Written)), Rehearsals: t.Rehearsals}
	out.Grid.Chords = make([]analysis.Change, len(t.Grid.Chords))
	for i, ch := range t.Grid.Chords {
		if !ch.Silent {
			ch.Chord.Root = ch.Chord.Root.Transpose(by)
			ch.Bass = ch.Bass.Transpose(by)
		}
		out.Grid.Chords[i] = ch
	}
	move := func(n naming.SpelledNote) (naming.SpelledNote, error) {
		m, ok := naming.SpellAbove(n, steps, n.Class().Transpose(by))
		if !ok {
			return m, fmt.Errorf("walk: %s: %w %v moved", t.Title, errSpelling, n)
		}
		return m, nil
	}
	for i, w := range t.Written {
		var err error
		if !w.NoChord {
			if w.Root, err = move(w.Root); err != nil {
				return Tune{}, err
			}
		}
		if w.HasBass {
			if w.Bass, err = move(w.Bass); err != nil {
				return Tune{}, err
			}
		}
		out.Written[i] = w
	}
	return out, nil
}
