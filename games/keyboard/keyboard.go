// Package keyboard is where key presses come from.
//
// # Why this port exists
//
// The analysis engine has no port: the caller drives it, and an
// interface that exists only to be mocked is worth less than the
// parameter it replaces. This one is different. There really are
// several sources of key presses, and a game has no reason to know
// which one it is reading.
//
// A MIDI keyboard is one. A recorded sequence replayed in training
// mode is another, and it is not a test double: a player working a
// phrase without touching the keys is a real mode of a real game.
//
// # What it buys us
//
// Containment. gitlab.com/gomidi/midi is imported by midi.go and by
// nothing else in the repository, the same way oto is imported by
// synth/device.go and nowhere else. Two files, two surfaces.
//
// That is what keeps two doors open. Trading rtmidi, which is C++ and
// drags a toolchain behind it, for a reader of /dev/snd/midiC*D* in
// plain Go is then one file. Compiling to the browser, where the same
// library offers a Web MIDI driver, is a build tag. Neither is planned,
// both stay cheap, and that is the whole point of paying for an
// interface here.
package keyboard

import "time"

// An Event is a key going down or coming up.
//
// Key is a MIDI key number, which is also what synth.Tuning takes. It
// is deliberately not a harmony.Pitch: what a key means musically is
// the analysis layer's business, and this package is a wire.
type Event struct {
	Key      int
	Velocity float64
	Down     bool

	// At is when the event was received, stamped by the source.
	//
	// Not read from the driver, whose own timestamps have a different
	// origin on every platform. What a game needs is a clock it shares
	// with everything else it does, so the source stamps on arrival
	// and every implementation agrees on what the number means.
	At time.Time
}

// A Source produces key events.
//
// Listen hands them over on whatever goroutine the source runs on, and
// the callback must return promptly: it sits on the path between a
// finger and a sound, and anything slow in it is latency the player
// feels. Queue the work, do not do it here.
type Source interface {
	// Name is what to show the player, and what to write in a log when
	// the wrong device was picked.
	Name() string

	// Listen starts delivery. Calling it twice is an error.
	Listen(recv func(Event)) error

	// Close stops delivery. Safe to call without having listened.
	Close() error
}
