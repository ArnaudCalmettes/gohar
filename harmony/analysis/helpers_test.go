package analysis_test

import (
	"strings"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// changesOf builds changes of one bar each from chords written as a
// root and a pattern, a nil pattern for a silence.
func changesOf(loops bool, chords ...any) analysis.Changes {
	c := analysis.Changes{Loops: loops}
	for i := 0; i+1 < len(chords); i += 2 {
		root := chords[i].(harmony.PitchClass)
		ch := analysis.Change{Bass: root, Length: 4 * analysis.TicksPerBeat}
		if p, ok := chords[i+1].(harmony.ChordPattern); ok {
			ch.Chord = harmony.Chord{Root: root, Pattern: p}
		} else {
			ch.Silent = true
		}
		ch.Start = analysis.Ticks(len(c.Chords)) * ch.Length
		c.Chords = append(c.Chords, ch)
	}
	return c
}

// The names below write the expectations of the tests as a musician
// reads a chart, until naming names chords and tonalities.

var noteNames = [12]string{"C", "D♭", "D", "E♭", "E", "F", "F♯", "G", "A♭", "A", "B♭", "B"}

// changeName names a change as a chart writes it: Cm7, B♭/D, N.C. for a
// silence.
func changeName(ch analysis.Change) string {
	if ch.Silent {
		return "N.C."
	}
	name := chordName(ch.Chord)
	if ch.Inverted() {
		name += "/" + noteNames[ch.Bass]
	}
	return name
}

// chordName names the chords the tests use, "?" for any other quality.
func chordName(ch harmony.Chord) string {
	flatNine, _ := harmony.NewChordPattern(0, 4, 7, 10, 13)
	nine, _ := harmony.NewChordPattern(0, 4, 7, 10, 14)
	quality, ok := map[harmony.ChordPattern]string{
		harmony.ChordMajorTriad:          "",
		harmony.ChordMinorSixth:          "m6",
		harmony.ChordMajorSeventh:        "maj7",
		harmony.ChordMinorSeventh:        "m7",
		harmony.ChordDominantSeventh:     "7",
		harmony.ChordHalfDiminished:      "m7♭5",
		harmony.ChordDiminishedSeventh:   "dim7",
		harmony.ChordDominantSeventhSus4: "7sus4",
		flatNine:                         "7♭9",
		nine:                             "9",
	}[ch.Pattern]
	if !ok {
		quality = "?"
	}
	return noteNames[ch.Root] + quality
}

// tonalityName names tonalities on one tonic as the analyse command
// does: "C" for the major alone, "Dm" for the natural minor alone, "Fm
// harm", "Dm harm/mel", and "D♭ M/m mel" when the major and a minor are
// both announced.
func tonalityName(ts []harmony.Tonality) string {
	if len(ts) == 0 {
		return "none"
	}
	tonic := noteNames[ts[0].Tonic()]
	major := false
	var minors []string
	for _, t := range ts {
		switch s, _ := harmony.NamedScaleOf(t.Pattern()); s {
		case harmony.NamedMajor:
			major = true
		case harmony.NamedNaturalMinor:
			minors = append(minors, "nat")
		case harmony.NamedHarmonicMinor:
			minors = append(minors, "harm")
		case harmony.NamedMelodicMinor:
			minors = append(minors, "mel")
		}
	}
	switch {
	case len(minors) == 0:
		return tonic
	case major:
		return tonic + " M/m " + strings.Join(minors, "/")
	case len(minors) == 1 && minors[0] == "nat":
		return tonic + "m"
	}
	return tonic + "m " + strings.Join(minors, "/")
}

// sharpened spells the root of a chord name with a sharp, as a chart
// writes a chord whose bass walks up: G♯dim7 rather than A♭dim7.
func sharpened(name string) string {
	for flat, sharp := range map[string]string{"D♭": "C♯", "E♭": "D♯", "A♭": "G♯", "B♭": "A♯"} {
		if strings.HasPrefix(name, flat) {
			return sharp + strings.TrimPrefix(name, flat)
		}
	}
	return name
}
