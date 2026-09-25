// Package games holds the playable examples.
//
// This is where the dependencies live: the game engine, the entity
// system, the MIDI adapter. Its own module for exactly that reason, so
// that importing the harmony library to transpose three chords does not
// pull in a graphics stack.
package games
