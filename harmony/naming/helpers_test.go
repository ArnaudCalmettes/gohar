package naming_test

import (
	"testing"

	"github.com/ArnaudCalmettes/gohar/harmony"
)

// collectOffsets reads a pattern back as the offsets it holds, so that
// a failed assertion prints something a musician can read instead of a
// bit mask.
func collectOffsets(t *testing.T, p harmony.ScalePattern) []harmony.Semitones {
	t.Helper()

	var out []harmony.Semitones
	for _, n := range p.Offsets() {
		out = append(out, n)
	}
	return out
}
