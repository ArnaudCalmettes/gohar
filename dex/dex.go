package dex

import (
	"cmp"
	"iter"
	"maps"
	"slices"
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

// earn records one more demonstration and reports whether it was the
// first.
//
// Reports are expected in order. One arriving late only fails to move
// Last backwards; First and FirstIn stay those of the first report
// applied, which is the one the player lived first on this dex.
func (m *Mark) earn(at time.Time, in GameID) bool {
	m.Count++
	if m.Count == 1 {
		m.First, m.Last, m.FirstIn = at, at, in
		return true
	}
	m.refresh(at)
	return false
}

// refresh moves Last forward on a mark already held, and does nothing
// else: freshness is not a demonstration, so Count stays put.
func (m *Mark) refresh(at time.Time) {
	if m.Held() && at.After(m.Last) {
		m.Last = at
	}
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
	// however well played. Whether the colour was open is the game's
	// judgement, not the dex's: see [FactChosen].
	Used Mark

	// Overheard records that the notion rang out under the player's
	// hands or in their ears without anyone naming it.
	//
	// Not a demonstration, and not one of the four marks: a notion only
	// overheard is not met, stays out of the collection and shows as a
	// silhouette. It is the musician who plays something because it
	// sounds classy, long before learning it has a name. The silhouette
	// is the moment a game can say "this one is called phrygian natural
	// six, want to catch it?", and the naming that follows is a
	// FactHeard, which is the discovery.
	//
	// Count here is how often it rang unnamed, which a game may want to
	// tell: "you have played this twelve times without knowing".
	Overheard Mark
}

// Tonics returns the set of tonics the notion has been produced in.
//
// A twelve bit mask for display, the circle of fifths filling in.
// Derived on demand: the dates and counts live per tonic.
func (e Entry) Tonics() harmony.PitchSet {
	var s harmony.PitchSet
	for tonic := range e.Produced {
		s = s.With(tonic)
	}
	return s
}

// held reports whether any mark exists on the entry.
func (e *Entry) held() bool {
	return e.Met.Held() || e.Recognized.Held() || e.Used.Held() ||
		len(e.Produced) > 0
}

// refresh moves every held mark forward, one tonic at a time for
// production: what rang out was played somewhere, and the dex cannot
// tell which topography, so they all count as warm.
func (e *Entry) refresh(at time.Time) {
	e.Met.refresh(at)
	e.Recognized.refresh(at)
	e.Used.refresh(at)
	for tonic, m := range e.Produced {
		m.refresh(at)
		e.Produced[tonic] = m
	}
}

// produce earns the mark of one tonic and reports whether that tonic is
// new.
func (e *Entry) produce(tonic harmony.PitchClass, at time.Time, in GameID) bool {
	if e.Produced == nil {
		e.Produced = map[harmony.PitchClass]Mark{}
	}
	m := e.Produced[tonic]
	first := m.earn(at, in)
	e.Produced[tonic] = m
	return first
}

// A FactKind says what a game observed.
type FactKind uint8

const (
	// FactHeard: the game sounded a notion and named it to the player.
	FactHeard FactKind = iota + 1

	// FactNamed: the player identified a notion they were made to hear.
	//
	// Sent on an identification only. A wrong answer is sent as
	// nothing at all: it gets corrected, it does not get learnt. What
	// is learnt is the correction the player produces afterwards, and
	// that one arrives as a FactNamed like any other.
	FactNamed

	// FactProduced: the player played a notion.
	FactProduced

	// FactChosen: the player produced a notion where the frame left the
	// choice open. Marks Used and the production of the tonic: who
	// chose, played.
	//
	// The game decides that the choice was open, the dex takes its word
	// for it. For a mode, that means the slot did not name the mode and
	// its characteristic degrees actually rang: a player sounding only
	// the notes of G7 chose no colour at all.
	FactChosen

	// FactSounded: a notion actually rang out, without anything asking
	// for it. Moves Last on the marks already held, never their Count,
	// and earns [Entry.Overheard]. On a notion never met, that makes a
	// silhouette appear, not a discovery.
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
// would punish, the second is a judgement in disguise: the player
// assesses themselves, and gets an honest account only when they ask
// for one.
type Fact struct {
	Kind   FactKind
	Notion Notion

	// Tonic carries the pitch class for the facts that need one.
	// Essential to FactProduced and FactChosen, since the hand has
	// twelve topographies; recorded for the others and unused.
	Tonic harmony.PitchClass
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
	//
	// FactSounded here means a silhouette just appeared: a notion never
	// met was overheard for the first time.
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
func New() *Dex {
	return &Dex{entries: map[Notion]*Entry{}}
}

// Apply folds a report into the dex and returns what changed.
//
// This is where the marking rule lives, and it lives here alone: games
// send what happened, never what it means. A rule that changes changes
// in one place.
//
// What the dex does not do is judge. Whether an answer was right,
// whether a colour was open, whether a detour held together: the game
// decides, and sends the fact or does not. The dex never reads a target.
//
// Only what is new comes back: a discovery, a mark earned for the first
// time, a tonic added. A refresh is not something to celebrate.
func (d *Dex) Apply(r Report) []Change {
	var changes []Change
	for _, f := range r.Facts {
		if f.Notion.IsZero() {
			continue
		}

		e, ok := d.entries[f.Notion]
		if !ok {
			e = &Entry{Notion: f.Notion}
		}

		if f.Kind == FactSounded {
			d.entries[f.Notion] = e
			e.refresh(r.At)
			first := e.Overheard.earn(r.At, r.Game)
			if first && !e.held() {
				changes = append(changes, Change{Notion: f.Notion, Marked: FactSounded})
			}
			continue
		}

		c := Change{Notion: f.Notion, Discovery: !e.held()}
		switch f.Kind {
		case FactHeard:
			if e.Met.earn(r.At, r.Game) {
				c.Marked = FactHeard
			}
		case FactNamed:
			if e.Recognized.earn(r.At, r.Game) {
				c.Marked = FactNamed
			}
		case FactProduced:
			c.Tonic = f.Tonic
			c.NewTonic = e.produce(f.Tonic, r.At, r.Game)
			if c.NewTonic && len(e.Produced) == 1 {
				c.Marked = FactProduced
			}
		case FactChosen:
			c.Tonic = f.Tonic
			c.NewTonic = e.produce(f.Tonic, r.At, r.Game)
			if e.Used.earn(r.At, r.Game) {
				c.Marked = FactChosen
			}
		default:
			// An unknown kind marks nothing. Opening the entry anyway
			// would show a notion the player never demonstrated.
			continue
		}

		d.entries[f.Notion] = e
		if c.Discovery || c.Marked != 0 || c.NewTonic {
			changes = append(changes, c)
		}
	}
	return changes
}

// Look returns what is known about a notion.
//
// The entry is a copy, its production map included: what a game does
// with it never reaches the collection.
func (d *Dex) Look(n Notion) (Entry, bool) {
	e, ok := d.entries[n]
	if !ok {
		return Entry{}, false
	}
	out := *e
	out.Produced = maps.Clone(e.Produced)
	return out, true
}

// Collection iterates over the entries the player has met.
//
// The view opened for pleasure: what I hold, what I have crossed. No
// notion of time in it.
//
// In the order of the written forms, so that two displays of the same
// dex agree. Musical order is the display's business.
func (d *Dex) Collection() iter.Seq[Entry] {
	return func(yield func(Entry) bool) {
		for _, n := range d.notions() {
			e, _ := d.Look(n)
			if !e.held() {
				// Overheard only: a silhouette, not something held.
				continue
			}
			if !yield(e) {
				return
			}
		}
	}
}

// notions returns every notion with an entry, overheard ones included,
// sorted by written form.
func (d *Dex) notions() []Notion {
	ns := slices.Collect(maps.Keys(d.entries))
	slices.SortFunc(ns, func(a, b Notion) int {
		return cmp.Compare(a.String(), b.String())
	})
	return ns
}

// A Due is a notion whose recall has come round, and how badly.
type Due struct {
	Notion Notion

	// Tonic is set when it is one tonic of a production that is
	// cooling rather than the whole entry.
	Tonic   harmony.PitchClass
	ByTonic bool
	Overdue time.Duration
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
//
// What was crossed, and as silhouettes what was only overheard. The
// other source of silhouettes, [Notion.Components], is undecided.
func (d *Dex) Visible() []Notion {
	return d.notions()
}
