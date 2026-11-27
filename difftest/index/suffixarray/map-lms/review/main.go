package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Length, Count int
	Backing       []int64
	Want          string
}

func main() {
	cases := []fixture{
		{0, 9, []int64{7, 8}, "bounds:7,8|bounds:7,8"},
		{1, -1, []int64{1, 99}, "ok:0,99|ok:0,99"},
		{2, 0, []int64{6, 9, 99}, "ok:5,8,99|ok:5,8,99"},
		{6, 0, []int64{4, -1, 2, 3, 19, 21, 99}, "ok:4,-1,2,3,1,2,99|ok:4,-1,2,3,1,2,99"},
		{6, math.MaxInt, []int64{4, -1, 2, 3, 19, 21, 99}, "ok:4,-1,2,3,1,2,99|ok:4,-1,2,3,1,2,99"},
		{5, 2, []int64{2147483648, 0, 2147483647, 7, 8, 99}, "ok:-2147483648,0,2147483647,7,2147483646,99|ok:2147483648,0,2147483647,2147483647,2147483646,99"},
	}
	mutants := []string{
		"ok:7,8|ok:7,8",
		"ok:1,99|ok:1,99",
		"ok:6,5,99|ok:6,5,99",
		"ok:4,-1,2,3,19,21,99|ok:4,-1,2,3,19,21,99",
		"ok:4,-1,2,3,19,21,99|ok:4,-1,2,3,19,21,99",
		"ok:-2147483648,0,2147483647,2147483647,2147483646,99|ok:2147483648,0,2147483647,2147483647,2147483646,99",
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d", c.Length, c.Count)
		for _, v := range c.Backing {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	rejects := 0
	for i, got := range check.LeanRows("godiff-suffixarray-map-lms", len(cases), &wire) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, cases[i].Want)
		}
		if got != mutants[i] {
			rejects++
		}
	}
	if rejects != len(mutants) {
		log.Fatalf("rejected %d of %d changed outcomes", rejects, len(mutants))
	}
	f, err := os.CreateTemp("", "golean-map-lms-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_SUFFIXARRAY_MAP_LMS_CASES", f.Name()))
	must(check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayMapLms", "difftest/index/suffixarray/map-lms/testdata/map_test.go"))
	fmt.Printf("Suffix array map: %d paired independent rows, %d mutated results rejected\n", len(cases), rejects)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
