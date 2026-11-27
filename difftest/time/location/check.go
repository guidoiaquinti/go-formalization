package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Wall, Ext uint64
	Old, Next int
	Want      string
}

func expected(c fixture) string {
	wall, ext := c.Wall, c.Ext
	if wall >= 1<<63 {
		ext = (wall/(1<<30))%(1<<33) + 59453308800
		wall %= 1 << 30
	}
	output := func(loc int) string { return fmt.Sprintf("%d,%d,%d", wall, ext, loc) }
	set := c.Next
	if set == 1 {
		set = 0
	}
	in := "nil"
	if c.Next != 0 {
		in = output(set)
	}
	exposed := c.Old
	if exposed == 0 {
		exposed = 1
	}
	return strings.Join([]string{output(c.Old), output(set), output(0), output(2), in, fmt.Sprint(exposed)}, "|")
}

func main() {
	f, err := os.CreateTemp("", "golean-time-location")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.Wall, c.Ext, c.Old, c.Next)
		}
		for i, row := range check.LeanRows("godiff-time-location", len(cases), &input) {
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
	for _, wall := range []uint64{0, 1, (1 << 30) - 1, 1 << 30, math.MaxInt64, 1 << 63, (1 << 63) + 1, math.MaxUint64} {
		for _, ext := range []uint64{0, 1, math.MaxInt64, 1 << 63, math.MaxUint64} {
			for old := range 5 {
				for next := range 5 {
					fixed = append(fixed, fixture{Wall: wall, Ext: ext, Old: old, Next: next})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{Wall: r.Uint64(), Ext: r.Uint64(), Old: r.Intn(5), Next: r.Intn(5)}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_TIME_LOCATION_ROWS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanLocation", "difftest/time/location/testdata/export_test.go", "difftest/time/location/testdata/location_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
