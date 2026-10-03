// Package synth turns note events into a stream of audio samples.
//
// What it produces is an io.Reader over linear PCM. Two libraries
// only: oto, behind Device, to reach the sound card, and go-meltysynth,
// behind Sampler, to play soundfonts; everything else is the standard
// library.
//
// Nothing here knows about harmony. A pitch arrives as a MIDI note
// number and leaves as a frequency, and the conversion lives here
// rather than in the harmony library on purpose: that library exists to
// be free of frequencies and temperament, and handing it a tuning
// standard would undo the abstraction it was built for.
package synth
