// Package harmony provides primitives for music theory computation.
//
// It knows no note names, no language and no notation, and no
// frequencies either: a pitch here is a MIDI note number, and what it
// sounds like at a given tuning is the synth's business. This package
// manipulates numbers and sets. Everything that turns a value into
// something a human reads or hears named lives in the naming
// subpackage.
//
// This separation is the organising principle of the library. A pitch
// class in this package is a number between 0 and 11; whether it is
// written C sharp or D flat is a question of spelling, and spelling
// needs a tonal context that the core has no business guessing.
package harmony
