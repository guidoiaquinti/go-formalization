package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Left, Right                    []byte
	LeftOffset, RightOffset, Order int
	Alias                          bool
}

func main() {
	f, err := os.CreateTemp("", "golean-bytealg-compare-")
	must(err)
	defer os.Remove(f.Name())
	buffer := bufio.NewWriter(f)
	rows := json.NewEncoder(buffer)
	values := [][]byte{nil, {}, {0}, {0, 1}, {0, 0}, {1, 0}, {127}, {128}, {255}, {0xc0, 0x80}}
	var cases []fixture
	for _, left := range values {
		for _, right := range values {
			cases = append(cases, fixture{Left: left, Right: right})
		}
	}
	for a := range 256 {
		for b := range 256 {
			cases = append(cases, fixture{Left: []byte{byte(a)}, Right: []byte{byte(b)}, LeftOffset: a % 64, RightOffset: b % 64})
		}
	}
	for offset := range 64 {
		for _, width := range []int{15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 129} {
			left := bytes.Repeat([]byte{0x80}, width)
			right := bytes.Clone(left)
			right[width-1] = 0x7f
			cases = append(cases, fixture{Left: left, Right: right, LeftOffset: offset, RightOffset: 63 - offset},
				fixture{Left: right, Right: left, LeftOffset: offset, RightOffset: 63 - offset},
				fixture{Left: left, Right: left, LeftOffset: offset, RightOffset: 63 - offset})
		}
	}
	for first := 0; first < len(cases); first += 5000 {
		last := min(first+5000, len(cases))
		compare(fmt.Sprintf("boundary rows %d–%d", first, last-1), cases[first:last], rows)
	}
	fmt.Printf("internal/bytealg comparison boundaries: %d cases.\n", len(cases))
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				a, b := make([]byte, rng.Intn(257)), make([]byte, rng.Intn(257))
				rng.Read(a)
				rng.Read(b)
				alias := false
				switch i % 6 {
				case 0:
					b = a
					alias = true
				case 1:
					b = bytes.Clone(a)
				case 2:
					b = a[:rng.Intn(len(a)+1)]
				case 3:
					b = bytes.Clone(a)
					if len(b) > 0 {
						b[rng.Intn(len(b))] ^= 255
					}
				case 4:
					b = append(bytes.Clone(a), 0)
				}
				cases[i] = fixture{a, b, rng.Intn(64), rng.Intn(64), 0, alias}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_BYTEALG_COMPARE_CASES", f.Name()))
	const data = "difftest/internal/bytealg/compare/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgCompare", data+"bridge_test.go", data+"compare_test.go"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var wire bytes.Buffer
	for i := range cases {
		c := &cases[i]
		size := func(raw []byte) int {
			if raw == nil {
				return -1
			}
			return len(raw)
		}
		fmt.Fprintf(&wire, "%d %d", size(c.Left), size(c.Right))
		for _, raw := range [][]byte{c.Left, c.Right} {
			for _, b := range raw {
				fmt.Fprintf(&wire, " %d", b)
			}
		}
		fmt.Fprintln(&wire)
		c.Order = slices.Compare(c.Left, c.Right)
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-compare", len(cases), &wire) {
		c := cases[i]
		want := fmt.Sprintf("%d,%d,%d,%t", c.Order, c.Order, c.Order, c.Order == 0)
		if got != want {
			log.Fatalf("%s row %d a=%x b=%x: Lean=%s want=%s", label, i, c.Left, c.Right, got, want)
		}
		must(rows.Encode(c))
	}
	fmt.Printf("internal/bytealg comparison %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
