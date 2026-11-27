package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Mode int
	Word uint64
	Want string
}

func extract(w uint64, start, width int) uint64 {
	return (w >> start) & (math.MaxUint64 >> (64 - width))
}

func main() {
	cases := []fixture{{Word: 0}, {Word: math.MaxUint64}}
	for i := range 64 {
		cases = append(cases, fixture{Word: uint64(1) << i})
	}
	cases = append(cases, fixture{Mode: 1}, fixture{Mode: 2}, fixture{Mode: 3})
	var wire bytes.Buffer
	var rejects [4]int
	for i := range cases {
		c := &cases[i]
		fmt.Fprintf(&wire, "%d %d %d\n", strconv.IntSize, c.Mode, c.Word)
		switch c.Mode {
		case 1:
			c.Want = "source:0" + strings.Repeat("|callback:0:1", 8)
		case 2:
			c.Want = "nil:0" + strings.Repeat("|nil:0:0", 8)
		case 3:
			c.Want = "source:0" + strings.Repeat("|nil:0:0", 8)
		default:
			w := c.Word
			vals := [8]uint64{
				extract(w, 0, 63), extract(w, 32, 32), w, extract(w, 33, 31),
				extract(w, 0, strconv.IntSize-1), extract(w, 0, strconv.IntSize),
				extract(w, 0, 53), extract(w, 32, 24),
			}
			c.Want = "source:0"
			for _, v := range vals {
				c.Want += fmt.Sprintf("|ok:%d:1", v)
			}
			for j, alternative := range [4]uint64{extract(w, 0, 32), extract(w, 0, 31), extract(w, 11, 53), extract(w, 40, 24)} {
				if vals[[4]int{1, 3, 6, 7}[j]] != alternative {
					rejects[j]++
				}
			}
		}
	}
	if rejects != [4]int{64, 62, 64, 32} {
		log.Fatalf("unexpected interval mutant counts %v", rejects)
	}
	for i, got := range check.LeanRows("godiff-rand-v2-word", len(cases), &wire) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, cases[i].Want)
		}
	}
	file, err := os.CreateTemp("", "golean-rand-word-review-")
	must(err)
	defer os.Remove(file.Name())
	rows := json.NewEncoder(file)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_RAND_WORD_CASES", file.Name()))
	const dir = "difftest/math/rand/v2/word/testdata/"
	must(check.TestInPackage("math/rand/v2", "TestGoLeanRandWord", dir+"bridge_test.go", dir+"word_test.go"))
	fmt.Printf("Rand one-word review: %d independent native and Lean rows, incorrect interval counts %v\n", len(cases), rejects)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
