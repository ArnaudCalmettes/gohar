package ireal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// Changes reads the chords of the timeline into the changes the
// analysis works on: each chord with its bass, when it starts and how
// long it sounds, looping as the app plays it.
//
// A chord that cannot be read, a typo in a quality typed by hand or a
// "W" with no chord before it, becomes a silent change, and the error
// says which: the analysis goes on around it rather than stopping, and
// the caller decides whether to report it.
func (t Timeline) Changes() (analysis.Changes, error) {
	c := analysis.Changes{
		Chords: make([]analysis.Change, len(t.Spans)),
		Loops:  true,
		Coda:   t.Coda,
		Bars:   t.Bars,
	}
	var errs []error
	for i, s := range t.Spans {
		ch := analysis.Change{Start: s.Start, Length: s.Length, Silent: true}
		if !s.NoChord {
			r, err := s.Chord.Read()
			if err != nil {
				errs = append(errs, fmt.Errorf("bar %d: %v: %w", s.Bar+1, s.Chord, err))
			} else {
				ch.Chord, ch.Silent = r.Chord(), false
				ch.Bass = ch.Chord.Root
				if r.HasBass {
					ch.Bass = r.Bass.Class()
				}
			}
		}
		c.Chords[i] = ch
	}
	return c, errors.Join(errs...)
}

// Changes unfolds the chart, times it and reads its chords: see
// [Timeline.Changes].
func (c Chart) Changes() (analysis.Changes, error) {
	return c.Timeline().Changes()
}

// HomeTonalities reads the tonality the app gives the song, its Key
// field: "Eb" for E flat major, "A-" for A minor, read in all three
// minors. It is the key signature of the chart, a clue to the
// tonality of the tune that its cadences confirm or not, and false
// when the app gives none that can be read.
func (s Song) HomeTonalities() ([]harmony.Tonality, bool) {
	root, minor := strings.CutSuffix(s.Key, "-")
	n, ok := spell(root)
	if !ok {
		return nil, false
	}
	if minor {
		return analysis.MinorTonalities(n.Class()), true
	}
	return analysis.MajorTonalities(n.Class()), true
}
