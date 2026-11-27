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
	Kind, Shape, Count, Order, Length, Offset int
	Words                                     []uint64
	Backing                                   []byte
	Want                                      string
}

func main() {
	cases := []fixture{
		{2, 3, 0, 0, 1, 0, []uint64{3}, []byte{5}, "true:3:5|false:1:3"},
		{2, 3, 1, 0, 3, 0, []uint64{3, 4}, []byte{5, 99, 98, 97, 77}, "false:3:5,4|false:3:3,0,0,97,77"},
		{4, 0, 0, 1, 0, 1, []uint64{258}, []byte{70, 3, 4, 77}, "true:3:772|true:3:70,1,2,77"},
		{2, 0, 0, 0, 0, 1, []uint64{3}, []byte{70, 5, 77}, "false:1:3|false:1:70,5,77"},
		{11, 0, 0, 0, 0, 2, []uint64{0x3f800000, 0x80000000}, []byte{77, 78, 1, 2, 3, 4}, "false:6:1065353216,2147483648|false:6:77,78,0,0,128,63"},
		{11, 0, 0, 1, 0, 1, []uint64{0x80000000, 0x3f800000}, []byte{77, 0, 128, 0, 7, 191, 128, 0, 0, 99}, "true:9:8388615,3212836864|true:9:77,128,0,0,0,63,128,0,0,99"},
		{13, 0, 0, 0, 0, 8, nil, []byte{77}, "true:8:|true:8:77"},
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d %d", c.Kind, c.Shape, c.Count, c.Order, c.Offset, c.Length)
		for _, v := range c.Words {
			fmt.Fprintf(&wire, " %d", v)
		}
		for _, v := range c.Backing {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	mutants := []string{
		"false:1:5|false:1:3",
		"false:3:5,4|false:3:3,99,98,97,77",
		"false:1:258|false:1:70,3,4,77",
		"true:2:5|true:2:70,3,77",
		"false:2:1065353216,2147483648|false:2:77,78,1,2,3,4",
		"true:9:8388615,3212836864|true:9:77,0,0,0,0,63,128,0,0,99",
	}
	rejected := 0
	for i, got := range check.LeanRows("godiff-binary-coder-value", len(cases), &wire) {
		if got != cases[i].Want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, cases[i].Want)
		}
		if i < len(mutants) && got != mutants[i] {
			rejected++
		}
	}
	if rejected != len(mutants) {
		log.Fatalf("rejected %d of %d variants", rejected, len(mutants))
	}
	f, err := os.CreateTemp("", "golean-coder-value-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_BINARY_CODER_VALUE_CASES", f.Name()))
	const path = "difftest/encoding/binary/coder-value/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanCoderValue", path+"bridge_test.go", path+"coder_value_test.go"))
	fmt.Printf("Coder values: %d independent native and Lean cases, %d distinct incorrect variants rejected\n", len(cases), rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
