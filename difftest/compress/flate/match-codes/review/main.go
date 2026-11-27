package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		offset uint32
		length byte
		want   string
	}{
		{255, 7, "15 7"},
		{256, 8, "16 8"},
		{32767, 9, "29 8"},
		{32768, 10, "14 9"},
		{32895, 254, "14 27"},
		{32896, 255, "15 28"},
		{33024, 0, "16 0"},
		{math.MaxUint32, 255, "29 28"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-flate-match-codes-")
	must(err)
	defer os.Remove(truth.Name())
	for _, c := range cases {
		fmt.Fprintln(&input, c.offset, c.length)
		_, err := fmt.Fprintln(truth, c.offset, c.length, c.want)
		must(err)
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_FLATE_MATCH_ROWS", truth.Name()))
	must(check.TestInPackage("compress/flate", "TestGoLeanMatchCodes", "difftest/compress/flate/match-codes/testdata/match_test.go"))

	rows := check.LeanRows("godiff-flate-match-codes", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].want {
			log.Fatalf("row %d = %q, want %q", i, row, cases[i].want)
		}
	}
	reject := func(row, column int, wrong string) {
		if strings.Fields(rows[row])[column] == wrong {
			log.Fatalf("row %d column %d accepted %q", row, column, wrong)
		}
	}
	reject(0, 0, "14")    // switch to the long table one byte early
	reject(1, 0, "0")     // switch to the long table one byte late
	reject(3, 0, "0")     // apply the window modulus before deciding which table to use
	reject(7, 0, "65535") // keep increasing the code outside the window
	reject(2, 1, "9")     // increment the code for every byte past the first group
	reject(5, 1, "27")    // omit the special final length code
	log.Print("eight independent pinned Go and Lean rows passed and rejected six wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
