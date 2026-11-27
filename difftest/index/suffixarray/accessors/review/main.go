package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Tag, Len32, Len64, ByteLen, I, J int
	Value                            int64
	Words32, Words64                 []int64
	Bytes                            []byte
	Want                             string
}

func main() {
	raw := []byte{10, 11, 12, 13}
	cases := []fixture{
		{7, 1, 2, 2, 0, 3, 4294967294, []int64{-1, 2, 3}, []int64{77, 88, 99}, raw,
			"ok:3|ok:-1|ok:1:-2,2,3/2:77,88,99|ok:3:-1,2,3/nil|ok:2:10,11,12,13|bounds"},
		{7, 0, 1, 2, 0, 1, 9, []int64{3, 4}, []int64{1}, raw,
			"ok:1|bounds|bounds|ok:1:3,4/nil|ok:2:10,11,12,13|bounds"},
		{6, 0, 1, 2, 0, 1, -6, nil, []int64{2, 6}, raw,
			"ok:1|ok:2|ok:nil/1:-6,6|ok:nil/1:2,6|ok:2:10,11,12,13|ok:0:12,13"},
		{6, 0, 1, 2, 0, 1, -6, nil, []int64{3, 6}, raw,
			"ok:1|ok:3|ok:nil/1:-6,6|ok:nil/1:3,6|ok:2:10,11,12,13|bounds"},
		{8, 0, 0, 0, 0, 0, 0, nil, nil, nil, "nilReceiver|nilReceiver|nilReceiver|nilReceiver|nilReceiver|nilReceiver"},
		{0, 0, 0, 0, 0, 0, 0, nil, nil, nil, "ok:0|bounds|bounds|ok:nil/nil|ok:nil|bounds"},
	}
	reject := func(row, observation int, mutated string) {
		if strings.Split(cases[row].Want, "|")[observation] == mutated {
			panic(fmt.Sprintf("mutation survived on row %d observation %d", row, observation))
		}
	}
	first := cases[0]
	reject(0, 0, fmt.Sprintf("ok:%d", first.Len32))
	reject(0, 1, fmt.Sprintf("ok:%d", uint32(first.Words32[0])))
	reject(0, 2, fmt.Sprintf("ok:1:%d,2,3/2:77,88,99", first.Value))
	reject(0, 3, "bounds")
	reject(0, 4, "ok:2:10,11")
	reject(3, 5, "ok:1:13")
	reject(1, 1, fmt.Sprintf("ok:%d", cases[1].Words64[0]))
	fmt.Println("rejected one-limb length, zero extension, full-width store, length-clamped slice, clipped byte cap, and capacity-based at mutations")

	file, err := os.CreateTemp("", "review-suffixarray-accessors-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	var rows bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&rows, c.Tag, " ", c.Len32, " ", len(c.Words32), " ", c.Len64, " ", len(c.Words64), " ", c.ByteLen, " ", c.I, " ", c.J, " ", c.Value)
		for _, words := range [][]int64{c.Words32, c.Words64} {
			for _, word := range words {
				fmt.Fprint(&rows, " ", word)
			}
		}
		for _, b := range c.Bytes {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
		must(encoder.Encode(c))
	}
	must(file.Close())
	for i, got := range check.LeanRows("godiff-suffixarray-accessors", len(cases), &rows) {
		if got != cases[i].Want {
			panic(fmt.Sprintf("fixture %d got %q want %q", i, got, cases[i].Want))
		}
	}
	must(os.Setenv("GOLEAN_SUFFIXARRAY_ACCESSORS_CASES", file.Name()))
	must(check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayAccessors", "difftest/index/suffixarray/accessors/testdata/accessors_test.go"))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
