package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Variant, A, B, N, Width int
	Trace                         string
}

func main() {
	file, err := os.CreateTemp("", "golean-sort-range-moves-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	for kind := range 3 {
		var fixed []fixture
		for width := 0; width <= 12; width++ {
			for a := 0; a <= width; a++ {
				for b := 0; b <= width; b++ {
					switch kind {
					case 0:
						if a <= b {
							fixed = add(fixed, kind, a, b, 0, width)
						}
					case 1:
						for n := 0; n <= width-max(a, b); n++ {
							fixed = add(fixed, kind, a, b, n, width)
						}
					case 2:
						if a == b {
							fixed = add(fixed, kind, a, b, b, width)
						}
						if a < b {
							for end := b + 1; end <= width; end++ {
								fixed = add(fixed, kind, a, b, end, width)
							}
						}
					}
				}
			}
		}
		compare(fmt.Sprintf("helper=%d fixed", kind), fixed, rows)
		for _, seed := range []int64{17, 2718, 65537} {
			rng := rand.New(rand.NewSource(seed))
			for block := range 4 {
				var cases []fixture
				for range 5000 {
					width := 2 + rng.Intn(95)
					a, b := rng.Intn(width+1), rng.Intn(width+1)
					n := 0
					switch kind {
					case 0:
						a, b = min(a, b), max(a, b)
					case 1:
						n = rng.Intn(width - max(a, b) + 1)
					case 2:
						a = rng.Intn(width - 1)
						b = a + 1 + rng.Intn(width-a-1)
						n = b + 1 + rng.Intn(width-b)
					}
					cases = add(cases, kind, a, b, n, width)
				}
				compare(fmt.Sprintf("helper=%d seed=%d block=%d", kind, seed, block), cases, rows)
			}
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_RANGE_MOVES_CASES", file.Name()))
	const data = "difftest/sort/range-moves/testdata/"
	must(check.TestInPackage("sort", "TestGoLeanRangeMoves", data+"range_moves_test.go"))
}

func add(cases []fixture, kind, a, b, n, width int) []fixture {
	return append(cases, fixture{kind, 0, a, b, n, width, ""}, fixture{kind, 1, a, b, n, width, ""})
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d\n", c.Kind, c.Variant, c.A, c.B, c.N)
	}
	for i, got := range check.LeanRows("godiff-sort-range-moves", len(cases), &wire) {
		c := cases[i]
		c.Trace = got
		must(rows.Encode(c))
	}
	fmt.Printf("sort range moves %s: %d Lean traces.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
