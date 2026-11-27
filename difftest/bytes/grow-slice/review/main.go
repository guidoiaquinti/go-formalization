package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Raw             []byte
	Capacity, Count int
	NilInput        bool
}

type fixture struct {
	Input input
	Want  string
}

func main() {
	raw := []byte{0, 255, 77}
	cases := []fixture{
		{input{NilInput: true}, "nil"},
		{input{}, "nil"},
		{input{NilInput: true, Count: math.MinInt}, "nil"},
		{input{Capacity: 256}, result(nil, 512)},
		{input{Raw: raw, Capacity: 4, Count: 5}, result(raw, 8)},
		{input{Raw: raw, Capacity: 4, Count: 6}, result(raw, 16)},
		{input{Raw: raw, Capacity: 4, Count: math.MaxInt}, result(raw, 8)},
		{input{Raw: raw, Capacity: 6, Count: math.MinInt}, result(raw, 16)},
		{input{Raw: raw, Capacity: 6, Count: -5}, result(raw, 16)},
		{input{NilInput: true, Count: 1}, result(nil, 8)},
		{input{Raw: raw, Capacity: 16, Count: 509}, result(raw, 512)},
	}
	var model bytes.Buffer
	for _, row := range cases {
		c := row.Input
		fmt.Fprintf(&model, "%d %d %d", c.Capacity, c.Count, math.MaxInt)
		for _, b := range c.Raw {
			fmt.Fprintf(&model, " %d", b)
		}
		fmt.Fprintln(&model)
	}
	for i, got := range check.LeanRows("godiff-bytes-grow-slice", len(cases), &model) {
		if got != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%q want=%q", i, got, cases[i].Want)
		}
	}
	if cases[3].Want == "nil" || cases[5].Want == result(raw, 9) ||
		cases[6].Want == "too-large" || strings.HasSuffix(cases[5].Want, "153") {
		log.Fatal("one incorrect variant escaped")
	}
	f, err := os.CreateTemp("", "golean-grow-slice-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_GROW_SLICE_CASES", f.Name()))
	const data = "difftest/bytes/grow-slice/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanGrowSlice", data+"bridge_test.go", data+"grow_slice_test.go"))
	fmt.Printf("growSlice: %d pinned native and Lean fixtures passed\n", len(cases))
}

func result(raw []byte, capacity int) string {
	wire := func(data []byte) string {
		parts := make([]string, len(data))
		for i, b := range data {
			parts[i] = strconv.Itoa(int(b))
		}
		return strings.Join(parts, ",")
	}
	return fmt.Sprintf("%d:%s:%s", capacity, wire(raw), wire(make([]byte, capacity-len(raw))))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
