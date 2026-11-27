package rand_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type globalShuffleFixture struct {
	Length, Mode, FailAt int
	Words                []uint64
	Want                 string
}

type shuffleSource struct {
	words []uint64
	pulls int
}

var shuffleExhausted, shuffleCallback = new(int), new(int)

func (s *shuffleSource) Uint64() uint64 {
	i := s.pulls
	s.pulls++
	if i >= len(s.words) {
		panic(shuffleExhausted)
	}
	return s.words[i]
}

func observeShuffle(call func()) (label string) {
	defer func() {
		switch p := recover(); p {
		case nil:
			label = "ok"
		case shuffleExhausted:
			label = "source"
		case shuffleCallback:
			label = "swap"
		case "invalid argument to Shuffle":
			label = "negative"
		default:
			if err, ok := p.(runtime.Error); ok {
				if strings.Contains(err.Error(), "makeslice: len out of range") {
					label = "negative"
					return
				}
				if strings.Contains(err.Error(), "nil pointer dereference") {
					label = "nil"
					return
				}
			}
			panic(p)
		}
	}()
	call()
	return
}

func observeGlobalShuffle(c globalShuffleFixture, inject bool) string {
	withSource := func(s *shuffleSource, call func(*rand.Rand)) {
		if inject {
			rand.GoLeanWithGlobalSource(s, func() { call(nil) })
		} else {
			call(rand.New(s))
		}
	}
	first, second := &shuffleSource{words: c.Words}, &shuffleSource{words: c.Words}
	var trace []string
	var swap func(int, int)
	if c.Mode != 2 {
		swap = func(i, j int) {
			trace = append(trace, fmt.Sprintf("%d,%d", i, j))
			if c.Mode == 1 && len(trace)-1 == c.FailAt {
				panic(shuffleCallback)
			}
		}
	}
	status := observeShuffle(func() {
		withSource(first, func(r *rand.Rand) {
			if inject {
				rand.Shuffle(c.Length, swap)
			} else {
				r.Shuffle(c.Length, swap)
			}
		})
	})
	shuffled := fmt.Sprintf("%s:%d:%s", status, first.pulls, strings.Join(trace, ";"))
	var indices []int
	status = observeShuffle(func() {
		withSource(second, func(r *rand.Rand) {
			if inject {
				indices = rand.Perm(c.Length)
			} else {
				indices = r.Perm(c.Length)
			}
		})
	})
	order := "none"
	if indices != nil {
		fields := make([]string, len(indices))
		for i, value := range indices {
			fields[i] = strconv.Itoa(value)
		}
		order = "slice:" + strings.Join(fields, ",")
	}
	return fmt.Sprintf("%s|%s:%d:%s", shuffled, status, second.pulls, order)
}

func TestGoLeanGlobalShuffle(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_RAND_GLOBAL_SHUFFLE_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	rows := 0
	for scanner.Scan() {
		var c globalShuffleFixture
		if err := json.Unmarshal(scanner.Bytes(), &c); err != nil {
			t.Fatal(err)
		}
		got, want := observeGlobalShuffle(c, true), observeGlobalShuffle(c, false)
		if got != want || got != c.Want {
			t.Fatalf("fixture %+v: native %q, Rand reference %q", c, got, want)
		}
		rows++
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	if rows == 0 {
		t.Fatal("missing cases")
	}
	t.Logf("%d package-level shuffle and permutation observations matched Go delegate and Lean", rows)
}
