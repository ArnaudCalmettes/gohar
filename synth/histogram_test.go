package synth

import (
	"testing"
	"time"
)

func TestHistogramBuckets(t *testing.T) {
	var h Histogram
	h.add(-time.Millisecond)
	h.add(2500 * time.Microsecond)
	h.add(time.Second)

	if h.Counts[0] != 1 || h.Counts[2] != 1 || h.Counts[histogramBuckets-1] != 1 {
		t.Errorf("counts %v", h.Counts)
	}
	if h.Total() != 3 {
		t.Errorf("total %d, want 3", h.Total())
	}
}

// A spread and a spike give the same worst case and different shapes,
// which is the whole reason for the histogram.
func TestHistogramQuantile(t *testing.T) {
	var h Histogram
	for range 99 {
		h.add(4 * time.Millisecond)
	}
	h.add(40 * time.Millisecond)

	if got := h.Quantile(0.5); got != 5*time.Millisecond {
		t.Errorf("median %v, want 5ms", got)
	}
	if got := h.Quantile(0.99); got != 5*time.Millisecond {
		t.Errorf("p99 %v, want 5ms: one spike in a hundred", got)
	}
	if got := h.Quantile(1); got != 41*time.Millisecond {
		t.Errorf("max %v, want 41ms", got)
	}
	if got := (Histogram{}).Quantile(0.5); got != 0 {
		t.Errorf("empty %v, want 0", got)
	}
}

// Read fills the histogram, and does not allocate while doing it.
func TestEngineRecordsDelays(t *testing.T) {
	var e Engine
	buf := make([]byte, 256*BytesPerFrame)

	e.NoteOn(60, 1, time.Now().Add(-3500*time.Microsecond))
	if _, err := e.Read(buf); err != nil {
		t.Fatal(err)
	}
	h := e.Histogram()
	if h.Total() != 1 || h.Counts[3] != 1 {
		t.Errorf("counts %v, want one event in the 3 ms bucket", h.Counts)
	}

	allocs := testing.AllocsPerRun(100, func() {
		e.NoteOn(62, 1, time.Now())
		_, _ = e.Read(buf)
	})
	if allocs != 0 {
		t.Errorf("Read allocates %v times per call", allocs)
	}
}
