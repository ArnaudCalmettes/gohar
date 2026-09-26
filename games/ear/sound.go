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

// modeNotes lays one mode out from `start`, and returns where it ends.
//
// Eight notes, the octave included: a scale that stops on its seventh
// leaves the ear hanging on the most characteristic note of some modes.
func modeNotes(tonic harmony.PitchClass, mode harmony.Degree, start time.Duration) ([]keyboard.Note, time.Duration) {
	pattern, _ := system.Mode(mode)
	first := scaleStart(tonic)
	end := start + lead + 7*step + lastHold

	notes := []keyboard.Note{{
		Key: first - 24, Velocity: pedalVelocity, Start: start, Length: end - start,
	}}
	at := start + lead
	for _, offset := range pattern.Offsets() {
		notes = append(notes, keyboard.Note{
			Key: first + int(offset), Velocity: scaleVelocity, Start: at, Length: step,
		})
		at += step
	}
	notes = append(notes, keyboard.Note{
		Key: first + 12, Velocity: scaleVelocity, Start: at, Length: lastHold,
	})
	return notes, end
}

// questionNotes is what the player hears when a question is asked.
func questionNotes(q Question) []keyboard.Note {
	notes, _ := modeNotes(q.Tonic, q.Mode, 0)
	return notes
}

// comparisonNotes plays the right mode, then the one chosen by
// mistake, on the same tonic. It is the correction being shown.
func comparisonNotes(q Question, chosen harmony.Degree) []keyboard.Note {
	right, end := modeNotes(q.Tonic, q.Mode, 0)
	wrong, _ := modeNotes(q.Tonic, chosen, end+gap)
	return append(right, wrong...)
}
