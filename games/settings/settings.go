// Package settings keeps what the games remember between two runs, as
// JSON files in one folder of the user's configuration: the dex of ear,
// the calibrated latency shared by every game (games/calibrate), and
// soon the language (see "Les scènes" in docs/architecture.md).
//
// The desktop keeps files (files.go); the browser keeps local storage
// behind the same three calls (storage_js.go).
package settings
