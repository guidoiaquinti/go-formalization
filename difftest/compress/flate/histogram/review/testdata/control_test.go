package flate

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"slices"
	"testing"
)

func TestReviewGoLeanHistogram(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_REVIEW_FLATE_HISTOGRAM"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	for row := 0; ; row++ {
		var c struct {
			Data          []byte
			Backing, Want []uint16
			Length        int
			Panic         bool
		}
		if err := decoder.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row == 5 {
				return
			}
			t.Fatal(err)
		}
		for _, run := range []func([]byte, []uint16){histogram, histogramSplit} {
			backing := slices.Clone(c.Backing)
			panicked := func() (panicked bool) {
				defer func() { panicked = recover() != nil }()
				run(c.Data, backing[:c.Length:len(backing)])
				return false
			}()
			if panicked != c.Panic || !slices.Equal(backing, c.Want) {
				t.Fatalf("row %d panic %t want %t, backing match %t", row,
					panicked, c.Panic, slices.Equal(backing, c.Want))
			}
		}
	}
}
