package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Present bool
	Words   []uint32
	Freq    []uint16
	Data    []byte
	Want    string
}
type reviewCase struct {
	width int
	fixture
	mutant string
}

func main() {
	w := strconv.IntSize
	cases := []reviewCase{
		{w, fixture{false, nil, []uint16{0, 0}, nil, "0 0 0"}, "nil nil 0"},
		{w, fixture{false, nil, []uint16{0, 1}, nil, "nil nil 0"}, "0 0 0"},
		{w, fixture{true, []uint32{0}, []uint16{1, 1}, []byte{0}, "bounds 2147483647 1"}, "bounds bounds 0"},
		{w, fixture{true, []uint32{256, 3}, []uint16{1, 2}, []byte{0, 1, 0}, "6 6 5"}, "6 2147483647 3"},
		{w, fixture{true, []uint32{7}, []uint16{3, 0, 0}, []byte{1}, "21 21 bounds"}, "bounds bounds 1"},
	}
	words := slices.Repeat([]uint32{255}, 258)
	frequencies := slices.Repeat([]uint16{65535}, 258)
	full := uint64(258 * 65535 * 255)
	narrow := uint64(uint32(full))
	wideText := strconv.FormatUint(full, 10) + " " + strconv.FormatUint(full, 10) + " 255"
	narrowText := strconv.FormatUint(narrow, 10) + " " + strconv.FormatUint(narrow, 10) + " 255"
	cases = append(cases,
		reviewCase{32, fixture{true, words, frequencies, []byte{0}, narrowText}, wideText},
		reviewCase{64, fixture{true, words, frequencies, []byte{0}, wideText}, narrowText})
	var wire bytes.Buffer
	for _, c := range cases {
		p := 0
		if c.Present {
			p = 1
		}
		fmt.Fprintf(&wire, "%d %d %d %d", c.width, p, len(c.Words), len(c.Freq))
		for _, v := range c.Words {
			fmt.Fprintf(&wire, " %d", v)
		}
		for _, v := range c.Freq {
			fmt.Fprintf(&wire, " %d", v)
		}
		for _, v := range c.Data {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	f, err := os.CreateTemp("", "golean-flate-cost-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	rejected, native := 0, 0
	for i, got := range check.LeanRows("godiff-flate-encoder-cost", len(cases), &wire) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, c.Want)
		}
		if got != c.mutant {
			rejected++
		}
		if c.width == w {
			must(rows.Encode(c.fixture))
			native++
		}
	}
	if rejected != len(cases) {
		log.Fatalf("rejected %d of %d altered outcomes", rejected, len(cases))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_COST_ROWS", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanEncoderCost", "difftest/compress/flate/encoder-cost/testdata/cost_test.go"))
	fmt.Printf("Encoder cost: %d Lean rows, %d native rows, %d altered outcomes rejected\n", len(cases), native, rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
