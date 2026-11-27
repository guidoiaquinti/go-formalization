package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	creator       uint16
	attrs, legacy uint32
	wide          uint64
	present       bool
	name          string
}

func main() {
	file, err := os.CreateTemp("", "golean-zip-header-info-")
	must(err)
	defer os.Remove(file.Name())
	out := bufio.NewWriter(file)
	fixed := []fixture{{}}
	for _, name := range []string{"", "/", "//", "a", "a/b", "a/b/", "a//.", "a/..", "a\\b/", "a/\xff\x00/"} {
		for _, creator := range []uint16{0, 0x300, 0xb00, 0xe00, 0x1300, 0xff00} {
			for _, attrs := range []uint32{0, 1, 17, 0x44200000, 0xa7770011, math.MaxUint32} {
				for _, wide := range []uint64{0, 1, math.MaxUint32, 1 << 32, math.MaxInt64, 1 << 63, math.MaxUint64} {
					fixed = append(fixed, fixture{creator | 42, attrs, math.MaxUint32, wide, true, name})
				}
			}
		}
	}
	compare("fixed", fixed, out)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				name := make([]byte, rng.Intn(38))
				rng.Read(name)
				for j := range name {
					if rng.Intn(7) == 0 {
						name[j] = '/'
					}
				}
				wide := rng.Uint64()
				if i%3 == 0 {
					wide = 0
				}
				rows[i] = fixture{uint16(rng.Uint32()), rng.Uint32(), rng.Uint32(), wide, true, string(name)}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows, out)
		}
	}
	must(out.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_ZIP_HEADER_INFO", file.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanHeaderInfo", "difftest/archive/zip/header-info/testdata/header_info_test.go"))
}

func compare(label string, rows []fixture, out *bufio.Writer) {
	var input bytes.Buffer
	for _, row := range rows {
		present := 0
		if row.present {
			present = 1
		}
		fmt.Fprint(&input, row.creator, " ", row.attrs, " ", row.legacy, " ", row.wide, " ", present)
		for _, b := range []byte(row.name) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	for i, observed := range check.LeanRows("godiff-zip-header-info", len(rows), &input) {
		row := rows[i]
		_, err := fmt.Fprintf(out, "%d %d %d %d %t %q\t%s\n", row.creator, row.attrs, row.legacy, row.wide, row.present, row.name, observed)
		must(err)
	}
	fmt.Printf("ZIP header info %s: %d Lean rows\n", label, len(rows))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
