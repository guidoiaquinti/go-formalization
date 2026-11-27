package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, A, B int
	State      uint64
	Want       string
}

func main() {
	file, err := os.CreateTemp("", "golean-sort-break-patterns-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	for kind := range 4 {
		compare(fmt.Sprintf("helper=%d boundaries", kind), boundaries(kind), rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				cases := make([]fixture, 5000)
				for i := range cases {
					c := fixture{Kind: kind}
					switch kind {
					case 0:
						c.State = rng.Uint64()
					case 1:
						c.A = int(rng.Uint64())
					default:
						size := int(rng.Uint64() & uint64(math.MaxInt))
						if i%3 == 0 {
							size = rng.Intn(1025)
						}
						room := math.MaxInt - size
						c.A = min(room, int(rng.Uint64()&uint64(math.MaxInt)))
						c.B = c.A + size
					}
					cases[i] = c
				}
				compare(fmt.Sprintf("helper=%d seed=%d block=%d", kind, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_BREAK_PATTERNS_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanBreakPatterns", "difftest/sort/break-patterns/testdata/break_patterns_test.go"))
}

func boundaries(kind int) []fixture {
	var cases []fixture
	switch kind {
	case 0:
		cases = append(cases, fixture{Kind: kind}, fixture{Kind: kind, State: math.MaxUint64})
		for i := range 64 {
			mask := uint64(1) << i
			cases = append(cases, fixture{Kind: kind, State: mask}, fixture{Kind: kind, State: ^mask})
		}
	case 1:
		for a := -1; a <= 8; a++ {
			cases = append(cases, fixture{Kind: kind, A: a})
		}
		for i := range strconv.IntSize - 1 {
			power := 1 << i
			for _, a := range []int{power - 1, power, power + 1, -power - 1, -power, -power + 1} {
				cases = append(cases, fixture{Kind: kind, A: a})
			}
		}
		cases = append(cases, fixture{Kind: kind, A: math.MinInt}, fixture{Kind: kind, A: math.MaxInt})
	default:
		appendSize := func(size int) {
			for _, a := range []int{0, min(19, math.MaxInt-size), math.MaxInt - size} {
				cases = append(cases, fixture{Kind: kind, A: a, B: a + size})
			}
		}
		for size := range 257 {
			appendSize(size)
		}
		for i := 9; i < strconv.IntSize-1; i++ {
			power := 1 << i
			appendSize(power - 1)
			appendSize(power)
			appendSize(power + 1)
		}
		appendSize(math.MaxInt)
	}
	return cases
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		if c.Kind == 0 {
			fmt.Fprintf(&wire, "%d %d %d 0\n", c.Kind, strconv.IntSize, c.State)
		} else {
			fmt.Fprintf(&wire, "%d %d %d %d\n", c.Kind, strconv.IntSize, c.A, c.B)
		}
	}
	for i, result := range check.LeanRows("godiff-sort-break-patterns", len(cases), &wire) {
		c := cases[i]
		c.Want = result
		must(rows.Encode(c))
	}
	fmt.Printf("sort break patterns %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
