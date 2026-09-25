package synth

import "time"

// HistogramBucket is the width of one bucket of a [Histogram].
const HistogramBucket = time.Millisecond

// histogramBuckets covers 0 to 63 ms, and the last bucket takes
// everything above. A delay past 60 ms is already a failure the ear
// reports on its own; how far past matters less than how often.
const histogramBuckets = 64

// A Histogram counts delays between a key event and the samples, one
// millisecond per bucket.
//
// # Why not only the worst case
//
// A worst case says that something went wrong once. It cannot tell a
// single spike, a garbage collection or a page fault, from a regular
// spread that the ear hears as a soft keyboard. The shape tells them
// apart, and it is what the test under load has to read.
//
// A fixed array: filled on the audio goroutine, which must never
// allocate, and copied by value to whoever wants to read it.
type Histogram struct {
	Counts [histogramBuckets]uint64
}

func (h *Histogram) add(d time.Duration) {
	i := int(d / HistogramBucket)
	switch {
	case i < 0:
		// A clock that went backwards between two goroutines. Counted
		// as immediate rather than dropped, so that the total stays
		// the number of events.
		i = 0
	case i >= histogramBuckets:
		i = histogramBuckets - 1
	}
	h.Counts[i]++
}

// Total returns the number of events counted.
func (h Histogram) Total() uint64 {
	var n uint64
	for _, c := range h.Counts {
		n += c
	}
	return n
}

// Quantile returns the upper bound of the bucket holding the q
// quantile, q between 0 and 1: Quantile(0.99) is the delay 99 events
// in 100 stayed under. Zero when nothing was counted.
//
// An upper bound and not an interpolation: with one millisecond
// buckets, pretending to more precision would only be noise.
func (h Histogram) Quantile(q float64) time.Duration {
	total := h.Total()
	if total == 0 {
		return 0
	}
	target := uint64(q*float64(total) + 0.5)
	if target == 0 {
		target = 1
	}
	var seen uint64
	for i, c := range h.Counts {
		seen += c
		if seen >= target {
			return time.Duration(i+1) * HistogramBucket
		}
	}
	return histogramBuckets * HistogramBucket
}
