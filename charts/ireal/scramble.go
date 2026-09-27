package ireal

// block is the size of the pieces the app scrambles a chart in.
const block = 50

// unscramble undoes the scrambling of a chart, the prefix removed.
//
// The chart is cut into blocks of fifty characters, and in each the
// first five are swapped with the last five, and characters 10 to 23
// with their mirrors. A block is left as is when fewer than two
// characters would follow it, and so is the tail shorter than a block.
//
// Every swap is its own inverse, so the same function scrambles: that
// is how the tests make charts of their own.
func unscramble(s string) string {
	out := make([]byte, 0, len(s))
	for len(s) > block {
		b := []byte(s[:block])
		s = s[block:]
		if len(s) >= 2 {
			swapMirrors(b)
		}
		out = append(out, b...)
	}
	return string(append(out, s...))
}

// swapMirrors swaps the scrambled positions of one block with their
// mirrors across its middle.
func swapMirrors(b []byte) {
	for i := 0; i < 5; i++ {
		b[i], b[block-1-i] = b[block-1-i], b[i]
	}
	for i := 10; i < 24; i++ {
		b[i], b[block-1-i] = b[block-1-i], b[i]
	}
}
