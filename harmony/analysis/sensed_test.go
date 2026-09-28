package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// sensedOf writes what the ear expects after each chord, one chord to
// a line: the ground it is heard in, the local tonic a cadence has
// just tonicised, and what the cadence being played awaits.
//
//	Gm7♭5 in E♭, awaiting Fm harm
func sensedOf(c analysis.Changes) string {
	var out []string
	for i, s := range analysis.Sense(c, analysis.Blocks(c, analysis.Approaches(c))) {
		line := changeName(c.Chords[i]) + " in "
		if s.Ground == nil {
			line += "?"
		} else {
			line += tonalityName(s.Ground)
		}
		if s.Local != nil {
			line += ", on " + tonalityName(s.Local)
		}
		if s.Awaited != nil {
			line += ", awaiting " + tonalityName(s.Awaited)
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

func TestSense(t *testing.T) {
	const (
		c, db, d, eb, f, g, ab, a, bb harmony.PitchClass = 0, 1, 2, 3, 5, 7, 8, 9, 10
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		min6    = harmony.ChordMinorSixth
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	flatNine, _ := harmony.NewChordPattern(0, 4, 7, 10, 13)
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    []string
	}{
		// Tenderly, bars 1 to 9, as docs/grilles.md hears them. E♭m7
		// A♭7 awaits nothing: its two is the tonic, borrowed. Plagal
		// cadences do await their tonic: A♭7 E♭m7, the IV7-Im of the
		// melodic minor (En Harmonie borrows A♭7 from E flat melodic
		// minor), and D♭7, the ♭VII7 of a minor plagal. The
		// cadence onto Fm7♭5 tonicises F minor for one chord, and
		// Fm7♭5 is already the two of E flat minor.
		"Tenderly": {
			changesOf(false, eb, maj7, ab, dom7, eb, min7, ab, dom7, f, min7, db, dom7, eb, maj7,
				g, halfDim, c, flatNine, f, halfDim, bb, flatNine),
			[]string{
				"E♭maj7 in E♭",
				"A♭7 in E♭, awaiting E♭m mel",
				"E♭m7 in E♭",
				"A♭7 in E♭",
				"Fm7 in E♭",
				"D♭7 in E♭, awaiting E♭m",
				"E♭maj7 in E♭",
				"Gm7♭5 in E♭, awaiting Fm harm",
				"C7♭9 in E♭, awaiting Fm harm",
				"Fm7♭5 in E♭, on Fm harm, awaiting E♭m harm",
				"B♭7♭9 in E♭, awaiting E♭m harm",
			},
		},
		// The same two five on the tonic, resolving this time: it
		// tonicises D flat.
		"a two five on the tonic that resolves": {
			changesOf(false, eb, maj7, eb, min7, ab, dom7, db, maj7),
			[]string{
				"E♭maj7 in E♭",
				"E♭m7 in E♭, awaiting D♭",
				"A♭7 in E♭, awaiting D♭",
				"D♭maj7 in E♭, on D♭",
			},
		},
		// Autumn Leaves, heard once from its start: a two chord installs
		// nothing, the first cadence installs B flat, and G minor comes
		// as a local tonic. Setting when it becomes the ground is the
		// threshold of a modulation, still to come.
		"Autumn Leaves, the relative first": {
			changesOf(false, c, min7, f, dom7, bb, maj7, eb, maj7, a, halfDim, d, dom7, g, min6),
			[]string{
				"Cm7 in ?, awaiting B♭",
				"F7 in ?, awaiting B♭",
				"B♭maj7 in B♭",
				"E♭maj7 in B♭",
				"Am7♭5 in B♭, awaiting Gm harm",
				"D7 in B♭, awaiting Gm harm",
				"Gm6 in B♭, on Gm harm",
			},
		},
	} {
		if got, want := sensedOf(tc.changes), strings.Join(tc.want, "\n"); got != want {
			t.Errorf("%s:\ngot\n%s\nwant\n%s", name, got, want)
		}
	}
}

// The tonality of the tune is read at its end: Autumn Leaves, begun in
// B flat, ends in G minor.
func TestTune(t *testing.T) {
	const (
		c, d, eb, f, g, a, bb harmony.PitchClass = 0, 2, 3, 5, 7, 9, 10
	)
	leaves := changesOf(false, c, harmony.ChordMinorSeventh, f, harmony.ChordDominantSeventh,
		bb, harmony.ChordMajorSeventh, eb, harmony.ChordMajorSeventh,
		a, harmony.ChordHalfDiminished, d, harmony.ChordDominantSeventh, g, harmony.ChordMinorSixth)
	sensed := analysis.Sense(leaves, analysis.Blocks(leaves, analysis.Approaches(leaves)))
	if got := tonalityName(analysis.Tune(leaves, sensed)); got != "Gm nat/harm/mel" {
		t.Errorf("%s, want Gm nat/harm/mel", got)
	}
}

// groundsOf writes the ground of each chord, one chord to a line: the
// one the ear has installed there, and the one heard afterwards when
// the region started earlier.
//
//	Dm7 in D, in C afterwards
func groundsOf(c analysis.Changes) string {
	sensed := analysis.Sense(c, analysis.Blocks(c, analysis.Approaches(c)))
	grounds := analysis.Grounds(c, sensed)
	var out []string
	name := func(ts []harmony.Tonality) string {
		if ts == nil {
			return "?"
		}
		return tonalityName(ts)
	}
	for i, s := range sensed {
		line := changeName(c.Chords[i]) + " in " + name(s.Ground)
		if after := name(grounds[i]); after != name(s.Ground) {
			line += ", in " + after + " afterwards"
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// Tune Up modulates at each phrase: C major installs on the second bar
// of Cmaj7, B flat on Gm7, which still holds in it, and D comes home
// with one cadence. Black Orpheus modulates to C major on its second
// cadence there, A7♭9 Dm7 between them keeping C, and comes home to A
// minor. Along Came Betty, heard as a chorus that comes round, dances
// between A flat and its neighbours, each tonic held a single bar, and
// never modulates.
func TestModulation(t *testing.T) {
	const (
		c, d, eb, e, f, gb, g, ab, a, bb, b harmony.PitchClass = 0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		minor   = harmony.ChordMinorTriad
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	flatNine, _ := harmony.NewChordPattern(0, 4, 7, 10, 13)
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    []string
	}{
		"Tune Up": {
			changesOf(false, e, min7, a, dom7, d, maj7, d, maj7, d, min7, g, dom7, c, maj7, c, maj7,
				c, min7, f, dom7, bb, maj7, g, min7, e, min7, a, dom7, d, maj7, d, maj7),
			[]string{
				"Em7 in ?, in D afterwards",
				"A7 in ?, in D afterwards",
				"Dmaj7 in D",
				"Dmaj7 in D",
				"Dm7 in D, in C afterwards",
				"G7 in D, in C afterwards",
				"Cmaj7 in D, in C afterwards",
				"Cmaj7 in C",
				"Cm7 in C, in B♭ afterwards",
				"F7 in C, in B♭ afterwards",
				"B♭maj7 in C, in B♭ afterwards",
				"Gm7 in B♭",
				"Em7 in B♭, in D afterwards",
				"A7 in B♭, in D afterwards",
				"Dmaj7 in D",
				"Dmaj7 in D",
			},
		},
		"Black Orpheus": {
			changesOf(false, a, minor, b, halfDim, e, flatNine, a, minor, b, halfDim, e, flatNine, a, minor,
				d, min7, g, dom7, c, maj7, a, flatNine, d, min7, g, dom7, c, maj7, f, maj7,
				b, halfDim, e, dom7, a, minor),
			[]string{
				"Am in Am nat/harm/mel",
				"Bm7♭5 in Am nat/harm/mel",
				"E7♭9 in Am nat/harm/mel",
				"Am in Am nat/harm/mel",
				"Bm7♭5 in Am nat/harm/mel",
				"E7♭9 in Am nat/harm/mel",
				"Am in Am nat/harm/mel",
				"Dm7 in Am nat/harm/mel, in C afterwards",
				"G7 in Am nat/harm/mel, in C afterwards",
				"Cmaj7 in Am nat/harm/mel, in C afterwards",
				"A7♭9 in Am nat/harm/mel, in C afterwards",
				"Dm7 in Am nat/harm/mel, in C afterwards",
				"G7 in Am nat/harm/mel, in C afterwards",
				"Cmaj7 in C",
				"Fmaj7 in C",
				"Bm7♭5 in C, in Am nat/harm/mel afterwards",
				"E7 in C, in Am nat/harm/mel afterwards",
				"Am in Am nat/harm/mel",
			},
		},
		"Along Came Betty": {
			changesOf(true, bb, min7, b, min7, e, dom7, bb, min7, b, min7, e, dom7, a, maj7, ab, dom7,
				g, maj7, gb, dom7, gb, min7, g, min7, c, dom7, gb, min7, g, min7, c, dom7,
				f, maj7, a, flatNine, d, min7, g, dom7,
				c, halfDim, f, dom7, bb, halfDim, eb, dom7, ab, maj7),
			[]string{
				"B♭m7 in A♭",
				"Bm7 in A♭",
				"E7 in A♭",
				"B♭m7 in A♭",
				"Bm7 in A♭",
				"E7 in A♭",
				"Amaj7 in A♭",
				"A♭7 in A♭",
				"Gmaj7 in A♭",
				"F♯7 in A♭",
				"F♯m7 in A♭",
				"Gm7 in A♭",
				"C7 in A♭",
				"F♯m7 in A♭",
				"Gm7 in A♭",
				"C7 in A♭",
				"Fmaj7 in A♭",
				"A7♭9 in A♭",
				"Dm7 in A♭",
				"G7 in A♭",
				"Cm7♭5 in A♭",
				"F7 in A♭",
				"B♭m7♭5 in A♭",
				"E♭7 in A♭",
				"A♭maj7 in A♭",
			},
		},
	} {
		if got, want := groundsOf(tc.changes), strings.Join(tc.want, "\n"); got != want {
			t.Errorf("%s:\ngot\n%s\nwant\n%s", name, got, want)
		}
	}
}

// A looping chart is heard as its second chorus: Tune Up starts on Em7
// in D, where its chorus ended, and not in an unknown tonic.
func TestSenseLoops(t *testing.T) {
	const c, d, e, f, g, a, bb harmony.PitchClass = 0, 2, 4, 5, 7, 9, 10
	const (
		maj7 = harmony.ChordMajorSeventh
		min7 = harmony.ChordMinorSeventh
		dom7 = harmony.ChordDominantSeventh
	)
	tuneUp := changesOf(true, e, min7, a, dom7, d, maj7, d, maj7, d, min7, g, dom7, c, maj7, c, maj7,
		c, min7, f, dom7, bb, maj7, g, min7, e, min7, a, dom7, d, maj7, d, maj7)
	if got, want := strings.SplitN(groundsOf(tuneUp), "\n", 2)[0], "Em7 in D"; got != want {
		t.Errorf("%s, want %s", got, want)
	}
}

// A chorus that ends on a turnaround ends in its ground: G7 C7 after F
// is not a tune in C.
func TestTuneAfterATurnaround(t *testing.T) {
	const c, d, f, g harmony.PitchClass = 0, 2, 5, 7
	ch := changesOf(false, f, harmony.ChordMajorSeventh, d, harmony.ChordMinorSeventh,
		g, harmony.ChordDominantSeventh, c, harmony.ChordDominantSeventh)
	sensed := analysis.Sense(ch, analysis.Blocks(ch, analysis.Approaches(ch)))
	if got := tonalityName(analysis.Tune(ch, sensed)); got != "F" {
		t.Errorf("%s, want F", got)
	}
}

// A blues is recognised by its form, and in its tonic from bar 1, even
// on a I7: Sonnymoon For Two is in B flat, where F7 B♭7 would otherwise
// install E flat, the IV.
func TestBlues(t *testing.T) {
	const (
		c, db, d, eb, e, f, g, ab, a, bb harmony.PitchClass = 0, 1, 2, 3, 4, 5, 7, 8, 9, 10
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		"Sonnymoon For Two": {
			barsOf(true, bar{bb, dom7}, bar{eb, dom7}, bar{bb, dom7}, bar{bb, dom7},
				bar{eb, dom7}, bar{eb, dom7}, bar{bb, dom7}, bar{bb, dom7},
				bar{c, min7}, bar{f, dom7}, bar{bb, dom7}, bar{f, dom7}),
			"B♭",
		},
		"Blues For Alice, a Bird blues": {
			barsOf(true, bar{f, maj7}, bar{e, halfDim, a, dom7}, bar{d, min7, g, dom7}, bar{c, min7, f, dom7},
				bar{bb, dom7}, bar{bb, min7, eb, dom7}, bar{a, min7, d, dom7}, bar{ab, min7, db, dom7},
				bar{g, min7}, bar{c, dom7}, bar{f, dom7, d, dom7}, bar{g, min7, c, dom7}),
			"F",
		},
		"Mr. P.C., a minor blues": {
			barsOf(true, bar{c, min7}, bar{c, min7}, bar{c, min7}, bar{c, min7},
				bar{f, min7}, bar{f, min7}, bar{c, min7}, bar{c, min7},
				bar{ab, dom7}, bar{g, dom7}, bar{c, min7}, bar{g, dom7}),
			"Cm nat/harm/mel",
		},
		"twelve bars, not a blues": {
			barsOf(true, bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{c, maj7},
				bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{c, maj7},
				bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{c, maj7}),
			"none",
		},
	} {
		got := "none"
		if ts, ok := analysis.Blues(tc.changes); ok {
			got = tonalityName(ts)
		}
		if got != tc.want {
			t.Errorf("%s: %s, want %s", name, got, tc.want)
		}
	}
}

// Sonnymoon For Two is heard in B flat from its first bar.
func TestSenseBlues(t *testing.T) {
	const c, eb, f, bb harmony.PitchClass = 0, 3, 5, 10
	const (
		min7 = harmony.ChordMinorSeventh
		dom7 = harmony.ChordDominantSeventh
	)
	type bar = []any
	sonnymoon := barsOf(true, bar{bb, dom7}, bar{eb, dom7}, bar{bb, dom7}, bar{bb, dom7},
		bar{eb, dom7}, bar{eb, dom7}, bar{bb, dom7}, bar{bb, dom7},
		bar{c, min7}, bar{f, dom7}, bar{bb, dom7}, bar{f, dom7})
	got := strings.SplitN(groundsOf(sonnymoon), "\n", 3)
	if want := []string{"B♭7 in B♭", "E♭7 in B♭"}; got[0] != want[0] || got[1] != want[1] {
		t.Errorf("%q, want %q", got[:2], want)
	}
}

// The end of My Way. After a modulation to F, F/C is a six-four over
// the tonic pedal, not a tonic: with C it makes a plagal cadence, an
// amen that brings home back, and heard afterwards C starts with it.
func TestHomeOnSight(t *testing.T) {
	const c, f, g harmony.PitchClass = 0, 5, 7
	ch := changesOf(false, c, harmony.ChordMajorSeventh, g, harmony.ChordMinorSeventh,
		c, harmony.ChordDominantSeventh, f, harmony.ChordMajorSeventh, f, harmony.ChordMajorSeventh,
		f, harmony.ChordMajorTriad, c, harmony.ChordMajorTriad)
	ch.Chords[5].Bass = c
	want := []string{
		"Cmaj7 in C",
		"Gm7 in C, in F afterwards",
		"C7 in C, in F afterwards",
		"Fmaj7 in C, in F afterwards",
		"Fmaj7 in F",
		"F/C in F, in C afterwards",
		"C in C",
	}
	if got := groundsOf(ch); got != strings.Join(want, "\n") {
		t.Errorf("got\n%s\nwant\n%s", got, strings.Join(want, "\n"))
	}
}

// Where a tune ends, and so which tonality it is in.
func TestTuneEndings(t *testing.T) {
	const (
		c, d, e, f, gb, g, a, b harmony.PitchClass = 0, 2, 4, 5, 6, 7, 9, 11
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		// Softly, As In A Morning Sunrise: the minor tonic is written m7,
		// and a minor cadence resolves on it.
		"a tonic written m7": {
			barsOf(true, bar{c, min7}, bar{d, halfDim, g, dom7}, bar{c, min7}, bar{d, halfDim, g, dom7}),
			"Cm nat/harm/mel",
		},
		// Sugar: G♭7 Fm7 tonicises the IV, and the last G7 goes back to
		// Cm7 across the loop.
		"the last cadence goes back to the start": {
			barsOf(true, bar{c, min7}, bar{d, halfDim, g, dom7}, bar{c, min7}, bar{g, dom7},
				bar{c, min7}, bar{gb, dom7}, bar{f, min7}, bar{d, halfDim}, bar{g, dom7}),
			"Cm nat/harm/mel",
		},
		// Fly Me To The Moon: the turnaround goes back to Am7, a VI, after
		// the tune has ended on its I.
		"a turnaround toward a VI": {
			barsOf(true, bar{a, min7}, bar{d, min7}, bar{g, dom7}, bar{c, maj7},
				bar{f, maj7}, bar{b, halfDim}, bar{e, dom7}, bar{a, min7},
				bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{b, halfDim, e, dom7}),
			"C",
		},
	} {
		sensed := analysis.Sense(tc.changes, analysis.Blocks(tc.changes, analysis.Approaches(tc.changes)))
		if got := tonalityName(analysis.Tune(tc.changes, sensed)); got != tc.want {
			t.Errorf("%s: %s, want %s", name, got, tc.want)
		}
	}
}

// A minor tune that ends on its tonic made major stays minor: the
// tierce picarde. A tune that stays major longer is major.
func TestPicardy(t *testing.T) {
	const c, d, g harmony.PitchClass = 0, 2, 7
	const (
		maj     = harmony.ChordMajorTriad
		maj7    = harmony.ChordMajorSeventh
		min6    = harmony.ChordMinorSixth
		dom7    = harmony.ChordDominantSeventh
		min7    = harmony.ChordMinorSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
		picardy bool
	}{
		"a minor tune ending major": {
			barsOf(false, bar{c, min6}, bar{d, halfDim, g, dom7}, bar{c, min6}, bar{d, halfDim, g, dom7}, bar{c, maj}),
			"Cm nat/harm/mel", true,
		},
		"a major tune with a minor bar": {
			barsOf(false, bar{c, maj7}, bar{d, min7, g, dom7}, bar{c, maj7}, bar{d, halfDim, g, dom7}, bar{c, min6},
				bar{d, min7, g, dom7}, bar{c, maj7}),
			"C", false,
		},
	} {
		c := tc.changes
		sensed := analysis.Sense(c, analysis.Blocks(c, analysis.Approaches(c)))
		if got := tonalityName(analysis.Tune(c, sensed)); got != tc.want {
			t.Errorf("%s: %s, want %s", name, got, tc.want)
		}
		if got := analysis.Picardy(c, sensed); got != tc.picardy {
			t.Errorf("%s: tierce picarde %v, want %v", name, got, tc.picardy)
		}
	}
}
