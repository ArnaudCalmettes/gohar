package ireal

import (
	"regexp"
	"strconv"
	"strings"
)

// A Played is one measure as it is played, once the form is unfolded.
type Played struct {
	Index int // of the written measure, in Chart.Measures

	// Events are what sounds in the bar. For a measure that repeats one
	// or two bars before ("x", "r"), they are those bars' events; for
	// an empty measure, none: the chord before goes on.
	Events []Event

	// From is the written measure the events come from: Index, unless
	// the bar repeats another. Their cells count in that measure, which
	// may not be as wide as this one: "Kcl" is two cells, the bar it
	// repeats usually four.
	From int

	// Coda marks the first bar of the coda, reached by the jump of a
	// "D.S." or "D.C. al Coda". The coda concludes the tune: the app
	// plays it on the last chorus only.
	Coda bool
}

// Unfold returns the measures in the order they are played: repeats
// taken, endings chosen pass by pass, D.C. and D.S. followed to the
// coda or to the fine.
//
// # The conventions
//
//   - A repeat is played twice, or as many times as a comment inside it
//     says ("3x"), or once per ending when it has some.
//   - Endings are taken in the order they are written, whatever their
//     number: charts typed by hand hold "N1" twice, or "N3" right after
//     "N1", and the order is what the typist meant.
//   - A repeat closed but never opened goes back to the start of the
//     chart, or to just after the repeat before it. A repeat opened but
//     never closed is not a repeat.
//   - "D.C." goes back to the start, "D.S." to the segno, wherever the
//     comment writes it ("After solos, D.C. al Coda"). On the way back,
//     repeats are not taken again and the last ending is played, as is
//     the custom on a lead sheet. "al Coda" jumps from the first coda
//     sign to the last one; "al Fine" stops after the measure marked
//     "Fine". "al 2nd ending" takes that ending instead of the last,
//     wherever it is written, even after the direction (a 3rd ending
//     often is), and stops at a "Fine" when there is one, as a bare
//     "D.C." does.
//
// A chart whose directions cannot be followed is played as written:
// Unfold never fails, and never plays more than a bounded number of
// bars.
func (c Chart) Unfold() []Played {
	f := newForm(c.Measures)
	return f.play()
}

// A repeat is the span of a repeat sign, from s to e inclusive.
type repeat struct {
	s, e    int
	times   int
	endings []int // where each ending starts, in written order
}

type form struct {
	ms      []Measure
	closing map[int]*repeat // by the index of the measure closing it
	ending  map[int]*repeat // by the index of a measure opening an ending
}

var timesComment = regexp.MustCompile(`^\s*(\d+)\s*[xX]\s*$`)

func newForm(ms []Measure) form {
	f := form{ms: ms, closing: map[int]*repeat{}, ending: map[int]*repeat{}}
	var opens []int
	after := 0 // where a repeat closed but never opened goes back to
	for i, m := range ms {
		if m.Open == RepeatOpen {
			opens = append(opens, i)
		}
		if m.Close != RepeatClose {
			continue
		}
		r := &repeat{s: after, e: i, times: 2}
		if n := len(opens); n > 0 {
			r.s, opens = opens[n-1], opens[:n-1]
		}
		after = i + 1

		for j := r.s; j <= r.e; j++ {
			if ms[j].Ending > 0 {
				r.endings = append(r.endings, j)
			}
			for _, mark := range ms[j].Marks {
				if g := timesComment.FindStringSubmatch(mark.Comment.Text); mark.Kind == Comment && g != nil {
					r.times, _ = strconv.Atoi(g[1])
				}
			}
		}
		// The last ending is written after the repeat sign.
		if len(r.endings) > 0 && i+1 < len(ms) && ms[i+1].Ending > 0 {
			r.endings = append(r.endings, i+1)
		}
		if len(r.endings) > 0 {
			r.times = len(r.endings)
		}
		for _, j := range r.endings {
			f.ending[j] = r
		}
		f.closing[i] = r
	}
	return f
}

// direction is a D.C. or D.S. found in a comment.
type direction struct {
	at     int  // the measure it is written in
	segno  bool // D.S. rather than D.C.
	coda   bool // al Coda: jump to the coda, else stop at a Fine
	ending int  // al Nth ending: the ending to take on the way back
}

var (
	daCapo   = regexp.MustCompile(`\bd\.([cs])\.(.*)$`)
	nthEnd   = regexp.MustCompile(`^\s*al\s+(\d)(st|nd|rd|th)\s+end`)
	toFine   = regexp.MustCompile(`^\s*al\s+fine`)
	toCodaAt = regexp.MustCompile(`^\s*al\s+coda`)
)

// direction finds the D.C. or D.S. of a chart and what it goes to: the
// coda, the fine, an ending, or the end of the chart when it says
// nothing more. One it cannot read ("D.C. on cue") is not followed.
func (f form) direction() (direction, bool) {
	for i, m := range f.ms {
		for _, mark := range m.Marks {
			if mark.Kind != Comment {
				continue
			}
			g := daCapo.FindStringSubmatch(strings.ToLower(mark.Comment.Text))
			if g == nil {
				continue
			}
			d := direction{at: i, segno: g[1] == "s"}
			rest := g[2]
			switch {
			case toCodaAt.MatchString(rest):
				d.coda = true
			case toFine.MatchString(rest):
			case nthEnd.MatchString(rest):
				d.ending, _ = strconv.Atoi(nthEnd.FindStringSubmatch(rest)[1])
			case strings.TrimSpace(rest) != "":
				continue
			}
			return d, true
		}
	}
	return direction{}, false
}

// endingAt returns where the nth ending of a repeat starts: among its
// endings, or else the first measure after the repeat marked so.
func (f form) endingAt(r *repeat, n int) int {
	for _, j := range r.endings {
		if f.ms[j].Ending == n {
			return j
		}
	}
	for j := r.e + 1; j < len(f.ms); j++ {
		if f.ms[j].Ending == n {
			return j
		}
	}
	return -1
}

func (f form) marked(k Kind, i int) bool {
	for _, mark := range f.ms[i].Marks {
		if mark.Kind == k {
			return true
		}
	}
	return false
}

func (f form) fine(i int) bool {
	for _, mark := range f.ms[i].Marks {
		if mark.Kind == Comment && strings.EqualFold(strings.TrimSpace(mark.Comment.Text), "fine") {
			return true
		}
	}
	return false
}

func (f form) play() []Played {
	ms := f.ms
	var out []Played
	dir, hasDir := f.direction()
	back := false // on the way back after the D.C. or D.S.
	toCoda, coda := -1, -1
	for i := range ms {
		if f.marked(Coda, i) {
			if toCoda < 0 {
				toCoda = i
			}
			coda = i
		}
	}
	segno := 0
	for i := range ms {
		if f.marked(Segno, i) {
			segno = i
			break
		}
	}

	pass := map[*repeat]int{}
	chosen := false  // the ending a D.C. al Nth ending asks for, taken
	twoBack := false // the bar after an "r" repeats the one two back too
	jumped := false  // to the coda: the next bar played opens it
	limit := 16*len(ms) + 64

	for i := 0; i < len(ms) && len(out) < limit; {
		if r, ok := f.ending[i]; ok && back && dir.ending > 0 {
			if !chosen {
				chosen = true
				if j := f.endingAt(r, dir.ending); j >= 0 && j != i {
					i = j
					continue
				}
			}
		} else if ok {
			k := indexOf(r.endings, i)
			want := pass[r]
			if back {
				want = len(r.endings) - 1
			}
			switch {
			case k > want && !back:
				// Our ending is over: back to the start of the repeat.
				pass[r]++
				i = r.s
				continue
			case k < want:
				i = r.endings[want]
				continue
			}
		}

		out = append(out, f.bar(i, out, &twoBack))
		if jumped {
			out[len(out)-1].Coda, jumped = true, false
		}

		switch {
		case hasDir && !back && i == dir.at:
			back = true
			i = 0
			if dir.segno {
				i = segno
			}
			continue
		case back && !dir.coda && f.fine(i):
			return out
		case back && dir.coda && i == toCoda && coda > toCoda:
			i, jumped = coda, true
			continue
		}

		if r, ok := f.closing[i]; ok && !back && pass[r]+1 < r.times {
			pass[r]++
			i = r.s
			continue
		}
		i++
	}
	return out
}

// bar resolves what sounds in measure i given the bars played before.
func (f form) bar(i int, out []Played, twoBack *bool) Played {
	m := f.ms[i]
	b := Played{Index: i, Events: m.Events, From: i}
	switch {
	case m.Repeat == 1 && len(out) >= 1:
		b.Events, b.From = out[len(out)-1].Events, out[len(out)-1].From
	case m.Repeat == 2 && len(out) >= 2:
		b.Events, b.From = out[len(out)-2].Events, out[len(out)-2].From
		*twoBack = true
		return b
	case *twoBack && len(m.Events) == 0 && len(out) >= 2:
		b.Events, b.From = out[len(out)-2].Events, out[len(out)-2].From
	}
	*twoBack = false
	return b
}

func indexOf(s []int, v int) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}
