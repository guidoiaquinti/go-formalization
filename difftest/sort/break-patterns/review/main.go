package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Kind, A, B int
	State      uint64
	Want       string
}

func main() {
	cases := []fixture{
		{Kind: 0, State: 1, Want: "1082269761"},
		{Kind: 0, State: 1 << 63, Want: "9295429630892703744"},
		{Kind: 0, State: math.MaxUint64, Want: "1065361344"},
		{Kind: 1, A: -1, Want: "0"},
		{Kind: 1, A: 0, Want: "1"},
		{Kind: 1, A: 8, Want: "16"},
		{Kind: 1, A: math.MaxInt, Want: strconv.FormatUint(1<<(strconv.IntSize-1), 10)},
	}
	for _, variant := range []int{2, 3} {
		cases = append(cases,
			fixture{Kind: variant, A: 0, B: 7, Want: "swaps:"},
			fixture{Kind: variant, A: 0, B: 8, Want: "swaps:3,0;4,4;5,0"},
			fixture{Kind: variant, A: 19, B: 28, Want: "swaps:22,19;23,23;24,20"},
		)
	}
	// Witnesses for skipping the last xor, rounding an exact power down,
	// starting at span seven, and using random % length instead of folding.
	mutants := map[int]string{
		0:  "8257",
		5:  "8",
		7:  "swaps:1,0;2,4;3,0",
		9:  "swaps:22,19;23,21;24,22",
		10: "swaps:1,0;2,4;3,0",
		12: "swaps:22,19;23,21;24,22",
	}
	var input bytes.Buffer
	for _, c := range cases {
		if c.Kind == 0 {
			fmt.Fprintf(&input, "%d %d %d %d\n", c.Kind, strconv.IntSize, c.State, c.B)
		} else {
			fmt.Fprintf(&input, "%d %d %d %d\n", c.Kind, strconv.IntSize, c.A, c.B)
		}
	}
	rows := check.LeanRows("godiff-sort-break-patterns", len(cases), &input)
	file, err := os.CreateTemp("", "golean-sort-break-patterns-review")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	rejected := 0
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %+v: Lean got %q", c, rows[i])
		}
		must(encoder.Encode(c))
		if broken, ok := mutants[i]; ok && broken != rows[i] {
			rejected++
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_BREAK_PATTERNS_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanBreakPatterns", "difftest/sort/break-patterns/testdata/break_patterns_test.go"))
	if rejected != 6 {
		log.Fatalf("expected six mutant rejections, got %d", rejected)
	}
	fmt.Println("13 independent cases passed; six mutant witnesses rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
