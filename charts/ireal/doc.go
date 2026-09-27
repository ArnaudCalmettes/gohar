// Package ireal reads the chord charts of iReal Pro, as the app shares
// them: an irealb:// URL holding one song or a whole playlist.
//
// It reads, it does not interpret. A song comes out with its metadata
// and its chart as the app wrote it, and [Lex] cuts the chart into
// tokens without losing a character: joined back, the tokens give the
// chart again, which is what will let an edited chart be written back
// one day. Chord symbols stay strings here; turning them into chords is
// the business of a chord symbol parser, not of a file format.
//
// # The format
//
// Nothing of it is documented by its author. What follows comes from
// the free readers that worked it out (ireal-reader in JavaScript, and
// the Lua parser it credits) and was checked against real playlists.
//
//   - The URL is percent encoded. Songs are separated by "===", and a
//     playlist ends with its name.
//   - A song is ten fields separated by "=": title, composer, an empty
//     field, style, key, transposition, chart, accompaniment style,
//     tempo, repeats. Older writers left some of them out; the chart is
//     found by its prefix, and the rest is read around it.
//   - The chart starts with the prefix "1r34LbKcu7" and is scrambled by
//     blocks of fifty characters (see [unscramble]).
//   - Once unscrambled, three abbreviations remain: "XyQ" for empty
//     cells, "LZ" for a bar line, "Kcl" for a bar repeated. The lexer
//     keeps them as they are, each as its own token.
package ireal
