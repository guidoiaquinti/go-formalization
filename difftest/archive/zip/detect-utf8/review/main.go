package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type observation struct {
	Raw            []byte
	Valid, Require bool
}

func main() {
	cases := []observation{
		{nil, true, false},
		{[]byte{0x20, 0x5b, 0x5d, 0x7d}, true, false},
		{[]byte{0x1f}, true, true},
		{[]byte{0x7e}, true, true},
		{[]byte{0x5c}, true, true},
		{[]byte{0x7f}, true, true},
		{[]byte{0xef, 0xbf, 0xbd}, true, true},
		{[]byte{0x5c, 0x80}, false, false},
		{[]byte{0x5c, 0xe2, 0x82}, false, false},
		{[]byte{0xed, 0xa0, 0x80}, false, false},
		{[]byte{0xf4, 0x8f, 0xbf, 0xbf}, true, true},
	}
	mutants := map[int]string{
		1: "true true", 2: "true false", 3: "true false", 4: "true false", 6: "false false",
		7: "false true", 8: "false true", 9: "true true", 10: "false false",
	}
	var wire bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&wire, 0)
		for _, v := range c.Raw {
			fmt.Fprintf(&wire, " %d", v)
		}
		fmt.Fprintln(&wire)
	}
	rejected := 0
	for i, got := range check.LeanRows("godiff-zip-detect-utf8", len(cases), &wire) {
		c := cases[i]
		want := fmt.Sprintf("%t %t", c.Valid, c.Require)
		if got != want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, want)
		}
		if alt, ok := mutants[i]; ok && alt != got {
			rejected++
		}
	}
	if rejected != len(mutants) {
		log.Fatalf("rejected %d of %d mutants", rejected, len(mutants))
	}
	f, err := os.CreateTemp("", "golean-zip-utf8-review-")
	must(err)
	defer os.Remove(f.Name())
	rows := json.NewEncoder(f)
	for _, c := range cases {
		must(rows.Encode(c))
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_DETECT_UTF8_CASES", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanDetectUTF8", "difftest/archive/zip/detect-utf8/testdata/detect_test.go"))
	fmt.Printf("ZIP UTF-8 detection: %d fixed native and Lean inputs, %d altered results rejected\n", len(cases), rejected)
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
