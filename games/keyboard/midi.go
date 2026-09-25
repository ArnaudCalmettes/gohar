package keyboard

// The only file in the repository that imports a MIDI library. Keep it
// that way: everything above talks to Source.

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

// Ports lists the MIDI inputs, in the order OpenMIDI numbers them.
func Ports() []string {
	ins := midi.GetInPorts()
	out := make([]string, len(ins))
	for i, p := range ins {
		out[i] = p.String()
	}
	return out
}

// Shutdown releases the MIDI driver. Call it once, when the program is
// done with MIDI altogether, after closing every source.
//
// Separate from Source.Close because the driver is process wide while a
// source is one device.
func Shutdown() {
	midi.CloseDriver()
}

// OpenMIDI opens a MIDI input.
//
// spec is a number as Ports lists them, or a fragment of a name, or
// empty for the first input. A number is tried first, because that is
// what a listing shows and therefore what anyone will type: matching
// "1" against the names instead quietly picks whichever port happens to
// contain that digit, which on Linux is usually Midi Through, which
// never sends anything.
func OpenMIDI(spec string) (Source, error) {
	ins := midi.GetInPorts()
	if len(ins) == 0 {
		return nil, errors.New("keyboard: no MIDI input found")
	}

	in := ins[0]
	switch {
	case spec == "":
	case isNumber(spec):
		n, _ := strconv.Atoi(spec)
		if n < 0 || n >= len(ins) {
			return nil, fmt.Errorf("keyboard: no MIDI input numbered %d", n)
		}
		in = ins[n]
	default:
		found := false
		for _, p := range ins {
			if strings.Contains(strings.ToLower(p.String()), strings.ToLower(spec)) {
				in, found = p, true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("keyboard: no MIDI input matching %q", spec)
		}
	}

	return &midiSource{in: in}, nil
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

type midiSource struct {
	in   drivers.In
	stop func()
}

func (m *midiSource) Name() string { return m.in.String() }

func (m *midiSource) Listen(recv func(Event)) error {
	if m.stop != nil {
		return errors.New("keyboard: already listening")
	}

	stop, err := midi.ListenTo(m.in, func(msg midi.Message, _ int32) {
		at := time.Now()

		var channel, key, velocity uint8
		switch {
		case msg.GetNoteStart(&channel, &key, &velocity):
			recv(Event{
				Key:      int(key),
				Velocity: float64(velocity) / 127,
				Down:     true,
				At:       at,
			})
		case msg.GetNoteEnd(&channel, &key):
			recv(Event{Key: int(key), At: at})
		}

		// Everything else is dropped on purpose. Control changes, pitch
		// bend, aftertouch and the pads of a controller are real events
		// that a game will want one day, and each will need a decision
		// about what it means. Passing them through untyped now would
		// make that decision by accident.
	})
	if err != nil {
		return fmt.Errorf("keyboard: listening to %s: %w", m.in, err)
	}

	m.stop = stop
	return nil
}

func (m *midiSource) Close() error {
	if m.stop != nil {
		m.stop()
		m.stop = nil
	}
	return nil
}
