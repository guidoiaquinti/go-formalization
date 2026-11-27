package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	type fixture struct {
		Bytes        []byte
		Target       byte
		Offset, Want int
	}
	cases := []fixture{
		{nil, 0, -1, 0}, {[]byte{}, 255, 15, 0},
		{[]byte{0xff, 0, 0xff, 0x80}, 0xff, 31, 2},
		{[]byte("éf\xc3"), 0xc3, 1, 2},
		{bytes.Repeat([]byte{0xff}, 31), 0xff, 3, 31},
		{bytes.Repeat([]byte{0x80}, 32), 0x80, 15, 32},
		{append(bytes.Repeat([]byte{0xff}, 32), 0x80), 0x80, 31, 1},
		{append([]byte{0}, bytes.Repeat([]byte{0xff}, 63)...), 0, 31, 1},
		{append(bytes.Repeat([]byte{0xff}, 64), 0), 0xff, 33, 64},
		{bytes.Repeat([]byte{0x80}, 129), 0x81, 63, 0},
	}
	f, err := os.CreateTemp("", "golean-review-bytealg-count-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var model bytes.Buffer
	for _, c := range cases {
		must(enc.Encode(c))
		fmt.Fprint(&model, c.Target)
		for _, b := range c.Bytes {
			fmt.Fprintf(&model, " %d", b)
		}
		fmt.Fprintln(&model)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_BYTEALG_COUNT_CASES", f.Name()))
	const dir = "difftest/internal/bytealg/count/testdata/"
	must(check.TestInPackage("internal/bytealg", "TestGoLeanBytealgCount", dir+"bridge_test.go", dir+"count_test.go"))
	for i, got := range check.LeanRows("godiff-internal-bytealg-count", len(cases), &model) {
		want := strings.TrimPrefix(strings.Repeat(fmt.Sprintf(",%d", cases[i].Want), 4), ",")
		if got != want {
			log.Fatalf("row %d got %q want %q", i, got, want)
		}
	}
	fmt.Printf("bytealg count passed %d fixed inputs through four pinned Go and Lean entry points\n", len(cases))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
