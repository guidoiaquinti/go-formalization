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
	Kind, Variant, A, B, C int
	Values                 []int
	Want                   string
}

func main() {
	fixed := []fixture{
		{0, 0, 1, 4, 0, []int{9, 3, 2, 1, 8}, "9,1,2,3,8:L2,1;S2,1;L3,2;S3,2;L2,1;S2,1"},
		{0, 0, 0, 2, 0, []int{4, 4}, "4,4:L1,0"},
		{1, 0, 0, 3, 1, []int{9, 0, 3, 3, 8}, "9,3,0,3,8:L2,3;L1,2;S1,2"},
		{1, 0, 0, 3, 1, []int{9, 6, 3, 5, 8}, "9,6,3,5,8:L2,3;L1,3"},
		{2, 0, 1, 3, 0, []int{9, 2, 1, 8}, "9,1,2,8:L1,2;S1,2;S1,1"},
		{2, 0, 1, 1, 0, []int{9, 2}, "9,2:"},
	}
	broken := map[int]string{
		1: "4,4:L1,0;S1,0",
		2: "9,3,3,0,8:L2,3;L1,3;S1,3",
		4: "9,1,2,8:L1,2;S1,2",
	}
	var cases []fixture
	var input bytes.Buffer
	for _, f := range fixed {
		for variant := range 2 {
			f.Variant = variant
			cases = append(cases, f)
			fmt.Fprintf(&input, "%d %d %d %d %d", f.Kind, variant, f.A, f.B, f.C)
			for _, v := range f.Values {
				fmt.Fprintf(&input, " %d", v)
			}
			fmt.Fprintln(&input)
		}
	}
	rows := check.LeanRows("godiff-sort-local-sorts", len(cases), &input)
	file, err := os.CreateTemp("", "golean-sort-local-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	rejected := 0
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %+v: Lean %q", c, rows[i])
		}
		must(enc.Encode(c))
		if m, ok := broken[i/2]; ok && rows[i] != m {
			rejected++
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_LOCAL_SORT_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanLocalSorts", "difftest/sort/local-sorts/testdata/local_sorts_test.go"))
	if rejected != 6 {
		log.Fatalf("expected six rejected traces, got %d", rejected)
	}
	fmt.Println("12 independent native and Lean traces passed; six variant-specific mutant traces rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
