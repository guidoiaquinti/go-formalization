package flate

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"slices"
	"testing"
)

func TestGoLeanHistogram(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_HISTOGRAM_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := gob.NewDecoder(bufio.NewReader(f))
	for row := 0; ; row++ {
		var c struct {
			Data    []byte
			Backing []uint16
			Length  int
			Nil     bool
		}
		if err := input.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d rows agree with independent oracle for both pinned Go entry points", row)
				return
			}
			t.Fatal(err)
		}
		short := len(c.Backing) < 256
		want := slices.Clone(c.Backing)
		if !short {
			frequencies := [256]uint64{}
			for _, b := range c.Data {
				frequencies[b]++
			}
			for i, n := range frequencies {
				want[i] = uint16(uint64(want[i]) + n)
			}
		}
		for _, entry := range []struct {
			name string
			run  func([]byte, []uint16)
		}{{"histogram", histogram}, {"histogramSplit", histogramSplit}} {
			var backing []uint16
			if !c.Nil {
				backing = make([]uint16, len(c.Backing))
				copy(backing, c.Backing)
			}
			panicked := func() (panicked bool) {
				defer func() { panicked = recover() != nil }()
				entry.run(c.Data, backing[:c.Length])
				return false
			}()
			if panicked != short || !slices.Equal(backing, want) {
				t.Fatalf("row %d %s (len %d cap %d data len %d): panic=%t, want %t, backing match=%t",
					row, entry.name, c.Length, len(backing), len(c.Data), panicked, short, slices.Equal(backing, want))
			}
		}
	}
}
