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
	Word         uint64
	Base, Length int
	Cells        []byte
	Want         string
}

func main() {
	f, err := os.CreateTemp("", "golean-integer-format-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	var fixed []fixture
	for _, word := range []uint64{0, 1, 9, 10, 99, 100, 101, 255, 256, 1 << 31, 1 << 32, math.MaxInt64, 1 << 63, 1<<63 | 1, math.MaxUint64 - 1, math.MaxUint64} {
		for base := -1; base <= 37; base++ {
			for _, spare := range []int{0, 1, 2, 9, 17, 64, 65} {
				fixed = append(fixed, fixture{word, base, 2, bytes.Repeat([]byte{165}, 2+spare), ""})
			}
		}
	}
	for _, base := range []int{math.MinInt, 0, 2, 10, 16, 36, math.MaxInt} {
		fixed = append(fixed, fixture{Word: math.MaxUint64, Base: base}, fixture{Word: math.MaxUint64, Base: base, Cells: []byte{}})
	}
	compare("fixed", fixed, rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				word := rng.Uint64()
				if i%7 == 0 {
					word = uint64(rng.Intn(200))
				}
				base := rng.Intn(43) - 3
				if i%5 == 0 {
					base = 10
				}
				length := rng.Intn(8)
				cells := make([]byte, length+rng.Intn(70))
				rng.Read(cells)
				if i%11 == 0 {
					length, cells = 0, nil
				}
				cases[i] = fixture{word, base, length, cells, ""}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_STRCONV_FORMAT_CASES", f.Name()))
	must(check.TestInPackage("internal/strconv", "TestGoLeanIntegerFormat", "difftest/strconv/integer-format/testdata/format_test.go"))
}
func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for _, c := range cases {
		p := 0
		if c.Cells != nil {
			p = 1
		}
		fmt.Fprintf(&wire, "%d %d %d %d %d", strconv.IntSize, c.Word, c.Base, c.Length, p)
		for _, x := range c.Cells {
			fmt.Fprintf(&wire, " %d", x)
		}
		fmt.Fprintln(&wire)
	}
	for i, got := range check.LeanRows("godiff-strconv-integer-format", len(cases), &wire) {
		c := cases[i]
		c.Want = got
		must(rows.Encode(c))
	}
	fmt.Printf("integer formatting %s: %d Lean rows.\n", label, len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
