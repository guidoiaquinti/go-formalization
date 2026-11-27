package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Profile, Limit int
	Text, Needle   []byte
	Offsets        []int
	Want           string
}

func main() {
	text, needle := []byte("banana"), []byte("ana")
	sa := []int{5, 3, 1, 0, 4, 2, 99}
	cases := []fixture{
		{1, -1, text, needle, sa, "2:3,1,0,4,2,99/nil|2:3,1"},
		{2, 1, text, needle, sa, "nil/2:3,1,0,4,2,99|1:3"},
		{1, 3, text, []byte("m"), sa, "0:4,2,99/nil|nil"},
		{1, 4, text, nil, sa, "6:5,3,1,0,4,2,99/nil|nil"},
		{0, 0, nil, needle, nil, "panic|nil"},
		{0, 1, nil, needle, nil, "panic|panic"},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&input, "%d %d %d %d %d", strconv.IntSize, c.Profile, len(c.Text), len(c.Needle), c.Limit)
		for _, values := range [][]byte{c.Text, c.Needle} {
			for _, b := range values {
				fmt.Fprintf(&input, " %d", b)
			}
		}
		for _, i := range c.Offsets {
			fmt.Fprintf(&input, " %d", i)
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-suffixarray-lookup", len(cases), &input)
	file, err := os.CreateTemp("", "golean-suffixarray-lookup-review")
	must(err)
	defer os.Remove(file.Name())
	enc := json.NewEncoder(file)
	for i, c := range cases {
		if rows[i] != c.Want {
			log.Fatalf("fixture %d: Lean %q want %q", i, rows[i], c.Want)
		}
		must(enc.Encode(c))
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_SUFFIXARRAY_LOOKUP_CASES", file.Name()))
	must(check.TestInPackage("index/suffixarray", "TestGoLeanSuffixarrayLookup",
		"difftest/index/suffixarray/lookup/testdata/lookup_test.go"))
	if strings.HasSuffix(rows[0], "|2:1,3") || strings.HasPrefix(rows[0], "2:3,1/nil") ||
		strings.HasSuffix(rows[1], "|2:3,1") || strings.HasSuffix(rows[4], "|panic") {
		log.Fatal("accepted sorted result, clipped capacity, ignored limit or eager receiver dereference")
	}
	fmt.Println("Six independent native/Lean fixtures passed; four mutation witnesses rejected.")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
