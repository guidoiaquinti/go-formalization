package main

import (
	"bytes"
	"cmp"
	"fmt"
	"iter"
	"math/rand"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	words                             []int64
	tag, producer, mode, failAt, mask int
}

func main() {
	edges := []fixture{{[]int64{513, 256, 512}, 1, 1, 1, -1, 0}, {[]int64{513, 256, 512}, 1, 2, 1, -1, 0}}
	for _, size := range []int{0, 1, 2, 3, 19, 20, 21, 39, 40, 41, 79, 80, 81} {
		for pattern := range 4 {
			words := make([]int64, size)
			for i := range words {
				keys := [...]int{i / 2, -i / 2, i%5 - 2, 7}
				words[i] = int64(keys[pattern])*256 + int64(size-i)
			}
			for producer := range 3 {
				for mode := range 3 {
					for _, stop := range []int{-1, 0, 1, 3, 19, 20, 41} {
						edges = append(edges, fixture{words, 1, producer, mode, stop, 2730})
					}
				}
			}
		}
	}
	for producer := range 3 {
		for mode := range 3 {
			edges = append(edges, fixture{producer: producer, mode: mode, failAt: -1})
		}
	}
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		samples := make([]fixture, 20_000)
		for i := range samples {
			words := make([]int64, rng.Intn(60))
			for j := range words {
				words[j] = int64(rng.Intn(13)-6)*256 + int64(j)
			}
			samples[i] = fixture{words, 1, rng.Intn(3), rng.Intn(3), rng.Intn(80) - 15, rng.Intn(4096)}
		}
		compare(fmt.Sprintf("seed=%d", seed), samples)
	}
}

func compare(label string, cases []fixture) {
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, c.tag, " ", c.producer, " ", c.mode, " ", c.failAt, " ", c.mask)
		for _, word := range c.words {
			fmt.Fprint(&rows, " ", word)
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-slices-public-stable", len(cases), &rows) {
		c := cases[i]
		if want := observe(c); got != want {
			panic(fmt.Sprintf("%s row %d %+v: Lean %q Go %q", label, i, c, got, want))
		}
	}
	fmt.Printf("public stable sort %s: %d matched across both pinned entries\n", label, len(cases))
}

func checkStable(input, output []int64) {
	counts := make(map[int64][]int64)
	for _, word := range input {
		key := word >> 8
		counts[key] = append(counts[key], word)
	}
	var previous int64
	for i, word := range output {
		key := word >> 8
		if i > 0 && previous > key || len(counts[key]) == 0 || counts[key][0] != word {
			panic("result changed stable key order")
		}
		previous = key
		counts[key] = counts[key][1:]
	}
}

type producerPanic int
type comparePanic int

func observe(c fixture) string {
	values := slices.Clone(c.words)
	var results []string
	for kind := range 2 {
		input := slices.Clone(values)
		if c.tag == 0 {
			input = nil
		}
		var collected []int64
		var calls []string
		acknowledged := 0
		finished := false
		var seq iter.Seq[int64]
		if c.producer != 0 {
			seq = func(yield func(int64) bool) {
				for _, value := range values {
					if !yield(value) {
						panic("collection stopped before producer completed")
					}
					acknowledged++
				}
				if c.producer == 2 {
					panic(producerPanic(7))
				}
				finished = true
			}
		}
		var compare func(int64, int64) int
		if c.mode != 0 {
			compare = func(x, y int64) int {
				if kind == 1 && !finished {
					panic("comparison preceded iterator completion")
				}
				call := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", x, y))
				if call == c.failAt {
					panic(comparePanic(call))
				}
				if c.mode == 1 {
					return cmp.Compare(x>>8, y>>8)
				}
				if c.mask&(1<<(call%12)) != 0 {
					return -1
				}
				return 0
			}
		}
		status := capture(kind == 1 && seq == nil, func() {
			switch kind {
			case 0:
				slices.SortStableFunc(input, compare)
			case 1:
				collected = slices.SortedStableFunc(seq, compare)
			}
		})
		if kind == 1 {
			input = collected
		}
		if status == "ok" || kind == 0 {
			frequencies := make(map[int64]int)
			for _, word := range c.words {
				frequencies[word]++
			}
			for _, x := range input {
				frequencies[x]--
			}
			for word, n := range frequencies {
				if n != 0 {
					panic(fmt.Sprintf("method %d changed frequency of %d", kind, word))
				}
			}
			if status == "ok" && c.mode == 1 {
				checkStable(values, input)
			}
		}
		data := "nil"
		if input != nil {
			words := make([]string, len(input))
			for i, value := range input {
				words[i] = strconv.FormatInt(value, 10)
			}
			data = "slice:" + strings.Join(words, ",")
		}
		results = append(results, fmt.Sprintf("%s/%s/%d/%s", status, data, acknowledged, strings.Join(calls, ";")))
	}
	return strings.Join(results, "|")
}

func capture(nilSequence bool, call func()) (status string) {
	defer func() {
		switch failure := recover().(type) {
		case nil:
		case producerPanic:
			status = fmt.Sprintf("producer:%d", failure)
		case comparePanic:
			status = fmt.Sprintf("compare:%d", failure)
		case runtime.Error:
			if failure.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			status = "nilCompare"
			if nilSequence {
				status = "nilSequence"
			}
		default:
			panic(failure)
		}
	}()
	call()
	return "ok"
}
