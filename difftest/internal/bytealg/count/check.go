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
	Bytes        []byte
	Target       byte
	Offset, Want int
}

func main() {
	file, err := os.CreateTemp("", "golean-bytealg-count-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	rows := json.NewEncoder(buffer)
	cases := []fixture{{Offset: -1}, {Offset: 0}, {Target: 255, Offset: 31}}
	for target := range 256 {
		raw := []byte{byte(target), byte(target - 1), byte(target), byte(target + 1)}
		cases = append(cases, fixture{raw, byte(target), target % 64, 0})
	}
	for offset := range 64 {
		for _, width := range []int{0, 1, 2, 15, 16, 17, 30, 31, 32, 33, 47, 48, 63, 64, 65, 95, 96, 127, 128, 129, 255, 256, 257, 511, 512, 513} {
			raw := make([]byte, width)
			for i := range raw {
				if i%3 != 1 {
					raw[i] = 255
				}
			}
			cases = append(cases, fixture{raw, 255, offset, 0}, fixture{raw, 0, offset, 0})
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
				cases[i] = fixture{raw, target, rng.Intn(64), 0}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, rows)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_BYTEALG_COUNT_CASES", file.Name()))
	const data = "difftest/internal/bytealg/count/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgCount", data+"bridge_test.go", data+"count_test.go"))
}

func compare(label string, cases []fixture, rows *json.Encoder) {
	var input bytes.Buffer
	for i := range cases {
		c := &cases[i]
		fmt.Fprint(&input, c.Target)
		var frequencies [256]int
		for _, byteValue := range c.Bytes {
			frequencies[byteValue]++
			fmt.Fprintf(&input, " %d", byteValue)
		}
		c.Want = frequencies[c.Target]
		fmt.Fprintln(&input)
	}
	for i, got := range check.LeanRows("godiff-internal-bytealg-count", len(cases), &input) {
		c := cases[i]
		want := fmt.Sprintf("%d,%d,%d,%d", c.Want, c.Want, c.Want, c.Want)
		if got != want {
			log.Fatalf("%s row %d bytes=%x target=%02x Lean=%s want=%s", label, i, c.Bytes, c.Target, got, want)
		}
		must(rows.Encode(c))
	}
	fmt.Printf("internal/bytealg count %s: %d Lean rows agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
