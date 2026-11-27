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
	Kind, Shape, A, B int
	Replies           []int
	Want              string
}

func main() {
	cases := []fixture{
		{0, 0, 4, 0, []int{2, 1, 0}, "false:3,2;2,1"},
		{0, 1, 2, 0, []int{0}, "true:1,0"},
		{0, 2, 3, 0, []int{0, 2}, "callback:2,1"},
		{0, 3, 3, 0, []int{0, 0}, "true:2,1;1,0"},
		{0, 4, 2, 0, []int{0}, "true:1,0"},
		{0, 5, 1, 0, nil, "invalid:"},
		{0, 6, 1, 0, nil, "unsupported:"},
		{1, 0, 9, 0, []int{1}, "true:1,9"},
		{1, 0, 0, 2, []int{8}, "callback:8,0"},
		{1, 1, 1, 0, []int{9}, "nil:"},
	}
	mutants := map[int]string{
		0: "callback:1,0",
		1: "unsupported:",
		7: "false:9,1",
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d", c.Kind, c.Shape, c.A, c.B)
		for _, reply := range c.Replies {
			fmt.Fprintf(&input, " %d", reply)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-sort-comparison-entrypoints", len(cases), &input)
	file, err := os.CreateTemp("", "golean-sort-comparison-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	rejected := 0
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
		if mutant, ok := mutants[i]; ok && rows[i] != mutant {
			rejected++
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_COMPARISON_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanComparisonEntrypoints",
		"difftest/sort/comparison-entrypoints/testdata/comparison_test.go"))
	if rejected != 3 {
		log.Fatalf("expected three mutation witnesses, got %d", rejected)
	}
	fmt.Println("Ten independent observations matched native Go and Lean; three mutant witnesses were rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
