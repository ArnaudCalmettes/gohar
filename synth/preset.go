package synth

// A Preset names one instrument of a soundfont the way General MIDI
// does: a bank and a patch number, 0 to 127. Bank 128 holds the drum
// kits, where each key is a different percussion.
type Preset struct {
	Name  string
	Bank  int
	Patch int
}
