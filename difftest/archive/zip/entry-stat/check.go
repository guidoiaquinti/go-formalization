package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Flags uint8
	Name  string
	Want  string
}

func expected(flags uint8, name string) string {
	header := "nil"
	if flags&4 != 0 {
		header = "header:29"
	}
	entry := header
	switch flags & 3 {
	case 1, 3:
		message := name + ": duplicate entries in zip file"
		cells := make([]string, len(message))
		for i := range len(message) {
			cells[i] = strconv.Itoa(int(message[i]))
		}
		entry = "duplicate:" + strings.Join(cells, ",")
	case 2:
		entry = "directory:17"
	}
	dir := "nil"
	if flags&8 != 0 {
		dir = entry
	}
	return entry + "|" + dir + "|" + header
}

func main() {
	f, err := os.CreateTemp("", "golean-zip-entry-stat")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, c.Flags&1, " ", c.Flags>>1&1, " ", c.Flags>>2&1, " ", c.Flags>>3&1)
			for i := range len(c.Name) {
				fmt.Fprint(&input, " ", c.Name[i])
			}
			fmt.Fprintln(&input)
		}
		for i, got := range check.LeanRows("godiff-zip-entry-stat", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c.Flags, c.Name)
			if got != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, got, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d ZIP stat observations agree\n", label, len(cases))
	}
	names := []string{"", "a", "a/", "/", "//", "foo/bar", "foo//bar/", "a\\b", "\x00", "\xff", "\x80\x00", "ä", "日本", "\n", "a\t", " : "}
	fixed := make([]fixture, 0, len(names)*16)
	for _, name := range names {
		for flags := range 16 {
			fixed = append(fixed, fixture{Flags: uint8(flags), Name: name})
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				name := make([]byte, r.Intn(49))
				_, err := r.Read(name)
				must(err)
				cases[i] = fixture{Flags: uint8(r.Intn(16)), Name: string(name)}
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_ENTRY_STAT_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanEntryStat", "difftest/archive/zip/entry-stat/testdata/stat_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
