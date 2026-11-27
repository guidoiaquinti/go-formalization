package flate

import (
	"bufio"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"testing"
)

func TestGoLeanTokenState(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_FLATE_TOKEN_STATE_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	input := gob.NewDecoder(bufio.NewReader(f))
	var initial tokens
	for i := range initial.tokens {
		initial.tokens[i] = token(65536 + i)
	}
	for row := 0; ; row++ {
		var c struct {
			Case struct {
				N, Hist uint16
				Filled  int
				Lit     byte
				Data    []byte
				Full    bool
			}
			Indices []uint16
			States  []struct {
				N             uint16
				Filled        int
				Literal       [256]uint16
				Extra, Offset [32]uint16
				Writes        map[uint16]uint32
			}
		}
		if err := input.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d rows agree for all five pinned Go declarations", row)
				return
			}
			t.Fatal(err)
		}
		initial.n, initial.nFilled = c.Case.N, c.Case.Filled
		for i := range initial.litHist {
			initial.litHist[i] = c.Case.Hist
		}
		for i := range initial.extraHist {
			initial.extraHist[i], initial.offHist[i] = c.Case.Hist, c.Case.Hist
		}
		for op, want := range c.States {
			got := initial
			switch op {
			case 0:
				got.Reset()
			case 1:
				got.AddLiteral(c.Case.Lit)
			case 2:
				got.AddEOB()
			case 3:
				emitLiterals(&got, c.Case.Data)
			}
			if got.n != want.N || got.nFilled != want.Filled || got.litHist != want.Literal ||
				got.extraHist != want.Extra || got.offHist != want.Offset {
				t.Fatalf("row %d operation %d: count or full histogram mismatch", row, op)
			}
			checkCell := func(i uint16) {
				value, changed := want.Writes[i]
				if !changed {
					value = uint32(initial.tokens[i])
				}
				if actual := uint32(got.tokens[i]); actual != value {
					t.Fatalf("row %d operation %d backing index %d: got %d want %d", row, op, i, actual, value)
				}
			}
			if c.Case.Full {
				for i := range got.tokens {
					checkCell(uint16(i))
				}
			} else {
				for _, i := range c.Indices {
					checkCell(i)
				}
			}
		}
		got := initial
		view := got.Slice()
		if view == nil || len(view) != int(c.Case.N) || cap(view) != len(initial.tokens) {
			t.Fatalf("row %d Slice: nil=%t len=%d cap=%d", row, view == nil, len(view), cap(view))
		}
		full := view[:cap(view)]
		if full[len(full)-1] != 131071 {
			t.Fatalf("row %d Slice: final backing word mismatch", row)
		}
		full[len(full)-1] = 73
		if got.tokens[len(full)-1] != 73 {
			t.Fatalf("row %d Slice: slice write not visible in receiver", row)
		}
		got.tokens[len(full)-1] = 29
		if full[len(full)-1] != 29 {
			t.Fatalf("row %d Slice: receiver write not visible in slice", row)
		}
	}
}
