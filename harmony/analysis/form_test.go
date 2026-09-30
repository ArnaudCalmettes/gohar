package analysis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
	"github.com/ArnaudCalmettes/gohar/harmony/analysis"
)

// formName writes a form as a musician reads it: each section with its
// letter and its length in bars, and its transposition when there is
// one, "A8 A8 B8 A8".
func formName(sections []analysis.Section) string {
	var parts []string
	for _, s := range sections {
		p := fmt.Sprintf("%s%d", s.Label, s.Bars)
		if s.Shift != 0 {
			p += fmt.Sprintf("+%d", s.Shift)
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, " ")
}

func TestSections(t *testing.T) {
	const (
		c, db, d, eb, e, f, g, ab, a, bb, b harmony.PitchClass = 0, 1, 2, 3, 4, 5, 7, 8, 9, 10, 11
	)
	const (
		maj7    = harmony.ChordMajorSeventh
		min7    = harmony.ChordMinorSeventh
		min6    = harmony.ChordMinorSixth
		six     = harmony.ChordMajorSixth
		dom7    = harmony.ChordDominantSeventh
		halfDim = harmony.ChordHalfDiminished
	)
	type bar = []any
	join := func(parts ...[]bar) []bar {
		var out []bar
		for _, p := range parts {
			out = append(out, p...)
		}
		return out
	}

	// An A of 8 bars with its turnaround, and the same A whose last
	// two bars lead to the bridge instead, as in Billy Boy.
	a1 := []bar{{c, maj7}, {a, min7}, {d, min7}, {g, dom7}, {e, min7}, {a, dom7}, {d, min7}, {g, dom7}}
	a2 := []bar{{c, maj7}, {a, min7}, {d, min7}, {g, dom7}, {e, min7}, {a, dom7}, {g, min7}, {c, dom7}}
	bridge := []bar{{f, maj7}, {f, maj7}, {f, min6}, {f, min6}, {e, min7}, {a, dom7}, {d, min7}, {g, dom7}}

	// An A of 14 bars, as in Alone Together, the second one with other
	// last bars, then a bridge and a last section that come back
	// nowhere.
	long1 := []bar{{d, min6}, {e, halfDim, a, dom7}, {d, min6}, {e, halfDim, a, dom7}, {a, halfDim}, {d, dom7},
		{g, min7}, {g, min7}, {b, halfDim}, {e, dom7}, {a, min7}, {d, dom7}, {g, maj7}, {c, dom7}}
	long2 := append(append([]bar{}, long1[:12]...), bar{e, halfDim}, bar{a, dom7})
	longB := []bar{{f, maj7}, {f, maj7}, {bb, dom7}, {bb, dom7}, {eb, maj7}, {ab, dom7}, {db, maj7}, {a, dom7}}
	longC := []bar{{d, min6}, {e, halfDim, a, dom7}, {d, min6}, {bb, dom7}, {e, halfDim}, {a, dom7}, {d, min6},
		{e, halfDim, a, dom7}}

	// The two halves of an ABAC start alike, and part at the middle.
	abacB := []bar{{e, dom7}, {e, dom7}, {a, min7}, {a, min7}, {d, dom7}, {d, dom7}, {d, min7}, {g, dom7}}
	abacC := []bar{{f, maj7}, {bb, dom7}, {e, min7}, {a, dom7}, {d, min7}, {g, dom7}, {c, six}, {c, six}}

	// An A whose second time matches on 5 bars only, then goes its own
	// way for 3.
	a1Short := append(append([]bar{}, a1[:5]...), bar{bb, dom7}, bar{eb, maj7}, bar{ab, dom7})

	// A tune of 16 bars in one breath, nothing coming back.
	breath := append(append([]bar{}, long1...), bar{f, maj7}, bar{bb, dom7})

	for name, tc := range map[string]struct {
		bars []bar
		want string
	}{
		"AABA, the second A leading to the bridge": {join(a1, a2, bridge, a1), "A8 A8 B8 A8"},
		"ABAC": {join(a1, abacB, a1, abacC), "A8 B8 A8 C8"},
		"sections of 14 bars, as in Alone Together": {join(long1, long2, longB, longC), "A14 A14 B8 C8"},
		"an A matching on 5 bars covers 8":          {join(a1, abacB, a1Short, abacC), "A8 B8 A8 C8"},
		"a tune in one breath, cut in 8":            {breath, "A8 B8"},
		// A pickup of 2 bars: the chart ends 2 bars into its last A,
		// which says nothing of the length of an A.
		"AABA after a pickup": {
			join([]bar{{bb, min7}, {eb, dom7}}, a1, a2, bridge, a1[:6]),
			"-2 A8 A8 B8 A6",
		},
		// A blues: three phrases of 4 bars, none coming back.
		"a blues": {
			[]bar{{f, dom7}, {bb, dom7}, {f, dom7}, {f, dom7}, {bb, dom7}, {bb, dom7}, {f, dom7}, {d, dom7},
				{g, min7}, {c, dom7}, {f, dom7}, {c, dom7}},
			"A12",
		},
	} {
		if got := formName(analysis.Sections(barsOf(true, tc.bars...))); got != tc.want {
			t.Errorf("%s: %s, want %s", name, got, tc.want)
		}
	}
	if got := analysis.Sections(changesOf(false, c, maj7, g, dom7)); got != nil {
		t.Errorf("changes without bars: %s, want no form", formName(got))
	}
}
