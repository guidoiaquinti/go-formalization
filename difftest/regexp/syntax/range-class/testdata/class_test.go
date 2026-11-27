package syntax

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func rangeWire(raw []rune) string {
	cells := make([]string, len(raw))
	for i, value := range raw {
		cells[i] = strconv.Itoa(int(value))
	}
	return strings.Join(cells, ",")
}

func runRangeOperation(original []rune, size int, operation func([]rune) []rune) string {
	backing := slices.Clone(original)
	result := operation(backing[:size:len(backing)])
	reused := cap(result) == 0 && len(backing) == 0
	if cap(result) > 0 && len(backing) > 0 {
		reused = &result[:cap(result)][0] == &backing[0]
	}
	return fmt.Sprintf("%s:%s:%t", rangeWire(result), rangeWire(backing), reused)
}

type rangeRef struct {
	out, old []rune
	reused   bool
}

func (s *rangeRef) add(lo, hi rune) {
	for i := len(s.out) - 2; i >= max(0, len(s.out)-4); i -= 2 {
		if s.out[i+1]+1 < lo || hi+1 < s.out[i] {
			continue
		}
		s.out[i], s.out[i+1] = min(s.out[i], lo), max(s.out[i+1], hi)
		if s.reused {
			copy(s.old[i:], s.out[i:i+2])
		}
		return
	}
	s.push(lo, hi)
}
func (s *rangeRef) push(lo, hi rune) {
	if s.reused && len(s.out)+2 <= len(s.old) {
		copy(s.old[len(s.out):], []rune{lo, hi})
	} else {
		s.reused = false
	}
	s.out = append(s.out, lo, hi)
}
func (s rangeRef) wire() string {
	return fmt.Sprintf("%s:%s:%t", rangeWire(s.out), rangeWire(s.old), s.reused)
}

func TestGoLeanRangeClass(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_REGEXP_RANGE_CLASS_ROWS"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	in := gob.NewDecoder(bufio.NewReader(f))
	indexedRows := 0
	for row := 0; ; row++ {
		var c struct {
			Backing, Clean []rune
			Size, I, J     int
			Lo, Hi, Rune   rune
			Want           string
		}
		if err := in.Decode(&c); err != nil {
			if errors.Is(err, io.EOF) && row > 0 {
				t.Logf("%d pinned range rows, including %d valid-index Less and Swap rows", row, indexedRows)
				return
			}
			t.Fatal(err)
		}
		makeRef := func() rangeRef { return rangeRef{slices.Clone(c.Backing[:2*c.Size]), slices.Clone(c.Backing), true} }
		gotBacking := slices.Clone(c.Backing)
		visible := gotBacking[:2*c.Size]
		pairs := ranges{&visible}
		gotCompare, wantCompare := "skip", "skip"
		if c.Size > 0 {
			indexedRows++
			less := pairs.Less(c.I, c.J)
			pairs.Swap(c.I, c.J)
			gotCompare = fmt.Sprintf("%t:%s", less, rangeWire(gotBacking))
			old := slices.Clone(c.Backing)
			i, j := 2*c.I, 2*c.J
			less = old[i] < old[j] || old[i] == old[j] && old[j+1] < old[i+1]
			a, b := old[i], old[i+1]
			old[i], old[i+1] = old[j], old[j+1]
			old[j], old[j+1] = a, b
			wantCompare = fmt.Sprintf("%t:%s", less, rangeWire(old))
		}
		refRange, refClass, refAppendNeg := makeRef(), makeRef(), makeRef()
		refRange.add(c.Lo, c.Hi)
		refNeg := rangeRef{nil, slices.Clone(c.Clean), true}
		has := false
		prev := rune(-1)
		for i := 0; i < len(c.Clean); i += 2 {
			lo, hi := c.Clean[i], c.Clean[i+1]
			refClass.add(lo, hi)
			if lo <= c.Rune && c.Rune <= hi {
				has = true
			}
			if prev+1 < lo {
				refAppendNeg.add(prev+1, lo-1)
				refNeg.push(prev+1, lo-1)
			}
			prev = hi
		}
		if prev < unicode.MaxRune {
			refAppendNeg.add(prev+1, unicode.MaxRune)
			refNeg.push(prev+1, unicode.MaxRune)
		}
		got := fmt.Sprintf("%d|%s|%s|%s|%s|%s|%t", pairs.Len(), gotCompare,
			runRangeOperation(c.Backing, 2*c.Size, func(r []rune) []rune { return appendRange(r, c.Lo, c.Hi) }),
			runRangeOperation(c.Backing, 2*c.Size, func(r []rune) []rune { return appendClass(r, c.Clean) }),
			runRangeOperation(c.Backing, 2*c.Size, func(r []rune) []rune { return appendNegatedClass(r, c.Clean) }),
			runRangeOperation(c.Clean, len(c.Clean), negateClass), inCharClass(c.Rune, c.Clean))
		ref := fmt.Sprintf("%d|%s|%s|%s|%s|%s|%t", c.Size, wantCompare, refRange.wire(), refClass.wire(), refAppendNeg.wire(), refNeg.wire(), has)
		if got != ref || got != c.Want {
			t.Fatalf("row %d: Go %q Lean %q oracle %q", row, got, c.Want, ref)
		}
	}
}
