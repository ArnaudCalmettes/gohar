package analysis

import (
	"slices"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// A Section is one stretch of the form of a tune: from which bar, over
// how many, and as which section, a letter given in the order the
// sections first sound, "A", "B", or "-" for a pickup of a bar or two
// before the first section. Shift is its transposition from the first
// time that section sounds: the second A of All The Things You Are is
// the first a fifth up.
type Section struct {
	Label string
	From  int // the first bar, an index in Changes.Bars
	Bars  int
	Shift harmony.Semitones
}

// SectionOptions tune how [SectionsWith] reads a form.
type SectionOptions struct {
	// MinRun is the shortest passage that counts as coming back, in
	// bars.
	MinRun int

	// Carrure prefers sections of 4 and 8 bars, the one piece of
	// knowledge of a style the form uses (see [Sections]).
	Carrure bool

	// Tolerance lets a recurrence of 6 bars or more hold one bar that
	// differs, and compares chords by family (see [chordFamily]).
	Tolerance bool

	// Loose compares the roots of the chords only.
	Loose bool
}

// DefaultSections are the options of [Sections].
var DefaultSections = SectionOptions{MinRun: 4, Carrure: true}

// Sections finds the sections of a tune from its chords alone, with
// the options [DefaultSections]. It reads no [A] [B] mark and knows no
// AABA nor blues: the sections are the passages that come back.
//
// # What the sources say
//
// The time of a tune is made of levels nested in one another, beats,
// bars, groups of bars, and above them the forms: sections, tunes
// (Siron, La partition intérieure, p. 147). The standards group their
// bars in « carrure », symmetric fragments of 4, 8 or 16 bars, most
// often phrases of 8 bars, the blues three phrases of 4 (p. 146 and
// 147). The song form AABA has sections of 8 bars and its variants,
// other lengths (Alone Together, 14 + 14 + 8 + 8) or other orders,
// ABAC, AABC (p. 408). The A's of an AABA often differ only by their
// cadence-boucle, the turnaround that leads to the next section (p. 353).
//
// # The method
//
// Each bar becomes what sounds in it: its chords, where they start,
// their root and their tetrad. Two passages are the same when bar after
// bar their chords are, at one transposition: a bridge that repeats a
// motif a fourth up is found too. Every passage of MinRun bars or more
// that comes back later is a recurrence.
//
// The form grows from its first bar: the longest recurrence that starts
// where a section is known to start is placed first, and each
// occurrence placed tells where more sections start. Only when none is
// left does a recurrence start anywhere, and never a transposed one:
// short chains of dominants look alike a fifth apart.
//
// When one occurrence is followed right away by the next, as A1 by A2,
// the chart tells how long the section is, and every occurrence covers
// that length and no more: the A's that differ by their last bars are
// both A, and the bridge after an A' whose last bars prepare it is
// another section (Billy Boy). A stretch that comes back nowhere is a
// section of its own.
//
// # The carrure
//
// The carrure only breaks ties, so that the form stays one that any
// style could have. A recurrence within 2 bars of a multiple of 8 is
// given that length; one shorter than 8 bars may not start off the grid
// of 8 bars, and one that starts on it covers 8 bars when it has room,
// its last bars a turnaround that differs (the A's that match on 5
// bars only in There Is No Greater Love). A stretch of 16 bars or more,
// a multiple of 8, that comes back nowhere is cut in sections of 8: the
// B and the C of Autumn Leaves, and a tune of 16 bars in one breath,
// Blue Bossa, whose charts mark two sections of 8. The grid starts
// where most recurrences do, past a pickup or an introduction. A length
// the chart tells is never rounded.
//
// On 1350 standards, it reads an AABA in about a third of them, and an
// ABAC in one in seven; a blues, whose three phrases of 4 bars do not
// come back, is one section.
func Sections(c Changes) []Section {
	return SectionsWith(c, DefaultSections)
}

// SectionsWith finds the sections of a tune with the options `o` (see
// [Sections]). It returns nil when the changes have no bars.
func SectionsWith(c Changes, o SectionOptions) []Section {
	f := former{o: o, bars: formBarsOf(c)}
	if len(f.bars) == 0 {
		return nil
	}
	return f.form()
}

// A formSlot is one chord sounding in a bar: where it starts in the
// bar, or 0 for the chord held from the bar before, and what it is.
type formSlot struct {
	at      Ticks
	silent  bool
	root    harmony.PitchClass
	quality harmony.ChordPattern
}

// A formBar is what sounds in one bar, in order.
type formBar []formSlot

// formBarsOf cuts the changes into bars.
func formBarsOf(c Changes) []formBar {
	if len(c.Bars) == 0 || len(c.Chords) == 0 {
		return nil
	}
	last := c.Chords[len(c.Chords)-1]
	end := last.Start + last.Length
	out := make([]formBar, len(c.Bars))
	for b, from := range c.Bars {
		to := end
		if b+1 < len(c.Bars) {
			to = c.Bars[b+1]
		}
		for _, ch := range c.Chords {
			if ch.Start+ch.Length <= from || ch.Start >= to {
				continue
			}
			s := formSlot{at: max(ch.Start, from) - from, silent: ch.Silent}
			if !ch.Silent {
				s.root, s.quality = ch.Chord.Root, ch.Chord.Pattern.Tetrad()
			}
			out[b] = append(out[b], s)
		}
	}
	return out
}

// A recurrence is a passage of `length` bars starting at bar `a` that
// sounds again from bar `b`, transposed by `shift` semitones.
type recurrence struct {
	a, b, length int
	shift        harmony.Semitones
}

// An occurrence is one place a section sounds: from which bar, over
// how many bars the recurrence matched, as which section and at which
// transposition from the section's first occurrence.
type occurrence struct {
	start, length int
	label         int
	shift         harmony.Semitones
}

// former reads one form.
type former struct {
	grid   int // the bar the grid of 8 bars starts from (see phase)
	o      SectionOptions
	bars   []formBar
	placed []occurrence // sorted by start
}

// form places the recurrences one by one, growing the form from its
// first bar: first the best recurrence that starts on a boundary
// already known, else the best untransposed one anywhere. The best is
// the longest, then the untransposed, then the earliest.
func (f *former) form() []Section {
	recs := f.recurrences()
	f.grid = phase(recs)
	grid := f.grid
	done := make([]bool, len(recs))
	labels := 0
	for {
		bounds := f.boundaries()
		i := f.pick(recs, done, bounds, grid, true)
		if i < 0 {
			i = f.pick(recs, done, bounds, grid, false)
		}
		if i < 0 {
			break
		}
		done[i] = true
		labels = f.place(recs[i], labels)
	}
	return f.sections()
}

// same reports whether bar `b` sounds as bar `a` transposed by `t`
// semitones: the same roots, and the same tetrads, or the same family
// of them with Tolerance, or anything with Loose.
func (f *former) same(a, b formBar, t harmony.Semitones) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.at != y.at || x.silent != y.silent {
			return false
		}
		if x.silent {
			continue
		}
		if x.root.Transpose(t) != y.root {
			return false
		}
		switch {
		case f.o.Loose:
		case f.o.Tolerance && chordFamily(x.quality) != chordFamily(y.quality),
			!f.o.Tolerance && x.quality != y.quality:
			return false
		}
	}
	return true
}

// recurrences finds every passage of MinRun bars or more that comes
// back, at every distance and transposition, the best first. A passage
// that repeats over and over (A A A) is cut into as many recurrences as
// the distance allows, each no longer than the distance: two
// occurrences never overlap. With Tolerance, two runs of bars that
// match, parted by one bar that does not, make one recurrence when it
// is 6 bars long or more.
func (f *former) recurrences() []recurrence {
	var out []recurrence
	n := len(f.bars)
	for lag := f.o.MinRun; lag < n; lag++ {
		for t := harmony.Semitones(0); t < 12; t++ {
			var runs [][2]int // [start, end) of the bars that match
			run := 0
			for i := 0; i+lag <= n; i++ {
				if i+lag < n && f.same(f.bars[i], f.bars[i+lag], t) {
					run++
					continue
				}
				if run > 0 {
					runs = append(runs, [2]int{i - run, i})
				}
				run = 0
			}
			spans := runs
			if f.o.Tolerance {
				for k := 0; k+1 < len(runs); k++ {
					if from, to := runs[k][0], runs[k+1][1]; runs[k+1][0] == runs[k][1]+1 && to-from >= 6 {
						spans = append(spans, [2]int{from, to})
					}
				}
			}
			for _, sp := range spans {
				for start := sp[0]; sp[1]-sp[0] >= f.o.MinRun && start+f.o.MinRun <= sp[1]; start += lag {
					out = append(out, recurrence{a: start, b: start + lag, length: min(lag, sp[1]-start), shift: t})
				}
			}
		}
	}
	slices.SortStableFunc(out, func(x, y recurrence) int {
		switch {
		case x.length != y.length:
			return y.length - x.length
		case (x.shift == 0) != (y.shift == 0):
			if x.shift == 0 {
				return -1
			}
			return 1
		}
		return x.a - y.a
	})
	return out
}

// phase returns the bar the grid of 8 bars starts from, between 0 and
// 7: the one most recurrences start on, weighed by their length. A tune
// with a pickup bar, or an introduction of 4, has its grid shifted as
// much.
func phase(recs []recurrence) int {
	var weight [8]int
	for _, r := range recs {
		if r.a%8 == r.b%8 {
			weight[r.a%8] += r.length
		}
	}
	best := 0
	for p, w := range weight {
		if w > weight[best] {
			best = p
		}
	}
	return best
}

// pick returns the first recurrence not `done` that can be placed: one
// that starts on a boundary when `anchored`, an untransposed one
// otherwise. A recurrence that no longer fits is marked done: what is
// placed only grows, so it never will. So is, with Carrure, one shorter
// than 8 bars that starts off the grid of 8 bars from bar `grid`.
func (f *former) pick(recs []recurrence, done []bool, bounds map[int]bool, grid int, anchored bool) int {
	for i, r := range recs {
		if done[i] {
			continue
		}
		if f.fit(r) < f.o.MinRun || f.o.Carrure && r.length < 8 && ((r.a-grid)%8 != 0 || (r.b-grid)%8 != 0) {
			done[i] = true
			continue
		}
		if anchored && (bounds[r.a] || bounds[r.b]) || !anchored && r.shift == 0 {
			return i
		}
	}
	return -1
}

// startsAt returns the occurrence placed at bar `start`, or -1.
func (f *former) startsAt(start int) int {
	return slices.IndexFunc(f.placed, func(o occurrence) bool { return o.start == start })
}

// room returns how many bars from `start` are free, up to `length`: 0
// when `start` falls inside an occurrence already placed.
func (f *former) room(start, length int) int {
	for _, o := range f.placed {
		switch {
		case o.start == start:
		case o.start < start && start < o.start+o.length:
			return 0
		case start < o.start && o.start < start+length:
			length = o.start - start
		}
	}
	return length
}

// fit returns how many bars of `r` can be placed: both occurrences
// must fit, a new one in free bars, or on the start of one placed.
func (f *former) fit(r recurrence) int {
	ia, ib := f.startsAt(r.a), f.startsAt(r.b)
	if ia >= 0 && ib >= 0 && f.placed[ia].label == f.placed[ib].label {
		return 0 // nothing new
	}
	length := r.length
	if ia < 0 {
		length = min(length, f.room(r.a, length))
	}
	if ib < 0 {
		length = min(length, f.room(r.b, length))
	}
	return length
}

// place adds the occurrences of `r`, under a new label or the label of
// the occurrence already there, and returns the number of labels used.
func (f *former) place(r recurrence, labels int) int {
	length := f.fit(r)
	ia, ib := f.startsAt(r.a), f.startsAt(r.b)
	switch {
	case ia < 0 && ib < 0:
		f.placed = append(f.placed,
			occurrence{start: r.a, length: length, label: labels},
			occurrence{start: r.b, length: length, label: labels, shift: r.shift})
		labels++
	case ib < 0:
		a := f.placed[ia]
		f.placed = append(f.placed, occurrence{start: r.b, length: length, label: a.label, shift: mod12(a.shift + r.shift)})
	case ia < 0:
		b := f.placed[ib]
		f.placed = append(f.placed, occurrence{start: r.a, length: length, label: b.label, shift: mod12(b.shift - r.shift)})
	default:
		// Two sections found apart are one: relabel the second, its
		// transpositions measured from the first.
		from, delta := f.placed[ib].label, mod12(f.placed[ia].shift+r.shift-f.placed[ib].shift)
		for i := range f.placed {
			if f.placed[i].label == from {
				f.placed[i].label, f.placed[i].shift = f.placed[ia].label, mod12(f.placed[i].shift+delta)
			}
		}
	}
	slices.SortFunc(f.placed, func(x, y occurrence) int { return x.start - y.start })
	return labels
}

// span returns how many bars there are from occurrence `i` to the next
// one, or to the end.
func (f *former) span(i int) int {
	if i+1 < len(f.placed) {
		return f.placed[i+1].start - f.placed[i].start
	}
	return len(f.bars) - f.placed[i].start
}

// lengths returns the length of each section whose length the chart
// tells: an occurrence followed by the next one within half its matched
// length again, as A1 by A2, whose last bars differ, spans its whole
// section. The span found most often wins, the shortest at equal
// counts.
func (f *former) lengths() map[int]int {
	// The span of the last occurrence runs to the end of the chart,
	// which may cut its section short, and tells nothing: the last A of
	// a tune with a pickup of 2 bars is 6 bars long, the others 8.
	inner := map[int]map[int]int{}
	for i, o := range f.placed[:max(0, len(f.placed)-1)] {
		s := f.span(i)
		if s*2 > o.length*3 {
			continue
		}
		if inner[o.label] == nil {
			inner[o.label] = map[int]int{}
		}
		inner[o.label][s]++
	}
	out := map[int]int{}
	for label, c := range inner {
		best := 0
		for s, k := range c {
			if best == 0 || k > c[best] || k == c[best] && s < best {
				best = s
			}
		}
		out[label] = best
	}
	return out
}

// extent returns how many bars occurrence `i` covers: the length of its
// section when the chart tells it and it fits before the next
// occurrence, what matched otherwise. With Carrure, a matched length
// within 2 bars of a multiple of 8 becomes that multiple, and one
// shorter than 8 bars that starts on the grid becomes 8, when it fits.
func (f *former) extent(i int, known map[int]int) int {
	o := f.placed[i]
	s := f.span(i)
	if k, ok := known[o.label]; ok && k >= o.length && k <= s {
		return k
	}
	l := o.length
	if f.o.Carrure {
		switch c := (l + 4) / 8 * 8; {
		case c > 0 && c-l <= 2 && l-c <= 2 && c <= s:
			l = c
		case l < 8 && (o.start-f.grid)%8 == 0 && s >= 8:
			l = 8
		}
	}
	return l
}

// boundaries returns the bars where a section is known to start: the
// first, the start of every occurrence, and its end when the length of
// its section is known. With Carrure, the first bar of the grid too:
// after a pickup, the tune starts there, and a longer recurrence that
// starts elsewhere, a turnaround and the next A alike to the end of the
// bridge and the last A, must not come first.
func (f *former) boundaries() map[int]bool {
	out := map[int]bool{0: true}
	if f.o.Carrure {
		out[f.grid] = true
	}
	known := f.lengths()
	for i, o := range f.placed {
		out[o.start] = true
		if _, ok := known[o.label]; ok {
			if end := o.start + f.extent(i, known); end < len(f.bars) {
				out[end] = true
			}
		}
	}
	return out
}

// sections turns the occurrences into sections. An occurrence covers
// its section when its length is known, and no more, unless the tune
// ends there: the last A of All The Things You Are is 12 bars long.
// Otherwise it runs to the next one when what follows is no longer
// than half of it, a turnaround that differs. A longer stretch is a
// section of its own, and what sounds before the first occurrence too,
// unless it is shorter than MinRun: a pickup, labelled "-".
func (f *former) sections() []Section {
	n := len(f.bars)
	type raw struct {
		label      int // -1 for a stretch that comes back nowhere, -2 for a pickup
		from, bars int
		shift      harmony.Semitones
	}
	var raws []raw
	stretch := func(from, to int) {
		if to <= from {
			return
		}
		if from == 0 && to < n && to < f.o.MinRun {
			raws = append(raws, raw{label: -2, from: from, bars: to})
			return
		}
		if l := to - from; f.o.Carrure && l > 8 && l%8 == 0 {
			for k := 0; k < l; k += 8 {
				raws = append(raws, raw{label: -1, from: from + k, bars: 8})
			}
			return
		}
		raws = append(raws, raw{label: -1, from: from, bars: to - from})
	}
	known := f.lengths()
	pos := 0
	for i, o := range f.placed {
		stretch(pos, o.start)
		end := o.start + f.extent(i, known)
		_, told := known[o.label]
		if next := o.start + f.span(i); (!told || next == n) && (next-end)*2 <= end-o.start {
			end = next
		}
		raws = append(raws, raw{label: o.label, from: o.start, bars: end - o.start, shift: o.shift})
		pos = end
	}
	stretch(pos, n)

	// Letters in the order the sections first sound, transpositions
	// from that first time.
	var out []Section
	letter := map[int]string{}
	first := map[int]harmony.Semitones{}
	next := 'A'
	for _, r := range raws {
		s := Section{From: r.from, Bars: r.bars}
		switch l, ok := letter[r.label]; {
		case r.label == -2:
			s.Label = "-"
		case r.label < 0:
			s.Label = string(next)
			next++
		case ok:
			s.Label, s.Shift = l, mod12(r.shift-first[r.label])
		default:
			letter[r.label], first[r.label] = string(next), r.shift
			s.Label = string(next)
			next++
		}
		out = append(out, s)
	}
	return out
}

// chordFamily returns the chord that stands for `p` with Tolerance: the
// major tonic chords are one (triad, 6, maj7), the minor ones another
// (triad, m6, m7, m(maj7)), and the dominants a third, whatever their
// fifth.
func chordFamily(p harmony.ChordPattern) harmony.ChordPattern {
	switch p {
	case harmony.ChordMajorTriad, harmony.ChordMajorSixth, harmony.ChordMajorSeventh, harmony.ChordMajorSeventhNo5:
		return harmony.ChordMajorSeventh
	case harmony.ChordMinorTriad, harmony.ChordMinorSixth, harmony.ChordMinorSeventh, harmony.ChordMinorSeventhNo5,
		harmony.ChordMinorMajorSeventh, harmony.ChordMinorMajorSeventhNo5:
		return harmony.ChordMinorSeventh
	case harmony.ChordDominantSeventh, harmony.ChordDominantSeventhNo5, harmony.ChordDominantSeventhFlat5, harmony.ChordDominantSeventhSharp5:
		return harmony.ChordDominantSeventh
	}
	return p
}

// mod12 brings a transposition between 0 and 11.
func mod12(t harmony.Semitones) harmony.Semitones {
	return (t%12 + 12) % 12
}
