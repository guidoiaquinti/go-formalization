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
	Kind, Size, Offset, Length int
	Backing                    []byte
	Want                       string
}

func expected(dec, enc int, status, backing string) string {
	return fmt.Sprintf("%d|%s:%d:%s", dec, status, enc, backing)
}

func main() {
	cells := []byte{1, 2, 3, 4, 5}
	const wireCells = "1,2,3,4,5"
	cases := []fixture{
		{1, 0, 2, 1, cells, expected(6, 2, "bounds", wireCells)},
		{1, 0, 1, 1, cells, expected(5, 5, "ok", "1,0,0,0,0")},
		{3, 2, 2, 1, cells, expected(4, 4, "ok", "1,2,0,0,5")},
		{5, 0, 5, 0, cells, expected(5, 5, "ok", wireCells)},
		{5, 0, 6, 0, cells, expected(6, 6, "bounds", wireCells)},
		{6, 0, 2, 4, cells, expected(1, 2, "bounds", wireCells)},
		{0, 0, math.MinInt, 4, cells, expected(math.MaxInt, math.MinInt, "bounds", wireCells)},
		{1, 0, math.MaxInt, 4, cells, expected(math.MinInt+3, math.MaxInt, "bounds", wireCells)},
		{7, 0, 0, 0, nil, expected(-1, 0, "bounds", "")},
		{5, 0, 0, 0, nil, expected(0, 0, "ok", "")},
	}
	var input bytes.Buffer
	var mutants [5]int
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d", strconv.IntSize, c.Kind, c.Size, c.Offset)
		for _, b := range c.Backing {
			fmt.Fprintf(&input, " %d", b)
		}
		fmt.Fprintln(&input)
		dec, enc, _ := strings.Cut(c.Want, "|")
		if (c.Kind == 0 || c.Kind == 6 || c.Kind == 7) && dec != strconv.Itoa(c.Offset) {
			mutants[0]++
		}
		if (c.Offset == math.MinInt || c.Offset == math.MaxInt) && dec != strconv.Itoa(c.Offset) {
			mutants[1]++
		}
		if strings.HasPrefix(enc, "ok:") && c.Offset > c.Length {
			mutants[2]++
		}
		if strings.HasPrefix(enc, "bounds:") && dec != strconv.Itoa(c.Offset) {
			mutants[3]++
		}
		if c.Kind == 5 && strings.HasPrefix(enc, "bounds:") {
			mutants[4]++
		}
	}
	for i, got := range check.LeanRows("godiff-binary-coder-skip", len(cases), &input) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, cases[i].Want)
		}
	}
	if mutants != [5]int{3, 2, 2, 5, 1} {
		log.Fatalf("incorrect mutant counts: %v", mutants)
	}
	file, err := os.CreateTemp("", "golean-coder-skip-review-")
	must(err)
	defer os.Remove(file.Name())
	rows := json.NewEncoder(file)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_BINARY_CODER_SKIP_CASES", file.Name()))
	const data = "difftest/encoding/binary/coder-skip/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanCoderSkip", data+"bridge_test.go", data+"coder_skip_test.go"))
	fmt.Printf("coder skip: %d independent native and Lean rows, mutant rejections %v\n", len(cases), mutants)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
