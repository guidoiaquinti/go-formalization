package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/bits"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Variant, A, B, Limit int
	Values               []int
	Want                 string
}

func main() {
	f, err := os.CreateTemp("", "golean-pdqsort-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	compare("boundaries", boundaries(), rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 0, 10000)
			for i := range 5000 {
				n := rng.Intn(65)
				if i%127 == 0 {
					n = 65 + rng.Intn(144)
				}
				a := rng.Intn(5)
				limit := []int{0, 1, 2, 4, bits.Len(uint(n)), math.MaxInt}[rng.Intn(6)]
				values := sample(rng, a, n, i%8)
				cases = add(cases, a, n, limit, values)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_PDQ_CASES", f.Name()))
	must(check.TestInPackage("sort", "TestGoLeanPdqsort", "difftest/sort/pdqsort/testdata/pdqsort_test.go"))
}

func sample(rng *rand.Rand, a, n, pattern int) []int {
	values := make([]int, a+n+2)
	for j := range values {
		if j < a {
			values[j] = math.MinInt
			continue
		}
		switch pattern {
		case 0:
			values[j] = int(rng.Uint64())
		case 1:
			values[j] = rng.Intn(3)
		case 2:
			values[j] = j
		case 3:
			values[j] = len(values) - j
		case 4:
			values[j] = 1
		case 5:
			values[j] = min(j, len(values)-j)
		case 6:
			values[j] = (j % 11) * 2
		case 7:
			values[j] = (j + 7) * (j + 7) % 97
		}
	}
	return values
}
func add(cases []fixture, a, n, limit int, values []int) []fixture {
	return append(cases, fixture{0, a, a + n, limit, values, ""}, fixture{1, a, a + n, limit, values, ""})
}
func boundaries() []fixture {
	var cases []fixture
	rng := rand.New(rand.NewSource(1))
	for _, n := range []int{0, 1, 2, 7, 8, 11, 12, 13, 14, 31, 32, 48, 49, 50, 51, 64, 65, 127, 128, 129, 255, 256} {
		for _, a := range []int{0, 1, 5} {
			for _, limit := range []int{0, 1, 2, bits.Len(uint(n)), math.MaxInt} {
				for pattern := range 8 {
					cases = add(cases, a, n, limit, sample(rng, a, n, pattern))
				}
			}
		}
	}
	for _, n := range []int{24, 31, 32, 48, 49} {
		values := sample(rng, 0, n, 2)
		for i := 1; i <= 3; i++ {
			values[i*(n/4)] = -1
		}
		cases = add(cases, 0, n, 8, values)
	}
	return cases
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Variant, strconv.IntSize, c.A, c.B, c.Limit)
		for _, v := range c.Values {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-sort-pdqsort", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("sort pdqsort %s: %d Lean traces.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
