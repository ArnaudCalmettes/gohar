package ireal

import (
	"strings"
	"testing"
)

// A reference computed by an independent implementation of the
// algorithm the free readers describe: 52 characters, one block
// scrambled and a tail of two left alone.
func TestUnscrambleReference(t *testing.T) {
	in := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	want := "XWVUTfghijNMLKJIHGFEDCBAyzxwvutsrqponmlkOPQRSedcbaYZ"
	if got := unscramble(in); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

// Every swap undoes itself: the same function scrambles, whatever the
// length, around the edge cases of the last block.
func TestUnscrambleIsItsOwnInverse(t *testing.T) {
	for _, n := range []int{0, 1, 49, 50, 51, 52, 99, 100, 101, 102, 177} {
		s := strings.Repeat("0123456789abcdefghijKLMNOPQRSTuvwxyz", 6)[:n]
		if got := unscramble(unscramble(s)); got != s {
			t.Errorf("length %d: round trip lost the chart", n)
		}
	}
}

// A block with fewer than two characters after it is left alone.
func TestUnscrambleSparesTheLastBlockBeforeATinyTail(t *testing.T) {
	s := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY"
	if got := unscramble(s); got != s {
		t.Errorf("51 characters: the block was scrambled into %s", got)
	}
}
