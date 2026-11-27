package main

import (
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

	backing := []int{81, 1, 2, 3, 4, 94}
	view := backing[1:5]
	got := slices.Delete(view, 1, 3)
	if !slices.Equal(got, []int{1, 4}) || cap(got) != cap(view) || &got[0] != &view[0] ||
		!slices.Equal(backing, []int{81, 1, 4, 0, 0, 94}) {
		log.Fatalf("Delete returned %v with backing %v", got, backing)
	}

	uncleared := []int{1, 2, 3, 4}
	noClear := append(uncleared[:1], uncleared[3:]...)
	if !slices.Equal(noClear, got) || slices.Equal(uncleared, view) {
		log.Fatalf("no-clear mutant returned %v with backing %v", noClear, uncleared)
	}

	before := slices.Clone(backing)
	expectBoundsPanic(func() { slices.Delete(view, 0, 5) })
	if !slices.Equal(backing, before) {
		log.Fatalf("invalid Delete changed backing to %v, want %v", backing, before)
	}
	withinCapacity := slices.Clone(before)
	short := withinCapacity[1:5]
	expanded := short[:cap(short)]
	mutant := slices.Delete(expanded, 0, 5)
	if len(mutant) != 0 || slices.Equal(withinCapacity, before) {
		log.Fatal("capacity-bounds mutant did not accept and mutate past-length input")
	}
	if got := slices.Delete([]int(nil), 0, 0); got != nil {
		log.Fatalf("empty deletion changed a nil slice to %v", got)
	}
	fmt.Printf("%s %s/%s: Delete controls passed, no-clear and capacity-bounds mutants rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func expectBoundsPanic(call func()) {
	defer func() {
		if p := recover(); p == nil {
			log.Fatal("expected a bounds panic")
		} else if _, ok := p.(runtime.Error); !ok {
			log.Fatalf("unexpected panic %T: %v", p, p)
		}
	}()
	call()
}
