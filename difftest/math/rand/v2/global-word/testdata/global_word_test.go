package rand_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strings"
	"testing"
)

var globalWordCalls = []func(*rand.Rand, uint64) uint64{
	func(_ *rand.Rand, _ uint64) uint64 { return uint64(rand.Int64()) },
	func(_ *rand.Rand, _ uint64) uint64 { return uint64(rand.Uint32()) },
	func(_ *rand.Rand, _ uint64) uint64 { return rand.Uint64() },
	func(_ *rand.Rand, _ uint64) uint64 { return uint64(rand.Int32()) },
	func(_ *rand.Rand, _ uint64) uint64 { return uint64(rand.Int()) },
	func(_ *rand.Rand, _ uint64) uint64 { return uint64(rand.Uint()) },
	func(_ *rand.Rand, w uint64) uint64 {
		value := rand.Float64()
		want := math.Ldexp(float64(selectBits(w, 0, 53)), -53)
		if math.Float64bits(value) != math.Float64bits(want) {
			panic("incorrect global Float64 bit pattern")
		}
		return uint64(math.Ldexp(value, 53))
	},
	func(_ *rand.Rand, w uint64) uint64 {
		value := rand.Float32()
		want := float32(math.Ldexp(float64(selectBits(w, 32, 24)), -24))
		if math.Float32bits(value) != math.Float32bits(want) {
			panic("incorrect global Float32 bit pattern")
		}
		return uint64(math.Ldexp(float64(value), 24))
	},
}

func TestGoLeanGlobalWord(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_RAND_GLOBAL_WORD_CASES"))
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
		src := &scriptedWord{word: c.Word, panic: c.Mode == 1}
		var source rand.Source = src
		if c.Mode == 2 {
			source = nil
		}
		var results []string
		rand.GoLeanWithGlobalSource(source, func() {
			for _, method := range globalWordCalls {
				prior := src.calls
				status, value := callWord(method, nil, c.Word)
				results = append(results, fmt.Sprintf("%s:%d:%d", status, value, src.calls-prior))
			}
		})
		want := wordOracle(c)
		_, want, _ = strings.Cut(want, "|")
		if got := strings.Join(results, "|"); got != want || got != c.Want {
			t.Fatalf("fixture %+v: native %q, oracle %q", c, got, want)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing fixtures")
	}
	t.Logf("%d direct global-to-source observations per declaration matched Lean and bit oracle", count)
}
