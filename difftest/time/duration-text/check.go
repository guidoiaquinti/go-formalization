package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Value   int64
	Present bool
	Cells   [32]byte
	Want    string
}

func render(n int64) string {
	if n == 0 {
		return "0s"
	}
	var b strings.Builder
	unit := uint64(n)
	if n < 0 {
		b.WriteByte('-')
		unit = uint64(-(n + 1)) + 1
	}
	padded := func(value uint64, digits int) string {
		if value == 0 {
			return ""
		}
		return "." + strings.TrimRight(fmt.Sprintf("%0*d", digits, value), "0")
	}
	switch {
	case unit < 1000:
		fmt.Fprintf(&b, "%dns", unit)
	case unit < 1000000:
		fmt.Fprintf(&b, "%d%sµs", unit/1000, padded(unit%1000, 3))
	case unit < 1000000000:
		fmt.Fprintf(&b, "%d%sms", unit/1000000, padded(unit%1000000, 6))
	default:
		hours, remainder := unit/3600000000000, unit%3600000000000
		minutes, remainder := remainder/60000000000, remainder%60000000000
		if hours > 0 {
			fmt.Fprintf(&b, "%dh", hours)
		}
		if hours > 0 || minutes > 0 {
			fmt.Fprintf(&b, "%dm", minutes)
		}
		fmt.Fprintf(&b, "%d%ss", remainder/1000000000, padded(remainder%1000000000, 9))
	}
	return b.String()
}

func wire(data []byte) string {
	out := make([]string, len(data))
	for i, b := range data {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}

func main() {
	f, err := os.CreateTemp("", "golean-duration-text")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			present := 0
			if c.Present {
				present = 1
			}
			fmt.Fprint(&input, c.Value, " ", present)
			for _, b := range c.Cells {
				fmt.Fprint(&input, " ", b)
			}
			fmt.Fprintln(&input)
		}
		for i, row := range check.LeanRows("godiff-time-duration-text", len(cases), &input) {
			c := cases[i]
			text := render(c.Value)
			if got := time.Duration(c.Value).String(); got != text {
				log.Fatalf("%s row %d: Go %q want %q", label, i, got, text)
			}
			c.Want = "nil"
			if c.Present {
				storage := c.Cells
				start := len(storage) - len(text)
				copy(storage[start:], text)
				c.Want = fmt.Sprintf("%d:%s", start, wire(storage[:]))
			}
			if want := wire([]byte(text)) + "|" + c.Want; row != want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d pinned Go and Lean rows agree\n", label, len(cases))
	}
	var fixed []fixture
	for _, base := range []int64{0, 1, 999, 1000, 1000000, 1000000000, 60000000000, 3600000000000, 86400000000000} {
		for _, delta := range []int64{-101, -10, -1, 0, 1, 10, 101} {
			for _, sign := range []int64{-1, 1} {
				for _, present := range []bool{false, true} {
					c := fixture{Value: sign * (base + delta), Present: present}
					for i := range c.Cells {
						c.Cells[i] = byte(192 + i)
					}
					fixed = append(fixed, c)
				}
			}
		}
	}
	fixed = append(fixed, fixture{Value: math.MinInt64, Present: true}, fixture{Value: math.MaxInt64, Present: true})
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			c := fixture{Value: int64(r.Uint64()), Present: r.Intn(7) != 0}
			if i%3 == 0 {
				c.Value %= 2_000_000_000
			}
			r.Read(c.Cells[:])
			cases[i] = c
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_DURATION_TEXT_ROWS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanDurationText", "difftest/time/duration-text/testdata/export_test.go", "difftest/time/duration-text/testdata/duration_text_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
