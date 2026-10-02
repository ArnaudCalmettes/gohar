package main

import (
	"fmt"
	"slices"
	"strings"
)

// phrasing reports how the rests of the phrases (see analysis.Phrases)
// agree with the conclusive cadences of the sections found (see
// analysis.Conclusions), where the phrases are meant to rest: a check
// that the two stay together. With `list`, the charts where they part,
// bar by bar.
func phrasing(rs []reading, list bool) string {
	var sections, concluding, both, alone, inside int
	var parted []reading
	for _, r := range rs {
		if len(r.found) == 0 {
			continue
		}
		rests := map[int]bool{}
		for _, bar := range r.rests {
			rests[bar] = true
		}
		ends := map[int]bool{}
		differs := false
		for i, bar := range r.conclusions {
			if r.found[i].Label == "-" {
				continue // a pickup
			}
			sections++
			if bar < 0 {
				continue
			}
			concluding++
			ends[bar] = true
			if rests[bar] {
				both++
			} else {
				alone++
				differs = true
			}
		}
		for _, bar := range r.rests {
			if !ends[bar] {
				inside++
				differs = true
			}
		}
		if differs {
			parted = append(parted, r)
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "\nthe phrases, where the analysis rests, compared with where the sections found conclude\n\n")
	fmt.Fprintf(&b, "%5d  sections\n", sections)
	fmt.Fprintf(&b, "%5d  concluding on a cadence (Siron)\n", concluding)
	fmt.Fprintf(&b, "%5d    where a phrase rests too\n", both)
	fmt.Fprintf(&b, "%5d    where no phrase rests\n", alone)
	fmt.Fprintf(&b, "%5d  rests where no section concludes\n", inside)
	fmt.Fprintf(&b, "%5d  charts where they part\n", len(parted))
	if !list || len(parted) == 0 {
		return b.String()
	}
	slices.SortFunc(parted, func(a, b reading) int { return strings.Compare(a.title, b.title) })
	fmt.Fprintf(&b, "\n  %s %s\n", pad("chart", 44), "form found / rests / conclusions, by bar")
	for _, r := range parted {
		fmt.Fprintf(&b, "  %s %s\n", pad(cut(r.title, 44), 44), lengths(r.found))
		fmt.Fprintf(&b, "  %s rests       %s\n", pad("", 44), bars(r.rests))
		fmt.Fprintf(&b, "  %s conclusions %s\n", pad("", 44), bars(r.conclusions))
	}
	return b.String()
}

// bars writes bar indices from 1, "-" for none.
func bars(bs []int) string {
	var parts []string
	for _, bar := range bs {
		if bar < 0 {
			parts = append(parts, "-")
			continue
		}
		parts = append(parts, fmt.Sprint(bar+1))
	}
	if len(parts) == 0 {
		return "-"
	}
	return strings.Join(parts, " ")
}
