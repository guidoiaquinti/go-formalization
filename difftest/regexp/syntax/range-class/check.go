package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Backing, Clean []rune
	Size, I, J     int
	Lo, Hi, Rune   rune
	Want           string
}

func main() {
	f, err := os.CreateTemp("", "golean-regexp-range-class")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, c.Size, " ", len(c.Backing)/2, " ", c.Lo, " ", c.Hi, " ", c.Rune, " ", c.I, " ", c.J)
			for _, class := range [][]rune{c.Backing, c.Clean} {
				for _, value := range class {
					fmt.Fprint(&input, " ", value)
				}
			}
			fmt.Fprintln(&input)
		}
		for i, got := range check.LeanRows("godiff-regexp-range-class", len(cases), &input) {
			c := cases[i]
			c.Want = got
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean range rows\n", label, len(cases))
	}
	var fixed []fixture
	prefixes := [][]rune{nil, {0, 0}, {1, 2, 9, 10}, {4, 7, 4, 8}, {0, 2, 4, 8, 9, 10}, {0, 0x10ffff}}
	classes := [][]rune{nil, {0, 0}, {1, 7}, {0, 0x10ffff}, {0x10ffff, 0x10ffff}, {0, 4, 8, 9, 0x10fffe, 0x10ffff}, {2, 2, 4, 7}}
	requests := [][2]rune{{0, 0}, {3, 4}, {7, 9}, {3, 18}, {0x10ffff, 0x10ffff}, {0, 0x10ffff}}
	for _, prefix := range prefixes {
		for _, clean := range classes {
			for _, request := range requests {
				for spare := range 4 {
					back := append([]rune(nil), prefix...)
					for i := range spare {
						back = append(back, rune(401+i*3), rune(402+i*3))
					}
					for _, r := range []rune{0, 3, 8, 9, 0x10ffff} {
						fixed = append(fixed, fixture{Backing: back, Clean: clean, Size: len(prefix) / 2,
							I: 0, J: max(0, len(prefix)/2-1), Lo: request[0], Hi: request[1], Rune: r})
					}
				}
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				size := r.Intn(7)
				backing := make([]rune, 2*(size+r.Intn(5)))
				for j := 0; j < len(backing); j += 2 {
					lo := r.Intn(160)
					backing[j], backing[j+1] = rune(lo), rune(lo+r.Intn(12))
				}
				var clean []rune
				cursor := r.Intn(3)
				base := []int{0, 0xd7ff, 0x10ff70}[r.Intn(3)]
				for range r.Intn(6) {
					lo := base + cursor
					hi := lo + r.Intn(8)
					clean = append(clean, rune(lo), rune(hi))
					cursor = hi - base + 2 + r.Intn(5)
				}
				lo := r.Intn(160)
				c := fixture{Backing: backing, Clean: clean, Size: size, Lo: rune(lo), Hi: rune(lo + r.Intn(12)), Rune: rune(base + r.Intn(80))}
				if size > 0 {
					c.I, c.J = r.Intn(size), r.Intn(size)
				}
				if len(clean) > 0 && i%3 == 0 {
					c.Rune = clean[r.Intn(len(clean))]
				}
				cases[i] = c
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_REGEXP_RANGE_CLASS_ROWS", f.Name()))
	must(check.TestInPackage("regexp/syntax", "TestGoLeanRangeClass", "difftest/regexp/syntax/range-class/testdata/class_test.go"))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
