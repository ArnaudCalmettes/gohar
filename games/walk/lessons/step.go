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
	read() bool                        // the bubble is read: a key, or its reading time; no side effect
	beat(s Stage, n int) bool          // beat `n` of the pulse begins, counted from its start
}

// still is the step that waits for nothing it does not name.
type still struct{}

func (still) pressed(Stage, []int) bool  { return false }
func (still) released(Stage, []int) bool { return false }
func (still) phraseEnded(Stage) bool     { return false }
func (still) read() bool                 { return false }
func (still) beat(Stage, int) bool       { return false }

// The steps out of rhythm, whose own methods leave no room for still,
// let the pulse go by.
func (*Ask) beat(Stage, int) bool      { return false }
func (*Guess) beat(Stage, int) bool    { return false }
func (*Repeat) beat(Stage, int) bool   { return false }
func (*Hang) beat(Stage, int) bool     { return false }
func (*PlayLine) beat(Stage, int) bool { return false }

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

// Band is "L'orchestre": the band keeps the pulse at `BPM`, or stops
// at 0, and the count of the bar shows when `Count`, as the walker says
// what they are. Over once its bubble is read; the pulse and the count
// stay as they are for the steps after it.
type Band struct {
	still
	Phrase string
	BPM    float64
	Count  bool
}

func (st Band) said() string { return st.Phrase }

func (st Band) begin(s Stage) {
	s.Say(st.Phrase)
	s.Pulse(st.BPM)
	s.Count(st.Count)
}

func (Band) read() bool { return true }

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
