package main

import (
	"bytes"
	"fmt"
	"log"
	"math"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := [][2]int64{{math.MinInt64, -1152921504606846974}, {-25, -1}, {-24, -1}, {-23, 0},
		{-17, 0}, {-16, 0}, {-9, 0}, {-8, 1}, {0, 2}, {7, 2}, {8, 3},
		{math.MaxInt64 - 16, 1152921504606846975}, {math.MaxInt64 - 15, -1152921504606846976},
		{math.MaxInt64, -1152921504606846974}}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, c[0])
	}
	for i, row := range check.LeanRows("godiff-internal-bytealg-index-amd64-profile", len(cases), &input) {
		want := fmt.Sprintf("%d,31,63", cases[i][1])
		if row != want {
			log.Fatalf("count %d: Lean %s, want %s", cases[i][0], row, want)
		}
	}
	fmt.Printf("%d signed cutover cases and both initializer branches passed.\n", len(cases))
}
