package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Words                                 []int64
	Kind, A, B, Pivot, Mode, FailAt, Mask int
	Want                                  string
}

func main() {
	cases := []fixture{
		{[]int64{91, 4, 1, 2, 8, 92}, 0, 1, 5, 3, 1, -1, 0,
			"ok:2,true/91,1,2,4,8,92/|ok:2,true/91,1,2,4,8,92/1,2;4,2;8,2;4,2"},
		{[]int64{91, 4, 8, 2, 1, 92}, 0, 1, 5, 3, 1, 2, 0,
			"ok:2,false/91,1,2,4,8,92/|panic:2/91,2,1,4,8,92/8,2;1,2;4,2"},
		{[]int64{91, 4, 8, 2, 2, 92}, 1, 1, 5, 3, 1, 2, 0,
			"ok:3/91,2,2,4,8,92/|panic:2/91,2,2,4,8,92/2,8;2,2;2,4"},
		{[]int64{91, 2, 1, 3, 92}, 1, 1, 4, 1, 1, -1, 0,
			"ok:3/91,2,1,3,92/|ok:3/91,2,1,3,92/2,1;2,3;2,3"},
		{[]int64{91, 4, 92}, 0, 1, 2, 1, 0, -1, 0, "ok:1,true/91,4,92/|ok:1,true/91,4,92/"},
		{[]int64{91, 4, 92}, 1, 1, 2, 1, 0, -1, 0, "ok:2/91,4,92/|ok:2/91,4,92/"},
	}
	flag := cases[0]
	if flag.A == flag.Pivot || flag.A == 2 {
		panic("miscounting pivot movement as a cross-swap survived")
	}
	panicCase := cases[1]
	withoutCross := slices.Clone(panicCase.Words)
	withoutCross[panicCase.A], withoutCross[panicCase.Pivot] = withoutCross[panicCase.Pivot], withoutCross[panicCase.A]
	if slices.Equal(withoutCross, []int64{91, 2, 1, 4, 8, 92}) {
		panic("cross-swap rollback mutation survived")
	}
	equal := cases[2]
	backwards := fmt.Sprintf("%d,%d", equal.Words[equal.A+1], equal.Words[equal.Pivot])
	if backwards == "2,8" {
		panic("equal-partition reversed operands survived")
	}
	fmt.Printf("rejected pivot-only flag change, cross-swap rollback %v, and reversed first equal operands %s\n", withoutCross, backwards)

	f, err := os.CreateTemp("", "review-partition-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.Kind, " ", c.A, " ", c.B, " ", c.Pivot, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
		for _, word := range c.Words {
			fmt.Fprint(&input, " ", word)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(f.Close())
	for i, row := range check.LeanRows("godiff-slices-partition", len(cases), &input) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_PARTITION_CASES", f.Name()))
	const data = "difftest/slices/partition/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanPartition", data+"bridge_test.go", data+"partition_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
