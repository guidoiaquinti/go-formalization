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

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Bool, Nil   bool
	Tag, Length int
	Backing     []byte
	Raw         string
	Want        string
}

func main() {
	f, err := os.CreateTemp("", "golean-internal-strconv-bool")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := gob.NewEncoder(w)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			value := 0
			if c.Bool {
				value = 1
			}
			fmt.Fprint(&input, value, " ", c.Tag, " ", c.Length, " ", len(c.Backing))
			for _, raw := range [][]byte{c.Backing, []byte(c.Raw)} {
				for _, b := range raw {
					fmt.Fprint(&input, " ", b)
				}
			}
			fmt.Fprintln(&input)
		}
		for i, got := range check.LeanRows("godiff-internal-strconv-bool", len(cases), &input) {
			c := cases[i]
			c.Want = got
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean internal strconv observations\n", label, len(cases))
	}
	inputs := []string{"1", "t", "T", "true", "TRUE", "True", "0", "f", "F", "false", "FALSE", "False",
		"", "tRue", "FALSe", "yes", "2", "01", "+1", "true ", " false", "1\x00", "\xff", "fa\x00lse"}
	tags := []int{-2147483648, -1, 0, 1, 2, 3, 4, 5, 2147483647}
	if strconv.IntSize == 64 {
		var min64, max64 int64 = -9223372036854775808, 9223372036854775807
		tags = append(tags, int(min64), int(max64))
	}
	var fixed []fixture
	for j, raw := range inputs {
		for _, value := range []bool{false, true} {
			fixed = append(fixed, fixture{Bool: value, Nil: true, Raw: raw, Tag: tags[j%len(tags)]})
			for capacity := range 13 {
				backing := make([]byte, capacity)
				for i := range backing {
					backing[i] = byte(240 + i)
				}
				for length := 0; length <= capacity; length++ {
					fixed = append(fixed, fixture{Bool: value, Tag: tags[(j+length+capacity)%len(tags)],
						Length: length, Backing: backing, Raw: raw})
				}
			}
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				backing := make([]byte, r.Intn(33))
				raw := make([]byte, r.Intn(25))
				_, err := r.Read(backing)
				must(err)
				_, err = r.Read(raw)
				must(err)
				text := string(raw)
				if i%4 == 0 {
					text = inputs[r.Intn(len(inputs))]
				}
				tag := int(r.Uint64())
				if i%4 == 1 {
					tag = tags[r.Intn(len(tags))]
				}
				cases[i] = fixture{Bool: r.Intn(2) == 0, Tag: tag, Length: r.Intn(len(backing) + 1), Backing: backing, Raw: text}
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
	must(w.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_INTERNAL_STRCONV_BOOL_ROWS", f.Name()))
	must(check.TestInPackage("internal/strconv", "TestGoLeanInternalBool", "difftest/internal/strconv/bool/testdata/bool_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
