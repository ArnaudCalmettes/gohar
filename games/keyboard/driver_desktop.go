//go:build !js

package keyboard

// On the desktop, MIDI goes through RtMidi: ALSA on Linux, CoreMIDI on
// macOS, WinMM on Windows (see "Deux surfaces" in docs/architecture.md).
import _ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
