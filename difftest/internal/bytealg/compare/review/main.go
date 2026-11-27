package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	type fixture struct {
		Left, Right                    []byte
		LeftOffset, RightOffset, Order int
		Alias                          bool
	}
	prefix := bytes.Repeat([]byte{0x80}, 64)
	cases := []fixture{
		{nil, []byte{}, 0, 15, 0, false},
		{[]byte{255}, []byte{128}, 31, 0, 1, false},
		{[]byte{1, 0, 255}, []byte{1, 0}, 1, 2, 1, false},
		{[]byte{0x7f}, []byte{0x80}, 7, 6, -1, false},
		{[]byte("é"), []byte{0xc3}, 0, 63, 1, false},
		{[]byte{1, 255}, []byte{2, 0, 0}, 7, 2, -1, false},
		{prefix, prefix, 15, 33, 0, false},
		{prefix, prefix, 63, 63, 0, true},
		{append(bytes.Clone(prefix), 0), append(bytes.Clone(prefix), 255), 1, 47, -1, false},
		{append(bytes.Clone(prefix), 255), prefix, 47, 1, 1, false},
	}
	f, err := os.CreateTemp("", "golean-review-bytealg-compare-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var model bytes.Buffer
	width := func(data []byte) int {
		if data == nil {
			return -1
		}
		return len(data)
	}
	for _, c := range cases {
		must(enc.Encode(c))
		fmt.Fprintf(&model, "%d %d", width(c.Left), width(c.Right))
		for _, part := range [][]byte{c.Left, c.Right} {
			for _, b := range part {
				fmt.Fprintf(&model, " %d", b)
			}
		}
		fmt.Fprintln(&model)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_BYTEALG_COMPARE_CASES", f.Name()))
	const dir = "difftest/internal/bytealg/compare/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgCompare", dir+"bridge_test.go", dir+"compare_test.go"))
	for i, got := range check.LeanRows("godiff-internal-bytealg-compare", len(cases), &model) {
		c := cases[i]
		want := fmt.Sprintf("%d,%d,%d,%t", c.Order, c.Order, c.Order, c.Order == 0)
		if got != want {
			log.Fatalf("row %d: got %q want %q", i, got, want)
		}
	}
	fmt.Printf("bytealg comparison passed %d fixed inputs through four native and Lean entries\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
