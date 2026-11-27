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
	Bytes       []byte
	Target      byte
	Offset      int
	First, Last int
}

func main() {
	file, err := os.CreateTemp("", "golean-bytealg-byte-index-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	cases := []fixture{{Offset: -1}, {Offset: 0}, {Target: 255, Offset: 31}}
	for target := range 256 {
		raw := []byte{byte(target), byte(target - 1), byte(target), byte(target + 1)}
		cases = append(cases, fixture{Bytes: raw, Target: byte(target), Offset: target % 64})
	}
	for offset := range 64 {
		for _, width := range []int{0, 1, 2, 15, 16, 17, 30, 31, 32, 33, 47, 48, 63, 64, 65, 95, 96, 127, 128, 129, 255, 256, 257, 511, 512, 513} {
			raw := make([]byte, width)
			for i := range raw {
				if i%3 != 1 {
					raw[i] = 255
				}
			}
			cases = append(cases, fixture{Bytes: raw, Target: 255, Offset: offset}, fixture{Bytes: raw, Target: 0, Offset: offset})
		}
	}
	compare("boundaries", cases, rows)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				target := byte(rng.Intn(256))
				raw := make([]byte, rng.Intn(129))
				rng.Read(raw)
				if i%3 == 0 {
					for j := range raw {
						if rng.Intn(2) == 0 {
							raw[j] = target
						}
					}
				}
				if i%7 == 0 && len(raw) > 0 {
					target = raw[rng.Intn(len(raw))]
				}
				cases[i] = fixture{Bytes: raw, Target: target, Offset: rng.Intn(64)}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BYTEALG_BYTE_INDEX_CASES", file.Name()))
	const data = "difftest/internal/bytealg/byte-index/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgByteIndex", data+"bridge_test.go", data+"byte_index_test.go"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var input bytes.Buffer
	for i := range cases {
		c := &cases[i]
		fmt.Fprint(&input, c.Target)
		var positions [256][]int
		for position, b := range c.Bytes {
			positions[b] = append(positions[b], position)
			fmt.Fprintf(&input, " %d", b)
		}
		hits := positions[c.Target]
		c.First, c.Last = -1, -1
		if len(hits) > 0 {
			c.First, c.Last = hits[0], hits[len(hits)-1]
		}
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-byte-index", len(cases), &input) {
		c := cases[i]
		want := fmt.Sprintf("%d,%d,%d,%d", c.First, c.First, c.Last, c.Last)
		if got != want {
			log.Fatalf("%s row %d bytes=%x target=%02x Lean=%s want=%s", label, i, c.Bytes, c.Target, got, want)
		}
		must(rows.Encode(c))
	}
	fmt.Printf("internal/bytealg byte index %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
