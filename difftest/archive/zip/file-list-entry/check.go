package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Name, Other                   string
	Present, Directory, Duplicate bool
	Seconds                       int64
	Nanos, Offset                 int
	Want                          string
}

func main() {
	f, err := os.CreateTemp("", "golean-zip-file-list")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			present := 0
			if c.Present {
				present = 1
			}
			fmt.Fprint(&input, len(c.Name), " ", present, " ", c.Seconds, " ", c.Nanos)
			for _, raw := range []string{c.Name, c.Other} {
				for i := range len(raw) {
					fmt.Fprint(&input, " ", raw[i])
				}
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-zip-file-list-entry", len(cases), &input) {
			c := cases[i]
			c.Want = row
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean metadata rows\n", label, len(cases))
	}
	names := []string{"", "/", "//", "///", ".", "./", "../", "a", "a/", "a//", "a/b", "a/z", "b/a", "/x/", "a\\b", "a/\x00\xff/", "\x80", "\x7f", "a/."}
	var fixed []fixture
	for i, a := range names {
		for j, b := range names {
			for flags := range 8 {
				fixed = append(fixed, fixture{Name: a, Other: b, Present: flags&1 != 0, Directory: flags&2 != 0, Duplicate: flags&4 != 0,
					Seconds: []int64{-62135596800, -1, 0, 946684800, 1719792000, 13569465599}[(i+j)%6],
					Nanos:   []int{0, 1, 987654321}[(i+2*j)%3], Offset: []int{-18000, 0, 19800}[(2*i+j)%3]})
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				a, b := randomName(r), randomName(r)
				switch i % 4 {
				case 0:
					b = a + "/"
				case 1:
					a, b = "common/"+a, "common/"+b
				case 2:
					b = a
				}
				flags := r.Intn(8)
				cases[i] = fixture{Name: a, Other: b, Present: flags&1 != 0, Directory: flags&2 != 0, Duplicate: flags&4 != 0,
					Seconds: r.Int63n(100000000000) - 65000000000, Nanos: r.Intn(1000000000), Offset: r.Intn(100801) - 50400}
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_FILE_LIST_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanZipFileList", "difftest/archive/zip/file-list-entry/testdata/entry_test.go"))
}

func randomName(r *rand.Rand) string {
	raw := make([]byte, r.Intn(40))
	_, err := r.Read(raw)
	must(err)
	for i := range raw {
		if r.Intn(5) == 0 {
			raw[i] = '/'
		}
	}
	return string(raw)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
