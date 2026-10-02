package harmony

// An ApproachKind says how a chord prepares the chord it leads to: as
// its dominant, its chromatic dominant, its diminished chord, its
// suspension or its two. The zero value means it does not prepare it.
//
// "Approach" is the English for what the French teaching calls a
// préparation (see docs/grilles.md). An [Approach] is the path a player
// took through a slot; its kinds are what each step of such a path, or
// of a chart, does for the step after it.
//
// # Between two chords, not in a key
//
// Whether a dominant is the dominant of the key or a secondary one
// depends on the key, which only the analysis of a whole chart knows.
// A kind only says what one chord does for the next, whatever their
// qualities: En Harmonie's secondary dominant prepares "any chord of
// arrival, whatever its quality or its function".
//
// # What two chords cannot tell
//
// A diminished chord between two chords a tone apart is also a passing
// chord, its bass walking by semitones. That reading needs three chords
// and their basses, and is the analysis's to make: a diminished
// approach here only says the chord is a dominant without its root.
// When the bass does walk, the two readings stand together, each from
// its own calculation, and both are exact.
//
// # A bit field
//
// For the same reason as [Function]: a reading is not exclusive, and a
// kind to come, a parallel chord, will sit beside others.
type ApproachKind uint16

// NoApproach records that a chord does not prepare the next.
const NoApproach ApproachKind = 0

const (
	// DominantApproach is a chord with the tritone of a dominant, a
	// fifth above the chord it leads to: A7 before Dm7.
	DominantApproach ApproachKind = 1 << iota

	// ChromaticApproach is a chord with the tritone of a dominant, a
	// semitone above the chord it leads to: E♭7 before Dm7. It shares
	// its tritone with the dominant a tritone away, which is what makes
	// it one.
	ChromaticApproach

	// DiminishedApproach is a diminished seventh chord with a note a
	// semitone under the root of the chord it leads to: C♯dim7 before
	// Dm7. It is the dominant seventh flat nine of that chord without
	// its root: C♯dim7 is A7(♭9) with no A.
	DiminishedApproach

	// SuspensionApproach is a suspended chord before the chord on the
	// same root that has its third: A7sus4 before A7. The fourth waits
	// and falls to the third.
	SuspensionApproach

	// TwoApproach is a minor seventh or half diminished chord a fifth
	// above a dominant: Em7 before A7, the two of a two five. Also a
	// semitone above it, as the two of its tritone twin resolving
	// straight onto it: B♭m7 before A7, En Harmonie's chromatic
	// subdominant.
	TwoApproach

	// PlagalApproach is a subdominant before the tonic it concludes on:
	// the IV of any quality, Fmaj7, F7, Fm6 before C, the plagal amen;
	// and the ♭VII7, B♭7 before C, the minor plagal in disguise (with F
	// in the bass it is Fm6); and the IIm7♭5, Dm7♭5 before C, « se
	// confondant avec » Fm6, the same four notes over its sixth (En
	// Harmonie, tome 2, p. 34, on the Gm7♭5 Fmaj7 of I'm Old Fashioned).
	// The tonic must be a chord that can be one.
	// A plagal cadence concludes as a V-I does, but draws less: what it
	// may do in the analysis of a chart is the analysis's to say.
	//
	// The other ♭VII-I of En Harmonie, B♭maj7 from the mixolydian and
	// B♭m7 from the phrygian, are not read from two chords: without
	// their mode, they are the stepwise motion of any tonal chart, Em7
	// Fmaj7 in C, Dm7 Em7. They wait for a chart that carries its
	// modes.
	PlagalApproach
)

// Has reports whether `k` includes the given kind.
func (k ApproachKind) Has(other ApproachKind) bool {
	return k&other != 0
}

// String writes the kinds as the analyse command labels them, the
// degree the chord stands on for the chord it prepares: V, ♭II (the
// chromatic dominant), ° (the diminished chord), sus, II, IV (the
// plagal subdominant, the ♭VII7 included). Several kinds are joined
// by a slash; no approach is "-".
func (k ApproachKind) String() string {
	if k == NoApproach {
		return "-"
	}
	var out string
	for _, n := range []struct {
		kind ApproachKind
		name string
	}{
		{TwoApproach, "II"},
		{SuspensionApproach, "sus"},
		{DominantApproach, "V"},
		{ChromaticApproach, "♭II"},
		{DiminishedApproach, "°"},
		{PlagalApproach, "IV"},
	} {
		if k.Has(n.kind) {
			if out != "" {
				out += "/"
			}
			out += n.name
		}
	}
	return out
}

// ApproachOf returns how `from` prepares `to` when it comes right
// before it.
//
// # The surface, not the form beneath
//
// Between two dominants, the kind reads what is played. In Sophisticated
// Lady, F7 before E7 is a chromatic approach; En Harmonie hears it as the
// dominant of the B♭7 that E7 stands in for. Both are true, since E7 and
// B♭7 share their tritone. Recovering the form beneath, a chain of
// dominants by fifths, is the analysis's business, and it can: a
// chromatic approach to a dominant is a dominant approach to its
// tritone twin.
//
// Likewise Fm7 before B♭7 is a two, whether the B♭7 then goes to E♭ or
// to C: which of the two five of E flat or the IVm7-♭VII7 of C it was,
// the chord of arrival decides, and the analysis reads it.
func ApproachOf(from, to Chord) ApproachKind {
	k := approachTo(from, to.Root)
	switch {
	case isDiminishedSeventh(from.Pattern) && from.Set().Contains(to.Root.Transpose(-1)):
		k |= DiminishedApproach
	case isSuspension(from.Pattern) && from.Root == to.Root && hasThird(to.Pattern):
		k |= SuspensionApproach
	case isTwo(from.Pattern) && isDominant(to.Pattern) &&
		(from.Root == to.Root.Transpose(7) || from.Root == to.Root.Transpose(1)):
		k |= TwoApproach
	}
	if isTonicChord(to.Pattern) && hasThird(from.Pattern) && !isDiminishedSeventh(from.Pattern) &&
		(from.Root == to.Root.Transpose(5) ||
			from.Root == to.Root.Transpose(10) && isDominant(from.Pattern) ||
			from.Root == to.Root.Transpose(2) && from.Pattern.Tetrad() == ChordHalfDiminished) {
		k |= PlagalApproach
	}
	return k
}

// isTonicChord reports whether a pattern can be a tonic: a major or
// minor triad, maj7, 6, m6, m(maj7), and the m7 charts write for a
// minor tonic.
func isTonicChord(p ChordPattern) bool {
	switch p.Tetrad() {
	case ChordMajorTriad, ChordMajorSeventh, ChordMajorSixth,
		ChordMinorTriad, ChordMinorSixth, ChordMinorMajorSeventh, ChordMinorSeventh:
		return true
	}
	return false
}

// approachTo reads the dominant approaches of a chord to a target known
// by its root alone, as a [Target] may be.
func approachTo(from Chord, to PitchClass) ApproachKind {
	if !isDominant(from.Pattern) {
		return NoApproach
	}
	switch from.Root {
	case to.Transpose(7):
		return DominantApproach
	case to.Transpose(1):
		return ChromaticApproach
	}
	return NoApproach
}

// isDominant reports whether a pattern holds the tritone of a dominant,
// between a major third and a minor seventh.
func isDominant(p ChordPattern) bool {
	return p.HasAll(IntMajorThird, IntMinorSeventh)
}

func hasThird(p ChordPattern) bool {
	return p.HasAny(IntMajorThird, IntMinorThird)
}

// isSuspension reports whether a pattern holds a fourth or a second in
// place of its third.
func isSuspension(p ChordPattern) bool {
	return !hasThird(p) && p.HasAny(IntPerfectFourth, IntMajorSecond)
}

// isDiminishedSeventh reports whether the tetrad of a pattern is the
// diminished seventh. The major seventh of the diminished chord of the
// scale 2-1 sits at fourteen, outside the tetrad, and does not change
// that.
func isDiminishedSeventh(p ChordPattern) bool {
	return p.Tetrad() == ChordDiminishedSeventh
}

// isTwo reports whether the tetrad of a pattern is a minor seventh or a
// half diminished seventh: the qualities of a two.
func isTwo(p ChordPattern) bool {
	t := p.Tetrad()
	return t == ChordMinorSeventh || t == ChordMinorSeventhNo5 || t == ChordHalfDiminished
}
