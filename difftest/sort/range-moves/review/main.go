package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, Variant, A, B, N, Width int
	Trace                         string
}

func main() {
	fixed := []fixture{
		{0, 0, 1, 5, 0, 6, "swaps:1,4;2,3"},
		{0, 0, 3, 4, 0, 6, "swaps:"},
		{1, 0, 0, 1, 3, 5, "swaps:0,1;1,2;2,3"},
		{1, 0, 2, 2, 2, 5, "swaps:2,2;3,3"},
		{2, 0, 0, 2, 5, 6, "swaps:0,3;1,4;0,2;1,2"},
		{2, 0, 0, 3, 5, 6, "swaps:0,3;1,4;2,4;2,3"},
		{2, 0, 1, 3, 5, 6, "swaps:1,3;2,4"},
		{2, 0, 2, 2, 2, 6, "swaps:"},
	}
	// Reversing transposition order preserves a reversal permutation, but
	// changes the callback trace. Overlap makes stack-order swaps observable.
	// Triple reversal yields the same rotation but calls different swaps.
	mutants := map[int]string{
		0: "swaps:2,3;1,4",
		2: "swaps:2,3;1,2;0,1",
		4: "swaps:0,1;2,4;0,4;1,3",
	}
	var input bytes.Buffer
	for _, f := range fixed {
		for variant := range 2 {
			fmt.Fprintf(&input, "%d %d %d %d %d\n", f.Kind, variant, f.A, f.B, f.N)
		}
	}
	rows := check.LeanRows("godiff-sort-range-moves", 2*len(fixed), &input)
	file, err := os.CreateTemp("", "golean-sort-range-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	rejected := 0
	for i, f := range fixed {
		for variant := range 2 {
			f.Variant = variant
			if got := rows[2*i+variant]; got != f.Trace {
				log.Fatalf("fixture %+v: Lean got %q", f, got)
			}
			must(enc.Encode(f))
			if broken, ok := mutants[i]; ok && broken != rows[2*i+variant] {
				rejected++
			}
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_RANGE_MOVES_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanRangeMoves", "difftest/sort/range-moves/testdata/range_moves_test.go"))
	if rejected != 6 {
		log.Fatalf("expected six variant-specific mutant rejections, got %d", rejected)
	}
	fmt.Println("16 independent traces matched pinned Go and Lean; all 6 mutant variants rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
