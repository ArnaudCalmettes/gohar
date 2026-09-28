package analysis

import "github.com/ArnaudCalmettes/gohar/harmony"

// Blues recognises a twelve bar blues, played once or twice, or one of
// twenty four bars in double time, by its skeleton: a chord on its
// tonic at bar 1, the IV at bar 5, the tonic again at bar 11. It returns the tonalities of its
// tonic: major, or the three minors when bar 1 has a minor third (Mr.
// P.C.).
//
// The skeleton is loose on purpose. Bar 1 may hold any chord on the
// tonic: the I7 of the blues, but also the Imaj7 of the Bird blues
// (Blues For Alice) or the Im7 of the minor blues. Bars 7 to 10, where
// the variants part ways, are not looked at.
//
// The form is what tells that the I7 of a blues is its tonic: a
// seventh of kind (septième d'espèce), not a dominant, which its sound
// alone does not tell. Without the form, F7 at bar 1 would wait for a
// cadence, and the first one in a blues, F7 B♭7, installs the IV.
func Blues(c Changes) ([]harmony.Tonality, bool) {
	switch len(c.Bars) {
	case 12:
		return blues(c, 0, 1)
	case 24:
		// Twice twelve bars, each with the skeleton, or twelve in double
		// time.
		if ts, ok := blues(c, 0, 1); ok {
			if _, again := blues(c, 12, 1); again {
				return ts, true
			}
		}
		return blues(c, 0, 2)
	}
	// Longer forms that happen to hold the skeleton are songs, not
	// blues: If I Loved You, I Remember You.
	return nil, false
}

// blues reads the skeleton from bar `from`, k bars to a bar of the
// blues.
func blues(c Changes, from, k int) ([]harmony.Tonality, bool) {
	at := func(bar int) (Change, bool) {
		start := c.Bars[from+(bar-1)*k]
		for _, ch := range c.Chords {
			if ch.Start <= start && start < ch.Start+ch.Length {
				return ch, !ch.Silent
			}
		}
		return Change{}, false
	}
	one, ok1 := at(1)
	five, ok5 := at(5)
	eleven, ok11 := at(11)
	if !ok1 || !ok5 || !ok11 {
		return nil, false
	}
	root := one.Chord.Root
	if five.Chord.Root != root.Transpose(5) || eleven.Chord.Root != root {
		return nil, false
	}
	p := one.Chord.Pattern
	switch {
	case p.HasOffset(4):
		return MajorTonalities(root), true
	case p.HasOffset(3):
		return MinorTonalities(root), true
	}
	return nil, false
}
