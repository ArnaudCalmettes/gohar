package ireal

// A Chart is a chart read into measures, in the order they are written.
//
// Repeats, endings, segno and coda are marked where they are written
// and not unfolded: the order the song is played in is another step,
// which reads a Chart.
type Chart struct {
	Measures []Measure
}

// A Measure is what lies between two bar lines.
//
// # Cells, not beats
//
// The app lays a chart out on a grid of cells, and a measure is most
// often four cells wide whatever its time signature: on real playlists,
// four cells hold nine measures out of ten in 4/4 as in 3/4. Cells is
// that width, and each event keeps the cell it sits in. Turning cells
// into beats is left to the reader who needs it: a chord in cell 2 of
// 4 in a 4/4 measure is on beat 3, but a 3/4 measure of four cells
// says nothing that sure.
type Measure struct {
	Open  Kind // Bar, DoubleOpen or RepeatOpen
	Close Kind // Bar, DoubleClose, RepeatClose or FinalBar

	Section byte    // the rehearsal mark starting here, 0 if none
	Time    TimeSig // in force in this measure
	NewTime bool    // the time signature changes here
	Ending  int     // the numbered ending starting here, 0 if none

	// Repeat is 1 when the measure repeats the one before ("x"), 2
	// when it starts a repeat of the two before ("r").
	Repeat int

	Cells  int
	Events []Event
	Marks  []Mark
}

// An Event is what sounds from a cell on: a chord, no chord, or a slash
// that plays the chord before again.
type Event struct {
	Kind  Kind // Chord, NoChord or Slash
	Cell  int
	Chord ChordSymbol

	// Alternate is the chord written small above this one, if any.
	Alternate    ChordSymbol
	HasAlternate bool

	Small bool // written small, usually a passing chord
}

// A Mark is a sign written in a measure at a given cell: a segno, a
// coda, a fermata, where the player stops, or a comment. Its cell tells
// a coda that opens a measure from one that closes it.
type Mark struct {
	Kind    Kind // Segno, Coda, Fermata, PlayerEnd or Comment
	Cell    int
	Comment CommentText
}

// defaultTime is what a chart that never says is in.
var defaultTime = TimeSig{4, 4}

// Structure reads tokens into measures.
//
// Bar lines end a measure only once it has something in it: a run of
// bar lines between two measures, such as "]|" or "}[", closes the one
// before with the strongest of them and opens the next with the last
// opening one. "LZ" and "Kcl" carry cells of their own, as the app
// writes them: "LZ" is a space then a bar line, "Kcl" a bar line then a
// space and an "x".
//
// # Gaps are not measures
//
// Charts use empty cells to lay things out: after a repeat, to put the
// second ending under the first; after the final bar line, trailing
// spaces. Empty cells that follow a double bar line, a repeat or the
// final bar line, with nothing written in them, are such a gap and are
// dropped. Empty cells after a plain bar line are kept: that is a
// measure where the chord before goes on.
func Structure(tokens []Token) Chart {
	r := reader{time: defaultTime}
	r.cur = r.start(Bar)
	for _, t := range tokens {
		r.read(t)
	}
	if r.cur.filled() && !r.gap() {
		r.cur.Close = Bar
		r.out = append(r.out, r.cur)
	}
	return Chart{Measures: r.out}
}

type reader struct {
	out   []Measure
	cur   Measure
	time  TimeSig
	small bool
}

func (m *Measure) filled() bool { return m.Cells > 0 || m.Repeat > 0 }

// gap reports whether the current measure is only layout: nothing in
// it but comments, after a bar line that ends a passage. The comments
// of a gap, typically a credit after the final bar line, go to the
// measure before it.
func (r *reader) gap() bool {
	m := &r.cur
	if len(m.Events) > 0 || m.Repeat > 0 || m.Section != 0 || m.Ending != 0 ||
		m.NewTime || len(r.out) == 0 || r.out[len(r.out)-1].Close == Bar {
		return false
	}
	for _, mark := range m.Marks {
		if mark.Kind != Comment {
			return false
		}
	}
	last := &r.out[len(r.out)-1]
	for _, mark := range m.Marks {
		mark.Cell = last.Cells
		last.Marks = append(last.Marks, mark)
	}
	m.Marks = nil
	return true
}

func (r *reader) start(open Kind) Measure {
	return Measure{Open: open, Time: r.time}
}

func (r *reader) read(t Token) {
	m := &r.cur
	switch t.Kind {
	case Bar:
		if t.Raw == "LZ" {
			m.Cells++
		}
		r.bar(Bar)
	case DoubleOpen, DoubleClose, RepeatOpen, RepeatClose, FinalBar:
		r.bar(t.Kind)
	case BarRepeated:
		r.bar(Bar)
		r.cur.Cells += 2
		r.cur.Repeat = 1

	case Space:
		m.Cells++
	case EmptyCells:
		m.Cells += 3
	case Chord, NoChord, Slash:
		m.Events = append(m.Events, Event{Kind: t.Kind, Cell: m.Cells, Chord: t.Chord, Small: r.small})
		m.Cells++
	case Alternate:
		if n := len(m.Events); n > 0 {
			m.Events[n-1].Alternate, m.Events[n-1].HasAlternate = t.Chord, true
		}
	case RepeatOne, RepeatTwo:
		m.Repeat = 1
		if t.Kind == RepeatTwo {
			m.Repeat = 2
		}
		m.Cells++

	case Section:
		m.Section = t.Section
	case TimeSignature:
		r.time = t.Time
		m.Time, m.NewTime = t.Time, true
	case Ending:
		m.Ending = t.Ending
	case Segno, Coda, Fermata, PlayerEnd, Comment:
		m.Marks = append(m.Marks, Mark{Kind: t.Kind, Cell: m.Cells, Comment: t.Comment})

	case Small:
		r.small = true
	case Large:
		r.small = false
	}
	// Comma, Spacer and Unknown place nothing.
}

// bar handles one bar line: it ends the current measure if there is
// one, and otherwise strengthens the bar lines around the gap.
func (r *reader) bar(k Kind) {
	closing := k == DoubleClose || k == RepeatClose || k == FinalBar
	opening := k == DoubleOpen || k == RepeatOpen

	if r.cur.filled() {
		next := r.start(Bar)
		if opening {
			next.Open = k
		}
		if r.gap() {
			// The gap's own bar lines belong to the measures around
			// it: an opening one to the next.
			r.cur = next
			return
		}
		r.cur.Close = Bar
		if closing {
			r.cur.Close = k
		}
		r.out = append(r.out, r.cur)
		r.cur = next
		return
	}

	switch {
	case opening:
		r.cur.Open = k
	case closing && len(r.out) > 0:
		if last := &r.out[len(r.out)-1]; last.Close == Bar {
			last.Close = k
		}
	}
}
