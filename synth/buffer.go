package synth

import "time"

// DefaultBuffer is the starting point, not a recommendation. Raise it
// when a machine crackles, and expect to have to. It is the desktop's:
// in the browser, the sound goes through synth/web, whose buffers the
// browser sizes itself.
const DefaultBuffer = 10 * time.Millisecond
