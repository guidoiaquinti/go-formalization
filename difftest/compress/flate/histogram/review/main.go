package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	Data          []byte
	Backing, Want []uint16
	Length        int
	Panic         bool
}

func main() {
	short := make([]uint16, 255)
	short[0], short[254] = 51, 52
	full := make([]uint16, 257)
	full[0], full[255], full[256] = 65535, 65534, 777
	empty := slices.Clone(full)
	small := slices.Clone(full)
	small[0], small[255] = 0, 2
	long := slices.Clone(full)
	long[0], long[7], long[255] = 0, 1, 8191
	data := bytes.Repeat([]byte{255}, 8195)
	data[0], data[1] = 0, 7
	cases := []fixture{
		{nil, short, short, 0, true},
		{[]byte{0, 254, 255}, short, short, 255, true},
		{nil, full, empty, 0, false},
		{[]byte{0, 255, 255, 255, 255}, full, small, 0, false},
		{data, full, long, 1, false},
	}
	var input bytes.Buffer
	truth, err := os.CreateTemp("", "review-flate-histogram-")
	must(err)
	defer os.Remove(truth.Name())
	enc := gob.NewEncoder(truth)
	for _, c := range cases {
		fmt.Fprint(&input, len(c.Backing))
		for _, w := range c.Backing {
			fmt.Fprint(&input, " ", w)
		}
		for _, b := range c.Data {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		must(enc.Encode(c))
	}
	must(truth.Close())
	must(os.Setenv("GOLEAN_REVIEW_FLATE_HISTOGRAM", truth.Name()))
	must(check.TestInPackage("compress/flate", "TestReviewGoLeanHistogram",
		"difftest/compress/flate/histogram/review/testdata/control_test.go"))

	got := check.LeanRows("godiff-flate-histogram", len(cases), &input)
	for i, row := range got {
		c := cases[i]
		fields := []string{"0"}
		if c.Panic {
			fields[0] = "1"
		}
		for _, w := range c.Want {
			fields = append(fields, strconv.Itoa(int(w)))
		}
		want := strings.Join(fields, " ")
		if row != want+" "+want {
			log.Fatalf("row %d differs from independent full backing", i)
		}
	}
	reject := func(row, column int, mutant string) {
		fields := strings.Fields(got[row])
		stride := len(cases[row].Backing) + 1
		for index := range 2 {
			if fields[index*stride+column] == mutant {
				log.Fatalf("row %d column %d accepted mutant %q", row, column, mutant)
			}
		}
	}
	reject(0, 0, "0")       // return before checking cap on empty input
	reject(1, 1, "52")      // increment before detecting short cap
	reject(2, 0, "1")       // check visible length instead of capacity
	reject(3, 256, "65535") // saturate rather than wrap
	reject(4, 256, "2047")  // deduplicate identical lane updates
	reject(4, 1, "65535")   // drop the initial partial group
	reject(4, 257, "0")     // clear counters beyond the byte domain
	log.Print("five independent pinned Go and Lean rows passed and rejected seven wrong observations each")
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
