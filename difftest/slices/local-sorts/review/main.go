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
	Words                                []int64
	Kind, A, B, Root, Mode, FailAt, Mask int
	Want                                 string
}

func main() {
	cases := []fixture{
		{[]int64{91, 3, 2, 1, 92}, 0, 1, 4, 0, 1, 1, 0, "ok/91,1,2,3,92/|panic:1/91,2,3,1,92/2,3;1,3"},
		{[]int64{91, 2, 92}, 0, 1, 2, 0, 0, -1, 0, "ok/91,2,92/|ok/91,2,92/"},
		{[]int64{91, 0, 2, 3, 1, 92}, 1, 1, 5, 0, 1, 1, 0, "ok/91,3,2,0,1,92/|panic:1/91,0,2,3,1,92/2,3;0,3"},
		{[]int64{91, 0, 2, 92}, 1, 1, 3, 0, 1, -1, 0, "ok/91,2,0,92/|ok/91,2,0,92/0,2"},
		{[]int64{91, 4, 1, 3, 92}, 2, 1, 4, 0, 1, 2, 0, "ok/91,1,3,4,92/|panic:2/91,3,1,4,92/1,3;4,3;3,1"},
		{[]int64{91, 4, 1, 3, 92}, 2, 2, 2, 0, 0, -1, 0, "ok/91,4,1,3,92/|ok/91,4,1,3,92/"},
	}
	insertion := cases[0]
	// A transaction-style insertion loses the completed first swap.
	if slices.Equal(insertion.Words, []int64{91, 2, 3, 1, 92}) {
		panic("insertion rollback mutation survived")
	}
	sift := cases[2]
	rootFirst := fmt.Sprintf("%d,%d", sift.Words[sift.A], sift.Words[sift.A+1])
	if rootFirst == "2,3" {
		panic("root-first sift mutation survived")
	}
	heap := cases[4]
	drain := slices.Clone(heap.Words)
	// The first drain exchanges endpoints before reaching its next comparison.
	drain[heap.A], drain[heap.B-1] = drain[heap.B-1], drain[heap.A]
	if !slices.Equal(drain, []int64{91, 3, 1, 4, 92}) || slices.Equal(drain, heap.Words) {
		panic("drain comparison-before-exchange mutation survived")
	}
	fmt.Printf("rejected rollback, root-first %s, and a missing heap drain swap\n", rootFirst)

	f, err := os.CreateTemp("", "review-local-sorts-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.Kind, " ", c.A, " ", c.B, " ", c.Root, " ", c.Mode, " ", c.FailAt, " ", c.Mask)
		for _, word := range c.Words {
			fmt.Fprint(&input, " ", word)
		}
		fmt.Fprintln(&input)
		must(encoder.Encode(c))
	}
	must(f.Close())
	for i, row := range check.LeanRows("godiff-slices-local-sorts", len(cases), &input) {
		if row != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, row, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SLICES_LOCAL_SORT_CASES", f.Name()))
	const data = "difftest/slices/local-sorts/testdata/"
	must(check.TestInPackage("slices", "TestGoLeanLocalSorts", data+"bridge_test.go", data+"local_sorts_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
