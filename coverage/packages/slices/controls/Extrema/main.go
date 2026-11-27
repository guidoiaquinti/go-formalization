package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}

	positive := []int{8, 3, 12}
	negative := []int{-8, -3, -12}
	if got := slices.Min(positive); got != 3 {
		log.Fatalf("Min(%v) = %d, want 3", positive, got)
	}
	if got := slices.Max(negative); got != -3 {
		log.Fatalf("Max(%v) = %d, want -3", negative, got)
	}
	if minimum := slices.Min(append([]int{0}, positive...)); minimum != 0 {
		log.Fatalf("zero-initialized minimum mutant = %d, want 0", minimum)
	}
	if maximum := slices.Max(append([]int{0}, negative...)); maximum != 0 {
		log.Fatalf("zero-initialized maximum mutant = %d, want 0", maximum)
	}

	type record struct{ key, label int }
	records := []record{{4, 0}, {-2, 1}, {4, 2}, {-2, 3}, {0, 4}}
	compare := func(x, y record) int { return cmp.Compare(x.key, y.key) }
	if got := slices.MinFunc(records, compare); got.label != 1 {
		log.Fatalf("MinFunc returned %+v, want first minimum", got)
	}
	if got := slices.MaxFunc(records, compare); got.label != 0 {
		log.Fatalf("MaxFunc returned %+v, want first maximum", got)
	}
	lastMinimum, lastMaximum := records[0], records[0]
	for _, next := range records[1:] {
		if compare(next, lastMinimum) <= 0 {
			lastMinimum = next
		}
		if compare(next, lastMaximum) >= 0 {
			lastMaximum = next
		}
	}
	if lastMinimum.label != 3 || lastMaximum.label != 2 {
		log.Fatalf("last-tie mutants returned %+v and %+v, want labels 3 and 2", lastMinimum, lastMaximum)
	}

	checkPanic("slices.Min: empty list", func() { slices.Min([]int{}) })
	checkPanic("slices.Max: empty list", func() { slices.Max([]int(nil)) })
	checkPanic("slices.MinFunc: empty list", func() { slices.MinFunc([]record{}, compare) })
	checkPanic("slices.MaxFunc: empty list", func() { slices.MaxFunc([]record(nil), compare) })
	fmt.Printf("%s %s/%s: extrema and exact empty panics passed; zero-seed and last-tie mutants rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func checkPanic(want string, call func()) {
	defer func() {
		if got := recover(); got != want {
			log.Fatalf("empty input panicked with %v, want %q", got, want)
		}
	}()
	call()
}
