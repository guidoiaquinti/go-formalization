package sort

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type goLeanPatternSwaps [][2]int

func (*goLeanPatternSwaps) Len() int           { panic("unexpected Len") }
func (*goLeanPatternSwaps) Less(int, int) bool { panic("unexpected Less") }
func (s *goLeanPatternSwaps) Swap(a, b int)    { *s = append(*s, [2]int{a, b}) }

func (s goLeanPatternSwaps) describe() string {
	fields := make([]string, len(s))
	for i, pair := range s {
		fields[i] = fmt.Sprintf("%d,%d", pair[0], pair[1])
	}
	return "swaps:" + strings.Join(fields, ";")
}

// Multiplication and division by fixed powers of two give the stage
// transitions on the finite ring, with xor combining successive words.
func goLeanPatternNext(word uint64) uint64 {
	word ^= word * 8192
	word ^= word / 128
	return word ^ (word * 131072)
}

func goLeanPatternPower(n uint) uint {
	for power := uint(1); power != 0; power *= 2 {
		if power > n {
			return power
		}
	}
	return 0
}

func goLeanPatternOracle(a, b int) goLeanPatternSwaps {
	size := b - a
	if size < 8 {
		return nil
	}
	word := uint64(size)
	center := a + (size/4)*2
	power := goLeanPatternPower(uint(size))
	var swaps goLeanPatternSwaps
	for _, offset := range []int{-1, 0, 1} {
		word = goLeanPatternNext(word)
		dst := (uint(word) % power) % uint(size)
		swaps.Swap(center+offset, a+int(dst))
	}
	return swaps
}

func TestGoLeanBreakPatterns(t *testing.T) {
	f, err := os.Open(os.Getenv("GOLEAN_SORT_BREAK_PATTERNS_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var counts [4]int
	for scanner.Scan() {
		var c struct {
			Kind, A, B int
			State      uint64
			Want       string
		}
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		var native, want string
		switch c.Kind {
		case 0:
			state := xorshift(c.State)
			result := state.Next()
			if uint64(state) != result {
				t.Fatalf("Next returned %d but stored %d", result, state)
			}
			native = strconv.FormatUint(result, 10)
			want = strconv.FormatUint(goLeanPatternNext(c.State), 10)
		case 1:
			native = strconv.FormatUint(uint64(nextPowerOfTwo(c.A)), 10)
			want = strconv.FormatUint(uint64(goLeanPatternPower(uint(c.A))), 10)
		default:
			var trace goLeanPatternSwaps
			if c.Kind == 2 {
				breakPatterns(&trace, c.A, c.B)
			} else {
				breakPatterns_func(lessSwap{Less: trace.Less, Swap: trace.Swap}, c.A, c.B)
			}
			native = trace.describe()
			want = goLeanPatternOracle(c.A, c.B).describe()
		}
		if native != want || c.Want != want {
			t.Fatalf("%+v: native %s, Lean %s, oracle %s", c, native, c.Want, want)
		}
		counts[c.Kind]++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	t.Logf("pinned Go and independent oracle rows per helper: %v", counts)
}
