package ireal

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarkedSections(t *testing.T) {
	for name, c := range map[string]struct{ chart, want string }{
		// The A under a repeat sounds twice.
		"AABA with a repeat": {
			"{*AC^7XyQ|D-7XyQ|E-7XyQ|A7XyQ}[*BF^7XyQ|G7XyQ|E-7XyQ|A7XyQ][*AC^7XyQ|D-7XyQ|G7XyQ|C^7XyQZ",
			"A4 A4 B4 A4",
		},
		"bars before the first mark": {
			"[C^7XyQ|G7XyQ][*AC^7XyQ|D-7XyQ|G7XyQ|C^7XyQZ",
			"-2 A4",
		},
		"no mark": {"[C^7XyQ|G7XyQ|C^7XyQZ", ""},
	} {
		var parts []string
		for _, s := range Structure(Lex(c.chart)).MarkedSections() {
			parts = append(parts, fmt.Sprintf("%s%d", s.Label, s.Bars))
		}
		if got := strings.Join(parts, " "); got != c.want {
			t.Errorf("%s: %q, want %q", name, got, c.want)
		}
	}
}
