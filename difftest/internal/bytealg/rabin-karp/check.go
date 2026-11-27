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
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct{ source, separator []byte }

func main() {
	file, err := os.CreateTemp("", "golean-rabin-karp-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	fixture := json.NewEncoder(buffer)
	first := []byte{0x28, 0x9e, 0x1a, 0x30, 0xb8, 0xe4, 0x2a}
	second := []byte{0x20, 0x0c, 0xf2, 0x59, 0x3d, 0x9d, 0x05}
	backA, backB := bytes.Clone(first), bytes.Clone(second)
	slices.Reverse(backA)
	slices.Reverse(backB)
	edges := []input{{first, second}, {append(append(bytes.Clone(first), 0), second...), second},
		{backA, backB}, {append(append(bytes.Clone(backB), 0), backA...), backB}}
	for _, s := range [][]byte{nil, {}, {0}, {0, 0}, {0, 255}, []byte("aaaaa"), []byte("abcabcabc"), {0xff, 0x80, 0xc0}, bytes.Repeat([]byte{255, 0}, 128)} {
		for _, sep := range [][]byte{nil, {}, {0}, {255}, {0, 255}, []byte("aa"), []byte("abc"), {0x80, 0xc0}, bytes.Repeat([]byte{255, 0}, 128)} {
			edges = append(edges, input{s, sep})
		}
	}
	for b := range 256 {
		octet := byte(b)
		edges = append(edges, input{[]byte{octet, octet + 1, octet, octet + 1}, []byte{octet, octet + 1}},
			input{[]byte{octet}, []byte{octet, octet + 1}})
	}
	compare("boundaries", edges, fixture)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5000)
			for i := range cases {
				source, sep := make([]byte, rng.Intn(129)), make([]byte, rng.Intn(53))
				rng.Read(source)
				rng.Read(sep)
				switch i % 5 {
				case 0:
					source = append(append(append(bytes.Clone(sep), source...), sep...), sep...)
				case 1:
					start := rng.Intn(len(source) + 1)
					sep = source[start : start+rng.Intn(len(source)-start+1)]
				case 2:
					for j := range source {
						source[j] &= 3
					}
					for j := range sep {
						sep[j] &= 3
					}
				case 3:
					source, sep = sep, source
				}
				cases[i] = input{source, sep}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, fixture)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_RABIN_KARP_CASES", file.Name()))
	const data = "difftest/internal/bytealg/rabin-karp/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanRabinKarp", data+"bridge_test.go", data+"rabin_karp_test.go"))
}

type observation struct {
	Input, Separator []byte
	Hashes           [4]uint32
	First, Last      string
}

func expected(c input) observation {
	o := observation{Input: c.source, Separator: c.separator, First: "panic", Last: "panic"}
	// Sum the individual terms from the least-significant position.
	var factor uint32 = 1
	for i := range len(c.separator) {
		o.Hashes[0] += uint32(c.separator[len(c.separator)-1-i]) * factor
		o.Hashes[2] += uint32(c.separator[i]) * factor
		factor *= 16777619
	}
	o.Hashes[1], o.Hashes[3] = factor, factor
	if len(c.source) >= len(c.separator) {
		first, last := -1, -1
		for i := 0; i+len(c.separator) <= len(c.source); i++ {
			if slices.Equal(c.source[i:i+len(c.separator)], c.separator) {
				if first < 0 {
					first = i
				}
				last = i
			}
		}
		o.First, o.Last = strconv.Itoa(first), strconv.Itoa(last)
	}
	return o
}

func compare(label string, cases []input, fixture *json.Encoder) {
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, len(c.separator))
		for _, raw := range [][]byte{c.separator, c.source} {
			for _, b := range raw {
				fmt.Fprintf(&rows, " %d", b)
			}
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-rabin-karp", len(cases), &rows) {
		c := cases[i]
		o := expected(c)
		h := o.Hashes
		want := fmt.Sprintf("%d,%d,%d,%d,%s,%s", h[0], h[1], h[2], h[3], o.First, o.Last)
		if got != want {
			log.Fatalf("%s row %d input=%x sep=%x Lean=%s want=%s", label, i, c.source, c.separator, got, want)
		}
		must(fixture.Encode(o))
	}
	fmt.Printf("bytealg Rabin-Karp %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
