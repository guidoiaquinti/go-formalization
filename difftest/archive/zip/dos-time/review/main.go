package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Date, Clock   uint16
	Seconds       int64
	Nanos, Offset int
	Local, UTC    [6]int
	Want          string
}

func main() {
	fixed := []fixture{
		{0, 0, 315532801, 999999999, -3600, [6]int{1979, 12, 31, 23, 0, 1}, [6]int{1980, 1, 1, 0, 0, 1},
			"312768000|65439:47104|312768000|315532801:999999999:33:0"},
		{65535, 65535, 951868799, 123, 19800, [6]int{2000, 3, 1, 5, 29, 59}, [6]int{2000, 2, 29, 23, 59, 59},
			"4362710642|10337:11197|4362710642|951868799:123:10333:49021"},
		{61533, 0, 4107542399, 0, 60, [6]int{2100, 3, 1, 0, 0, 59}, [6]int{2100, 2, 28, 23, 59, 59},
			"4107542400|61537:29|4107542400|4107542399:0:61532:49021"},
		{10333, 0, 315532799, 1, 3600, [6]int{1980, 1, 1, 0, 59, 59}, [6]int{1979, 12, 31, 23, 59, 59},
			"951782400|33:1917|951782400|315532799:1:65439:49021"},
		{22944, 53215, 4354819200, 10, -60, [6]int{2107, 12, 31, 23, 59, 0}, [6]int{2108, 1, 1, 0, 0, 0},
			"1735696982|65439:48992|1735696982|4354819200:10:33:0"},
		{33, 31, 1735689598, 5, 2, [6]int{2025, 1, 1, 0, 0, 0}, [6]int{2024, 12, 31, 23, 59, 58},
			"315532862|23073:0|315532862|1735689598:5:22943:49021"},
	}
	f, err := os.CreateTemp("", "golean-zip-dos-review")
	must(err)
	defer os.Remove(f.Name())
	enc := json.NewEncoder(f)
	var input bytes.Buffer
	var mutations [2]int
	for _, c := range fixed {
		fmt.Fprint(&input, c.Date, " ", c.Clock, " ", c.Seconds, " ", c.Nanos)
		for _, civil := range [][6]int{c.Local, c.UTC} {
			for _, cell := range civil {
				fmt.Fprint(&input, " ", cell)
			}
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
		fields := strings.Split(c.Want, "|")
		if prefix := fmt.Sprintf("%d:%d:", c.Seconds, c.Nanos); prefix+fields[1] != fields[3] {
			mutations[0]++
		}
		if sec := c.UTC[5]; sec%2 != 0 {
			t := c.UTC[3]*2048 + c.UTC[4]*32 + (sec+1)/2
			if !strings.HasSuffix(fields[3], fmt.Sprintf(":%d", t)) {
				mutations[1]++
			}
		}
	}
	must(f.Close())
	for i, got := range check.LeanRows("godiff-zip-dos-time", len(fixed), &input) {
		if got != fixed[i].Want {
			log.Fatalf("fixture %d: Lean %q, want %q", i, got, fixed[i].Want)
		}
	}
	must(os.Setenv("GOLEAN_ZIP_DOS_ROWS", f.Name()))
	must(check.TestInPackage("archive/zip", "TestGoLeanDOSTime", "difftest/archive/zip/dos-time/testdata/dos_test.go"))
	if want := [2]int{6, 4}; mutations != want {
		log.Fatalf("mutation disagreements %v, want %v", mutations, want)
	}
	fmt.Printf("%d literal ZIP DOS time rows; mutation disagreements %v\n", len(fixed), mutations)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
