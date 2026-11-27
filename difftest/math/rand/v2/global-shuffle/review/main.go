package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type source struct {
	words []uint64
	pulls int
}

var sourceStop, swapStop = new(int), new(int)

func (s *source) Uint64() uint64 {
	index := s.pulls
	s.pulls++
	if index == len(s.words) {
		panic(sourceStop)
	}
	return s.words[index]
}

func capture(call func()) (status string) {
	defer func() {
		if value := recover(); value != nil {
			switch value {
			case sourceStop:
				status = "source"
			case swapStop:
				status = "swap"
			case "invalid argument to Shuffle":
				status = "negative"
			default:
				if e, ok := value.(runtime.Error); ok {
					if strings.Contains(e.Error(), "nil pointer dereference") {
						status = "nil"
						return
					}
					if strings.Contains(e.Error(), "makeslice: len out of range") {
						status = "negative"
						return
					}
				}
				panic(value)
			}
		}
	}()
	call()
	return "ok"
}

type fixture struct {
	n, mode, failAt int
	words           []uint64
	want            string
}

func native(c fixture) string {
	s := &source{words: c.words}
	var calls []string
	swap := func(i, j int) {
		calls = append(calls, fmt.Sprintf("%d,%d", i, j))
		if c.mode == 1 && len(calls)-1 == c.failAt {
			panic(swapStop)
		}
	}
	if c.mode == 2 {
		swap = nil
	}
	status := capture(func() { rand.New(s).Shuffle(c.n, swap) })
	shuffle := fmt.Sprintf("%s:%d:%s", status, s.pulls, strings.Join(calls, ";"))

	s = &source{words: c.words}
	var p []int
	status = capture(func() { p = rand.New(s).Perm(c.n) })
	result := "none"
	if p != nil {
		fields := make([]string, len(p))
		for i, v := range p {
			fields[i] = strconv.Itoa(v)
		}
		result = "slice:" + strings.Join(fields, ",")
	}
	return fmt.Sprintf("%s|%s:%d:%s", shuffle, status, s.pulls, result)
}

func main() {
	const exact = uint64(0xaaaaaaaaaaaaaaab)
	cases := []fixture{
		{4, 0, 0, []uint64{1, 0, exact, 1}, "ok:4:3,1;2,2;1,1|ok:4:slice:0,3,2,1"},
		{4, 1, 1, []uint64{1, 0, exact, 1}, "swap:3:3,1;2,2|ok:4:slice:0,3,2,1"},
		{4, 0, 0, []uint64{1, 0}, "source:3:3,1|source:3:none"},
		{3, 2, 0, []uint64{0, exact, 1}, "nil:2:|ok:3:slice:0,1,2"},
		{2, 2, 0, nil, "source:1:|source:1:none"},
		{0, 2, 0, nil, "ok:0:|ok:0:slice:"},
		{1, 2, 0, nil, "ok:0:|ok:0:slice:0"},
		{-1, 0, 0, nil, "negative:0:|negative:0:none"},
	}
	mutants := []string{
		"ok:4:1,1;2,2;3,1|ok:4:slice:0,3,2,1",
		"swap:4:3,1;2,2;1,1|ok:4:slice:0,3,2,1",
		"source:2:3,1|source:2:none",
		"nil:0:|ok:3:slice:0,1,2",
		"nil:0:|source:1:none",
		"ok:0:|ok:0:none",
	}
	var wire strings.Builder
	for _, c := range cases {
		if got := native(c); got != c.want {
			log.Fatalf("native %q want %q", got, c.want)
		}
		fmt.Fprintf(&wire, "%d %d %d %d", strconv.IntSize, c.n, c.mode, c.failAt)
		for _, word := range c.words {
			fmt.Fprintf(&wire, " %d", word)
		}
		fmt.Fprintln(&wire)
	}
	rejects := 0
	for i, got := range check.LeanRows("godiff-rand-v2-global-shuffle", len(cases), strings.NewReader(wire.String())) {
		if got != cases[i].want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].want)
		}
		if i < len(mutants) && got != mutants[i] {
			rejects++
		}
	}
	if rejects != len(mutants) {
		log.Fatalf("incorrect rejection count: %d", rejects)
	}
	file, err := os.CreateTemp("", "golean-global-shuffle-review-")
	must(err)
	defer os.Remove(file.Name())
	rows := json.NewEncoder(file)
	for _, c := range cases {
		nativeInput := struct {
			Length, Mode, FailAt int
			Words                []uint64
			Want                 string
		}{c.n, c.mode, c.failAt, c.words, c.want}
		must(rows.Encode(nativeInput))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_RAND_GLOBAL_SHUFFLE_CASES", file.Name()))
	must(check.TestInPackage("math/rand/v2", "TestGoLeanGlobalShuffle",
		"difftest/math/rand/v2/global-word/testdata/bridge_test.go",
		"difftest/math/rand/v2/global-shuffle/testdata/global_shuffle_test.go"))
	fmt.Printf("Shuffle and Perm review: %d paired fixed cases, %d altered observations rejected\n", len(cases), rejects)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
