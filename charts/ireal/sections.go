package ireal

import "github.com/ArnaudCalmettes/gohar/harmony/analysis"

// MarkedSections returns the sections the chart declares by its
// rehearsal marks, "*A", "*B", bar by bar as the chart is played: a
// section under a repeat sounds twice, and each mark starts a section
// that runs to the next one. The label is the mark itself, "A" to "D",
// "i" for an introduction, "v" for a verse. Bars played before the
// first mark make a section labelled "-". It returns nil when the chart
// has no mark.
//
// The analysis does not read these marks: it finds the sections from
// the chords (see analysis.Sections). They are there to compare with,
// as the key the app declares, and are no more reliable.
func (c Chart) MarkedSections() []analysis.Section {
	var out []analysis.Section
	marked := false
	for bar, p := range c.Unfold() {
		switch mark := c.Measures[p.Index].Section; {
		case mark != 0:
			out = append(out, analysis.Section{Label: string(mark), From: bar})
			marked = true
		case len(out) == 0:
			out = append(out, analysis.Section{Label: "-", From: bar})
		}
		out[len(out)-1].Bars++
	}
	if !marked {
		return nil
	}
	return out
}
