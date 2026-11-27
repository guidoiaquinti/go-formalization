package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Value             uint64
	Length, Precision int
	Backing           []byte
	Nil               bool
	Want              string
}

func oracle(c fixture) string {
	digits := strconv.FormatUint(c.Value, 10)
	fraction := ""
	whole := c.Value
	if c.Precision > 0 {
		if padding := c.Precision - len(digits); padding > 0 {
			digits = strings.Repeat("0", padding) + digits
		}
		cut := len(digits) - c.Precision
		whole = 0
		if cut > 0 {
			var err error
			whole, err = strconv.ParseUint(digits[:cut], 10, 64)
			must(err)
		}
		if remain := strings.TrimRight(digits[cut:], "0"); remain != "" {
			fraction = "." + remain
		}
	}
	return output(c, strconv.FormatUint(c.Value, 10), 0) + "|" + output(c, fraction, whole)
}

func output(c fixture, text string, value uint64) string {
	backing := slices.Clone(c.Backing)
	state := "bounds"
	written := min(c.Length, len(text))
	copy(backing[c.Length-written:c.Length], text[len(text)-written:])
	if c.Length >= len(text) {
		state = fmt.Sprintf("%d:%d", c.Length-len(text), value)
	}
	digits := make([]string, len(backing))
	for i, b := range backing {
		digits[i] = strconv.Itoa(int(b))
	}
	return state + ":" + strings.Join(digits, ",")
}

func main() {
	f, err := os.CreateTemp("", "golean-time-format-helpers")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	writer := bufio.NewWriter(f)
	enc := gob.NewEncoder(writer)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprint(&input, c.Value, " ", c.Length, " ", c.Precision)
			for _, cell := range c.Backing {
				fmt.Fprint(&input, " ", cell)
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-time-format-helpers", len(cases), &input) {
			c := cases[i]
			c.Want = oracle(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q expected %q (value %d, length %d, precision %d)", label, i, row, c.Want, c.Value, c.Length, c.Precision)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	values := []uint64{0, 1, 9, 10, 11, 90, 99, 100, 101, 105, 120, 567, 999, 1000, 1234567000, 9999999999, 10000000000000000000, math.MaxUint64}
	lengths := []int{0, 1, 2, 3, 4, 9, 19, 20, 21, 32}
	precisions := []int{math.MinInt32, -1, 0, 1, 2, 3, 6, 9, 19, 20, 21, 32}
	for _, value := range values {
		for _, length := range lengths {
			for _, precision := range precisions {
				backing := make([]byte, length+3)
				for i := range backing {
					backing[i] = byte(250 - i)
				}
				fixed = append(fixed, fixture{Value: value, Length: length, Precision: precision, Backing: backing})
				if length == 0 {
					fixed = append(fixed, fixture{Value: value, Length: 0, Precision: precision, Nil: true},
						fixture{Value: value, Length: 0, Precision: precision, Backing: []byte{}})
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
				backing := make([]byte, r.Intn(65))
				_, err := r.Read(backing)
				must(err)
				value := r.Uint64()
				if i%3 == 0 {
					value = value / 1000 * 1000
				}
				cases[i] = fixture{Value: value, Length: r.Intn(len(backing) + 1), Precision: r.Intn(39) - 6, Backing: backing}
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
	must(writer.Flush())
	must(f.Close())
	must(os.Setenv("GOLEAN_TIME_FMT_ROWS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanFormatHelpers", "difftest/time/format-helpers/testdata/bridge_test.go",
		"difftest/time/format-helpers/testdata/format_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
