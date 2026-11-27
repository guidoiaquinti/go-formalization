package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"slices"
	"strings"
)

var marker = new(int)

func main() {
	pin, err := os.ReadFile(".go-version")
	if err != nil {
		log.Fatal(err)
	}
	if want := "go" + strings.TrimSpace(string(pin)); runtime.Version() != want {
		log.Fatalf("Go version %s, want %s", runtime.Version(), want)
	}

	input := []int{81, 1, 2, 3, 13, 6, 94}
	backing := slices.Clone(input)
	view := backing[1:6]
	var trace []int
	predicate := func(x int) bool {
		trace = append(trace, x)
		if x == 13 {
			panic(marker)
		}
		return x%2 == 0
	}
	catchMarker(func() { slices.DeleteFunc(view, predicate) })
	if !slices.Equal(trace, []int{1, 2, 3, 13}) ||
		!slices.Equal(backing, []int{81, 1, 3, 3, 13, 6, 94}) {
		log.Fatalf("panicking DeleteFunc produced trace %v and backing %v", trace, backing)
	}

	mutantBacking := slices.Clone(input)
	mutantView := mutantBacking[1:6]
	trace = nil
	catchMarker(func() {
		scratch := slices.Clone(mutantView)
		slices.DeleteFunc(scratch, predicate)
		copy(mutantView, scratch)
	})
	if !slices.Equal(trace, []int{1, 2, 3, 13}) || !slices.Equal(mutantBacking, input) ||
		slices.Equal(mutantBacking, backing) {
		log.Fatalf("buffered-write mutant produced trace %v and backing %v", trace, mutantBacking)
	}

	backing = slices.Clone(input)
	view = backing[1:6]
	got := slices.DeleteFunc(view, func(x int) bool { return x%2 == 0 })
	if !slices.Equal(got, []int{1, 3, 13}) || &got[0] != &view[0] || cap(got) != cap(view) ||
		!slices.Equal(backing, []int{81, 1, 3, 13, 0, 0, 94}) {
		log.Fatalf("successful DeleteFunc returned %v with backing %v", got, backing)
	}
	fmt.Printf("%s %s/%s: DeleteFunc success and panic controls passed, buffered-write mutant rejected\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func catchMarker(call func()) {
	defer func() {
		if got := recover(); got != marker {
			log.Fatalf("callback panic = %v, want original marker", got)
		}
	}()
	call()
}
