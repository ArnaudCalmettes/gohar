//go:build js && wasm

package keyboard

// In the browser, MIDI goes through the Web MIDI API. The driver asks
// for it as the program starts, which shows the browser's permission
// prompt, and waits for the answer. A browser without Web MIDI, Safari,
// would make it panic: the page stubs the API first (games/walk/web),
// and the game then starts with no MIDI input, as on a desktop without
// a keyboard.
import _ "gitlab.com/gomidi/midi/v2/drivers/webmididrv"
