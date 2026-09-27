package ireal

import "testing"

func structure(chart string) []Measure {
	return Structure(Lex(chart)).Measures
}

// Four measures of four cells, the first opening a repeat with a
// section and a time signature, the last closing it.
func TestStructureMeasures(t *testing.T) {
	ms := structure("{*AT34C^7XyQ|A-7XyQ|D-7 G7LZ x }")
	if len(ms) != 4 {
		t.Fatalf("%d measures, want 4", len(ms))
	}
	first, third, last := ms[0], ms[2], ms[3]
	if first.Open != RepeatOpen || first.Section != 'A' || first.Time != (TimeSig{3, 4}) || !first.NewTime {
		t.Errorf("first %+v", first)
	}
	if ms[1].Time != (TimeSig{3, 4}) || ms[1].NewTime {
		t.Errorf("the time signature is not carried: %+v", ms[1])
	}
	if third.Cells != 4 || len(third.Events) != 2 || third.Events[1].Cell != 2 || third.Events[1].Chord.Root != "G" {
		t.Errorf("two chords in cells 0 and 2 expected: %+v", third)
	}
	if last.Repeat != 1 || last.Close != RepeatClose {
		t.Errorf("last %+v", last)
	}
}

// A run of bar lines between two measures is one boundary: "]" closes,
// "[" opens, and no empty measure comes out of it.
func TestStructureBarRuns(t *testing.T) {
	ms := structure("[CXyQ]|[*BDXyQ}[EXyQZ")
	if len(ms) != 3 {
		t.Fatalf("%d measures, want 3", len(ms))
	}
	if ms[0].Close != DoubleClose || ms[1].Open != DoubleOpen || ms[1].Section != 'B' || ms[1].Close != RepeatClose || ms[2].Close != FinalBar {
		t.Errorf("%+v", ms)
	}
}

// "Kcl" is a bar line, then a measure repeating the one before; "LZ"
// carries a cell before its bar line.
func TestStructureAbbreviations(t *testing.T) {
	ms := structure("F^7XyQKcl LZC7LZ")
	if len(ms) != 3 {
		t.Fatalf("%d measures, want 3", len(ms))
	}
	if ms[1].Repeat != 1 || ms[1].Cells != 4 {
		t.Errorf("repeated measure %+v", ms[1])
	}
	if ms[2].Cells != 2 {
		t.Errorf("C7 then LZ is two cells: %+v", ms[2])
	}
}

// Endings, marks and alternates land on the measure they are written
// in, with the cell they sit at.
func TestStructureMarks(t *testing.T) {
	ms := structure("{SC7XyQ|N1D-7(F7) G7 Q}N2sE-7 lA7<D.S. al Coda>Z")
	if len(ms) != 3 {
		t.Fatalf("%d measures, want 3", len(ms))
	}
	if ms[0].Marks[0].Kind != Segno || ms[0].Marks[0].Cell != 0 {
		t.Errorf("segno %+v", ms[0].Marks)
	}
	m := ms[1]
	if m.Ending != 1 || !m.Events[0].HasAlternate || m.Events[0].Alternate.Root != "F" {
		t.Errorf("first ending %+v", m)
	}
	if len(m.Marks) != 1 || m.Marks[0].Kind != Coda || m.Marks[0].Cell != 4 {
		t.Errorf("coda at the end of the measure %+v", m.Marks)
	}
	last := ms[2]
	if last.Ending != 2 || !last.Events[0].Small || last.Events[1].Small || last.Marks[0].Comment.Text != "D.S. al Coda" {
		t.Errorf("second ending %+v", last)
	}
}

// Empty cells after a repeat or the final bar line are layout, not a
// measure; after a plain bar line, they are a measure where the chord
// goes on.
func TestStructureGaps(t *testing.T) {
	ms := structure("{CXyQ|N1DXyQ}XyQXyQ LZN2EXyQ|XyQZ ")
	if len(ms) != 4 {
		t.Fatalf("%d measures, want 4: %+v", len(ms), ms)
	}
	if ms[2].Ending != 2 || ms[3].Cells != 3 || len(ms[3].Events) != 0 || ms[3].Close != FinalBar {
		t.Errorf("%+v", ms)
	}
}

// A comment after the final bar line is not a measure: it goes to the
// last one.
func TestStructureCommentAfterTheEnd(t *testing.T) {
	ms := structure("[CXyQZ<*48credits> ")
	if len(ms) != 1 || len(ms[0].Marks) != 1 || ms[0].Marks[0].Comment.Text != "credits" {
		t.Errorf("%+v", ms)
	}
}
