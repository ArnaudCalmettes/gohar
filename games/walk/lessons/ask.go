package lessons

// The steps that ask for notes, one at a time: Ask, and Guess, its
// riddle.

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
