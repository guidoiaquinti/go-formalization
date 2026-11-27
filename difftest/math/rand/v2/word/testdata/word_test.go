package rand_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type randWordFixture struct {
	Mode int
	Word uint64
	Want string
}

var randWordPanic = new(int)

type scriptedWord struct {
	word  uint64
	panic bool
	calls int
}

func (s *scriptedWord) Uint64() uint64 {
	s.calls++
	if s.panic {
		panic(randWordPanic)
	}
	return s.word
}

var wordCalls = []func(*rand.Rand, uint64) uint64{
	func(r *rand.Rand, _ uint64) uint64 { return uint64(r.Int64()) },
	func(r *rand.Rand, _ uint64) uint64 { return uint64(r.Uint32()) },
	func(r *rand.Rand, _ uint64) uint64 { return r.Uint64() },
	func(r *rand.Rand, _ uint64) uint64 { return uint64(r.Int32()) },
	func(r *rand.Rand, _ uint64) uint64 { return uint64(r.Int()) },
	func(r *rand.Rand, _ uint64) uint64 { return uint64(r.Uint()) },
	func(r *rand.Rand, w uint64) uint64 {
		value := r.Float64()
		want := math.Ldexp(float64(selectBits(w, 0, 53)), -53)
		if math.Float64bits(value) != math.Float64bits(want) {
			panic("incorrect Float64 bit pattern")
		}
		return uint64(math.Ldexp(value, 53))
	},
	func(r *rand.Rand, w uint64) uint64 {
		value := r.Float32()
		want := float32(math.Ldexp(float64(selectBits(w, 32, 24)), -24))
		if math.Float32bits(value) != math.Float32bits(want) {
			panic("incorrect Float32 bit pattern")
		}
		return uint64(math.Ldexp(float64(value), 24))
	},
}

func callWord(method func(*rand.Rand, uint64) uint64, r *rand.Rand, word uint64) (state string, value uint64) {
	defer func() {
		if p := recover(); p != nil {
			if p == randWordPanic {
				state = "callback"
			} else if e, ok := p.(runtime.Error); ok && strings.Contains(e.Error(), "nil pointer dereference") {
				state = "nil"
			} else {
				panic(p)
			}
		}
	}()
	return "ok", method(r, word)
}

func wordNative(t *testing.T, c randWordFixture) string {
	s := &scriptedWord{word: c.Word, panic: c.Mode == 1}
	var iface rand.Source = s
	installed := "source:0"
	if c.Mode == 2 {
		iface = nil
		installed = "nil:0"
	}
	r := rand.New(iface)
	if r == nil || r == rand.New(iface) || rand.GoLeanRandSource(r) != iface || s.calls != 0 {
		t.Fatal("New did not construct an unused Rand retaining the input source")
	}
	if c.Mode == 3 {
		r = nil
	}
	results := []string{installed}
	for _, f := range wordCalls {
		previous := s.calls
		state, value := callWord(f, r, c.Word)
		results = append(results, fmt.Sprintf("%s:%d:%d", state, value, s.calls-previous))
	}
	return strings.Join(results, "|")
}

// Select a bit interval by weighting individual input digits.
func selectBits(word uint64, start, width int) uint64 {
	var result uint64
	for i := range width {
		mask := uint64(1) << (start + i)
		if word&mask != 0 {
			result += uint64(1) << i
		}
	}
	return result
}

func wordOracle(c randWordFixture) string {
	switch c.Mode {
	case 1:
		return "source:0" + strings.Repeat("|callback:0:1", len(wordCalls))
	case 2:
		return "nil:0" + strings.Repeat("|nil:0:0", len(wordCalls))
	case 3:
		return "source:0" + strings.Repeat("|nil:0:0", len(wordCalls))
	}
	w := c.Word
	fields := []string{"source:0"}
	for _, interval := range [][2]int{{0, 63}, {32, 32}, {0, 64}, {33, 31}, {0, strconv.IntSize - 1},
		{0, strconv.IntSize}, {0, 53}, {32, 24}} {
		fields = append(fields, fmt.Sprintf("ok:%d:1", selectBits(w, interval[0], interval[1])))
	}
	return strings.Join(fields, "|")
}

func TestGoLeanRandWord(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_RAND_WORD_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		var c randWordFixture
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		got, want := wordNative(t, c), wordOracle(c)
		if got != c.Want || got != want {
			t.Fatalf("fixture %+v: Go %q, oracle %q", c, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing fixtures")
	}
	t.Logf("%d native observations per declaration matched Lean and bit interval oracle", count)
}
