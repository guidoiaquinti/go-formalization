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

type input struct {
	NilReceiver, NilStorage bool
	Length, Spare, Offset   int
}

type fixture struct {
	Input input
	Want  string
}

func main() {
	cases := []fixture{
		{input{NilReceiver: true}, "panic"},
		{input{NilStorage: true}, "true"},
		{input{NilStorage: true, Offset: -1}, "false"},
		{input{NilStorage: true, Offset: math.MinInt}, "false"},
		{input{Spare: 8}, "true"},
		{input{Length: 8, Spare: 8, Offset: 7}, "false"},
		{input{Length: 8, Spare: 8, Offset: 8}, "true"},
		{input{Length: 8, Spare: 8, Offset: 9}, "true"},
		{input{Length: 8, Offset: math.MaxInt}, "true"},
		{input{Length: 8, Spare: 8, Offset: math.MinInt}, "false"},
	}
	var model bytes.Buffer
	for _, f := range cases {
		if f.Input.NilReceiver {
			fmt.Fprintln(&model, 0)
		} else {
			fmt.Fprintln(&model, 1, f.Input.Length, f.Input.Offset)
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-buffer-empty", len(cases), &model) {
		if row != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, row, cases[i].Want)
		}
	}
	misses := [3]int{}
	for _, f := range cases[1:] {
		c := f.Input
		for variant, got := range []bool{c.Offset == c.Length, c.Offset >= c.Length+c.Spare, uint(c.Offset) >= uint(c.Length)} {
			if fmt.Sprint(got) != f.Want {
				misses[variant]++
			}
		}
	}
	if misses[0] == 0 || misses[1] == 0 || misses[2] == 0 {
		log.Fatalf("an incorrect variant survived: %v", misses)
	}
	file, err := os.CreateTemp("", "golean-empty-review-")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(cases))
	must(file.Close())
	must(os.Setenv("GOLEAN_BUFFER_EMPTY_CASES", file.Name()))
	const dir = "difftest/bytes/buffer-empty/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferEmpty", dir+"bridge_test.go", dir+"buffer_empty_test.go"))
	fmt.Printf("Buffer.empty: %d independent native and Lean cases, equality/capacity/unsigned variants fail %v\n", len(cases), misses)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
