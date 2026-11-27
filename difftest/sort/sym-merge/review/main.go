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
	Variant, A, M, B int
	Keys             []int
	Want             string
}

func main() {
	fixed := []fixture{
		{0, 0, 1, 4, []int{2, 1, 2, 2}, "1,0,2,3:L2,0;L1,0;S0,1"},
		{0, 0, 3, 4, []int{1, 2, 2, 2}, "0,1,2,3:L3,1;L3,2"},
		{0, 1, 4, 5, []int{9, 1, 4, 7, 3, 8}, "0,1,4,2,3,5:L4,2;L4,1;S4,3;S3,2"},
		{0, 0, 2, 4, []int{2, 4, 1, 3}, "2,0,3,1:L2,1;L3,0;S1,2;L1,0;S0,1;L3,2;S2,3"},
	}
	mutants := map[int]string{
		0: "1,2,3,0:L2,0;L3,0;S0,1;S1,2;S2,3",
		1: "0,3,1,2:L3,1;L3,0;S3,2;S2,1",
		3: "2,0,3,1:L2,1;L3,0;L1,0;S0,1;L3,2;S2,3",
	}
	var cases []fixture
	var input bytes.Buffer
	for _, c := range fixed {
		for variant := range 2 {
			c.Variant = variant
			cases = append(cases, c)
			fmt.Fprintf(&input, "%d %d %d %d", variant, c.A, c.M, c.B)
			for _, key := range c.Keys {
				fmt.Fprintf(&input, " %d", key)
			}
			fmt.Fprintln(&input)
		}
	}
	rows := check.LeanRows("godiff-sort-sym-merge", len(cases), &input)
	file, err := os.CreateTemp("", "golean-sort-sym-merge-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	rejected := 0
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
		if bad, ok := mutants[i/2]; ok && rows[i] != bad {
			rejected++
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_SYM_MERGE_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanSymMerge", "difftest/sort/sym-merge/testdata/sym_merge_test.go"))
	if rejected != 6 {
		log.Fatalf("expected six rejected traces, got %d", rejected)
	}
	fmt.Println("Eight independent ordered callback traces matched native Go and Lean; six mutation witnesses rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
