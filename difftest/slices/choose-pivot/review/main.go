package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                    []int64
	A, B, Mode, FailAt, Mask int
	Want                     string
}

func main() {
	ascending, descending := make([]int64, 52), make([]int64, 52)
	for i := range ascending {
		ascending[i], descending[i] = int64(i), int64(len(ascending)-i)
	}
	cases := []fixture{
		{ascending, 3, 10, 0, -1, 0, "ok:5,1/|ok:5,1/"},
		{ascending, 3, 11, 1, -1, 0, "ok:7,1/|ok:7,1/7,5;9,7;7,5"},
		{ascending, 1, 50, 1, -1, 0, "ok:25,1/|ok:25,1/25,13;37,25;25,13"},
		{ascending, 1, 51, 2, 3, 0, "ok:25,1/|panic:3/13,12;14,13;13,12;25,24"},
		{ascending, 1, 51, 1, -1, 0, "ok:25,1/|ok:25,1/13,12;14,13;13,12;25,24;26,25;25,24;37,36;38,37;37,36;25,13;37,25;25,13"},
		{descending, 1, 51, 0, -1, 0, "ok:25,2/|nil/"},
		{descending, 1, 50, 0, -1, 0, "ok:25,0/|nil/"},
	}

	short := cases[0]
	midpoint := short.A + (short.B-short.A)/2
	if midpoint != 6 || midpoint == 5 {
		panic("ordinary midpoint mutation survived")
	}
	threshold := cases[3]
	quarter := (threshold.B - threshold.A) / 4
	medianOnlyFirstPair := fmt.Sprintf("%d,%d", threshold.Words[threshold.A+2*quarter], threshold.Words[threshold.A+quarter])
	if medianOnlyFirstPair != "25,13" || medianOnlyFirstPair == "13,12" {
		panic("ninther threshold mutation survived")
	}
	medianOnly := cases[6]
	a := medianOnly.A
	q := (medianOnly.B - a) / 4
	positions := []int{a + q, a + 2*q, a + 3*q}
	swaps := 0
	for _, pair := range [][2]int{{0, 1}, {1, 2}, {0, 1}} {
		x, y := pair[0], pair[1]
		if medianOnly.Words[positions[y]] < medianOnly.Words[positions[x]] {
			positions[x], positions[y] = positions[y], positions[x]
			swaps++
		}
	}
	badHint := 0
	if swaps >= 3 {
		badHint = 2
	}
	if swaps != 3 || fmt.Sprintf("ok:%d,%d/|nil/", positions[1], badHint) == medianOnly.Want {
		panic("descending-hint threshold mutation survived")
	}
	fmt.Printf("rejected ordinary midpoint %d, delayed ninther first pair %s and three-swap decreasing hint\n", midpoint, medianOnlyFirstPair)

	f, err := os.CreateTemp("", "review-choose-pivot-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, strconv.IntSize, " ", c.A, " ", c.B, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
		for _, word := range c.Words {
			fmt.Fprint(&input, " ", word)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(f.Close())
	for i, row := range check.LeanRows("godiff-slices-choose-pivot", len(cases), &input) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_PIVOT_CASES", f.Name()))
	const data = "difftest/slices/choose-pivot/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanChoosePivot", data+"bridge_test.go", data+"pivot_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
