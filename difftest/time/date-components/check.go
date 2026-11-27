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
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Year, Month, Day, Hi, Lo, Base int
	Want                           string
}

func expected(c fixture) string {
	leap := c.Year%4 == 0 && (c.Year%100 != 0 || c.Year%400 == 0)
	civil := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	before := 0
	for _, days := range civil[:c.Month-1] {
		before += days
	}
	days := civil[c.Month-1]
	if leap && c.Month == 2 {
		days++
	}
	shift := uint64(c.Year) + 292277022400
	offset := before - 59
	if c.Month <= 2 {
		shift--
		offset = before + 306
	}
	centuries, years := shift/100, shift%100
	centuryDays := (centuries * 146097) / 4
	date := centuryDays + 365*years + years/4 + uint64(offset) + uint64(c.Day) - 1
	var q, r big.Int
	q.DivMod(big.NewInt(int64(c.Lo)), big.NewInt(int64(c.Base)), &r)
	normalizedHigh := int(uint64(c.Hi) + uint64(q.Int64()))
	return fmt.Sprintf("%d|%t|%d|%d|%d,%d", date, leap, before, days, normalizedHigh, r.Int64())
}

func main() {
	file, err := os.CreateTemp("", "golean-time-date-components-*")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, strconv.IntSize, c.Year, c.Month, c.Day, c.Hi, c.Lo, c.Base)
		}
		for i, got := range check.LeanRows("godiff-time-date-components", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				log.Fatalf("%s row %d %+v: Lean %q, want %q", label, i, c, got, c.Want)
			}
			must(encoder.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows matched\n", label, len(cases))
	}
	years := []int{math.MinInt, math.MinInt + 1, -2000, -1900, -400, -100, -4, -1, 0, 1, 4,
		100, 400, 1900, 2000, 2100, math.MaxInt - 1, math.MaxInt}
	limits := []int{math.MinInt, math.MinInt + 1, -61, -60, -1, 0, 1, 59, 60, 61, math.MaxInt - 1, math.MaxInt}
	var fixed []fixture
	for y, year := range years {
		for month := 1; month <= 12; month++ {
			for i, lo := range limits {
				for _, base := range []int{1, 2, 7, 60, 1_000_000_000, math.MaxInt} {
					fixed = append(fixed, fixture{Year: year, Month: month, Day: limits[(i+y)%len(limits)],
						Hi: limits[(i+y+month)%len(limits)], Lo: lo, Base: base})
				}
			}
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20000)
		for i := range cases {
			cases[i] = fixture{Year: int(rng.Uint64()), Month: 1 + rng.Intn(12), Day: int(rng.Uint64()),
				Hi: int(rng.Uint64()), Lo: int(rng.Uint64()), Base: 1 + rng.Intn(math.MaxInt)}
		}
		compare(fmt.Sprint("seed ", seed), cases)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_TIME_DATE_COMPONENTS", file.Name()))
	must(check.TestInPackage("time", "TestGoLeanDateComponents", "difftest/time/date-components/testdata/export_test.go",
		"difftest/time/date-components/testdata/date_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
