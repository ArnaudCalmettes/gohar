package lessons

import "slices"

// A Step is one of the walker's gestures in a lesson (see the table of
// "Le bonhomme professeur" in docs/debutants.md). The Runner calls its
// methods; each but begin says whether the step is over. "Jouer
// ensemble", over the band, comes with the band's lessons.
type Step interface {
	said() string // the step's phrase, "" for none
	begin(s Stage)
	pressed(s Stage, held []int) bool  // a key went down, the last of `held`
	released(s Stage, held []int) bool // a key went up, `held` the keys still down
	phraseEnded(s Stage) bool          // the walker's phrase is over
	read() bool                        // the bubble is read: a key, or its reading time; no side effect
}

// still is the step that waits for nothing it does not name.
type still struct{}

func (still) pressed(Stage, []int) bool  { return false }
func (still) released(Stage, []int) bool { return false }
func (still) phraseEnded(Stage) bool     { return false }
func (still) read() bool                 { return false }

// Say is "Dire": a bubble, over once read.
type Say struct {
	still
	Phrase string
}

func (st Say) said() string  { return st.Phrase }
func (st Say) begin(s Stage) { s.Say(st.Phrase) }
func (Say) read() bool       { return true }

// Show is "Montrer": the walker lights keys, and says what they are.
// It prepares the next step: over once its bubble is read, it leaves
// the keys lit for the step that asks for them.
type Show struct {
	still
	Phrase string
	PCs    []int // the pitch classes lit
}

func (st Show) said() string { return st.Phrase }

func (st Show) begin(s Stage) {
	s.Say(st.Phrase)
	s.Light(st.PCs)
}

func (Show) read() bool { return true }

// Play is "Jouer": the walker plays a phrase, over when it ends.
type Play struct {
	still
	Phrase string // said while he plays, or ""
	Keys   []int  // MIDI numbers
}

func (st Play) said() string { return st.Phrase }

func (st Play) begin(s Stage) {
	if st.Phrase != "" {
		s.Say(st.Phrase)
	}
	s.Play(st.Keys)
}

func (Play) phraseEnded(Stage) bool { return true }

// Ask is "Demander": one target, over at the right answer. A miss sounds
// the casserole and shows the answer, which stays shown until played:
// the lesson waits, it never moves on past a wrong note. The keys a Show
// lit stay lit until the answer.
type Ask struct {
	Phrase string
	Want   Target
	Chords []string // written while it asks, if any: "the root of this one"

	// Teasing marks a request nobody misses in good faith, "press the
	// two black keys of a pair": misses in a row there are on purpose,
	// and the walker gets annoyed (see Stage.Tease).
	Teasing bool

	partial bool // a group begun, to judge again when its keys go up
	misses  int  // in a row, on a teasing request
}

// The misses in a row that annoy the walker, more and more (see the
// easter egg of "Rendre le cours amusant" in docs/debutants.md).
const (
	missesHey       = 4
	missesOnPurpose = 7
	missesFed       = 10
)

func (st *Ask) said() string { return st.Phrase }

func (st *Ask) begin(s Stage) {
	st.partial, st.misses = false, 0
	s.Say(st.Phrase)
	if st.Chords != nil {
		s.Chords(st.Chords)
	}
}

func (st *Ask) pressed(s Stage, held []int) bool {
	switch st.Want.Judge(held) {
	case Hit:
		st.misses = 0
		s.Light(nil)
		s.Hit()
		return true
	case Partial:
		st.partial = true
		return false
	}
	st.partial = false
	st.miss(s)
	return false
}

// released judges a group left unfinished: its keys all up, it is a
// miss.
func (st *Ask) released(s Stage, held []int) bool {
	if st.partial && len(held) == 0 {
		st.partial = false
		st.miss(s)
	}
	return false
}

// miss sounds the casserole, and on a teasing request counts the misses
// in a row.
func (st *Ask) miss(s Stage) {
	s.Miss(st.Want.Hint(), 0)
	if !st.Teasing {
		return
	}
	st.misses++
	switch st.misses {
	case missesHey:
		s.Tease(Hey)
	case missesOnPurpose:
		s.Tease(OnPurpose)
	case missesFed:
		s.Tease(FedUp)
		st.misses = 0 // he will be back, and count again
	}
}

func (*Ask) phraseEnded(Stage) bool { return false }
func (*Ask) read() bool             { return false }

// Repeat is "Répéter après lui": the walker plays a phrase, then the
// player plays it back, note after note, without a mistake. A wrong
// note sounds the casserole, and the walker plays the phrase again, to
// be played back again from its start: "do sol do" is not found among
// the notes of a scale.
//
// Without Keys, the walker only says the phrase, "C, E, C", and the
// player plays it from the words. A wrong note then shows every key of
// the phrase, lit until it is played through, from its start again.
type Repeat struct {
	Phrase string
	Keys   []int    // what the walker plays, MIDI numbers; nil to say it only
	Want   []Target // a target for each note

	at int // the note expected next, -1 while the walker plays
}

func (st *Repeat) said() string { return st.Phrase }

func (st *Repeat) begin(s Stage) {
	s.Say(st.Phrase)
	s.Light(nil)
	if st.Keys == nil {
		st.at = 0
		return
	}
	st.at = -1
	s.Play(st.Keys)
}

func (st *Repeat) phraseEnded(Stage) bool {
	if st.Keys != nil {
		st.at = 0
	}
	return false
}

func (st *Repeat) pressed(s Stage, held []int) bool {
	if st.at < 0 {
		return false // he is still playing: the keys sound, nothing counts
	}
	if st.Want[st.at].Judge(held) != Hit {
		if st.Keys == nil {
			s.Miss(st.hint(), st.at)
			st.at = 0
			return false
		}
		s.Miss(nil, st.at)
		st.at = -1
		s.Play(st.Keys)
		return false
	}
	if st.Keys != nil {
		s.Light(nil) // what a Show lit, out at the first note
	}
	st.at++
	if st.at < len(st.Want) {
		return false
	}
	if st.Keys == nil {
		s.Light(nil) // the phrase's keys, lit after a miss
	}
	s.Hit()
	return true
}

// hint is every pitch class the phrase wants, once each, in order.
func (st *Repeat) hint() []int {
	var pcs []int
	for _, w := range st.Want {
		for _, p := range w.Hint() {
			if !slices.Contains(pcs, p) {
				pcs = append(pcs, p)
			}
		}
	}
	return pcs
}

func (*Repeat) released(Stage, []int) bool { return false }
func (*Repeat) read() bool                 { return false }

// Hang is a phrase left hanging, the gag of lesson 1.2 (see
// docs/debutants/chapitre-1.md): said, and played back from the words
// as a Repeat without Keys is, but its last note is held by the game
// once the player lets it go. The walker frets, and the note that
// resolves it, the first one of its pitch class above, blinks; played,
// it lets go of the held note, and he thanks the player.
type Hang struct {
	Phrase  string   // what is asked, said
	Want    []Target // a target for each note, the last one left hanging
	Humpf   string   // said while it hangs
	Resolve Note     // what resolves it
	Thanks  string   // said once resolved

	at   int // the note expected next
	held int // the key left hanging, 0 before
}

func (st *Hang) said() string { return st.Phrase }

func (st *Hang) alsoSays() []string { return []string{st.Humpf, st.Thanks} }

func (st *Hang) begin(s Stage) {
	if st.held != 0 { // played over, R: nothing left hanging
		s.Release()
		s.Blink(0)
	}
	st.at, st.held = 0, 0
	s.Say(st.Phrase)
	s.Light(nil)
}

func (st *Hang) pressed(s Stage, held []int) bool {
	if st.held == 0 {
		r := Repeat{Want: st.Want}
		if st.Want[st.at].Judge(held) != Hit {
			s.Miss(r.hint(), st.at)
			st.at = 0
			return false
		}
		if st.at++; st.at < len(st.Want) {
			return false
		}
		s.Light(nil) // the phrase's keys, lit after a miss
		st.held = held[len(held)-1]
		s.Hold(st.held)
		s.Say(st.Humpf)
		up := (int(st.Resolve) - pc(st.held) + 12) % 12
		if up == 0 {
			up = 12
		}
		s.Blink(st.held + up)
		return false
	}
	if st.Resolve.Judge(held) != Hit {
		s.Miss(nil, 0)
		return false
	}
	s.Blink(0)
	s.Release()
	s.Say(st.Thanks)
	return true
}

func (*Hang) released(Stage, []int) bool { return false }
func (*Hang) phraseEnded(Stage) bool     { return false }
func (*Hang) read() bool                 { return false }

// Guess asks for a note the player has not been shown, to work out
// from what he knows: "Et mi♯ ?" A miss sounds the casserole and says
// `Nope`, showing nothing; at the `Chances`th, the walker gives
// `Answer` and shows the key, to be played as an Ask's answer is. The
// chances are not counted aloud.
type Guess struct {
	Phrase  string // the question
	Want    Note
	Nope    string // said at a miss, while chances are left
	Answer  string // said once they are spent
	Chances int

	misses int
}

func (st *Guess) said() string { return st.Phrase }

func (st *Guess) alsoSays() []string { return []string{st.Nope, st.Answer} }

func (st *Guess) begin(s Stage) {
	st.misses = 0
	s.Say(st.Phrase)
	s.Light(nil)
}

func (st *Guess) pressed(s Stage, held []int) bool {
	if st.Want.Judge(held) == Hit {
		s.Light(nil)
		s.Hit()
		return true
	}
	st.misses++
	switch {
	case st.misses < st.Chances:
		s.Miss(nil, 0)
		s.Say(st.Nope)
	case st.misses == st.Chances:
		s.Miss(st.Want.Hint(), 0)
		s.Say(st.Answer)
	default:
		s.Miss(st.Want.Hint(), 0)
	}
	return false
}

func (*Guess) released(Stage, []int) bool { return false }
func (*Guess) phraseEnded(Stage) bool     { return false }
func (*Guess) read() bool                 { return false }

// Write is "Écrire": the walker writes chord symbols, `Chords` in
// ChordPro, side by side, or as a line of a grid, a chord a bar, when
// `AsLine`, and says what they are. Over once its bubble is read; what
// it wrote stays until another step writes.
type Write struct {
	still
	Phrase string
	Chords []string
	AsLine bool
}

func (st Write) said() string { return st.Phrase }

func (st Write) begin(s Stage) {
	s.Say(st.Phrase)
	if st.AsLine {
		s.Chords(nil)
		s.Line(st.Chords, -1)
		return
	}
	s.Line(nil, -1)
	s.Chords(st.Chords)
}

func (Write) read() bool { return true }

// PlayLine plays a line of a grid on the bass, the roots of its chords,
// a note a bar. With Keys, the walker plays it, the bar he plays shaded,
// over when his phrase ends. Without, the player plays `Want`, a target
// a bar, the line waiting for him, without tempo, as the first palier
// does: a wrong root sounds the casserole and shows the right one, and
// the line stays on its bar.
type PlayLine struct {
	Phrase string
	Chords []string
	Keys   []int    // the walker's roots, MIDI numbers; nil for the player's turn
	Want   []Target // the player's, a target a bar

	at int // the bar played next
}

func (st *PlayLine) said() string { return st.Phrase }

func (st *PlayLine) begin(s Stage) {
	s.Chords(nil)
	s.Bass(true)
	s.Say(st.Phrase)
	s.Light(nil)
	if st.Keys != nil {
		s.Line(st.Chords, -1)
		s.PlayBars(st.Keys)
		return
	}
	st.at = 0
	s.Line(st.Chords, 0)
}

func (st *PlayLine) pressed(s Stage, held []int) bool {
	if st.Keys != nil {
		return false // he plays: the keys sound, nothing counts
	}
	if st.Want[st.at].Judge(held) != Hit {
		s.Miss(st.Want[st.at].Hint(), 0) // a bar of its own: shown on the next one
		return false
	}
	s.Light(nil)
	if st.at++; st.at < len(st.Want) {
		s.Line(st.Chords, st.at)
		return false
	}
	s.Line(st.Chords, -1)
	s.Hit()
	return true
}

func (st *PlayLine) phraseEnded(Stage) bool  { return st.Keys != nil }
func (*PlayLine) released(Stage, []int) bool { return false }
func (*PlayLine) read() bool                 { return false }
