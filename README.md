# gohar

Harmony in Go, and games to train the ear with it.

gohar is a music theory library built around the modal harmony taught
by Bernard Maury's school, and a set of small games that use it with a
MIDI keyboard. The library computes; the games make you listen and
play. Every notion a player meets, recognises or plays is collected in
a shared dex, which celebrates what is known and never shows what is
left.

## Repository layout

The repository is a Go workspace of five modules, so that using the
theory library never pulls in a graphics or audio stack.

| Module    | What it holds |
|-----------|---------------|
| `harmony` | The theory core: pitches, intervals, scales, chords, the five mother scales and their 35 modes, tonalities, functions. No note names, no frequencies. |
| `harmony/naming` | Words for the numbers: note spelling, mode, interval and scale names in French and English, as signs (`phrygien ♮6`, `ré♭ majeur`) or words (`phrygien bécarre 6`), and chord symbols in a chosen style (`Cmaj7`, `C-7`, `Cø`). |
| `harmony/analysis` | Deterministic chord recognition, without scoring, and the analysis of a chord chart after the book *En Harmonie*: preparations, cadences, cells, degrees, the tonic the ear senses and its modulations, the form, half cadences and pedals. |
| `dex`     | The player's collection of musical notions, shared by every game. |
| `synth`   | The sound: a small polyphonic synthesiser (sine and 8-bit console timbres), soundfonts, recorded sounds, a mixer, notes scheduled to the sample, and the audio output, tuned for low latency. `synth/sf2` slims a soundfont down to what a game plays. |
| `charts`  | Chord charts: gohar's open format, a profile of ChordPro, read and written; iReal Pro playlists read and converted into it; and commands to see their analysis and their form. |
| `games`   | The playable programs, Ebitengine and MIDI included, and what they share: the canvas and the keyboard on screen, musical time, settings, languages. |

Design notes, in French, live in [`docs/`](docs/): `architecture.md`
for the choices and their reasons, `dex.md` for the collection,
`oreille.md` for the ear trainer, `voicings.md` for chord positions,
`grilles.md` for chart analysis, `formats.md` for the open chart format
(ChordPro and ABC), `walk.md` for the design of *Walk with me*,
`debutants.md` for the course for complete beginners, `glossaire.md`
for the vocabulary, `chantiers.md` for what is open.

## Requirements

- Go 1.27 or later.
- On Linux, a C and C++ toolchain and the ALSA and X11/OpenGL
  development headers, needed by Ebitengine, oto and rtmidi. On Debian
  or Ubuntu:

  ```sh
  sudo apt install build-essential libasound2-dev libgl1-mesa-dev xorg-dev
  ```

- A MIDI keyboard is recommended, not required.
- Nothing to download for the sounds: *Walk with me* embeds a soundfont
  of its own, GeneralUser GS slimmed down to the instruments it plays
  (2.5 MB rather than 32).

Developed and tested on Linux. Other platforms supported by Ebitengine
should work but have not been tried.

## The games

Run them from the `games` directory.

### ear

The ear trainer. Something sounds over a pedal on its tonic, and you
name it. A menu offers four activities, all open from the start:

- **Degrees**: one note of a major scale, and you answer its degree.
  The tonic stays for the whole series.
- **Tetrachords**: four notes, and you name the shape among the four of
  the natural system.
- **Modes**: a mode of the natural system, to name among three.
- **All seven modes**: the same, with the seven offered every time, each
  in the place of its degree.

A mistake counts for nothing: the right answer and yours are played one
after the other, their names shown as they sound, and the question
comes back a little later. What you recognise goes into the dex, saved
in `~/.config/gohar/dex.json`.

```sh
cd games
go run ./ear                  # French, signs
go run ./ear -lang en         # English
go run ./ear -notation words  # si bémol rather than si♭
go run ./ear -port 1          # pick a MIDI input, see keys -list
go run ./ear -timbre square   # an 8-bit voice: pulse12, pulse25, square, triangle, noise
go run ./ear -timbre square -authentic   # with the consoles' raw aliasing
```

| Key | Action |
|-----|--------|
| `1` to `7`, click | pick an activity, answer |
| `R` | listen again, the whole correction after a mistake |
| `Space` | next question; at the end, the same activity again |
| `Enter` | back to the menu, at the end of a series |
| `P` | hide the piano during questions, to work by ear alone |
| `L` | switch between French and English |
| `N` | switch between signs and words |
| `H` | show audio delay figures |
| `Esc` | quit |

### walk

*Walk with me*, the game that teaches the left hand to replace the
bassist. The first palier: a chord chart, the hi-hat to count in, then
the ride to keep going, your left hand on a double bass, your right on
a piano, split at G3, and the roots you land on each change marked on
the chart. Land enough of them and the stick figure snaps his fingers
on 2 and 4: you hear it. He says a word now and then, when you rush,
drag, or swing. Without tempo, the chart waits for each root.

Four charts for now, chords only, written for gohar: a jazz blues in
F, *Satin Doll*, *Tune Up* and *Autumn Leaves*. At rest, ↑ and ↓ choose
the chart, ← and → the tempo, for the session, and D lets the stick
figure play the bass himself, the demo, on from the start without a
MIDI keyboard. A chart longer than
three rows turns its pages a row at a time. A run played through ends
on its review: what worked and what to consolidate, bar by bar on the
chart, how the notes sat on the beat, and one piece of advice.

The title screen plays the blues while you choose; Esc goes back to it
from the game. The options set the tempo and the language, kept from
one run to the next, and open the latency calibration.

```sh
cd games
go run ./walk                  # the title screen, then the space bar starts and stops
go run ./walk -demo            # the band walks the bass itself
go run ./walk -practice        # without tempo; T switches at rest
go run ./walk -record line.txt # writes down every note heard
go run ./walk -lang en         # in English; your session's language by default
```

`-sf2` plays another soundfont than the one embedded, `-list` shows
its presets and `-bass` picks one (`bank:patch`, `0:32` by default):
`make sounds` downloads the whole of GeneralUser GS into
`~/.cache/gohar` for that. `-choruses` sets how many
times the chart goes round, `-bpm` the tempo of the game (the one the
options keep, 120 at first),
`-split` the lowest key of the right hand (55, G3, by default).

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
make slim    # remake the soundfont of Walk with me, after make sounds
make wasm    # Walk with me in WebAssembly, into build/web
make serve   # serve it on http://localhost:8080, for Chrome or Edge
```

Every push to `master` publishes the browser version on GitHub Pages
(`.github/workflows/pages.yml`), at
<https://arnaudcalmettes.github.io/gohar/>. The repository's Pages
setting must name "GitHub Actions" as its source.

`make slim` is only needed when the game plays a new sound: it cuts
GeneralUser GS down again to the presets and drum keys listed in the
Makefile, into `games/walk/sounds/walk.sf2`.

The iReal Pro reader is checked against real playlists that are not
ours to publish. Export yours from the app as HTML into
`charts/ireal/testdata/local/`, which git ignores, and `go test -v
./...` in `charts` reads every chart in them and reports what it could
not read.

To see what the analysis makes of one of them, bar by bar:

```sh
cd charts
go run ./cmd/analyse ireal/testdata/local/playlist.html "tenderly"
```

The degrees are counted in the tonality the analysis hears, and the
heading says when the chart declares another; `-key declared` counts
them in the key the chart declares, `-key F` or `-key A-` in the one
you choose. `-legend` explains the marks before the chart, `-smells`
lists the chord spellings that still smell. The chord symbols are
written Cmaj7, Cm7, Cm7♭5, Cdim7, C6⁄9 by default; `-maj7 natural` or
`-maj7 delta`, `-minus`, `-halfdim`, `-dimsign` and `-slash69` write
C♮7 or CΔ7, C-7, Cø, C°7 and C6/9 instead.

And over whole playlists, where the tonality it hears differs from the
one the app declares, or the one checked by ear when a chart is listed
in `keys.txt`, grouped by how the two relate, the charts that lack what
tells a tonality set aside:

```sh
cd charts
go run ./cmd/corpus -aside ireal/testdata/set-aside.txt -keys ireal/testdata/keys.txt ireal/testdata/local/*.html
```

The open format of gohar's charts is a profile of ChordPro (see
`docs/formats.md`), read and written by `charts/chordpro`. To convert
playlists into it, one `.cho` file per song, for a private corpus kept
out of the repository as the playlists are:

```sh
cd charts
go run ./cmd/tochordpro -out ireal/testdata/local/cho ireal/testdata/local/*.html
```

## History

This repository used to be a single module. Its tagged versions remain
resolvable through the Go module proxy, and that code is kept in
[`gohar-archive`](https://github.com/ArnaudCalmettes/gohar-archive).