package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Left, Right []byte
	RightLength int
	Want        string
}

func expected(c fixture) string {
	// A mismatch can be observed only after the relevant word or tail slice
	// has passed its own bounds check.
	mismatch := 0
	for mismatch < min(len(c.Left), len(c.Right)) && c.Left[mismatch] == c.Right[mismatch] {
		mismatch++
	}
	whole := len(c.Left) / 8 * 8
	if mismatch < whole {
		if (mismatch/8+1)*8 > len(c.Right) {
			return "bounds"
		}
		return strconv.Itoa(mismatch)
	}
	if whole > c.RightLength || len(c.Left) > len(c.Right) {
		return "bounds"
	}
	return strconv.Itoa(mismatch)
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-match-length")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, len(c.Left), " ", c.RightLength)
			for _, cells := range [][]byte{c.Left, c.Right} {
				for _, b := range cells {
					fmt.Fprint(&input, " ", b)
				}
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-flate-match-length", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for _, left := range []int{0, 1, 7, 8, 9, 15, 16, 17, 23, 24, 25, 32, 33} {
		for _, right := range []int{0, 1, 7, 8, 9, 15, 16, 17, 23, 24, 25, 32, 33} {
			for _, visible := range []int{0, right / 2, right} {
				for _, flip := range []int{-1, 0, 1, 6, 7, 8, 15, 16, 23, 24, 31, 32} {
					c := fixture{Left: make([]byte, left), Right: make([]byte, right), RightLength: visible}
					if flip >= 0 && flip < right {
						c.Right[flip] = byte(1 << uint(flip%8))
					}
					fixed = append(fixed, c)
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Left: make([]byte, r.Intn(65)), Right: make([]byte, r.Intn(65))}
			c.RightLength = r.Intn(len(c.Right) + 1)
			r.Read(c.Left)
			r.Read(c.Right)
			common := r.Intn(min(len(c.Left), len(c.Right)) + 1)
			copy(c.Right, c.Left[:common])
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_MATCH_LENGTH_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanMatchLength", "difftest/compress/flate/match-length/testdata/match_length_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
