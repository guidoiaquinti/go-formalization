package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	backing                     []int64
	length, size, first, second int
	write                       bool
	value                       int64
}

func main() {
	var edges []fixture
	for n := range 12 {
		for _, size := range []int{math.MinInt, -1, 0, 1, 2, 3, 7, 16, math.MaxInt} {
			for first := range 5 {
				for second := range 5 {
					for _, write := range []bool{true, false} {
						backing := make([]int64, n+2)
						for i := range backing {
							backing[i] = int64(i + 1)
						}
						edges = append(edges, fixture{backing, n, size, first, second, write, 8})
					}
				}
			}
		}
	}
	edges = append(edges, fixture{nil, 0, 1, 5, 5, false, 0}, fixture{[]int64{}, 0, 1, 5, 5, false, 0})
	compare("boundaries", edges)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		samples := make([]fixture, 20_000)
		for i := range samples {
			n := rng.Intn(45)
			backing := make([]int64, n+rng.Intn(5))
			for j := range backing {
				backing[j] = int64(rng.Uint64())
			}
			samples[i] = fixture{backing, n, rng.Intn(n+15) - 3, rng.Intn(16), rng.Intn(16), i%2 == 0, int64(rng.Uint64())}
		}
		compare(fmt.Sprintf("seed=%d", seed), samples)
	}
}

func compare(label string, cases []fixture) {
	var rows bytes.Buffer
	for _, c := range cases {
		write := 0
		if c.write {
			write = 1
		}
		fmt.Fprint(&rows, c.length, " ", c.size, " ", c.first, " ", c.second, " ", write, " ", c.value)
		for _, value := range c.backing {
			fmt.Fprint(&rows, " ", value)
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-slices-chunk", len(cases), &rows) {
		c := cases[i]
		want := native(c)
		if independent := expected(c); independent != want || got != want {
			panic(fmt.Sprintf("%s row %d %+v: spec %q Go %q Lean %q", label, i, c, independent, want, got))
		}
	}
	fmt.Printf("Chunk %s: %d pinned inputs agree\n", label, len(cases))
}

func expected(c fixture) string {
	if c.size < 1 {
		return "cannot be less than 1"
	}
	backing := slices.Clone(c.backing)
	run := func(budget int, write bool) string {
		var parts []string
		for end := 0; end < c.length && len(parts) < budget; {
			start := end
			end += min(c.size, c.length-end)
			parts = append(parts, fmt.Sprintf("%d:%s", end-start, cells(backing[start:end])))
			if write {
				backing[start] = c.value
			}
		}
		return strings.Join(parts, ";") + "/" + cells(backing)
	}
	left := run(c.first, c.write)
	return left + "|" + run(c.second, false)
}

func native(c fixture) (report string) {
	backing := slices.Clone(c.backing)
	defer func() {
		if failure := recover(); failure != nil {
			if failure != "cannot be less than 1" {
				panic(failure)
			}
			report = "cannot be less than 1"
		}
	}()
	seq := slices.Chunk(backing[:c.length], c.size)
	run := func(budget int, write bool) string {
		var parts []string
		index := 0
		if budget != 0 {
			seq(func(part []int64) bool {
				if len(part) == 0 || &part[0] != &backing[index] {
					panic("chunk is not the original alias")
				}
				parts = append(parts, fmt.Sprintf("%d:%s", cap(part), cells(part)))
				saved := slices.Clone(backing)
				extended := append(part, math.MinInt64)
				extended[0] ^= 1
				if !slices.Equal(backing, saved) || &extended[0] == &part[0] {
					panic("append leaked into source backing")
				}
				if write {
					part[0] = c.value
				}
				index += len(part)
				return len(parts) < budget
			})
		}
		return strings.Join(parts, ";") + "/" + cells(backing)
	}
	left := run(c.first, c.write)
	return left + "|" + run(c.second, false)
}

func cells(values []int64) string {
	fields := make([]string, len(values))
	for i, value := range values {
		fields[i] = strconv.FormatInt(value, 10)
	}
	return strings.Join(fields, ",")
}
