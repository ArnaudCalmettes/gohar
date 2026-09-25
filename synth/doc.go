// Package synth turns note events into a stream of audio samples.
//
// It depends on nothing but the standard library. What it produces is
// an io.Reader over linear PCM, which is what a game hands to its audio
// context; the engine that consumes it is the game's business, not this
// package's.
//
// Nothing here knows about harmony. A pitch arrives as a MIDI note
// number and leaves as a frequency, and the conversion lives here
// rather than in the harmony library on purpose: that library exists to
// be free of frequencies and temperament, and handing it a tuning
// standard would undo the abstraction it was built for.
package synth
