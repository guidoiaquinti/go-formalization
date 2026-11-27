package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Variant, A, B int
	Values        []int
	CrossesStart  bool
	Want          string
}

func main() {
	values := make([]int, 54)
	for i := range values {
		values[i] = i
	}
	copy(values, []int{100, 101, 0, 1, -1})
	result := []string{"-1", "100", "101", "0", "1"}
	for i := 5; i < len(values); i++ {
		result = append(result, strconv.Itoa(i))
	}
	crossing := []string{"L3,2", "L4,3", "S4,3", "L3,2", "S3,2", "L2,1", "S2,1", "L1,0", "S1,0", "L5,4"}
	for i := 4; i < 52; i++ {
		crossing = append(crossing, fmt.Sprintf("L%d,%d", i, i-1))
	}
	fixed := []fixture{
		{0, 0, 0, nil, false, "false::"},
		{0, 0, 1, []int{4}, false, "true:4:"},
		{0, 1, 3, []int{9, 2, 1, 8}, false, "false:9,2,1,8:L2,1"},
		{0, 2, 52, values, true, "true:" + strings.Join(result, ",") + ":" + strings.Join(crossing, ";")},
	}
	var cases []fixture
	var input bytes.Buffer
	for _, f := range fixed {
		for variant := range 2 {
			f.Variant = variant
			cases = append(cases, f)
			fmt.Fprintf(&input, "%d %d %d", variant, f.A, f.B)
			for _, v := range f.Values {
				fmt.Fprintf(&input, " %d", v)
			}
			fmt.Fprintln(&input)
		}
	}
	rows := check.LeanRows("godiff-sort-partial-insertion", len(cases), &input)
	file, err := os.CreateTemp("", "golean-sort-partial-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SORT_PARTIAL_CASES", file.Name()))
	must(check.TestInPackage("sort", "TestGoLeanPartialInsertion", "difftest/sort/partial-insertion-sort/testdata/partial_insertion_test.go"))
	for variant := range 2 {
		// A generic sortedness check incorrectly calls the empty interval
		// successful. Bounding the left shift at a suppresses the witness
		// swaps S2,1 and S1,0 and preserves the original prefix.
		if strings.HasPrefix(rows[variant], "true:") ||
			!strings.HasPrefix(rows[6+variant], "true:-1,100,101,0,1,") ||
			!strings.Contains(rows[6+variant], ";S2,1;L1,0;S1,0;") {
			log.Fatal("mutation witness missing")
		}
	}
	fmt.Println("Eight independent native/Lean traces matched, including four variant-specific mutation witnesses.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
