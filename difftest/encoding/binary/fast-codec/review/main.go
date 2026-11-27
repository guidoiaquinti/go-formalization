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

type fixture struct {
	Kind, Form, Order, Length int
	Words                     []uint64
	Backing                   []byte
	Want                      string
}

func main() {
	cases := []fixture{
		{2, 2, 0, 2, []uint64{17, 18, 19}, []byte{13, 14, 99}, "ok:13,14,19|ok:17,18,99"},
		{1, 2, 0, 2, []uint64{17, 18, 19}, []byte{13, 14, 99}, "ok:13,14,19|bounds:17,18,99"},
		{1, 2, 0, 2, []uint64{17}, []byte{255, 14, 99}, "bounds:255|ok:17,14,99"},
		{0, 2, 0, 2, []uint64{0}, []byte{128, 14, 99}, "bounds:1|ok:0,14,99"},
		{4, 2, 1, 3, []uint64{257, 514}, []byte{4, 5, 6, 77}, "bounds:1029,514|bounds:1,1,6,77"},
		{9, 1, 1, 4, []uint64{0x80000000}, []byte{0x7f, 0x80, 0, 1, 77}, "ok:2139095041|ok:128,0,0,0,77"},
		{10, 1, 0, 8, []uint64{0x7ff0000000000001}, []byte{0, 0, 0, 0, 0, 0, 0, 0x80, 77}, "ok:9223372036854775808|ok:1,0,0,0,0,0,240,127,77"},
		{4, 0, 1, 2, []uint64{257}, []byte{4, 5, 77}, "unsupported:257|ok:1,1,77"},
		{11, 3, 0, 2, []uint64{257}, []byte{4, 5, 77}, "unsupported:257|ok:4,5,77"},
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprintf(&wire, "%d %d %d %d %d", c.Kind, c.Form, c.Order, len(c.Words), c.Length)
		for _, w := range c.Words {
			fmt.Fprintf(&wire, " %d", w)
		}
		for _, b := range c.Backing {
			fmt.Fprintf(&wire, " %d", b)
		}
		fmt.Fprintln(&wire)
	}
	var rejects [4]int
	for i, got := range check.LeanRows("godiff-binary-fast-codec", len(cases), &wire) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, c.Want)
		}
		// Each variant makes a distinct prediction on the relevant control row.
		switch i {
		case 0:
			if got != "ok:13,14,19|bounds:17,18,99" {
				rejects[0]++
			}
		case 2, 3:
			if !strings.HasPrefix(got, "ok:") {
				rejects[1]++
			}
		case 4:
			if got != "bounds:257,514|bounds:4,5,6,77" {
				rejects[2]++
			}
		case 5, 6:
			if got != []string{"ok:2143289344|ok:0,0,0,0,77", "ok:0|ok:0,0,0,0,0,0,248,127,77"}[i-5] {
				rejects[3]++
			}
		}
	}
	if rejects != [4]int{1, 2, 1, 2} {
		log.Fatalf("unexpected mutant rejections %v", rejects)
	}
	f, err := os.CreateTemp("", "golean-fast-codec-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_BINARY_FAST_CODEC_CASES", f.Name()))
	const base = "difftest/encoding/binary/fast-codec/testdata/"
	must(check.TestInPackage("encoding/binary", "TestGoLeanFastCodec", base+"bridge_test.go", base+"fast_codec_test.go"))
	fmt.Printf("Fast codec: %d independent native and Lean fixtures, mutant rejections %v\n", len(cases), rejects)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
