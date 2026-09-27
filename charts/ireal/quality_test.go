package ireal

import "testing"

// Every quality the lexer knows has a reading.
func TestEveryQualityIsRead(t *testing.T) {
	for _, q := range qualities {
		if _, ok := qualityOffsets[q]; !ok {
			t.Errorf("quality %q has no reading", q)
		}
	}
	if len(qualityOffsets) != len(qualities)+1 {
		t.Errorf("%d readings for %d qualities and the major triad", len(qualityOffsets), len(qualities))
	}
}

// Offsets stay inside the 24 positions of a chord pattern, strictly
// increasing, and never repeat the root.
func TestQualityOffsetsAreWellFormed(t *testing.T) {
	for q, o := range qualityOffsets {
		for i, n := range o {
			if n <= 0 || n >= 24 || (i > 0 && n <= o[i-1]) {
				t.Errorf("%q: %v", q, o)
			}
		}
	}
}

// A quality typed by hand is read once rewritten in the app's spelling.
func TestCustomQualities(t *testing.T) {
	for custom, like := range map[string]string{
		"m7":        "-7",
		"maj7":      "^7",
		"m7b5":      "-7b5",
		"ø7":        "h7",
		"dim7":      "o7",
		"7(b9,#11)": "7b9#11",
		"mMaj7":     "-^7",
		"7+":        "7#5",
	} {
		got, ok := offsets(custom, true)
		want := qualityOffsets[like]
		if !ok || len(got) != len(want) {
			t.Errorf("%q: %v, want %v like %q", custom, got, want, like)
		}
	}
	for _, dim := range []string{"o^7", "dim(maj7)"} {
		if o, ok := offsets(dim, true); !ok || o[len(o)-1] != 23 {
			t.Errorf("%s: %v, want the diminished seventh natural 14", dim, o)
		}
	}
	if _, ok := offsets("7us", true); ok {
		t.Errorf("a typo is not guessed")
	}
}
