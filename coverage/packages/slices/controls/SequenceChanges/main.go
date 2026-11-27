package main

import (
	"fmt"
	"log"
	"math"
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
	checkReverse()
	checkRepeat()
	checkConcat()
	checkCompact()
	fmt.Printf("%s %s/%s: reverse, repeat, concat and compact controls passed and rejected their mutants\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func checkReverse() {
	backing := []int{81, 2, 3, 4, 94}
	slices.Reverse(backing[1:4])
	if !slices.Equal(backing, []int{81, 4, 3, 2, 94}) {
		log.Fatalf("reverse changed backing to %v", backing)
	}
	unchanged := []int{81, 2, 3, 4, 94}
	copyMutant := slices.Clone(unchanged[1:4])
	slices.Reverse(copyMutant)
	if !slices.Equal(copyMutant, backing[1:4]) || slices.Equal(unchanged, backing) {
		log.Fatal("reverse-copy mutant did not leave the input backing unchanged")
	}
}

func checkRepeat() {
	values := []int{2, 3}
	want := []int{2, 3, 2, 3, 2, 3}
	if got := slices.Repeat(values, 3); !slices.Equal(got, want) || len(got) != cap(got) {
		log.Fatalf("Repeat returned %v with capacity %d", got, cap(got))
	}
	empty := slices.Repeat(values, 0)
	if empty == nil || len(empty) != 0 || cap(empty) != 0 {
		log.Fatalf("zero repetitions returned %v, want nonnil empty", empty)
	}
	repeatOrNil := func(values []int, count int) []int {
		if count == 0 {
			return nil
		}
		return slices.Repeat(values, count)
	}
	if mutant := repeatOrNil(values, 0); !slices.Equal(mutant, empty) || (mutant == nil) == (empty == nil) {
		log.Fatal("nil-result mutant did not isolate the nilness difference")
	}
	checkPanic("cannot be negative", func() { slices.Repeat([]int{}, -1) })
	checkPanic("the result of (len(x) * count) overflows", func() { slices.Repeat(values, math.MaxInt) })
}

func checkConcat() {
	backing := []int{2, 3, 0}
	first, second := backing[:2], []int{4}
	got := slices.Concat(first, nil, second)
	if !slices.Equal(got, []int{2, 3, 4}) || backing[2] != 0 {
		log.Fatalf("Concat returned %v with source backing %v", got, backing)
	}
	got[0], got[2] = 19, 20
	if first[0] != 2 || second[0] != 4 {
		log.Fatal("Concat result aliases an input")
	}
	if empty := slices.Concat([]int{}, nil); empty != nil {
		log.Fatalf("empty Concat returned %v, want nil", empty)
	}
	reuseMutant := append(first, second...)
	reuseMutant[0] = 19
	if backing[0] != 19 || backing[2] != 4 {
		log.Fatal("append-in-place mutant did not alter the first input backing")
	}
}

func checkCompact() {
	backing := []int{81, 1, 1, 2, 1, 1, 94}
	view := backing[1:6]
	got := slices.Compact(view)
	if !slices.Equal(got, []int{1, 2, 1}) || &got[0] != &view[0] ||
		!slices.Equal(backing, []int{81, 1, 2, 1, 0, 0, 94}) {
		log.Fatalf("Compact returned %v with backing %v", got, backing)
	}
	uncleared := []int{1, 1, 2, 1, 1}
	next := 1
	for i := 1; i < len(uncleared); i++ {
		if uncleared[i] != uncleared[i-1] {
			uncleared[next] = uncleared[i]
			next++
		}
	}
	if !slices.Equal(uncleared[:next], got) || slices.Equal(uncleared, view) {
		log.Fatalf("no-clear mutant returned prefix %v and backing %v", uncleared[:next], uncleared)
	}
}

func checkPanic(want string, call func()) {
	defer func() {
		if got := recover(); got != want {
			log.Fatalf("panic = %v, want %q", got, want)
		}
	}()
	call()
}
