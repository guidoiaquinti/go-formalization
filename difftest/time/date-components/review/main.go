package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Year, Month, Day, Hi, Lo, Base int
	Want                           string
}

func main() {
	cases := []fixture{
		{1900, 2, 28, 10, -1, 60, "106751991047891|false|31|28|9,59"},
		{2000, 2, 29, math.MaxInt, 60, 60, fmt.Sprintf("106751991084416|true|31|29|%d,0", math.MinInt)},
		{-400, 3, 1, -9, -121, 60, "106751990207835|true|59|31|-12,59"},
		{2024, 1, 1, 0, 0, 60, "106751991093123|true|0|31|0,0"},
		{1999, 12, 31, 0, math.MinInt, 2, fmt.Sprintf("106751991084356|false|334|31|%d,0", math.MinInt/2)},
		{2024, 3, 0, 1, 61, 60, "106751991093182|true|59|31|2,1"},
	}
	f, err := os.CreateTemp("", "golean-date-components-control-")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&input, strconv.IntSize, c.Year, c.Month, c.Day, c.Hi, c.Lo, c.Base)
		must(enc.Encode(c))
	}
	must(f.Close())
	var rejected [5]int
	for i, got := range check.LeanRows("godiff-time-date-components", len(cases), &input) {
		c := cases[i]
		if got != c.Want {
			log.Fatalf("row %d got %q want %q", i, got, c.Want)
		}
		fields := strings.Split(got, "|")
		absolute, err := strconv.ParseUint(fields[0], 10, 64)
		must(err)
		withoutDayZero := fields[0]
		if c.Day == 0 {
			withoutDayZero = strconv.FormatUint(absolute+1, 10)
		}
		civilLeap := strconv.FormatBool(c.Year%4 == 0)
		before := fields[2]
		days := fields[3]
		if c.Month >= 3 {
			old, err := strconv.Atoi(before)
			must(err)
			before = strconv.Itoa(old + 2)
		}
		if c.Month == 2 {
			days = "28"
		}
		naive := fmt.Sprintf("%d,%d", c.Hi+c.Lo/c.Base, c.Lo%c.Base)
		for j, w := range []string{withoutDayZero, civilLeap, before, days, naive} {
			if w != fields[j] {
				rejected[j]++
			}
		}
	}
	if want := [5]int{1, 1, 3, 1, 2}; rejected != want {
		log.Fatalf("mutation counts %v want %v", rejected, want)
	}
	must(os.Setenv("GOLEAN_TIME_DATE_COMPONENTS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanDateComponents", "difftest/time/date-components/testdata/export_test.go", "difftest/time/date-components/testdata/date_test.go"))
	fmt.Printf("date components review: %d pinned Go and Lean rows, mutations %v\n", len(cases), rejected)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
