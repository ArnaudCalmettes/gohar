# gohar

Harmony in Go, and games to train the ear with it.

gohar is a music theory library built around the modal harmony taught
by Bernard Maury's school, and a set of small games that use it with a
MIDI keyboard. The library computes; the games make you listen and
play. Every notion a player meets, recognises or plays is collected in
a shared dex, which celebrates what is known and never shows what is
left.

## Repository layout

The repository is a Go workspace of four modules, so that using the
theory library never pulls in a graphics or audio stack.

| Module    | What it holds |
|-----------|---------------|
| `harmony` | The theory core: pitches, intervals, scales, chords, the five mother scales and their 35 modes, tonalities, functions. No note names, no frequencies. |
| `harmony/naming` | Words for the numbers: note spelling, mode names in French and English, as signs (`phrygien ♮6`) or words (`phrygien bécarre 6`). |
| `harmony/analysis` | Deterministic chord recognition, without scoring. |
| `dex`     | The player's collection of musical notions, shared by every game. |
| `synth`   | A small polyphonic synthesiser (sine and 8-bit console timbres) and the audio output, tuned for low latency. |
| `games`   | The playable programs, Ebitengine and MIDI included. |

Design notes, in French, live in [`docs/`](docs/): `architecture.md`
for the choices and their reasons, `dex.md` for the collection,
`glossaire.md` for the vocabulary, `chantiers.md` for what is open.

## Requirements

- Go 1.27 or later.
- On Linux, a C and C++ toolchain and the ALSA and X11/OpenGL
  development headers, needed by Ebitengine, oto and rtmidi. On Debian
  or Ubuntu:

  ```sh
  sudo apt install build-essential libasound2-dev libgl1-mesa-dev xorg-dev
  ```

- A MIDI keyboard is recommended, not required.

Developed and tested on Linux. Other platforms supported by Ebitengine
should work but have not been tried.

## The games

Run them from the `games` directory.

### ear

The first ear training activity. A mode of the natural system sounds
over a pedal on its tonic, and you name it among three. A mistake counts
for nothing: the two modes are played one after the other, and the
mode comes back a little later on another tonic. What you recognise
goes into the dex, saved in `~/.config/gohar/dex.json`.

```sh
cd games
go run ./ear                  # French, signs
go run ./ear -lang en         # English
go run ./ear -notation words  # si bémol rather than si♭
go run ./ear -port 1          # pick a MIDI input, see keys -list
go run ./ear -timbre square   # an 8-bit voice: pulse12, pulse25, square, triangle
go run ./ear -timbre square -authentic   # with the consoles' raw aliasing
```

| Key | Action |
|-----|--------|
| `1` `2` `3`, click | answer |
| `R` | listen again |
| `Space` | next question |
| `L` | switch between French and English |
| `N` | switch between signs and words |
| `H` | show audio delay figures |
| `Q`, `Esc` | quit |

### keys

Plays your MIDI keyboard through the synthesiser, and nothing else. It
is the latency check: if it feels soft under the fingers, the games
will too.

```sh
cd games
go run ./keys -list           # list MIDI inputs
go run ./keys -port 1         # play
```

`-device` sets the audio buffer, `-a4` the tuning. `latency` and
`otolatency`, in the same directory, are the probes that measured the
audio path.

## Development

`go test ./...` does not cross module boundaries, so the Makefile runs
each module in turn:

```sh
make test    # every module
make vet
make fmt
make bench   # harmony benchmarks
```

## History

This repository used to be a single module. Its tagged versions remain
resolvable through the Go module proxy, and that code is kept in
[`gohar-archive`](https://github.com/ArnaudCalmettes/gohar-archive).