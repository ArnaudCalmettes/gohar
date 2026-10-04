package calibrate

import (
	"time"

	"github.com/ArnaudCalmettes/gohar/games/settings"
)

// File is where the offsets are kept, among the settings shared by the
// games: the latency belongs to the machine, not to a game.
const File = "latency.json"

// Key names the pair the offset belongs to: the MIDI keyboard, and the
// audio output as the game can describe it. Plugging another keyboard,
// or changing what `output` says, asks for a new calibration.
func Key(keyboard, output string) string {
	return keyboard + " | " + output
}

// offsets are the measures kept, by pair, in milliseconds: readable in
// the file, and finer than any player.
type offsets map[string]int64

// Load returns the offset measured for `key` in the file at `path`, and
// whether there is one.
func Load(path, key string) (time.Duration, bool, error) {
	o := offsets{}
	if _, err := settings.Load(path, &o); err != nil {
		return 0, false, err
	}
	ms, ok := o[key]
	return time.Duration(ms) * time.Millisecond, ok, nil
}

// Save keeps `d` for `key` in the file at `path`, beside the other
// pairs.
func Save(path, key string, d time.Duration) error {
	o := offsets{}
	if _, err := settings.Load(path, &o); err != nil {
		return err
	}
	o[key] = d.Milliseconds()
	return settings.Save(path, o)
}
