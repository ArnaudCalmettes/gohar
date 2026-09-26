package main

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/games/keyboard"
	"github.com/ArnaudCalmettes/gohar/harmony"
)

// How a question sounds: a pedal on the tonic, then the scale climbing
// over it. The colour is read against the pedal.
const (
	pedalVelocity = 0.45
	scaleVelocity = 0.7

	// lead lets the pedal settle before the scale starts, so that the
	// ear has the tonic before it has anything to compare with it.
	lead = 600 * time.Millisecond

	step     = 350 * time.Millisecond
	lastHold = 900 * time.Millisecond

	// gap separates the two modes of a comparison.
	gap = 500 * time.Millisecond
)

// scaleStart returns the MIDI key the scale starts on: the tonic
// between A3 and G sharp 4, so that every scale stays in the same
// comfortable register whatever the tonic.
func scaleStart(tonic harmony.PitchClass) int {
	return 57 + (int(tonic)+3)%12
}

// scaleNotes lays a shape out on `tonic` from `start`, over a pedal on
// the tonic, and returns where it ends.
//
// A scale closes on its octave: one that stops on its seventh leaves
// the ear hanging on the most characteristic note of some modes. A
// smaller shape, a tetrachord, stops where it stops. The last note is
// held.
func scaleNotes(tonic harmony.PitchClass, pattern harmony.ScalePattern, start time.Duration) ([]keyboard.Note, time.Duration) {
	var offsets []int
	for _, o := range pattern.Offsets() {
		offsets = append(offsets, int(o))
	}
	if pattern.IsHeptatonic() {
		offsets = append(offsets, 12)
	}

	first := scaleStart(tonic)
	end := start + lead + time.Duration(len(offsets)-1)*step + lastHold

	notes := []keyboard.Note{{
		Key: first - 24, Velocity: pedalVelocity, Start: start, Length: end - start,
	}}
	at := start + lead
	for i, offset := range offsets {
		length := step
		if i == len(offsets)-1 {
			length = lastHold
		}
		notes = append(notes, keyboard.Note{
			Key: first + offset, Velocity: scaleVelocity, Start: at, Length: length,
		})
		at += step
	}
	return notes, end
}
