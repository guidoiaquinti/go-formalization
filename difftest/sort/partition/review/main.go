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
	Kind, Variant, A, B, Pivot int
	Values                     []int
	Want                       string
}

func main() {
	fixed := []fixture{
		{0, 0, 1, 4, 1, []int{9, 2, 3, 1, 8}, "2:false:9,1,2,3,8:S1,1;L2,1;L3,1;S2,3;S2,1"},
		{0, 0, 0, 3, 1, []int{1, 2, 3}, "1:true:1,2,3:S0,1;L1,0;L2,0;L2,0;S1,0"},
		{0, 0, 0, 1, 0, []int{4}, "0:true:4:S0,0;S0,0"},
		{1, 0, 0, 4, 1, []int{2, 1, 3, 1}, "2:-:1,1,3,2:S0,1;L0,1;L0,3;S1,3;L0,2;L0,2"},
		{1, 0, 1, 4, 2, []int{9, 5, 5, 5, 8}, "4:-:9,5,5,5,8:S1,2;L1,2;L1,3"},
	}
	mutants := map[int]string{
		0: "2:true:9,1,2,3,8:S1,1;L2,1;L3,1;S2,3;S2,1",
		2: "0:true:4:S0,0",
		3: "1:-:1,1,3,2:S0,1;L0,1;L0,3;S1,3;L0,2;L0,2",
	}
	var cases []fixture
	var input bytes.Buffer
	for _, f := range fixed {
		for variant := range 2 {
			f.Variant = variant
			cases = append(cases, f)
			fmt.Fprintf(&input, "%d %d %d %d %d", f.Kind, variant, f.A, f.B, f.Pivot)
			for _, v := range f.Values {
				fmt.Fprintf(&input, " %d", v)
			}
			fmt.Fprintln(&input)
		}
	}
	rows := check.LeanRows("godiff-sort-partition", len(cases), &input)
	file, err := os.CreateTemp("", "golean-sort-partition-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	rejected := 0
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("%+v: Lean %q", c, rows[i])
		}
		must(enc.Encode(c))
		if broken, ok := mutants[i/2]; ok && broken != rows[i] {
			rejected++
		}
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_PARTITION_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanPartition", "difftest/sort/partition/testdata/partition_test.go"))
	if rejected != 6 {
		log.Fatalf("expected six mutant rejections, got %d", rejected)
	}
	fmt.Println("10 independent native and Lean traces passed; 6 variant-specific mutation witnesses rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
