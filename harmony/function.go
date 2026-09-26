package harmony

// A Function is the harmonic role a chord or a mode assumes: tonic,
// subdominant or dominant. The zero value means none.
//
// # Why it lives here
//
// Harmonic function is structure, not spelling. It first appeared on
// the mode catalogue in the naming package, because mode names belong
// there, and it moved here once it was clear the core needed it too:
// deciding whether an approach reached a target in the right role is a
// question about harmony, and a [Target] has to be able to ask it
// without depending on a package of names.
//
// # What it says and what it does not
//
// Nothing about the chord's quality. A minor tetrad in a tonic role is
// Tonic, not some minor variant of it: the minor is already in the chord
// symbol, and repeating it here would split one function into two
// values that behave alike.
//
// # Why a bit field
//
// A role is not exclusive. A plain major triad serves as tonic,
// subdominant and dominant depending on where it sits, and a mode can
// carry more than one recorded role. Three separate values would force
// a choice the harmony does not make.
type Function uint8

// NoFunction records that no role is known.
//
// Declared apart from the three roles on purpose. Placed in the same
// block, the iota would start at zero and shift every role up by one
// bit, which works until the day a value is persisted or compared to a
// literal.
const NoFunction Function = 0

// The three roles, one bit each.
const (
	Tonic Function = 1 << iota
	Subdominant
	Dominant
)

// Has reports whether `f` includes the given role.
func (f Function) Has(other Function) bool {
	return f&other != 0
}

// FunctionOf returns the roles a tetrad can assume, or NoFunction for
// one the table does not know.
//
// # Where this comes from
//
// Transcribed from the ChordType table of the article « Reconnaître les
// accords », which is the published source.
//
// # What it cannot say
//
// Which role a chord is playing right now. That depends on what it
// resolves to: a D flat seventh is a dominant of C only in relation to
// C. This table gives the roles a tetrad is capable of, and
// [Target.Resolve] decides which one it took.
//
// # Suspensions
//
// Both sus chords are subdominants, never dominants. The seventh sus
// four has no third and so no tritone; removing that interval is what
// turned it into a plagal chord rather than a dominant waiting to
// resolve.
func FunctionOf(tetrad ChordPattern) Function {
	return tetradFunctions[tetrad]
}

var tetradFunctions = map[ChordPattern]Function{
	ChordMajorSeventh:       Tonic | Subdominant,
	ChordMajorSeventhNo5:    Tonic | Subdominant,
	ChordMajorSeventhSharp5: Tonic | Subdominant,

	ChordDominantSeventh:      Dominant,
	ChordDominantSeventhNo5:   Dominant,
	ChordDominantSeventhFlat5: Dominant,

	// Not in the published table, which predates this tetrad. Filed with
	// the other sevenths on a major third: same tritone, same role.
	ChordDominantSeventhSharp5: Dominant,

	ChordMajorSixth: Tonic,
	ChordMajorTriad: Tonic | Subdominant | Dominant,

	// The published table lists only the fifthless form under this
	// type. The full chord is filed with it, the fifth adding no role.
	ChordMinorMajorSeventh:    Tonic,
	ChordMinorMajorSeventhNo5: Tonic,

	ChordMinorSeventh:    Subdominant,
	ChordMinorSeventhNo5: Subdominant,
	ChordHalfDiminished:  Subdominant,

	ChordMinorSixth: Tonic,
	ChordMinorTriad: Tonic | Subdominant,

	ChordDiminishedSeventh: Dominant,
	ChordDiminishedTriad:   Dominant | Subdominant,
	ChordAugmentedTriad:    Dominant,

	ChordSus2:                Subdominant,
	ChordSus4:                Subdominant,
	ChordDominantSeventhSus2: Subdominant,
	ChordDominantSeventhSus4: Subdominant,
}
