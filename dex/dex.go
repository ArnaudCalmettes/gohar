package dex

import (
	"iter"
	"time"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A GameID names a game. The only free string in this package, and it
// designates a program rather than a musical notion.
type GameID string

// A Mark records that the player has shown something, when, how often,
// and where it first happened.
//
// # Nothing ever removes one
//
// An acquis stays acquis. A fact that undid a mark would punish a
// mistake, and this collection makes no reproach. What passes is time,
// and [Mark.Last] is what says so.
//
// Count is what makes spaced repetition possible. A notion seen once a
// month ago and one seen eight times do not call for the same interval,
// and a date alone cannot tell them apart.
type Mark struct {
	First time.Time
	Last  time.Time
	Count int

	// FirstIn is the game where the mark was first earned, which is
	// what makes a collection tellable: met in one game, produced in
	// another.
	FirstIn GameID
}

// Held reports whether the mark exists at all.
func (m Mark) Held() bool {
	return m.Count > 0
}

// An Entry is what the dex knows about one notion.
//
// The four marks are independent, not steps on a ladder. Recognising a
// shape by ear and laying it under the hands are two different skills,
// and the second can come first: exploring at the keyboard is how many
// musicians meet a mode. A dex that stacked them would refuse to credit
// what the player can actually do.
type Entry struct {
	Notion Notion

	// Met records that the notion sounded and was named to the player.
	Met Mark

	// Recognized records that the player identified it by ear.
	//
	// Not detailed by tonic, on purpose. A relative ear recognises a
	// shape whatever the key, by definition, so counting tonics here
	// would measure nothing.
	Recognized Mark

	// Produced records that the player laid it under the hands, one
	// mark per tonic.
	//
	// Detailed by tonic because the keyboard has no such invariance:
	// twelve tonics, twelve topographies, twelve fingerings. Laying a
	// shape everywhere is precisely what builds the relative ear this
	// is all for.
	//
	// The per entry summary is computed from this; the reverse cannot
	// be invented, which is why the detail is what gets stored.
	Produced map[harmony.PitchClass]Mark

	// Used records that the player chose it where the frame left the
	// colour open.
	//
	// This does not measure free improvisation. It measures that there
	// was a choice: a slot asking for a dominant and getting a lydian
	// dominant counts, the same mode imposed by the game does not,
	// however well played.
	Used Mark
}

// Tonics returns the set of tonics the notion has been produced in.
//
// A twelve bit mask for display, the circle of fifths filling in.
// Derived on demand: the dates and counts live per tonic.
func (e Entry) Tonics() harmony.PitchSet { panic("TODO") }

// A FactKind says what a game observed.
type FactKind uint8

const (
	// FactHeard: the game sounded a notion and named it to the player.
	FactHeard FactKind = iota + 1

	// FactNamed: the player identified a notion they were made to hear.
	FactNamed

	// FactProduced: the player played a notion.
	FactProduced

	// FactChosen: the player produced a notion where the frame left the
	// choice open.
	FactChosen

	// FactSounded: a notion actually rang out, without anything asking
	// for it.
	//
	// The freshness fact, and the answer to a trap of spaced
	// repetition: notions are not independent. Working a two five one
	// sounds a minor seventh, a dominant and a major seventh, and
	// without this the player would be offered tomorrow what they just
	// practised without knowing it.
	//
	// The recognition engine fills this in, never the author of a game,
	// who would report what was programmed rather than what rang.
	FactSounded
)

// A Fact is one thing a game observed about one notion.
//
// # What no fact says
//
// That an attempt failed, or that a performance was good. The first
// would punish, the second is a mark in disguise and would put the
// decision back on the game's side.
type Fact struct {
	Kind   FactKind
	Notion Notion

	// Tonic carries the pitch class for the facts that need one.
	// Essential to FactProduced and FactChosen, since the hand has
	// twelve topographies; recorded for the others and unused.
	Tonic harmony.PitchClass

	// Correct says whether an identification was right. FactNamed only.
	Correct bool

	// Target is what the slot was asking for. FactChosen only, and the
	// one place the dex has to understand a target: it is how it checks
	// that the colour really was open.
	Target harmony.Target
}

// A Report is what a game sends once an activity is over.
//
// One report per activity, not a stream. A stream would be faithful and
// deafening, and would leave the dex to sort the noise.
//
// It carries no score, no overall success and no difficulty. A game's
// vocabulary does not come in here.
type Report struct {
	Game     GameID
	Activity string
	At       time.Time
	Facts    []Fact
}

// A Change is something a report made happen, so that a game can
// celebrate it.
//
// Discovery is the one worth showing: a wild lydian dominant appears.
type Change struct {
	Notion    Notion
	Discovery bool

	// Marked is the mark that was earned, or zero when the report only
	// refreshed what was already there.
	Marked FactKind

	// Tonic is the tonic that was added, for a production.
	Tonic    harmony.PitchClass
	NewTonic bool
}

// A Dex is one player's collection.
type Dex struct {
	entries map[Notion]*Entry
}

// New returns an empty dex.
func New() *Dex { panic("TODO") }

// Apply folds a report into the dex and returns what changed.
//
// This is where the decision lives, and it lives here alone: games send
// what happened, never what it means. A rule that changes changes in
// one place.
func (d *Dex) Apply(r Report) []Change { panic("TODO") }

// Look returns what is known about a notion.
func (d *Dex) Look(n Notion) (Entry, bool) { panic("TODO") }

// Collection iterates over the entries the player has met.
//
// The view opened for pleasure: what I hold, what I have crossed. No
// notion of time in it.
func (d *Dex) Collection() iter.Seq[Entry] { panic("TODO") }

// A Due is a notion whose recall has come round, and how badly.
type Due struct {
	Notion Notion

	// Tonic is set when it is one tonic of a production that is
	// cooling rather than the whole entry.
	Tonic    harmony.PitchClass
	ByTonic  bool
	Overdue  time.Duration
}

// Cooling returns what is due for review, most overdue first.
//
// The other view, the one the player goes looking for. Someone asking
// for an honest state of things accepts a frank answer; the same
// account shown unasked would be a telling off, which is why nothing
// here surfaces on its own.
//
// # Undecided
//
// The intervals. Spaced repetition stretches them with each success,
// and the right curve for a keyboard is a matter for the ear rather
// than for a formula copied from flashcards.
func (d *Dex) Cooling(now time.Time) []Due { panic("TODO") }

// Discoverable returns notions never met whose components the player
// already knows, so that a game can introduce them without losing
// anyone.
func (d *Dex) Discoverable() []Notion { panic("TODO") }

// Visible returns what the player is allowed to see, silhouettes
// included.
//
// A notion never crossed does not show. One does not chase what one
// cannot see, and a beginner facing thirty five empty slots reads
// everything they do not know.
func (d *Dex) Visible() []Notion { panic("TODO") }
