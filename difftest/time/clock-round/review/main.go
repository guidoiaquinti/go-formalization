package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	fixed := []struct{ input, want string }{
		{"5 0 3 10 5", "0,0,3|10,0,3|0,5|false"},
		{"4 0 2 10 4", "0,0,2|0,0,2|0,4|true"},
		{"999999995 -1 1 10 -9223372036854775808", "999999990,-1,1|0,0,1|1,5|true"},
		{"100000000 2 3 200000000 99999999", "0,2,3|200000000,2,3|0,100000000|true"},
		{"750000000 -2 2 1000000000 499999999", "0,-2,2|0,-1,2|0,750000000|true"},
		{"345678901 7 3 700000003 9223372036854775807", "30,7,3|30,7,3|0,345678871|false"},
		{"9223372040076001289 77 3 0 -1", "9,59453308803,3|9,59453308803,3|skip|false"},
		{"9223372040076001289 77 3 -1 9223372036854775807", "9,59453308803,3|9,59453308803,3|skip|true"},
		{"0 -9223372036854775808 0 1000000000 9223372036854775807", "0,-9223372036854775808,0|0,-9223372036854775808,0|0,0|false"},
		{"999999999 -1 2 1000000000 0", "0,-1,2|0,0,2|1,999999999|true"},
		{"9223372036854775818 12 3 20 10", "0,59453308800,3|20,59453308800,3|0,10|false"},
	}
	f, err := os.CreateTemp("", "golean-time-round-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	for _, c := range fixed {
		fmt.Fprintln(&input, c.input)
		var row struct {
			Wall                uint64
			Ext, Duration, Half int64
			Loc                 int
			Lean                string
		}
		_, err := fmt.Sscan(c.input, &row.Wall, &row.Ext, &row.Loc, &row.Duration, &row.Half)
		must(err)
		row.Lean = c.want
		must(enc.Encode(row))
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-time-clock-round", len(fixed), &input) {
		if got != fixed[i].want {
			log.Fatalf("fixture %d: Lean %q, expected %q", i, got, fixed[i].want)
		}
	}
	must(os.Setenv("GOLEAN_TIME_CLOCK_ROUND", f.Name()))
	must(check.TestInPackage("time", "TestGoLeanClockRound",
		"difftest/time/location/testdata/export_test.go",
		"difftest/time/clock-arithmetic/testdata/bridge_test.go",
		"difftest/time/clock-round/testdata/export_test.go",
		"difftest/time/clock-round/testdata/clock_round_test.go"))
	fmt.Printf("%d literal clock-rounding controls passed on pinned Go and Lean\n", len(fixed))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
