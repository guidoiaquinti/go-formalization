package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Offset, Flag int
	Stored, Spare     []byte
	Count             int
}

type fixture struct {
	Input input
	Want  string
}

func main() {
	data, spare := []byte{5, 0, 255}, []byte{7, 9}
	cases := []fixture{
		{input{Count: math.MinInt}, "nil-receiver"},
		{input{Tag: 1, Flag: -7, Count: 0}, "0:true:-7:0:0:nil::"},
		{input{Tag: 2, Flag: 127, Count: 0}, "0:true:127:0:0:slice::"},
		{input{Tag: 1, Flag: 4, Count: -1}, "bounds:4:0:0:nil::"},
		{input{Tag: 1, Flag: 4, Count: 1}, "0:false:4:0:0:nil::"},
		{input{Tag: 2, Offset: 2, Flag: -128, Stored: data, Spare: spare, Count: 2}, "3:true:-128:2:5:slice:5,0,255,7,9:"},
		{input{Tag: 2, Offset: 2, Flag: -128, Stored: data, Spare: spare, Count: 3}, "0:false:-128:2:5:slice:5,0,255:7,9"},
		{input{Tag: 2, Offset: 2, Flag: -128, Stored: data, Spare: spare, Count: math.MaxInt}, "0:false:-128:2:5:slice:5,0,255:7,9"},
		{input{Tag: 2, Offset: 2, Flag: -128, Stored: data, Spare: spare, Count: math.MinInt}, "bounds:-128:2:5:slice:5,0,255:7,9"},
		{input{Tag: 2, Offset: 2, Flag: -128, Stored: data, Spare: spare, Count: -3}, "3:true:-128:2:5:slice::5,0,255,7,9"},
		{input{Tag: 2, Offset: 2, Flag: -128, Stored: data, Spare: spare, Count: -4}, "bounds:-128:2:5:slice:5,0,255:7,9"},
		{input{Tag: 2, Offset: 3, Flag: 2, Stored: data, Spare: spare, Count: -2}, "3:true:2:3:5:slice:5:0,255,7,9"},
	}
	var model bytes.Buffer
	for _, fixture := range cases {
		c := fixture.Input
		fmt.Fprintf(&model, "%d %d", c.Tag, c.Count)
		if c.Tag != 0 {
			fmt.Fprintf(&model, " %d %d %d", len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&model, " %d", b)
				}
			}
		}
		fmt.Fprintln(&model)
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-try-grow-by-reslice", len(cases), &model) {
		if want := cases[i].Want; got != want {
			log.Fatalf("fixture %d: Lean=%q want=%q", i, got, want)
		}
	}
	if cases[5].Want == cases[6].Want || cases[9].Want == cases[7].Want ||
		cases[9].Want == "3:true:-128:0:5:slice::5,0,255,7,9" {
		log.Fatal("strict capacity, rejecting negative counts, or clamping offset escaped detection")
	}
	file, err := os.CreateTemp("", "golean-review-reslice-")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(cases))
	must(file.Close())
	must(os.Setenv("GOLEAN_BUFFER_TRY_GROW_BY_RESLICE_CASES", file.Name()))
	const dir = "difftest/bytes/buffer-try-grow-by-reslice/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferTryGrowByReslice", dir+"bridge_test.go", dir+"grow_test.go"))
	fmt.Printf("Buffer.tryGrowByReslice: %d native and Lean fixtures, three incorrect variants distinguished\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
