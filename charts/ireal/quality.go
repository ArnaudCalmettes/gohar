package ireal

import "strings"

// qualityOffsets gives, for each chord quality the app writes, the
// notes it stands for, in semitones above the root. Positions past the
// octave are kept as such, as harmony.ChordPattern wants them: a ninth
// is 14, not 2; a raised eleventh 18, not 6; a lowered thirteenth 20.
//
// # A table, not a grammar
//
// The app writes from a closed list, so each symbol is read by looking
// it up rather than by a parser that would have to guess. It is also
// what makes each reading checkable one line at a time.
//
// Every entry is already normalised: harmony.ChordPattern.Normalize
// leaves it as it is, and a test holds the table to that. Where a
// symbol and harmony's rules of construction disagree, the rules win,
// since recognition will compare what is played with them. The choices
// a grammar would hide are here in plain sight:
//
//   - The perfect fifth is written out unless the symbol alters it, as
//     in harmony's own tetrads: the symbol names a chord, what a
//     pianist leaves out is a voicing.
//   - A ninth, eleventh or thirteenth implies the seventh and the
//     extensions below it: 13 is 7, 9 and 13; -11 is -7, 9 and 11.
//   - 11 alone, a dominant eleventh without its third, is a 9sus: with
//     no third, the fourth stays a fourth.
//   - 7susadd3 is a 7 with an eleventh: with a third, the fourth
//     becomes an eleventh.
//   - 2 is a sus2: 1, 2, 5. 5 is the bare fifth.
//   - h is the half-diminished seventh, like h7; o is the diminished
//     triad, o7 the diminished seventh.
//   - 7alt is the pandiatonic chord of the altered mode, locrian ♭4:
//     7(♭9, ♭10, ♭5, ♭13). Its "3" is the diminished fourth, heard as a
//     major third by enharmony, which is what makes the quality a 7;
//     its ♭10 sounds where a ♯9 would, and its ♭5 is a fifth, at 6, not
//     a raised eleventh at 18.
var qualityOffsets = map[string][]int{
	"":     {4, 7},
	"-":    {3, 7},
	"o":    {3, 6},
	"+":    {4, 8},
	"5":    {7},
	"2":    {2, 7},
	"sus":  {5, 7},
	"add9": {4, 7, 14},

	"6":     {4, 7, 9},
	"69":    {4, 7, 9, 14},
	"^":     {4, 7, 11},
	"^7":    {4, 7, 11},
	"^9":    {4, 7, 11, 14},
	"^13":   {4, 7, 11, 14, 21},
	"^7#11": {4, 7, 11, 18},
	"^9#11": {4, 7, 11, 14, 18},
	"^7#5":  {4, 8, 11},

	"-6":    {3, 7, 9},
	"-69":   {3, 7, 9, 14},
	"-7":    {3, 7, 10},
	"-9":    {3, 7, 10, 14},
	"-11":   {3, 7, 10, 14, 17},
	"-^7":   {3, 7, 11},
	"-^9":   {3, 7, 11, 14},
	"-b6":   {3, 7, 8},
	"-#5":   {3, 8},
	"-add9": {3, 7, 14},
	"-7b5":  {3, 6, 10},
	"h":     {3, 6, 10},
	"h7":    {3, 6, 10},
	"h9":    {3, 6, 10, 14},
	"o7":    {3, 6, 9},

	"7":      {4, 7, 10},
	"9":      {4, 7, 10, 14},
	"11":     {5, 7, 10, 14},
	"13":     {4, 7, 10, 14, 21},
	"7b5":    {4, 6, 10},
	"7#5":    {4, 8, 10},
	"7b9":    {4, 7, 10, 13},
	"7#9":    {4, 7, 10, 15},
	"7#11":   {4, 7, 10, 18},
	"7b13":   {4, 7, 10, 20},
	"9b5":    {4, 6, 10, 14},
	"9#5":    {4, 8, 10, 14},
	"9#11":   {4, 7, 10, 14, 18},
	"13#11":  {4, 7, 10, 14, 18, 21},
	"13#9":   {4, 7, 10, 15, 21},
	"13b9":   {4, 7, 10, 13, 21},
	"7#9#5":  {4, 8, 10, 15},
	"7#9b5":  {4, 6, 10, 15},
	"7#9#11": {4, 7, 10, 15, 18},
	"7b9#11": {4, 7, 10, 13, 18},
	"7b9b5":  {4, 6, 10, 13},
	"7b9#5":  {4, 8, 10, 13},
	"7b9#9":  {4, 7, 10, 13, 15},
	"7b9b13": {4, 7, 10, 13, 20},
	"7alt":   {4, 6, 10, 13, 15, 20},

	"7sus":     {5, 7, 10},
	"9sus":     {5, 7, 10, 14},
	"13sus":    {5, 7, 10, 14, 21},
	"7b9sus":   {5, 7, 10, 13},
	"7b13sus":  {5, 7, 10, 20},
	"7susadd3": {4, 7, 10, 17},
}

// customSpellings turn the usual ways of writing a quality by hand into
// the app's own, so that a quality typed freely between stars can be
// looked up: "m7" is "-7", "maj7" is "^7". Longest first.
var customSpellings = strings.NewReplacer(
	"maj", "^", "Maj", "^", "MA", "^", "M7", "^7", "Δ", "^",
	"min", "-", "mi", "-", "m", "-",
	"ø", "h", "dim", "o", "°", "o", "aug", "+",
	"(", "", ")", "", ",", "", " ", "",
)

// customAliases are whole qualities typed by hand that no rewriting of
// their parts would reach: "7+" is a 7#5, not a 7 followed by an
// augmented triad.
var customAliases = map[string]string{
	"7+": "7#5",
	"+7": "7#5",
}

// handQualities are chords the app has no quality for, which players
// type by hand, in the app's spelling once rewritten.
//
// "o^7" or "dim(maj7)" is the diminished seventh ♮14, from the scale
// 2-1: in a chord holding a diminished triad, harmony puts the major
// seventh at 23, as a fourteenth.
var handQualities = map[string][]int{
	"o^7": {3, 6, 9, 23},
}

// offsets returns the notes a quality stands for. A quality typed
// freely is first rewritten in the app's spelling.
func offsets(quality string, custom bool) ([]int, bool) {
	if custom {
		if alias, ok := customAliases[quality]; ok {
			quality = alias
		}
		quality = customSpellings.Replace(quality)
		if o, ok := handQualities[quality]; ok {
			return o, true
		}
	}
	o, ok := qualityOffsets[quality]
	return o, ok
}
