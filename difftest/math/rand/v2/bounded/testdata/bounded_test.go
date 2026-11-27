package rand_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"math/bits"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"testing"
)

type boundedFixture struct {
	Bound uint64
	Words []uint64
	Want  string
}

type boundedSource struct {
	words []uint64
	pulls int
}

var exhausted = new(int)

func (s *boundedSource) Uint64() uint64 {
	i := s.pulls
	s.pulls++
	if i >= len(s.words) {
		panic(exhausted)
	}
	return s.words[i]
}

type signedByte int8

type boundedMethod struct {
	label string
	width int
	sign  bool
	call  func(*rand.Rand, uint64) uint64
}

var boundedMethods = []boundedMethod{
	{"", 64, false, rand.GoLeanRandUint64n},
	{"", 32, false, func(r *rand.Rand, n uint64) uint64 { return uint64(rand.GoLeanRandUint32n(r, uint32(n))) }},
	{"Int64N", 64, true, func(r *rand.Rand, n uint64) uint64 { return uint64(r.Int64N(int64(n))) }},
	{"Uint64N", 64, false, func(r *rand.Rand, n uint64) uint64 { return r.Uint64N(n) }},
	{"Int32N", 32, true, func(r *rand.Rand, n uint64) uint64 { return uint64(r.Int32N(int32(n))) }},
	{"Uint32N", 32, false, func(r *rand.Rand, n uint64) uint64 { return uint64(r.Uint32N(uint32(n))) }},
	{"IntN", strconv.IntSize, true, func(r *rand.Rand, n uint64) uint64 { return uint64(r.IntN(int(n))) }},
	{"UintN", strconv.IntSize, false, func(r *rand.Rand, n uint64) uint64 { return uint64(r.UintN(uint(n))) }},
	{"N", 8, true, func(r *rand.Rand, n uint64) uint64 { return uint64(r.N(signedByte(n))) }},
	{"N", 16, false, func(r *rand.Rand, n uint64) uint64 { return uint64(r.N(uint16(n))) }},
}

func nativeBounded(m boundedMethod, c boundedFixture) (wire string) {
	src := &boundedSource{words: c.Words}
	r := rand.New(src)
	defer func() {
		if failure := recover(); failure != nil {
			status := fmt.Sprint(failure)
			if failure == exhausted {
				status = "callback"
			}
			wire = fmt.Sprintf("%s:0:%d", status, src.pulls)
		}
	}()
	result := m.call(r, c.Bound)
	return fmt.Sprintf("ok:%d:%d", result, src.pulls)
}

var two64 = new(big.Int).Lsh(big.NewInt(1), 64)

func oracleBounded(m boundedMethod, c boundedFixture) string {
	mask := ^uint64(0) >> (64 - m.width)
	bound := c.Bound & mask
	if m.label != "" && (bound == 0 || m.sign && bound>>uint(m.width-1) != 0) {
		return "invalid argument to " + m.label + ":0:0"
	}
	if len(c.Words) == 0 {
		return "callback:0:1"
	}
	if bound == 0 {
		return fmt.Sprintf("ok:%d:1", c.Words[0]&(^uint64(0)>>(64-min(strconv.IntSize, m.width))))
	}
	if bits.OnesCount64(bound) == 1 {
		return fmt.Sprintf("ok:%d:1", c.Words[0]%bound)
	}
	n := new(big.Int).SetUint64(bound)
	cut := new(big.Int).Mod(two64, n).Uint64()
	for i, word := range c.Words {
		p := new(big.Int).SetUint64(word)
		p.Mul(p, n)
		if p.Uint64() >= cut {
			p.Rsh(p, 64)
			return fmt.Sprintf("ok:%d:%d", p.Uint64(), i+1)
		}
	}
	return fmt.Sprintf("callback:0:%d", len(c.Words)+1)
}

func TestGoLeanRandBounded(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_RAND_BOUNDED_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		var c boundedFixture
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		results := make([]string, len(boundedMethods))
		for i, m := range boundedMethods {
			got, want := nativeBounded(m, c), oracleBounded(m, c)
			if got != want {
				t.Fatalf("fixture %+v method %d: Go %q, oracle %q", c, i, got, want)
			}
			results[i] = got
		}
		if joined := strings.Join(results, "|"); joined != c.Want {
			t.Fatalf("fixture %+v: Go %q", c, joined)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("missing fixtures")
	}
	t.Logf("%d bounded draws per instantiation matched Go, Lean and big integer oracle", count)
}
