// Command tochordpro converts the songs of iReal Pro playlists into
// gohar's profile of ChordPro, one file per song (see docs/formats.md).
//
//	tochordpro -out dir playlist.html...
//
// A playlist is an export of the app: an HTML file holding irealb://
// links, or a text file holding them. The grids an iReal chart holds
// are their authors': what this command writes is a private corpus, to
// keep out of the repository, as testdata/local is.
//
// The form is written as it is played, the coda after {x_coda} (see
// chordpro.FromIReal). A song whose time signature changes is skipped,
// and said so; a chord that cannot be read is written N.C., and said so.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ArnaudCalmettes/gohar/charts/chordpro"
	"github.com/ArnaudCalmettes/gohar/charts/ireal"
)

var link = regexp.MustCompile(`irealb://[^"\s<]*`)

func main() {
	out := flag.String("out", ".", "the folder to write the .cho files in")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: tochordpro [-out dir] <playlist>...")
		os.Exit(2)
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	written, skipped := 0, 0
	for _, path := range flag.Args() {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, l := range link.FindAllString(string(data), -1) {
			p, err := ireal.Parse(l)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
				continue
			}
			for _, song := range p.Songs {
				if err := convert(song, *out); err != nil {
					fmt.Fprintf(os.Stderr, "%s: %v\n", song.Title, err)
					if errors.Is(err, errSkipped) {
						skipped++
						continue
					}
				}
				written++
			}
		}
	}
	fmt.Printf("%d songs written in %s, %d skipped\n", written, *out, skipped)
}

var errSkipped = errors.New("skipped")

// convert writes one song, even with chords it could not read: the
// error says which.
func convert(song ireal.Song, dir string) error {
	s, err := chordpro.FromIReal(song)
	if len(s.Body) == 0 {
		return fmt.Errorf("%w: %v", errSkipped, err)
	}
	f, ferr := os.Create(filepath.Join(dir, chordpro.Name(song.Title)))
	if ferr != nil {
		return ferr
	}
	if werr := chordpro.Write(f, s); werr != nil {
		f.Close()
		return werr
	}
	if cerr := f.Close(); cerr != nil {
		return cerr
	}
	return err
}
