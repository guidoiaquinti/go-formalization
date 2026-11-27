package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

const absoluteYears uint64 = 292277022400

type fixture struct {
	Days, Century                  uint64
	Year, JanFeb, Month, Day, Leap int
	Want                           string
}

func daysBefore(year int) int {
	return 365*year + year/4 - year/100 + year/400
}

func expected(row fixture) string {
	effective := row.Days % (1 << 62)
	cycles, rem := effective/146097, int(effective%146097)
	yearInCycle := sort.Search(400, func(i int) bool { return daysBefore(i+1) > rem })
	century, lowYear := 4*cycles+uint64(yearInCycle/100), yearInCycle%100
	day := rem - daysBefore(yearInCycle)
	absoluteYear := cycles*400 + uint64(yearInCycle)
	janFeb := day >= 306
	if janFeb {
		absoluteYear++
	}
	civilYear := int(absoluteYear - absoluteYears)
	month, monthDay := 3, day+1
	for _, length := range []int{31, 30, 31, 30, 31, 31, 30, 31, 30, 31, 31} {
		if monthDay <= length {
			break
		}
		monthDay -= length
		month++
	}
	yearDay := day + 60
	if janFeb {
		month -= 12
		yearDay = day - 305
	} else if absoluteYear%4 == 0 && (absoluteYear%100 != 0 || absoluteYear%400 == 0) {
		yearDay++
	}
	standaloneYear := int(row.Century*100 + uint64(row.Year+row.JanFeb) - absoluteYears)
	standaloneMonth := row.Month
	standaloneDay := row.Day + 60 + row.Leap
	if row.JanFeb == 1 {
		standaloneMonth -= 12
		standaloneDay = row.Day - 305
	}
	return fmt.Sprintf("%d,%d,%d|%d|%d|%d|%d,%d,%d|%d,%d", century, lowYear, day,
		standaloneYear, standaloneMonth, standaloneDay, civilYear, month, monthDay, civilYear, yearDay)
}

func main() {
	f, err := os.CreateTemp("", "golean-calendar-conversion-*")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	compare := func(name string, cases []fixture) {
		var input bytes.Buffer
		for _, c := range cases {
			fmt.Fprintln(&input, strconv.IntSize, c.Days, c.Century, c.Year, c.JanFeb, c.Month, c.Day, c.Leap)
		}
		for i, got := range check.LeanRows("godiff-time-calendar-conversion", len(cases), &input) {
			c := cases[i]
			c.Want = expected(c)
			if got != c.Want {
				log.Fatalf("%s row %d %+v: Lean %q, want %q", name, i, c, got, c.Want)
			}
			must(enc.Encode(c))
		}
		fmt.Printf("%s: %d Lean rows matched\n", name, len(cases))
	}
	var fixed []fixture
	for i, day := range []uint64{0, 1, 305, 306, 364, 365, 1460, 1461, 36523, 36524, 36525,
		146096, 146097, 146098, 1<<62 - 1, 1 << 62, 1<<62 + 1, math.MaxUint64 - 1, math.MaxUint64} {
		for d := range 366 {
			j, m := 0, 3+d%10
			if d >= 306 {
				j, m = 1, 13+(d-306)/31
			}
			fixed = append(fixed, fixture{Days: day, Century: day, Year: (d + i) % 100, JanFeb: j,
				Month: m, Day: d, Leap: (d + i) % 2})
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]fixture, 20000)
		for i := range cases {
			d, j, m := rng.Intn(366), 0, 3+rng.Intn(10)
			if d >= 306 {
				j, m = 1, 13+rng.Intn(2)
			}
			cases[i] = fixture{Days: rng.Uint64(), Century: rng.Uint64(), Year: rng.Intn(100),
				JanFeb: j, Month: m, Day: d, Leap: rng.Intn(2)}
		}
		compare(fmt.Sprint("seed ", seed), cases)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_TIME_CALENDAR", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanCalendarConversion", "difftest/time/calendar-conversion/testdata/export_test.go",
		"difftest/time/calendar-conversion/testdata/calendar_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
