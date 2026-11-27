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
	Backing                      []byte
	Length, Read, Written, Count int
	Full                         bool
	Value                        byte
	Want                         string
}

func expected(c fixture) string {
	wire := func(b []byte) string {
		cells := make([]string, len(b))
		for i, v := range b {
			cells[i] = strconv.Itoa(int(v))
		}
		return strings.Join(cells, ",")
	}
	state := func(r, w int, full bool, b []byte) string {
		return fmt.Sprintf("%d:%d:%t:%s", r, w, full, wire(b))
	}
	view := func(from, to int) string {
		return fmt.Sprintf("%t:%d:%d:%s", c.Backing != nil, len(c.Backing)-from, from, wire(c.Backing[from:to]))
	}
	size := c.Written
	if c.Full {
		size = c.Length
	}
	stored := "bounds"
	if c.Written < c.Length {
		b := slices.Clone(c.Backing)
		b[c.Written] = c.Value
		stored = state(c.Read, c.Written+1, c.Full, b)
	}
	next := state(c.Written, c.Written, c.Full, c.Backing)
	if c.Written == c.Length {
		next = state(0, 0, true, c.Backing)
	}
	return strings.Join([]string{strconv.Itoa(size), strconv.Itoa(c.Written - c.Read),
		strconv.Itoa(c.Length - c.Written), view(c.Written, c.Length),
		state(c.Read, c.Written+c.Count, c.Full, c.Backing), stored, view(c.Read, c.Written), next}, "|")
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-dict")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var in bytes.Buffer
		for _, c := range cases {
			present, full := 0, 0
			if c.Backing != nil {
				present = 1
			}
			if c.Full {
				full = 1
			}
			fmt.Fprint(&in, present, " ", c.Length, " ", c.Read, " ", c.Written, " ", full, " ", c.Count, " ", c.Value)
			for _, b := range c.Backing {
				fmt.Fprint(&in, " ", b)
			}
			fmt.Fprintln(&in)
		}
		for i, row := range check.LeanRows("godiff-flate-dictionary-state", len(cases), &in) {
			c := cases[i]
			c.Want = expected(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	fixed := []fixture{{}}
	for capacity := range 9 {
		for length := 0; length <= capacity; length++ {
			for w := 0; w <= length; w++ {
				for read := 0; read <= w; read++ {
					for _, full := range []bool{false, true} {
						c := fixture{Backing: make([]byte, capacity), Length: length, Read: read, Written: w, Count: length - w, Full: full, Value: 255}
						for i := range c.Backing {
							c.Backing[i] = byte(i + 1)
						}
						fixed = append(fixed, c)
					}
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Backing: make([]byte, r.Intn(65)), Full: r.Intn(2) == 1, Value: byte(r.Uint32())}
			c.Length = r.Intn(len(c.Backing) + 1)
			c.Written = r.Intn(c.Length + 1)
			c.Read = r.Intn(c.Written + 1)
			c.Count = r.Intn(c.Length - c.Written + 1)
			r.Read(c.Backing)
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_DICT_STATE", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanDictState", "difftest/compress/flate/dictionary-state/testdata/dict_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
