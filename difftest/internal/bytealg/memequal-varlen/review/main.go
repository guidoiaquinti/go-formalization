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
	Left, Right                   []byte
	Size, LeftOffset, RightOffset int
	Alias, Want                   bool
}

func main() {
	cases := []fixture{
		{[]byte{1, 2, 3, 9}, []byte{1, 2, 3, 8}, 3, 1, 3, false, true},
		{[]byte{1, 2, 3}, []byte{1, 2, 4}, 3, 7, 0, false, false},
		{[]byte{1, 2, 3}, []byte{1, 4, 3}, 3, 0, 15, false, false},
		{[]byte{0, 255, 0, 3, 4}, []byte{0, 255, 0, 3, 4}, 5, 11, 13, false, true},
		{[]byte{1, 2, 3}, []byte{1, 2, 3}, 3, 2, 2, true, true},
		{bytes.Repeat([]byte{0xa5}, 130), bytes.Repeat([]byte{0xa5}, 130), 129, 3, 17, false, true},
	}
	cases[5].Right[129] = 0xff
	left, right := bytes.Repeat([]byte{1}, 257), bytes.Repeat([]byte{1}, 257)
	right[256] = 2
	cases = append(cases, fixture{left, right, 257, 1, 17, false, false})
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d", c.Size, len(c.Left))
		for _, window := range [][]byte{c.Left, c.Right} {
			for _, b := range window {
				fmt.Fprintf(&input, " %d", b)
			}
		}
		fmt.Fprintln(&input)
	}
	fmt.Fprintln(&input, "0 0")
	for i, row := range check.LeanRows("godiff-internal-bytealg-memequal-varlen", len(cases)+1, &input) {
		want := true
		if i < len(cases) {
			want = cases[i].Want
		}
		if row != strconv.FormatBool(want) {
			log.Fatalf("case %d: Lean %s, want %t", i, row, want)
		}
	}
	out, err := os.CreateTemp("", "golean-memequal-varlen-review-*")
	must(err)
	defer os.Remove(out.Name())
	for _, c := range cases {
		must(json.NewEncoder(out).Encode(c))
	}
	must(out.Close())
	must(os.Setenv("GOLEAN_BYTEALG_MEMEQUAL_VARLEN_CASES", out.Name()))
	const data = "difftest/internal/bytealg/memequal-varlen/testdata/memequal_varlen_test.go"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanMemequalVarlen", data))
	fmt.Printf("%d native closure fixtures and %d Lean observations passed.\n", len(cases), len(cases)+1)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
