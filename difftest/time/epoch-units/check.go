package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const epoch = int64(62135596800)

type fixture struct {
	Wall, Ext        uint64
	Sec, Nsec, Units int64
	Want             string
}

func expected(c fixture) string {
	modulus := new(big.Int).Lsh(big.NewInt(1), 64)
	wrap := func(n *big.Int) uint64 { return new(big.Int).Mod(n, modulus).Uint64() }
	nanos := int64(c.Wall % 1073741824)
	absolute := int64(c.Ext)
	if c.Wall >= 1<<63 {
		absolute = 59453308800 + int64(c.Wall/1073741824%(1<<33))
	}
	extract := func(scale int64) string {
		value := new(big.Int).Sub(big.NewInt(absolute), big.NewInt(epoch))
		value.Mul(value, big.NewInt(scale))
		value.Add(value, big.NewInt(nanos/(1000000000/scale)))
		return fmt.Sprint(int64(wrap(value)))
	}
	raw := func(second *big.Int, nano uint64) string {
		ext := new(big.Int).Add(second, big.NewInt(epoch))
		return fmt.Sprintf("%d,%d", nano, wrap(ext))
	}
	construct := func(seconds int64, ticks *big.Int) string {
		var whole, remainder big.Int
		whole.DivMod(ticks, big.NewInt(1000000000), &remainder)
		whole.Add(&whole, big.NewInt(seconds))
		return raw(&whole, remainder.Uint64())
	}
	milli := new(big.Int).Mul(big.NewInt(c.Units), big.NewInt(1000000))
	micro := new(big.Int).Mul(big.NewInt(c.Units), big.NewInt(1000))
	return strings.Join([]string{extract(1000), extract(1000000), extract(1000000000),
		raw(big.NewInt(c.Sec), uint64(int64(int32(c.Nsec)))), construct(c.Sec, big.NewInt(c.Nsec)),
		construct(0, milli), construct(0, micro)}, "|")
}

func main() {
	f, err := os.CreateTemp("", "golean-time-epoch")
	must(err)
	defer os.Remove(f.Name())
	defer f.Close()
	enc := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.Wall, c.Ext, c.Sec, c.Nsec, c.Units)
		}
		for i, row := range check.LeanRows("godiff-time-epoch-units", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if row != c.Want {
				log.Fatalf("%s row %d: Lean %q want %q", label, i, row, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows agree\n", label, len(cases))
	}
	words := []uint64{0, 1, 999, 1000, 999999, 1000000, 999999999, 1 << 30, 1 << 63, math.MaxUint64}
	signed := []int64{math.MinInt64, math.MinInt64 + 1, -1000000001, -1000000000, -1000001, -1000000, -1001, -1000, -1, 0, 1, 999, 1000, 1000000000, 1000000001, math.MaxInt64}
	var fixed []fixture
	for _, wall := range words {
		for _, sec := range signed {
			for _, nsec := range signed {
				fixed = append(fixed, fixture{Wall: wall, Ext: uint64(sec + epoch), Sec: sec, Nsec: nsec, Units: nsec})
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20_000)
		for i := range cases {
			cases[i] = fixture{Wall: r.Uint64(), Ext: r.Uint64(), Sec: int64(r.Uint64()), Nsec: int64(r.Uint64()), Units: int64(r.Uint64())}
		}
		compare(fmt.Sprintf("seed %d", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_TIME_EPOCH_ROWS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanEpochUnits", "difftest/time/location/testdata/export_test.go",
		"difftest/time/epoch-units/testdata/export_test.go", "difftest/time/epoch-units/testdata/epoch_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
