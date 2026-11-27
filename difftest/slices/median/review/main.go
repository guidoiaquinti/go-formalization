package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                                []int64
	A, B, C, Center, Mode, FailAt, Count int
	Replies                              [3]int
	HasCounter                           bool
	Want                                 string
}

func main() {
	descending := []int64{90, 30, 20, 10, 80}
	cases := []fixture{
		{descending, 1, 2, 3, 2, 2, 2, 5, [3]int{-1, -1, -1}, true,
			"ok:2/8/|panic:2/7/20,30;10,30;10,20|ok:2/8/|panic:2/7/20,30;10,30;10,20"},
		{descending, 1, 2, 3, 2, 2, -1, 5, [3]int{-1, 1, 0}, true,
			"ok:2/8/|ok:1/6/20,30;10,30;30,20|ok:2/8/|ok:1/6/20,30;10,30;30,20"},
		{descending, 0, 3, 4, 1, 1, -1, math.MaxInt, [3]int{}, true,
			fmt.Sprintf("ok:4/%d/|ok:4/%d/10,90;80,90;80,10|ok:1/%d/|ok:1/%d/30,90;20,90;20,30",
				math.MinInt+1, math.MinInt+1, math.MinInt+2, math.MinInt+2)},
		{[]int64{7, 7, 7}, 0, 1, 2, 1, 1, -1, 4, [3]int{}, false,
			"ok:1/nil/|ok:1/nil/7,7;7,7;7,7|ok:1/nil/|ok:1/nil/7,7;7,7;7,7"},
		{descending, 1, 2, 3, 2, 2, 0, 5, [3]int{-1, -1, -1}, false,
			"nil/nil/|panic:0/nil/20,30|nil/nil/|panic:0/nil/20,30"},
	}
	panicCase := cases[0]
	priorCount := panicCase.Count
	for _, reply := range panicCase.Replies[:panicCase.FailAt] {
		if reply < 0 {
			priorCount++
		}
	}
	if priorCount != 7 || panicCase.Count == priorCount {
		panic("counter rollback mutation survived")
	}
	orderCase := cases[1]
	// With the third compare omitted, the first exchange followed by the
	// positive second reply yields the original left index and only two calls.
	withoutThird := fmt.Sprintf("ok:%d/%d/%d,%d;%d,%d", orderCase.A, orderCase.Count+1,
		orderCase.Words[orderCase.B], orderCase.Words[orderCase.A], orderCase.Words[orderCase.C], orderCase.Words[orderCase.A])
	if withoutThird == "ok:1/6/20,30;10,30;30,20" {
		panic("omitted third comparison survived")
	}
	forwardFirst := fmt.Sprintf("%d,%d", panicCase.Words[panicCase.A], panicCase.Words[panicCase.B])
	if forwardFirst == "20,30" {
		panic("forward argument order survived")
	}
	fmt.Printf("rejected rollback to %d, omitted third comparison %q and forward operands %q\n", panicCase.Count, withoutThird, forwardFirst)

	f, err := os.CreateTemp("", "review-median-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.HasCounter {
			present = 1
		}
		fmt.Fprint(&input, strconv.IntSize, " ", c.A, " ", c.B, " ", c.C, " ", c.Center, " ", present, " ", c.Count, " ", c.Mode, " ", c.FailAt, " ", c.Replies[0], " ", c.Replies[1], " ", c.Replies[2])
		for _, word := range c.Words {
			fmt.Fprint(&input, " ", word)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(f.Close())
	for i, row := range check.LeanRows("godiff-slices-median", len(cases), &input) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_MEDIAN_CASES", f.Name()))
	const data = "difftest/slices/median/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanMedian", data+"bridge_test.go", data+"median_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
