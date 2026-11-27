package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Flags      uint8
	Name, Want string
}

func main() {
	const duplicateTail = "58,32,100,117,112,108,105,99,97,116,101,32,101,110,116,114,105,101,115,32,105,110,32,122,105,112,32,102,105,108,101"
	duplicate := "duplicate:255,0," + duplicateTail
	cases := []fixture{
		{11, "\xff\x00", duplicate + "|" + duplicate + "|nil"},
		{10, "entry", "directory:17|directory:17|nil"},
		{12, "regular", "header:29|header:29|header:29"},
		{4, "missing entry", "header:29|nil|header:29"},
		{8, "missing file", "nil|nil|nil"},
	}
	var input bytes.Buffer
	f, err := os.CreateTemp("", "golean-zip-entry-review-")
	must(err)
	defer os.Remove(f.Name())
	enc := gob.NewEncoder(f)
	for _, c := range cases {
		for i := range 4 {
			fmt.Fprint(&input, (c.Flags>>i)&1)
			if i < 3 {
				fmt.Fprint(&input, " ")
			}
		}
		for _, b := range []byte(c.Name) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(f.Close())
	var rejected [4]int
	for i, got := range check.LeanRows("godiff-zip-entry-stat", len(cases), &input) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d: Lean %q want %q", i, got, c.Want)
		}
		parts := strings.Split(got, "|")
		if c.Flags&3 == 3 && parts[0] != "directory:17" {
			rejected[0]++
		}
		if c.Flags&8 == 0 && parts[1] != "header:29" {
			rejected[1]++
		}
		if parts[2] != "sticky-read-error" {
			rejected[2]++
		}
		if c.Flags&3 == 3 && parts[0] != "duplicate:"+duplicateTail {
			rejected[3]++
		}
	}
	if want := [4]int{1, 1, 5, 1}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_ZIP_ENTRY_STAT_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanEntryStat", "difftest/archive/zip/entry-stat/testdata/stat_test.go"))
	fmt.Printf("ZIP entry stat: %d pinned Go and Lean rows, mutation counts %v\n", len(cases), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
