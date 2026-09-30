package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// How the form the analysis finds relates to the one the chart marks,
// in the order the report lists them. The letters are compared once
// renamed in the order they first sound, so that a form marked with
// other letters is still the same form.
const (
	sameForm   = "same form"
	sameBounds = "same sections, other letters (ABAC marked, AABC found)"
	nearBounds = "sections within 2 bars (a pickup, a turnaround)"
	otherLevel = "same form at another level (A16 B16 marked, ABAC found)"
	otherForm  = "other form"
	unmarked   = "no marks, or a single one"
)

var formOrder = []string{sameBounds, nearBounds, otherLevel, otherForm, unmarked, sameForm}

// formRelation tells how the sections found relate to those marked.
// One form is the other at another level when all the sections of one
// start where sections of the other do: A16 B16 marked is ABAC found,
// its halves (Siron's levels nested in one another).
func formRelation(found, marked []analysis.Section) string {
	found, marked = withoutPickup(found), withoutPickup(marked)
	switch {
	case len(marked) < 2:
		return unmarked
	case len(found) != len(marked) && (startsWithin(found, marked) || startsWithin(marked, found)):
		return otherLevel
	case len(found) != len(marked):
		return otherForm
	}
	same, near := true, true
	for i := range found {
		d := found[i].From - marked[i].From
		same = same && d == 0
		near = near && d >= -2 && d <= 2
	}
	switch {
	case same && letters(found) == letters(marked):
		return sameForm
	case same:
		return sameBounds
	case near:
		return nearBounds
	}
	return otherForm
}

// withoutPickup returns the sections after a pickup of 3 bars or
// less, the first bar of the first section counted as 0: the analysis
// finds the A of I Should Care from its second bar, where the chart
// marks it from the first, the pickup included.
func withoutPickup(ss []analysis.Section) []analysis.Section {
	if len(ss) < 2 || ss[0].Bars > 3 {
		return ss
	}
	out := slices.Clone(ss[1:])
	for i := range out {
		out[i].From -= ss[0].Bars
	}
	return out
}

// startsWithin reports whether every section of `a` starts where a
// section of `b` does.
func startsWithin(a, b []analysis.Section) bool {
	starts := map[int]bool{}
	for _, s := range b {
		starts[s.From] = true
	}
	for _, s := range a {
		if !starts[s.From] {
			return false
		}
	}
	return true
}

// letters returns the letters of a form renamed in the order they first
// sound: "iAABA" becomes "ABBCB". A section transposed counts as
// another, as a chart marks it: the bridge of So What, its A a semitone
// up, is its B.
func letters(ss []analysis.Section) string {
	names := map[string]byte{}
	var b strings.Builder
	for _, s := range ss {
		key := fmt.Sprintf("%s%+d", s.Label, s.Shift)
		n, ok := names[key]
		if !ok {
			n = byte('A' + len(names))
			names[key] = n
		}
		b.WriteByte(n)
	}
	return b.String()
}

// lengths writes a form as the chart has it, each section with its
// label and its length in bars: "A8 A8 B8 A8".
func lengths(ss []analysis.Section) string {
	var parts []string
	for _, s := range ss {
		p := fmt.Sprintf("%s%d", s.Label, s.Bars)
		if s.Shift != 0 {
			p += fmt.Sprintf("+%d", s.Shift)
		}
		parts = append(parts, p)
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, " ")
}

// forms reports how the forms found agree with the marks of the charts:
// the count of each relation, and with `list`, the charts that
// disagree, both forms side by side.
func forms(rs []reading, list bool) string {
	groups := map[string][]reading{}
	for _, r := range rs {
		rel := formRelation(r.found, r.marked)
		groups[rel] = append(groups[rel], r)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "\nthe form, found from the chords and compared with the marks of the chart\n\n")
	for _, rel := range formOrder {
		fmt.Fprintf(&b, "%5d  %s\n", len(groups[rel]), rel)
	}
	if !list {
		return b.String()
	}
	for _, rel := range formOrder {
		g := groups[rel]
		if rel == sameForm || rel == unmarked || len(g) == 0 {
			continue
		}
		slices.SortFunc(g, func(a, b reading) int { return strings.Compare(a.title, b.title) })
		fmt.Fprintf(&b, "\n%s (%d)\n\n", rel, len(g))
		fmt.Fprintf(&b, "  %s %s\n", pad("chart", 44), "marked / found")
		for _, r := range g {
			fmt.Fprintf(&b, "  %s %s\n  %s %s\n", pad(cut(r.title, 44), 44), lengths(r.marked), pad("", 44), lengths(r.found))
		}
	}
	return b.String()
}
