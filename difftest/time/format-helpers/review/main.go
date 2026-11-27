package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Value             uint64
	Length, Precision int
	Backing           []byte
	Nil               bool
	Want              string
}

func main() {
	cases := []fixture{
		{567, 2, 0, []byte{7, 8, 9}, false, "bounds:54,55,9|2:567:7,8,9"},
		{120, 2, 3, []byte{7, 8, 9}, false, "bounds:50,48,9|bounds:49,50,9"},
		{120, 4, 3, []byte{7, 8, 9, 10, 11}, false, "1:0:7,49,50,48,11|1:0:7,46,49,50,11"},
		{18000, 6, 3, []byte{7, 8, 9, 10, 11, 12}, false, "1:0:7,49,56,48,48,48|6:18:7,8,9,10,11,12"},
		{1234506, 7, 4, []byte{7, 8, 9, 10, 11, 12, 13, 14}, false,
			"0:0:49,50,51,52,53,48,54,14|2:123:7,8,46,52,53,48,54,14"},
		{7, 0, math.MinInt, nil, true, "bounds:|0:7:"},
		{0, 2, -4, []byte{7, 8, 9}, false, "1:0:7,48,9|2:0:7,8,9"},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-time-format-helpers-")
	must(err)
	defer os.Remove(truth.Name())
	enc := gob.NewEncoder(truth)
	for _, c := range cases {
		fmt.Fprint(&input, c.Value, " ", c.Length, " ", c.Precision)
		for _, b := range c.Backing {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_TIME_FMT_ROWS", truth.Name()))
	must(check.TestInPackage("time", "TestGoLeanFormatHelpers", "difftest/time/format-helpers/testdata/bridge_test.go",
		"difftest/time/format-helpers/testdata/format_test.go"))
	rows := check.LeanRows("godiff-time-format-helpers", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].Want {
			log.Fatalf("row %d got %q, want %q", i, row, cases[i].Want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{0, 0, "bounds:7,8,9"},                 // pre-check every decimal write
		{1, 1, "bounds:7,8,9"},                 // pre-check decimal dot
		{1, 1, "0:0:46,49,50"},                 // write through spare capacity
		{2, 1, "0:0:46,49,50,48,11"},           // print trailing fractional zero
		{3, 1, "5:18:7,8,9,10,11,46"},          // add decimal point to empty fraction
		{4, 1, "0:123:7,8,46,52,53,48,54,14"},  // return wrong suffix index
		{4, 1, "2:1234:7,8,46,52,53,48,54,14"}, // divide one time too few
		{5, 1, "bounds:"},                      // try to write on negative precision
		{6, 0, "2:0:7,8,9"},                    // emit no digit for zero
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.value)
		}
	}
	log.Print("seven independent pinned Go and Lean rows passed and rejected nine wrong observations")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
