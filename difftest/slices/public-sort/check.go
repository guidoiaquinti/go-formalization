package main

import (
	"bytes"
	"cmp"
	"fmt"
	"iter"
	"math"
	"math/rand"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/binary64"
	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	words                                   []uint64
	kind, tag, producer, mode, failAt, mask int
}

func main() {
	boundaries := binary64.Boundaries()
	edges := []fixture{
		{words: []uint64{3, 2, 1}, tag: 1, producer: 1, mode: 1, failAt: 1},
		{words: []uint64{2, 1}, tag: 1, producer: 1, failAt: -1},
	}
	for _, size := range []int{0, 1, 2, 3, 11, 12, 13, 24, 48, 49, 50, 51, 67} {
		for pattern := range 4 {
			words := make([]uint64, size)
			for i := range words {
				switch pattern {
				case 0:
					words[i] = uint64(i)
				case 1:
					words[i] = uint64(-i)
				case 2:
					words[i] = boundaries[i%len(boundaries)]
				case 3:
					words[i] = boundaries[(size-i)%len(boundaries)]
				}
			}
			for kind := range 2 {
				for producer := range 3 {
					for mode := range 3 {
						for _, stop := range []int{-1, 0, 1, 3, 11, 12, 20} {
							edges = append(edges, fixture{words: words, kind: kind, tag: 1, producer: producer,
								mode: mode, failAt: stop, mask: 2730})
						}
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
			samples[i] = sample(rng, boundaries)
		}
		compare(fmt.Sprintf("seed=%d", seed), samples)
	}
}

func sample(rng *rand.Rand, boundaries []uint64) fixture {
	words := make([]uint64, rng.Intn(71))
	for i := range words {
		words[i] = rng.Uint64()
		if i%3 != 0 {
			words[i] = boundaries[rng.Intn(len(boundaries))]
		}
	}
	return fixture{words: words, kind: rng.Intn(2), tag: 1, producer: rng.Intn(3),
		mode: rng.Intn(3), failAt: rng.Intn(80) - 15, mask: rng.Intn(4096)}
}

func compare(label string, cases []fixture) {
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, strconv.IntSize, " ", c.kind, " ", c.tag, " ", c.producer, " ", c.mode, " ", c.failAt, " ", c.mask)
		for _, word := range c.words {
			fmt.Fprint(&rows, " ", word)
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-slices-public-sort", len(cases), &rows) {
		c := cases[i]
		var want string
		if c.kind == 1 {
			want = observe(c, math.Float64frombits, math.Float64bits)
		} else {
			want = observe(c, func(x uint64) int64 { return int64(x) }, func(x int64) uint64 { return uint64(x) })
		}
		if got != want {
			panic(fmt.Sprintf("%s row %d %+v: Lean %q Go %q", label, i, c, got, want))
		}
	}
	fmt.Printf("public sort %s: %d matched across four pinned entries\n", label, len(cases))
}

type producerPanic int
type comparePanic int

func observe[T cmp.Ordered](c fixture, decode func(uint64) T, encode func(T) uint64) string {
	values := make([]T, len(c.words))
	for i, word := range c.words {
		values[i] = decode(word)
	}
	var results []string
	for kind := range 4 {
		input := slices.Clone(values)
		if c.tag == 0 {
			input = nil
		}
		var collected []T
		var calls []string
		acknowledged := 0
		finished := false
		var seq iter.Seq[T]
		if c.producer != 0 {
			seq = func(yield func(T) bool) {
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
		var compare func(T, T) int
		if c.mode != 0 {
			compare = func(x, y T) int {
				if kind >= 2 && !finished {
					panic("comparison preceded iterator completion")
				}
				call := len(calls)
				calls = append(calls, fmt.Sprintf("%d,%d", encode(x), encode(y)))
				if call == c.failAt {
					panic(comparePanic(call))
				}
				if c.mode == 1 {
					return cmp.Compare(x, y)
				}
				if c.mask&(1<<(call%12)) != 0 {
					return -1
				}
				return 0
			}
		}
		status := capture(kind >= 2 && seq == nil, func() {
			switch kind {
			case 0:
				slices.Sort(input)
			case 1:
				slices.SortFunc(input, compare)
			case 2:
				collected = slices.Sorted(seq)
			case 3:
				collected = slices.SortedFunc(seq, compare)
			}
		})
		if kind >= 2 {
			input = collected
		}
		if status == "ok" || kind < 2 {
			frequencies := make(map[uint64]int)
			for _, word := range c.words {
				frequencies[word]++
			}
			for _, x := range input {
				frequencies[encode(x)]--
			}
			for word, n := range frequencies {
				if n != 0 {
					panic(fmt.Sprintf("method %d changed frequency of %d", kind, word))
				}
			}
			if status == "ok" && (kind%2 == 0 || c.mode == 1) && !slices.IsSorted(input) {
				panic("ordered result is not sorted")
			}
		}
		data := "nil"
		if input != nil {
			words := make([]string, len(input))
			for i, value := range input {
				words[i] = strconv.FormatUint(encode(value), 10)
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
