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
	Backing                                []byte
	Length, Read, Written, Distance, Count int
	Full                                   bool
	Want                                   string
}

func expected(c fixture) string {
	period := make([]byte, c.Distance)
	for i := range period {
		source := c.Written - c.Distance + i
		if source < 0 {
			source += c.Length
		}
		period[i] = c.Backing[source]
	}
	out := func(n int) string {
		b := slices.Clone(c.Backing)
		for i := range n {
			b[c.Written+i] = period[i%len(period)]
		}
		fields := make([]string, len(b))
		for i, v := range b {
			fields[i] = strconv.Itoa(int(v))
		}
		return fmt.Sprintf("%d:%d:%d:%t:%s", n, c.Read, c.Written+n, c.Full, strings.Join(fields, ","))
	}
	available := c.Length - c.Written
	fast := c.Count
	if c.Written < c.Distance || c.Count > available {
		fast = 0
	}
	return out(min(c.Count, available)) + "|" + out(fast)
}

func main() {
	f, err := os.CreateTemp("", "golean-flate-dict-copy")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var in bytes.Buffer
		for _, c := range cases {
			full := 0
			if c.Full {
				full = 1
			}
			fmt.Fprint(&in, c.Length, " ", c.Read, " ", c.Written, " ", full, " ", c.Distance, " ", c.Count)
			for _, b := range c.Backing {
				fmt.Fprint(&in, " ", b)
			}
			fmt.Fprintln(&in)
		}
		for i, row := range check.LeanRows("godiff-flate-dictionary-copy", len(cases), &in) {
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
	for length := 1; length <= 12; length++ {
		for w := 0; w <= length; w++ {
			for _, full := range []bool{false, true} {
				history := w
				if full {
					history = length
				}
				for dist := 1; dist <= history; dist++ {
					for _, count := range []int{0, 1, dist, dist + 1, length - w, length - w + 1, 2*length + 1} {
						c := fixture{Backing: make([]byte, length+3), Length: length, Read: w / 2, Written: w, Distance: dist, Count: count, Full: full}
						for i := range c.Backing {
							c.Backing[i] = byte(i + 11)
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
			length := 1 + r.Intn(65)
			c := fixture{Length: length, Backing: make([]byte, length+r.Intn(7)), Written: r.Intn(length + 1), Full: r.Intn(2) == 1}
			if c.Written == 0 {
				c.Full = true
			}
			history := c.Written
			if c.Full {
				history = length
			}
			c.Read = r.Intn(c.Written + 1)
			c.Distance = 1 + r.Intn(history)
			c.Count = r.Intn(2*length + 10)
			r.Read(c.Backing)
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_FLATE_DICT_COPY", f.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanDictionaryCopy", "difftest/compress/flate/dictionary-copy/testdata/copy_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
