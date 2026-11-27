package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Old, Preset                 []byte
	Length, Read, Written, Size int
	Full                        bool
	Want                        string
}

func expected(c fixture) string {
	wire := func(cells []byte) string {
		fields := make([]string, len(cells))
		for i, b := range cells {
			fields[i] = strconv.Itoa(int(b))
		}
		return strings.Join(fields, ",")
	}
	old := slices.Clone(c.Old)
	result := old
	reused := len(old) >= c.Size
	if !reused {
		result = make([]byte, c.Size)
	}
	keep := min(len(c.Preset), c.Size)
	for i := range keep {
		result[i] = c.Preset[len(c.Preset)-keep+i]
	}
	full := len(c.Preset) >= c.Size
	cursor := keep
	if full {
		cursor = 0
	}
	return fmt.Sprintf("%t|%d|%d|%d|%d|%t|%t|%s|%s", result != nil, c.Size, len(result), cursor, cursor,
		full, reused, wire(result), wire(old))
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-dict-init")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			present, full := 0, 0
			if c.Old != nil {
				present = 1
			}
			if c.Full {
				full = 1
			}
			fmt.Fprint(&input, present, " ", c.Length, " ", len(c.Old), " ", c.Read, " ", c.Written, " ", full, " ", c.Size)
			for _, p := range [][]byte{c.Old, c.Preset} {
				for _, b := range p {
					fmt.Fprint(&input, " ", b)
				}
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-flate-dictionary-init", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for _, old := range [][]byte{nil, {}, {41}, {42, 43, 44, 45}, {1, 2, 3, 4, 5, 6, 7, 8}} {
		for _, size := range []int{0, 1, 2, 3, 4, 7, 8, 9, 16} {
			for _, presetSize := range []int{0, 1, 2, 3, 4, 7, 8, 9, 16, 33} {
				preset := make([]byte, presetSize)
				for i := range preset {
					preset[i] = byte(128 + i)
				}
				for _, full := range []bool{false, true} {
					fixed = append(fixed, fixture{Old: old, Preset: preset, Length: len(old) / 2,
						Written: len(old) / 2, Size: size, Full: full})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Old: make([]byte, r.Intn(65)), Preset: make([]byte, r.Intn(66)), Size: r.Intn(66), Full: r.Intn(2) == 1}
			c.Length = r.Intn(len(c.Old) + 1)
			c.Written = r.Intn(c.Length + 1)
			c.Read = r.Intn(c.Written + 1)
			r.Read(c.Old)
			r.Read(c.Preset)
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_DICT_INIT", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanDictionaryInit", "difftest/compress/flate/dictionary-init/testdata/init_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
