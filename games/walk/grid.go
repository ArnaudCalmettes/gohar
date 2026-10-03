package main

import (
	"fmt"

	"github.com/ArnaudCalmettes/gohar/charts/ireal"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// jazzBlues is the grid of the first milestone, as iReal Pro exports
// it: a twelve bar jazz blues, the IV in bar 2, the VI7 in bar 8, a
// II-V in bars 9 and 10, and a turnaround, in F. Written for gohar,
// shared with it.
//
// The chords written small, A♭7 and G♭7, are chromatic dominants on the
// turnaround: the parser sets them aside, and they are the player's to
// add once advanced enough.
const jazzBlues = "irealb://12%20Bar%20Blues=Composer%20Unknown==Medium%20Swing=F=5=1r34LbKcu7X7D%7CQQ%7CBb7QyX%2C7bB%7CQyX7bBQ%7CyX7F%7CQyX7F%7CQyX%7CF6XyyX7F%5ByQ%7CG%2D7XyQ%7CC7XyQ%7CF6%20D7%28Ab7%29LZG7%20C7%28Gb7%29%20%5D%20=Jazz%2DMedium%20Swing=100=30"

// readGrid reads the first song of an iReal Pro link into changes, as
// written.
//
// The field the parser reads as a transposition (5 in jazzBlues) is
// ignored: the app shows this chart in F. iReal is only a way in for a
// grid; taking the player through the twelve keys is gohar's own job.
func readGrid(link string) (analysis.Changes, error) {
	p, err := ireal.Parse(link)
	if err != nil {
		return analysis.Changes{}, err
	}
	if len(p.Songs) == 0 {
		return analysis.Changes{}, fmt.Errorf("walk: no song in the link")
	}
	s := p.Songs[0]
	c, err := ireal.Structure(ireal.Lex(s.Chart)).Changes()
	if err != nil {
		return analysis.Changes{}, fmt.Errorf("walk: %s: %w", s.Title, err)
	}
	return c, nil
}
