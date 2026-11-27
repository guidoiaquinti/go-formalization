package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	version     uint16
	attrs, mode uint32
	flags       uint16
	slash       bool
}

func main() {
	f, err := os.CreateTemp("", "golean-zip-file-mode")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	output := bufio.NewWriter(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			slash := 0
			if c.slash {
				slash = 1
			}
			fmt.Fprintln(&input, c.version, c.attrs, c.mode, c.flags, slash)
		}
		lines := bytes.Split(bytes.TrimSuffix(input.Bytes(), []byte("\n")), []byte("\n"))
		for i, row := range check.LeanRows("godiff-zip-file-mode", len(cases), &input) {
			_, err := fmt.Fprintf(output, "%s %s\n", lines[i], row)
			must(err)
		}
		fmt.Printf("%s: %d Lean rows emitted\n", label, len(cases))
	}
	fixed := make([]fixture, 0, 139264)
	for n := range 65536 {
		fixed = append(fixed, fixture{uint16(n), uint32(n)<<16 | uint32(n), uint32(n), uint16(n), false})
		fixed = append(fixed, fixture{uint16(n), ^uint32(n), uint32(n)<<16 | 0x1ff, ^uint16(n), true})
	}
	for typ := range 128 {
		for mask := range 64 {
			mode := uint32(0)
			for index, bit := range []uint{31, 27, 26, 25, 24, 21, 19, 23, 22, 20, 8, 7, 6} {
				if (typ|(mask<<7))&(1<<index) != 0 {
					mode |= 1 << bit
				}
			}
			fixed = append(fixed, fixture{0x03fe, uint32(typ<<9 | mask), mode, 8, false})
		}
	}
	compare("exhaustive word projections and explicit type combinations", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			version := uint16(rng.Uint32())
			if i%2 == 0 {
				creators := [...]uint16{0, 3, 11, 14, 19}
				version = version&255 | creators[i%len(creators)]<<8
			}
			cases[i] = fixture{version, rng.Uint32(), rng.Uint32(), uint16(rng.Uint32()), rng.Intn(2) == 1}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(output.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_ZIP_MODE_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanFileMode", "difftest/archive/zip/file-mode/testdata/file_mode_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
