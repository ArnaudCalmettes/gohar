package ireal

// Ticks measure time in a chart: TicksPerBeat to a beat.
type Ticks int

// TicksPerBeat leaves room below the beat. Chords read from a chart
// always start on a beat (see [Chart.Timeline]), but what is timed
// later, a chord the player sounds, will not.
const TicksPerBeat Ticks = 2520

// A Span is one chord and how long it sounds, in the order played.
type Span struct {
	Chord   ChordSymbol
	NoChord bool // "n", or the silence before the first chord

	Bar    int   // the played measure it starts in, from 0
	Start  Ticks // from the start of the chorus
	Length Ticks
}

// A Timeline is a chorus as it sounds: chords with their durations,
// and where each bar starts.
//
// # A cycle
//
// The app plays a chart chorus after chorus, so the last span is
// followed by the first, and Next says so. It is what lets an analysis
// read a turnaround as the preparation of the first chord: the last
// bars look ahead past the final bar line.
//
// A coda is the conclusion, played on the last chorus only: the chorus
// loops back from the bar before it, and after the coda the tune is
// over.
type Timeline struct {
	Spans  []Span
	Bars   []Ticks // where each played measure starts
	Length Ticks

	// Coda is the first span of the coda, or 0 when there is none. A
	// chord held into the coda is split there, so that the coda starts
	// a span of its own.
	Coda int
}

// Next returns the index of the span that follows span i when the
// chorus goes round: the first after the last, or after the last one
// before the coda. It returns -1 after the coda: the tune has ended.
func (t Timeline) Next(i int) int {
	switch {
	case t.Coda > 0 && i == t.Coda-1:
		return 0
	case i == len(t.Spans)-1 && t.Coda > 0:
		return -1
	}
	return (i + 1) % len(t.Spans)
}

// Timeline unfolds the chart and times its chords.
//
// # Cells to beats
//
// A chord sits in a cell, and a measure is most often four cells wide
// whatever its time signature. A chord starts on the beat where its
// cell falls when the cells share the measure evenly, rounded up to
// the next beat, and never past the last one:
//
//   - two cells in 4/4 are two half measures, four in 12/8 are the four
//     dotted quarters;
//   - in 3/4 over four cells, cell 2 falls at a beat and a half and the
//     chord starts on beat 3, which is how real charts write a bar of
//     two chords: cell 2 ten times more often than cell 1;
//   - in 5/4 over four cells, cell 2 falls at two beats and a half and
//     the chord starts on beat 4: Take Five, three then two, as the app
//     plays it.
//
// The beat is the pulse: a quarter in 4/4 and 3/4, a dotted quarter in
// 6/8 and 12/8.
//
// # What goes on
//
// An empty measure, a slash and a chord written again the same all let
// the chord before go on: the timeline holds one span for them. A "W"
// root takes the root of the chord before.
func (c Chart) Timeline() Timeline {
	var t Timeline
	var now Ticks
	for bar, p := range c.Unfold() {
		m := c.Measures[p.Index]
		cells := c.Measures[p.From].Cells
		t.Bars = append(t.Bars, now)
		if p.Coda && len(t.Spans) > 0 {
			last := &t.Spans[len(t.Spans)-1]
			last.Length = now - last.Start
			t.Spans = append(t.Spans, Span{Chord: last.Chord, NoChord: last.NoChord, Bar: bar, Start: now})
			t.Coda = len(t.Spans) - 1
		}
		for _, e := range p.Events {
			if e.Kind == Slash {
				continue
			}
			t.add(Span{
				Chord:   e.Chord,
				NoChord: e.Kind == NoChord,
				Bar:     bar,
				Start:   now + at(e.Cell, cells, m.Time),
			})
		}
		now += Ticks(pulses(m.Time)) * TicksPerBeat
	}
	t.Length = now
	if n := len(t.Spans); n > 0 {
		t.Spans[n-1].Length = now - t.Spans[n-1].Start
	}
	return t
}

// add appends a span, closing the one before it.
func (t *Timeline) add(s Span) {
	if len(t.Spans) == 0 && s.Start > 0 {
		// A chart that starts on nothing, a pickup or a rest.
		t.Spans = append(t.Spans, Span{NoChord: true})
	}
	n := len(t.Spans)
	if n == 0 {
		t.Spans = append(t.Spans, s)
		return
	}
	last := &t.Spans[n-1]
	if s.Chord.Root == "W" && !last.NoChord {
		s.Chord.Root = last.Chord.Root
	}
	switch {
	case s.NoChord == last.NoChord && s.Chord == last.Chord:
		return // written again, goes on
	case s.Start <= last.Start:
		// Two chords on one beat, which the rule of a cell per beat
		// can make: the later one is what sounds.
		s.Bar, s.Start = last.Bar, last.Start
		*last = s
		return
	}
	last.Length = s.Start - last.Start
	t.Spans = append(t.Spans, s)
}

// pulses is how many beats a measure counts, the beat being the pulse:
// compound meters (6/8, 9/8, 12/8) count in dotted quarters.
func pulses(t TimeSig) int {
	if t.Unit == 8 && t.Beats > 3 && t.Beats%3 == 0 {
		return t.Beats / 3
	}
	return t.Beats
}

// at is where a cell starts in its measure: see [Chart.Timeline].
func at(cell, cells int, t TimeSig) Ticks {
	p := pulses(t)
	beat := min((cell*p+cells-1)/cells, p-1)
	return Ticks(beat) * TicksPerBeat
}
