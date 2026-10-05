// Command sfslim cuts a soundfont down to the presets a game plays,
// and, for a drum kit, to the keys it strikes (see synth/sf2).
//
//	go run ./cmd/sfslim -keep 0:32 -keep 0:0 -keep 128:32/44,51,76,77 -o walk.sf2 GeneralUser-GS.sf2
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ArnaudCalmettes/gohar/synth/sf2"
)

// keeps reads the repeated -keep flag: bank:patch, then for a kit a
// slash and the keys it strikes.
type keeps []sf2.Keep

func (k *keeps) String() string { return fmt.Sprint(*k) }

func (k *keeps) Set(s string) error {
	preset, keys, kit := strings.Cut(s, "/")
	bank, patch, ok := strings.Cut(preset, ":")
	if !ok {
		return fmt.Errorf("%q: bank:patch", s)
	}
	var keep sf2.Keep
	var err error
	if keep.Bank, err = strconv.Atoi(bank); err != nil {
		return err
	}
	if keep.Patch, err = strconv.Atoi(patch); err != nil {
		return err
	}
	if kit {
		for _, key := range strings.Split(keys, ",") {
			n, err := strconv.Atoi(key)
			if err != nil {
				return err
			}
			keep.Keys = append(keep.Keys, n)
		}
	}
	*k = append(*k, keep)
	return nil
}

func main() {
	var keep keeps
	flag.Var(&keep, "keep", "a preset to keep, bank:patch, and for a kit /key,key,…; repeated")
	out := flag.String("o", "slim.sf2", "the soundfont written")
	flag.Parse()
	if flag.NArg() != 1 || len(keep) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	in, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	f, err := sf2.Read(in)
	if err != nil {
		log.Fatal(err)
	}
	s, err := sf2.Slim(f, keep)
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	if err := sf2.Write(&b, s); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, b.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %d presets, %d samples, %.1f MB\n", *out, len(s.Presets)-1, len(s.SampleHdrs)-1, float64(b.Len())/1e6)
}
