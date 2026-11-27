package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Left, Right []byte
	RightLength int
	Want        string
}

func main() {
	fixed := []fixture{
		{nil, nil, 0, "0"},
		{[]byte("xy"), []byte("xy"), 0, "2"},
		{[]byte("01234567"), []byte("0123456X"), 0, "7"},
		{[]byte("01234567"), []byte("X123456"), 1, "bounds"},
		{[]byte("01234567"), []byte("01234567"), 7, "bounds"},
		{[]byte("01234567"), []byte("01234567"), 8, "8"},
		{[]byte("01234567x"), []byte("01234567y"), 8, "8"},
		{[]byte("01234567xy"), []byte("01234567z"), 8, "bounds"},
		{[]byte("0123456789abcdef"), []byte("0123456789abcdeX"), 8, "15"},
		{[]byte("0123456789abcdef"), []byte("0123456789abcdef"), 15, "bounds"},
		{[]byte("0123456789abcdef"), []byte("01234567X9abcde"), 8, "bounds"},
		{[]byte("0123456789abcdefx"), []byte("0123456789abcdefx"), 16, "17"},
		{[]byte{0, 255, 128}, []byte{0, 255, 128}, 1, "3"},
	}
	f, err := os.CreateTemp("", "golean-flate-match-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	var mutations [3]int
	for _, c := range fixed {
		fmt.Fprint(&input, len(c.Left), " ", c.RightLength)
		for _, data := range [][]byte{c.Left, c.Right} {
			for _, b := range data {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
		for i, mutant := range []string{
			observe(c, true, false, false),
			observe(c, false, true, false),
			observe(c, false, false, true),
		} {
			if mutant != c.Want {
				mutations[i]++
			}
		}
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-flate-match-length", len(fixed), &input) {
		if want := fixed[i].Want; got != want {
			log.Fatalf("control %d: Lean returned %q, want %q", i, got, want)
		}
	}
	must(os.Setenv("GOLEAN_FLATE_MATCH_LENGTH_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanMatchLength", "difftest/compress/flate/match-length/testdata/match_length_test.go"))
	if want := [3]int{6, 2, 3}; mutations != want {
		log.Fatalf("mutation disagreements %v, want %v", mutations, want)
	}
	fmt.Printf("%d literal control rows; mutation disagreements %v\n", len(fixed), mutations)
}

func observe(c fixture, clampToVisible, allowTailDrop, checkOnRead bool) string {
	capacity := len(c.Right)
	if clampToVisible {
		capacity = c.RightLength
	}
	for i, left := range c.Left {
		if !checkOnRead && i%8 == 0 {
			end := min(i+8, len(c.Left))
			if end > capacity {
				return "bounds"
			}
		}
		if !allowTailDrop && i >= len(c.Left)/8*8 && len(c.Left)/8*8 > c.RightLength {
			return "bounds"
		}
		if i >= capacity {
			return "bounds"
		}
		if left != c.Right[i] {
			return strconv.Itoa(i)
		}
	}
	if !allowTailDrop && len(c.Left)/8*8 > c.RightLength {
		return "bounds"
	}
	return strconv.Itoa(len(c.Left))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
