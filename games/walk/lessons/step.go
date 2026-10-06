package lessons

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
	read() bool                        // the bubble is read: a key, or its reading time
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
	missesHey      = 4
	missesOnPurpose = 7
	missesFed      = 10
)

func (st *Ask) said() string { return st.Phrase }

func (st *Ask) begin(s Stage) {
	st.partial, st.misses = false, 0
	s.Say(st.Phrase)
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
	s.Miss(st.Want.Hint())
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
type Repeat struct {
	Phrase string
	Keys   []int // what the walker plays, MIDI numbers
	Want   []Target // a target for each note

	at int // the note expected next, -1 while the walker plays
}

func (st *Repeat) said() string { return st.Phrase }

func (st *Repeat) begin(s Stage) {
	st.at = -1
	s.Say(st.Phrase)
	s.Light(nil)
	s.Play(st.Keys)
}

func (st *Repeat) phraseEnded(Stage) bool {
	st.at = 0
	return false
}

func (st *Repeat) pressed(s Stage, held []int) bool {
	if st.at < 0 {
		return false // he is still playing: the keys sound, nothing counts
	}
	if st.Want[st.at].Judge(held) != Hit {
		s.Miss(nil)
		st.at = -1
		s.Play(st.Keys)
		return false
	}
	s.Light(nil)
	st.at++
	if st.at < len(st.Want) {
		return false
	}
	s.Hit()
	return true
}

func (*Repeat) released(Stage, []int) bool { return false }
func (*Repeat) read() bool                  { return false }
