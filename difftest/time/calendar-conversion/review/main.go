package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Days, Century                  uint64
	Year, JanFeb, Month, Day, Leap int
	Want                           string
}

func main() {
	var origin int64 = -292277022400
	type control struct {
		days, century                  uint64
		year, janFeb, month, day, leap int
		civilYear                      int64
		civilMonth, civilDay, ordinal  int
	}
	fixed := []control{
		{0, 0, 0, 0, 3, 0, 1, origin, 3, 1, 61},
		{306, 0, 0, 1, 13, 306, 0, origin + 1, 1, 1, 1},
		{146097, 4, 0, 0, 3, 0, 1, origin + 400, 3, 1, 61},
		{1 << 62, 0, 0, 0, 3, 0, 1, origin, 3, 1, 61},
		{106751990354238, 2922770224, 0, 1, 13, 306, 0, 1, 1, 1, 1},
		{106751991047891, 2922770242, 99, 1, 14, 364, 0, 1900, 2, 28, 59},
		{106751991047892, 2922770243, 0, 0, 3, 0, 0, 1900, 3, 1, 60},
		{106751991073400, 2922770243, 69, 1, 13, 306, 0, 1970, 1, 1, 1},
		{106751991084416, 2922770243, 99, 1, 14, 365, 1, 2000, 2, 29, 60},
		{106751991084417, 2922770244, 0, 0, 3, 0, 1, 2000, 3, 1, 61},
		{106751991120941, 2922770245, 0, 0, 3, 0, 0, 2100, 3, 1, 60},
		{106751991230513, 2922770247, 99, 1, 14, 365, 1, 2400, 2, 29, 60},
	}
	f, err := os.CreateTemp("", "golean-calendar-conversion-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	wants := make([]string, 0, len(fixed)+1)
	record := func(c fixture) {
		fmt.Fprintln(&input, strconv.IntSize, c.Days, c.Century, c.Year, c.JanFeb, c.Month, c.Day, c.Leap)
		must(enc.Encode(c))
		wants = append(wants, c.Want)
	}
	for _, c := range fixed {
		y := int(c.civilYear)
		want := fmt.Sprintf("%d,%d,%d|%d|%d|%d|%d,%d,%d|%d,%d", c.century, c.year, c.day, y, c.civilMonth, c.ordinal, y, c.civilMonth, c.civilDay, y, c.ordinal)
		record(fixture{c.days, c.century, c.year, c.janFeb, c.month, c.day, c.leap, want})
	}
	record(fixture{0, math.MaxUint64, 98, 1, 14, 365, 1,
		fmt.Sprintf("0,0,0|%d|2|60|%d,3,1|%d,61", int(origin-1), int(origin), int(origin))})
	must(f.Close())
	for i, got := range check.LeanRows("godiff-time-calendar-conversion", len(wants), &input) {
		if got != wants[i] {
			log.Fatalf("fixture %d: Lean %q, expected %q", i, got, wants[i])
		}
	}
	must(os.Setenv("GOLEAN_TIME_CALENDAR", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanCalendarConversion",
		"difftest/time/calendar-conversion/testdata/export_test.go",
		"difftest/time/calendar-conversion/testdata/calendar_test.go"))
	fmt.Printf("%d literal Gregorian and native-width controls passed on pinned Go and Lean\n", len(wants))
}
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
