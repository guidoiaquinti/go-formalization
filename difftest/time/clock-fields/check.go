package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Seconds, Days, Century uint64
	Year, Day              uint32
	Want                   string
}

func expected(c fixture) string {
	days, seconds := c.Seconds/86400, c.Seconds%86400
	weekday := (c.Days%7 + 3) % 7
	if c.Days > math.MaxUint64-3 {
		weekday = c.Days - (math.MaxUint64 - 2)
	}
	janFeb := 0
	if c.Day >= 306 {
		janFeb = 1
	}
	leap := 0
	yearMod400 := (c.Century%4*100 + uint64(c.Year)) % 400
	if yearMod400%4 == 0 && (yearMod400%100 != 0 || yearMod400 == 0) {
		leap = 1
	}
	month := 3
	left := int(c.Day)
	for _, length := range []int{31, 30, 31, 30, 31, 31, 30, 31, 30, 31, 31} {
		if left < length {
			break
		}
		left -= length
		month++
	}
	return fmt.Sprintf("%d|%d,%d,%d|%d|%d|%d|%d,%d", days, seconds/3600, seconds/60%60,
		seconds%60, weekday, janFeb, leap, month, left+1)
}

func main() {
	f, err := os.CreateTemp("", "golean-time-fields-*")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	compare := func(label string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, c.Seconds, c.Days, c.Century, c.Year, c.Day)
		}
		for i, got := range check.LeanRows("godiff-time-clock-fields", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				log.Fatalf("%s row %d %+v: Lean %q, want %q", label, i, c, got, c.Want)
			}
			must(encoder.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows matched\n", label, len(cases))
	}
	var fixed []fixture
	boundaries := []uint64{0, 1, 59, 60, 61, 3599, 3600, 86399, 86400, 86401,
		math.MaxUint64 - 3, math.MaxUint64 - 2, math.MaxUint64 - 1, math.MaxUint64}
	for i, seconds := range boundaries {
		for day := range uint32(366) {
			fixed = append(fixed, fixture{Seconds: seconds, Days: boundaries[(i+int(day))%len(boundaries)],
				Century: uint64(i), Year: day % 100, Day: day})
		}
	}
	compare("boundaries", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20000)
		for i := range cases {
			cases[i] = fixture{Seconds: rng.Uint64(), Days: rng.Uint64(), Century: rng.Uint64(),
				Year: uint32(rng.Intn(100)), Day: uint32(rng.Intn(366))}
		}
		compare(fmt.Sprint("seed ", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_TIME_FIELDS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockFields",
		"difftest/time/clock-fields/testdata/export_test.go", "difftest/time/clock-fields/testdata/clock_fields_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
