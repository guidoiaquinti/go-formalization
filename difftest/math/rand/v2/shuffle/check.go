package main

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	randv2 "math/rand/v2"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	length, mode, failAt int
	words                []uint64
}

type scripted struct {
	words []uint64
	pulls int
}

var sourcePanic, swapPanic = new(int), new(int)

func (s *scripted) Uint64() uint64 {
	i := s.pulls
	s.pulls++
	if i >= len(s.words) {
		panic(sourcePanic)
	}
	return s.words[i]
}

func main() {
	var fixed []fixture
	ones, zeros := make([]uint64, 128), make([]uint64, 128)
	for i := range ones {
		ones[i] = ^uint64(0)
	}
	for n := -3; n <= 36; n++ {
		for mode := range 3 {
			for failAt := range 8 {
				for _, words := range [][]uint64{nil, {0}, {0, 0, 1}, {^uint64(0)}, ones, zeros} {
					fixed = append(fixed, fixture{n, mode, failAt, words})
				}
			}
		}
	}
	check.CompareSamples("godiff-rand-v2-shuffle", fixed, sample, encode, native, reference)
}

func sample(rng *rand.Rand, i int) fixture {
	length := rng.Intn(55) - 4
	mode := rng.Intn(3)
	if i%4 != 0 {
		mode = 0
	}
	words := make([]uint64, rng.Intn(70))
	for j := range words {
		if j%7 != 0 {
			words[j] = rng.Uint64()
		}
	}
	return fixture{length, mode, rng.Intn(9), words}
}

func encode(c fixture) string {
	var out strings.Builder
	fmt.Fprintf(&out, "%d %d %d %d", strconv.IntSize, c.length, c.mode, c.failAt)
	for _, w := range c.words {
		fmt.Fprintf(&out, " %d", w)
	}
	return out.String()
}

func classify(value any) string {
	switch value {
	case nil:
		return "ok"
	case sourcePanic:
		return "source"
	case swapPanic:
		return "swap"
	case "invalid argument to Shuffle":
		return "negative"
	}
	if e, ok := value.(runtime.Error); ok {
		switch {
		case strings.Contains(e.Error(), "nil pointer dereference"):
			return "nil"
		case strings.Contains(e.Error(), "makeslice: len out of range"):
			return "negative"
		}
	}
	panic(value)
}

func native(c fixture) string {
	s := &scripted{words: c.words}
	var calls [][2]int
	var swap func(i, j int)
	if c.mode != 2 {
		swap = func(i, j int) {
			calls = append(calls, [2]int{i, j})
			if c.mode == 1 && len(calls)-1 == c.failAt {
				panic(swapPanic)
			}
		}
	}
	status := func() (status string) {
		defer func() { status = classify(recover()) }()
		randv2.New(s).Shuffle(c.length, swap)
		return
	}()
	shuffle := shuffleWire(status, s.pulls, calls)
	p := &scripted{words: c.words}
	var perm []int
	status = func() (status string) {
		defer func() { status = classify(recover()) }()
		perm = randv2.New(p).Perm(c.length)
		return
	}()
	return shuffle + "|" + permWire(status, p.pulls, perm)
}

var radix = new(big.Int).Lsh(big.NewInt(1), 64)

func nextIndex(n int, words []uint64) (int, int, bool) {
	if bits.OnesCount(uint(n)) == 1 && len(words) != 0 {
		return int(words[0] % uint64(n)), 1, true
	}
	bound := big.NewInt(int64(n))
	minimum := new(big.Int).Mod(radix, bound).Uint64()
	for index, word := range words {
		scaled := new(big.Int).SetUint64(word)
		scaled.Mul(scaled, bound)
		if scaled.Uint64() >= minimum {
			return int(scaled.Rsh(scaled, 64).Uint64()), index + 1, true
		}
	}
	return 0, len(words) + 1, false
}

func oracleTrace(c fixture) (string, int, [][2]int) {
	if c.length < 0 {
		return "negative", 0, nil
	}
	pulls := 0
	var calls [][2]int
	for i := c.length; i > 1; i-- {
		index, took, accepted := nextIndex(i, c.words[pulls:])
		pulls += took
		if !accepted {
			return "source", pulls, calls
		}
		if c.mode == 2 {
			return "nil", pulls, calls
		}
		calls = append(calls, [2]int{i - 1, index})
		if c.mode == 1 && len(calls)-1 == c.failAt {
			return "swap", pulls, calls
		}
	}
	return "ok", pulls, calls
}

func reference(c fixture) string {
	status, pulls, trace := oracleTrace(c)
	shuffle := shuffleWire(status, pulls, trace)
	c.mode = 0
	status, pulls, trace = oracleTrace(c)
	var values []int
	if status == "ok" {
		values = make([]int, c.length)
		// Replay each final position backwards to its original index.
		for i := range values {
			source := i
			for j := len(trace) - 1; j >= 0; j-- {
				a, b := trace[j][0], trace[j][1]
				switch source {
				case a:
					source = b
				case b:
					source = a
				}
			}
			values[i] = source
		}
	}
	return shuffle + "|" + permWire(status, pulls, values)
}

func shuffleWire(status string, pulls int, trace [][2]int) string {
	fields := make([]string, len(trace))
	for i, call := range trace {
		fields[i] = fmt.Sprintf("%d,%d", call[0], call[1])
	}
	return fmt.Sprintf("%s:%d:%s", status, pulls, strings.Join(fields, ";"))
}

func permWire(status string, pulls int, values []int) string {
	if values == nil {
		return fmt.Sprintf("%s:%d:none", status, pulls)
	}
	fields := make([]string, len(values))
	for i, v := range values {
		fields[i] = strconv.Itoa(v)
	}
	return fmt.Sprintf("%s:%d:slice:%s", status, pulls, strings.Join(fields, ","))
}
