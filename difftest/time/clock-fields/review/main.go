package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixed := []struct{ input, want string }{
		{"0 0 0 0 0", "0|0,0,0|3|0|1|3,1"},
		{"59 1 1 0 30", "0|0,0,59|4|0|0|3,31"},
		{"60 2 4 0 31", "0|0,1,0|5|0|1|4,1"},
		{"3599 3 0 4 305", "0|0,59,59|6|0|1|12,31"},
		{"3600 18446744073709551612 1 1 306", "0|1,0,0|1|1|0|13,1"},
		{"86399 18446744073709551613 2 4 364", "0|23,59,59|0|1|1|14,28"},
		{"86400 18446744073709551614 3 99 365", "1|0,0,0|1|1|0|14,29"},
		{"18446744073709551615 18446744073709551615 18446744073709551615 0 365", "213503982334601|7,0,15|2|1|0|14,29"},
		{"86401 7 4 96 275", "1|0,0,1|3|0|1|12,1"},
		{"37321261 10 5 8 213", "431|23,1,1|6|0|1|9,30"},
	}
	f, err := os.CreateTemp("", "golean-clock-fields-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	var mutations [3]int
	for _, c := range fixed {
		fmt.Fprintln(&input, c.input)
		var row struct {
			Seconds, Days, Century uint64
			Year, Day              uint32
			Want                   string
		}
		_, err := fmt.Sscan(c.input, &row.Seconds, &row.Days, &row.Century, &row.Year, &row.Day)
		must(err)
		row.Want = c.want
		must(enc.Encode(row))
		fields := strings.Split(c.want, "|")
		if weekday := (row.Days%7 + 3) % 7; strconv.FormatUint(weekday, 10) != fields[2] {
			mutations[0]++
		}
		if row.Day >= 306 && (strings.HasPrefix(fields[5], "13,") || strings.HasPrefix(fields[5], "14,")) {
			mutations[1]++
		}
		if onlyFour := row.Year%4 == 0; onlyFour != (fields[4] == "1") {
			mutations[2]++
		}
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-time-clock-fields", len(fixed), &input) {
		if got != fixed[i].want {
			log.Fatalf("fixture %d: Lean %q, expected %q", i, got, fixed[i].want)
		}
	}
	must(os.Setenv("GOLEAN_TIME_FIELDS", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockFields",
		"difftest/time/clock-fields/testdata/export_test.go",
		"difftest/time/clock-fields/testdata/clock_fields_test.go"))
	if want := [3]int{3, 4, 2}; mutations != want {
		log.Fatalf("mutation disagreements %v, want %v", mutations, want)
	}
	fmt.Printf("%d literal absolute calendar and clock rows; mutation disagreements %v\n", len(fixed), mutations)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
