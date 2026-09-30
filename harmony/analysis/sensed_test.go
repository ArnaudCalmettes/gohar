package analysis_test

import (
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// sensedOf writes what the ear expects after each chord, one chord to
// a line: the ground it is heard in, the region a cadence has opened,
// the tonic a dominant makes of the chord when it is not the region's,
// and what the cadence being played awaits.
//
//	Gm7♭5 in E♭, awaiting Fm harm
func sensedOf(c analysis.Changes) string {
	var out []string
	for i, s := range sense(c) {
		line := changeName(c.Chords[i]) + " in "
		if s.Ground == nil {
			line += "?"
		} else {
			line += tonalityName(s.Ground)
		}
		if s.Region != nil {
			line += ", on " + tonalityName(s.Region)
		}
		if s.Tonicised != nil && !sameTonicName(s.Tonicised, s.Region) {
			line += ", as a I of " + tonalityName(s.Tonicised)
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
		minMaj7 = harmony.ChordMinorMajorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	flatNine, _ := harmony.NewChordPattern(0, 4, 7, 10, 13)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    []string
	}{
		// Tenderly, bars 1 to 9, as docs/grilles.md hears them. E♭m7
		// A♭7 awaits nothing: its two is the tonic, borrowed. Plagal
		// cadences do await their tonic: A♭7 E♭m7, the IV7-Im of the
		// melodic minor (En Harmonie borrows A♭7 from E flat melodic
		// minor), and D♭7, the ♭VII7 of a minor plagal. Gm7♭5 C7♭9
		// awaits F minor, and lands on Fm7♭5, no tonic but the two of E
		// flat minor.
		"Tenderly": {
			barsOf(false, bar{eb, maj7}, bar{ab, dom7}, bar{eb, min7}, bar{ab, dom7}, bar{f, min7}, bar{db, dom7},
				bar{eb, maj7}, bar{g, halfDim, c, flatNine}, bar{f, halfDim, bb, flatNine}),
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
				"Fm7♭5 in E♭, awaiting E♭m harm",
				"B♭7♭9 in E♭, awaiting E♭m harm",
			},
		},
		// My Lucky Star: the II tonicised stays local under Gm7 and
		// Gm(maj7) alike, G minor with a seventh that wanders.
		"a local tonic changing colour": {
			changesOf(false, f, maj7, a, min7, d, dom7, g, min7, g, minMaj7, g, min7, g, minMaj7),
			[]string{
				"Fmaj7 in F",
				"Am7 in F, awaiting Gm mel",
				"D7 in F, awaiting Gm mel",
				"Gm7 in F, on Gm mel",
				"Gm(maj7) in F, on Gm mel",
				"Gm7 in F, on Gm mel",
				"Gm(maj7) in F, on Gm mel",
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
	if got := tonalityName(analysis.Tune(leaves, phrasesOf(leaves))); got != "Gm nat/harm/mel" {
		t.Errorf("%s, want Gm nat/harm/mel", got)
	}
}

// groundsOf writes the ground of each chord, one chord to a line: the
// one the ear has installed there, and the one heard afterwards when
// the region started earlier.
//
//	Dm7 in D, in C afterwards
func groundsOf(c analysis.Changes) string {
	sensed := sense(c)
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
		minMaj7 = harmony.ChordMinorMajorSeventh
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
		// My Lucky Star: the II tonicised, and held four bars in the line
		// Gm7 Gm(maj7), G minor with a seventh that wanders, is still the
		// II of F, before Dm7 G7 C7 goes back to F6.
		"the two tonicised, not a modulation": {
			changesOf(false, f, maj7, a, min7, d, dom7, g, min7, c, dom7, g, min7, g, minMaj7, g, min7, g, minMaj7,
				d, min7, g, dom7, c, dom7, f, maj7),
			[]string{
				"Fmaj7 in F",
				"Am7 in F",
				"D7 in F",
				"Gm7 in F",
				"C7 in F",
				"Gm7 in F",
				"Gm(maj7) in F",
				"Gm7 in F",
				"Gm(maj7) in F",
				"Dm7 in F",
				"G7 in F",
				"C7 in F",
				"Fmaj7 in F",
			},
		},
		// Only Trust Your Heart: B7 Em7 tonicises the III of C, and Am7,
		// the VI of C, does not settle E minor before Dm7 G7 Cmaj7.
		"the three tonicised, not a modulation": {
			changesOf(false, c, maj7, b, dom7, e, min7, a, min7, d, min7, g, dom7, c, maj7),
			[]string{
				"Cmaj7 in C",
				"B7 in C",
				"Em7 in C",
				"Am7 in C",
				"Dm7 in C",
				"G7 in C",
				"Cmaj7 in C",
			},
		},
		// Every cadence defines the tonality, a plagal one included (En
		// Harmonie, tome 1, chapter 8 §3, p. 103): A♭maj7 E♭maj7, the IV
		// then the I, opens E flat, and two bars of it install it.
		"a plagal cadence opens a tonic": {
			changesOf(false, c, maj7, d, min7, g, dom7, c, maj7, ab, maj7, eb, maj7, eb, maj7),
			[]string{
				"Cmaj7 in C",
				"Dm7 in C",
				"G7 in C",
				"Cmaj7 in C",
				"A♭maj7 in C, in E♭ afterwards",
				"E♭maj7 in C, in E♭ afterwards",
				"E♭maj7 in E♭",
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
	if got := tonalityName(analysis.Tune(ch, phrasesOf(ch))); got != "F" {
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

// The end of My Way. The first phrase lands on C, held two bars after
// G7. After a modulation to F, F/C is the IV of C over a tonic pedal,
// not a tonic: with C it makes a plagal cadence, an amen that brings
// home back, and heard afterwards C starts with it.
func TestHomeOnSight(t *testing.T) {
	const c, f, g harmony.PitchClass = 0, 5, 7
	ch := changesOf(false, g, harmony.ChordDominantSeventh, c, harmony.ChordMajorSeventh, c, harmony.ChordMajorSeventh,
		g, harmony.ChordMinorSeventh,
		c, harmony.ChordDominantSeventh, f, harmony.ChordMajorSeventh, f, harmony.ChordMajorSeventh,
		f, harmony.ChordMajorTriad, c, harmony.ChordMajorTriad)
	ch.Chords[7].Bass = c
	want := []string{
		"G7 in ?, in C afterwards",
		"Cmaj7 in C",
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

// The imperfect cadence reaches the tonic on its fifth: C7 Fmaj7/C
// (En Harmonie, tome 1, chapter 8 §3.1, p. 103). Without a V before
// it, the same chord over the same bass is the IV of C over a tonic
// pedal, as at the end of My Way.
func TestTonicOnItsFifth(t *testing.T) {
	const c, f, g harmony.PitchClass = 0, 5, 7
	imperfect := changesOf(false, g, harmony.ChordMinorSeventh, c, harmony.ChordDominantSeventh, f, harmony.ChordMajorSeventh)
	imperfect.Chords[2].Bass = c
	if got := sense(imperfect)[2].Tonic; got == nil || tonalityName(got) != "F" {
		t.Errorf("Fmaj7/C after C7: tonic of %v, want F", got)
	}
	pedal := changesOf(false, c, harmony.ChordMajorSeventh, f, harmony.ChordMajorSeventh)
	pedal.Chords[1].Bass = c
	if got := sense(pedal)[1].Tonic; got != nil {
		t.Errorf("Fmaj7/C after Cmaj7: tonic of %s, want none", tonalityName(got))
	}
}

// The end of Only Trust Your Heart: Fmaj7/G before G7♭9 C6 is the V of
// C with its fourth, not the F the tune opened on coming back. The
// tune stops on the Cmaj7 that Dm7 G7 led to, and is in C.
func TestSlashChordEnding(t *testing.T) {
	const c, d, e, f, g, a, b harmony.PitchClass = 0, 2, 4, 5, 7, 9, 11
	const (
		maj7 = harmony.ChordMajorSeventh
		min7 = harmony.ChordMinorSeventh
		dom7 = harmony.ChordDominantSeventh
		six  = harmony.ChordMajorSixth
	)
	ch := changesOf(false, f, maj7, b, dom7, e, min7, a, min7, d, min7, g, dom7, c, maj7,
		f, maj7, g, dom7, c, six)
	ch.Chords[7].Bass = g
	if got := tonalityName(analysis.Tune(ch, phrasesOf(ch))); got != "C" {
		t.Errorf("%s, want C", got)
	}
}

// Where a tune stops, and so which tonality it is in.
func TestTuneEndings(t *testing.T) {
	const (
		c, db, d, eb, e, f, gb, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	const (
		min     = harmony.ChordMinorTriad
		minMaj7 = harmony.ChordMinorMajorSeventh
		min6    = harmony.ChordMinorSixth
		six     = harmony.ChordMajorSixth
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
		sus     = harmony.ChordDominantSeventhSus4
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
		// Sugar: its two sections of 8 bars end open on G7, the half
		// cadence, and come to rest nowhere. G♭7 Fm7 tonicises the IV in
		// the middle of the second, and the last G7 goes back to Cm7
		// across the loop.
		"the last cadence goes back to the start": {
			barsOf(true, bar{c, min7}, bar{d, halfDim, g, dom7}, bar{c, min7}, bar{g, dom7},
				bar{c, min7}, bar{c, min7}, bar{g, sus}, bar{g, dom7},
				bar{c, min7}, bar{gb, dom7}, bar{f, min7}, bar{eb, dom7},
				bar{d, halfDim}, bar{g, dom7}, bar{ab, dom7}, bar{g, dom7}),
			"Cm nat/harm/mel",
		},
		// Fly Me To The Moon opens on Am7 and comes back to it at bar 8,
		// but its first cadence goes to C: the Am7 is a VI, as in Blue
		// Skies. Its second section stops on Cmaj7 at bar 15, before the
		// turnaround.
		"a turnaround after the last tonic": {
			barsOf(true, bar{a, min7}, bar{d, min7}, bar{g, dom7}, bar{c, maj7},
				bar{f, maj7}, bar{b, halfDim}, bar{e, dom7}, bar{a, min7, a, dom7},
				bar{d, min7}, bar{g, dom7}, bar{c, maj7, f, dom7}, bar{e, min7, a, dom7},
				bar{d, min7}, bar{g, dom7}, bar{c, maj7}, bar{b, halfDim, e, dom7}),
			"C",
		},
		// Lullaby Of Birdland opens on F minor, and its section of 8
		// bars concludes on A♭maj7 at bar 7: Gm7♭5 C7 after it goes back
		// to the first chord, and is not played at the end. First and
		// last chords part, and F minor is heard the longer: it is in F
		// minor, where the chart declares A flat, as the ear hears it,
		// its A sections in F minor and its bridge in A flat.
		"the last tonic before the turnaround": {
			barsOf(true, bar{f, min, d, halfDim}, bar{g, dom7, c, dom7}, bar{f, min}, bar{bb, min7, eb, dom7},
				bar{c, min7, f, min7}, bar{bb, min7, eb, dom7}, bar{ab, maj7}, bar{g, halfDim, c, dom7}),
			"Fm nat/harm/mel",
		},
		// Yesterdays goes through B♭maj7 in its cycle of dominants, a V
		// with no two, at bar 14 of its 16: it stops in D minor.
		"a tonic passed through by a lone V": {
			barsOf(true, bar{d, min}, bar{e, halfDim, a, dom7}, bar{d, min}, bar{e, halfDim, a, dom7},
				bar{d, min, d, minMaj7}, bar{d, min7}, bar{b, halfDim}, bar{e, dom7},
				bar{a, dom7}, bar{d, dom7}, bar{g, dom7}, bar{c, dom7},
				bar{f, dom7}, bar{bb, maj7}, bar{e, halfDim}, bar{a, dom7}),
			"Dm nat/harm/mel",
		},
		// In a Sentimental Mood « est en Ré mineur pour se terminer dans
		// la tonalité de son relatif Fa majeur » (En Harmonie, tome 1,
		// chapter 10, p. 159): it holds Dm two bars, comes back to it by
		// A7, and concludes on Fmaj7. Its first and last chords part, and
		// D minor, heard the longer, is the tonality of the tune (tome 1,
		// chapter 8 §1.2, p. 100).
		"a tune that opens in its relative": {
			barsOf(true, bar{d, min, d, minMaj7}, bar{d, min7, d, min6}, bar{g, min7}, bar{g, min7, a, dom7},
				bar{d, min}, bar{d, dom7}, bar{g, min7, c, dom7}, bar{f, maj7}),
			"Dm nat/harm/mel",
		},
		// Guile's Theme comes back four times to C♯m7 by an aeolian
		// cadence, Amaj7 B7 C♯m7, and rests on it a bar each time. Its
		// one two five, to E, is the relative major tonicised in its
		// third section. (The helper spells C♯ as D♭.)
		"an aeolian cadence installs its I": {
			barsOf(true, bar{db, min7}, bar{a, maj7}, bar{gb, min7, ab, min7}, bar{db, min7},
				bar{a, maj7}, bar{gb, min7, ab, min7},
				bar{db, min7}, bar{a, maj7}, bar{b, dom7}, bar{ab, min7},
				bar{db, min7}, bar{a, maj7}, bar{b, dom7}, bar{ab, min7},
				bar{a, maj7}, bar{b, dom7}, bar{db, min7},
				bar{a, maj7}, bar{b, dom7}, bar{db, min7},
				bar{gb, min7}, bar{b, dom7}, bar{e, maj7},
				bar{a, maj7}, bar{b, dom7}, bar{db, min7},
				bar{a, maj7}, bar{b, dom7}, bar{db, min7}),
			"D♭m nat/harm/mel",
		},
		// The chart marks its end on C at bar 4, as Somewhere marks
		// its end two bars before the A♭ its last section concludes on:
		// the last chorus stops there.
		"the end the chart marks": {
			endingOn(barsOf(true, bar{d, min7, g, dom7}, bar{c, maj7}, bar{d, min7, g, dom7}, bar{c, maj7},
				bar{g, min7, c, dom7}, bar{f, maj7}, bar{c, dom7}, bar{f, maj7}), 5),
			"C",
		},
		// Unforgettable says its phrase again on the IV, as a blues goes
		// from I to IV: Cmaj7 at the end of the section is no tonic to
		// stop on, and Am7 D7 goes back to G.
		"a phrase said again on the IV": {
			barsOf(true, bar{g, maj7}, bar{a, min7, d, dom7}, bar{g, maj7}, bar{g, maj7},
				bar{c, maj7}, bar{d, min7, g, dom7}, bar{c, maj7}, bar{a, min7, d, dom7}),
			"G",
		},
		// Blue Skies opens on the same line from Am, and concludes on C6:
		// it is in C.
		"a tune that opens in its relative, again": {
			barsOf(true, bar{a, min, a, minMaj7}, bar{a, min7, a, min6}, bar{c, maj7, a, dom7}, bar{d, min7, g, dom7},
				bar{c, six}, bar{c, six}, bar{b, halfDim, e, dom7}),
			"C",
		},
	} {
		// A tune sets out from its first chord, when its first cadence
		// confirms it as the tonic (see analysis.FirstTonic).
		got := tonalityName(analysis.Tune(tc.changes, phrasesOf(tc.changes)))
		if first := analysis.FirstTonic(tc.changes, analysis.Blocks(tc.changes, analysis.Approaches(tc.changes))); first != nil && tonalityName(first) != got {
			got += ", setting out from " + tonalityName(first)
		}
		if got != tc.want {
			t.Errorf("%s: %s, want %s", name, got, tc.want)
		}
	}
}

// A minor tune that ends on its tonic made major stays minor: the
// picardy third. A tune that turns major before is major.
func TestPicardy(t *testing.T) {
	const c, d, e, f, g, a harmony.PitchClass = 0, 2, 4, 5, 7, 9
	const (
		six     = harmony.ChordMajorSixth
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
		// 'Round Midnight closes a section as it ends, on the major tonic,
		// and starts the next one again in the minor.
		"a section closing like the tune": {
			barsOf(false, bar{c, min6}, bar{d, halfDim, g, dom7}, bar{c, maj}, bar{c, min6},
				bar{d, halfDim, g, dom7}, bar{c, maj}),
			"Cm nat/harm/mel", true,
		},
		// Somewhere: the minor tonic comes by a IV-I in the middle of a
		// phrase, and a two five whose two is a m7, the II of the major,
		// has installed the major tonic before the end.
		"a minor chord in a major tune": {
			barsOf(false, bar{f, min6}, bar{c, min6}, bar{d, min7, g, dom7}, bar{c, six},
				bar{f, min6}, bar{c, min6}, bar{d, min7, g, dom7}, bar{c, six}),
			"C", false,
		},
		// Chega De Saudade: D minor for its first half, D major for its
		// second, and a stop on D6. Not a picardy third: the reading says
		// where it stops, the tune being as much minor as major.
		"a minor tune turning major": {
			barsOf(false, bar{d, min6}, bar{e, halfDim, a, dom7}, bar{d, min6}, bar{d, min6},
				bar{e, min7, a, dom7}, bar{d, six}, bar{d, six}, bar{e, min7, a, dom7}, bar{d, six}),
			"D", false,
		},
	} {
		c := tc.changes
		if got := tonalityName(analysis.Tune(c, phrasesOf(c))); got != tc.want {
			t.Errorf("%s: %s, want %s", name, got, tc.want)
		}
		if got := analysis.Picardy(c, phrasesOf(c)); got != tc.picardy {
			t.Errorf("%s: picardy third %v, want %v", name, got, tc.picardy)
		}
	}
}

// Home, where a tune sets out from, is where its first phrase comes to
// rest, however long that phrase is.
func TestHome(t *testing.T) {
	const (
		c, cs, d, eb, e, f, fs, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	const (
		maj     = harmony.ChordMajorTriad
		min     = harmony.ChordMinorTriad
		maj7    = harmony.ChordMajorSeventh
		six     = harmony.ChordMajorSixth
		min6    = harmony.ChordMinorSixth
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
		dim7    = harmony.ChordDiminishedSeventh
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    string
	}{
		// Autumn Leaves: B♭maj7 is passed through, Gm6 held two bars is
		// where the phrase rests.
		"Autumn Leaves": {
			barsOf(true, bar{c, min7}, bar{f, dom7}, bar{bb, maj7}, bar{eb, maj7},
				bar{a, halfDim}, bar{d, dom7}, bar{g, min6}, bar{g, min6}),
			"Gm nat/harm/mel",
		},
		// How Insensitive: a long sigh down from Dm, back to it.
		"How Insensitive": {
			barsOf(true, bar{d, min}, bar{d, min}, bar{cs, dim7}, bar{cs, dim7}, bar{c, min6}, bar{c, min6},
				bar{g, dom7}, bar{g, dom7}, bar{bb, maj7}, bar{bb, maj7}, bar{e, halfDim, a, dom7}, bar{d, min}),
			"Dm nat/harm/mel",
		},
		// There Will Never Be Another You opens on E♭maj7; Cm7,
		// held two bars after G7, is its VI.
		"There Will Never Be Another You": {
			barsOf(true, bar{eb, maj7}, bar{eb, maj7}, bar{d, halfDim, g, dom7}, bar{c, min7}, bar{c, min7},
				bar{f, min7}, bar{bb, dom7}, bar{eb, maj7}),
			"E♭",
		},
		// Just Friends opens on Cmaj7, the IV of G. Gmaj7 at bars 5 and
		// 11 confirms G in the middle of a section; the first conclusive
		// cadence is the last, G6 at bar 31.
		"Just Friends": {
			barsOf(true, bar{c, maj7}, bar{c, maj7}, bar{c, min7}, bar{f, dom7},
				bar{g, maj7}, bar{g, maj7}, bar{bb, min7}, bar{eb, dom7},
				bar{a, min7}, bar{d, dom7}, bar{g, maj7}, bar{e, min7},
				bar{a, dom7}, bar{a, dom7}, bar{a, min7}, bar{d, dom7, g, dom7},
				bar{c, maj7}, bar{c, maj7}, bar{c, min7}, bar{f, dom7},
				bar{g, maj7}, bar{g, maj7}, bar{bb, min7}, bar{eb, dom7},
				bar{a, min7}, bar{d, dom7}, bar{fs, halfDim, b, dom7}, bar{e, min7},
				bar{a, dom7}, bar{a, min7, d, dom7}, bar{g, six}, bar{d, min7, g, dom7}),
			"G",
		},
		// Cherokee: E♭maj7 held two bars is the IV, and the phrase comes
		// back to B♭6 by A♭7.
		"Cherokee": {
			barsOf(true, bar{bb, six}, bar{f, min7}, bar{bb, dom7}, bar{eb, maj7}, bar{eb, maj7},
				bar{ab, dom7}, bar{bb, six}, bar{c, min7, f, dom7}),
			"B♭",
		},
		// Honeysuckle Rose opens on a two, and C7 Gm7 is not a plagal
		// cadence to rest on: its first section concludes on F6 at bar 7.
		"Honeysuckle Rose": {
			barsOf(true, bar{g, min7, c, dom7}, bar{g, min7, c, dom7}, bar{g, min7, c, dom7}, bar{g, min7, c, dom7},
				bar{f, six, f, dom7}, bar{bb, six, c, dom7}, bar{f, six, bb, dom7}, bar{a, min7, d, dom7}),
			"F",
		},
	} {
		c := tc.changes
		if got := tonalityName(analysis.Home(c, phrasesOf(c))); got != tc.want {
			t.Errorf("%s: %s, want %s", name, got, tc.want)
		}
	}
}

// endingOn marks the end of the chart on change `i`, as a coda, a
// "Fine" or the player's end do.
func endingOn(c analysis.Changes, i int) analysis.Changes {
	c.End = i
	return c
}

// A dominant makes of the chord it resolves on a « degré I temporaire »,
// for that chord only (Siron, La partition intérieure, p. 341): Dm7
// after A7 in C. The tonic of the ground is not tonicised, nor the I of
// a plagal cadence, nor a chord with no major or minor triad: the
// Fm7♭5 of bar 9 of Tenderly, after Gm7♭5 C7♭9, is the II of E flat
// minor.
func TestTonicised(t *testing.T) {
	const c, db, d, eb, f, g, ab, a harmony.PitchClass = 0, 1, 2, 3, 5, 7, 8, 9
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	for name, tc := range map[string]struct {
		changes analysis.Changes
		want    []string
	}{
		"a II tonicised": {
			barsOf(false, bar{c, maj7}, bar{a, dom7}, bar{d, min7, g, dom7}, bar{c, maj7}),
			[]string{"Dm7 as a I of Dm harm/mel"},
		},
		"Tenderly, bars 1 to 9": {
			barsOf(false, bar{eb, maj7}, bar{ab, dom7}, bar{eb, min7}, bar{ab, dom7}, bar{f, min7},
				bar{db, dom7}, bar{eb, maj7}, bar{g, halfDim, c, dom7}, bar{f, halfDim}),
			nil,
		},
	} {
		var got []string
		for i, s := range sense(tc.changes) {
			if s.Tonicised != nil {
				got = append(got, changeName(tc.changes.Chords[i])+" as a I of "+tonalityName(s.Tonicised))
			}
		}
		if strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
			t.Errorf("%s:\ngot\n%s\nwant\n%s", name, strings.Join(got, "\n"), strings.Join(tc.want, "\n"))
		}
	}
}

// sameTonicName reports whether two sets of tonalities, nil or not,
// are on the same tonic.
func sameTonicName(a, b []harmony.Tonality) bool {
	return len(a) > 0 && len(b) > 0 && a[0].Tonic() == b[0].Tonic()
}
