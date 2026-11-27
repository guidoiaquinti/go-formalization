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
	stored := []byte{1, 2, 3, 4, 5, 6}
	spare := []byte{0, 255}
	cases := []fixture{
		{input{Count: -1}, "nil-receiver"},
		{input{Tag: 1, Flag: -7, Count: 0}, "ok:0:-7:0:0:nil:false::"},
		{input{Tag: 2, Offset: 6, Flag: -128, Stored: stored, Spare: spare, Count: 2}, "ok:0:0:0:8:slice:true:1,2:3,4,5,6,0,255"},
		{input{Tag: 2, Offset: 6, Flag: -128, Stored: stored, Spare: spare, Count: -1}, "bounds:0:0:8:slice:true::1,2,3,4,5,6,0,255"},
		{input{Tag: 2, Offset: 6, Flag: -128, Stored: stored, Spare: spare, Count: math.MaxInt}, "too-large:0:0:8:slice:true::1,2,3,4,5,6,0,255"},
		{input{Tag: 2, Offset: 4, Flag: 4, Stored: stored, Count: 1}, "ok:2:4:0:6:slice:true:5,6,3:4,5,6"},
		{input{Tag: 2, Offset: 3, Flag: 3, Stored: stored, Count: 1}, "ok:3:3:0:8:slice:false:4,5,6,0:0,0,0,0"},
		{input{Tag: 2, Offset: 2, Flag: 2, Stored: stored, Spare: spare, Count: 2}, "ok:6:2:2:8:slice:true:1,2,3,4,5,6,0,255:"},
		{input{Tag: 2, Offset: 2, Flag: 2, Stored: stored, Spare: spare, Count: -6}, "ok:6:2:2:8:slice:true::1,2,3,4,5,6,0,255"},
		{input{Tag: 2, Offset: 2, Flag: 2, Stored: stored, Spare: spare, Count: math.MinInt}, "bounds:2:2:8:slice:true:1,2,3,4,5,6:0,255"},
	}
	var model bytes.Buffer
	for _, fixture := range cases {
		c := fixture.Input
		fmt.Fprintf(&model, "%d %d %d", c.Tag, c.Count, math.MaxInt)
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
	for i, got := range check.LeanRows("godiff-bytes-buffer-grow-private", len(cases), &model) {
		if want := cases[i].Want; got != want {
			log.Fatalf("fixture %d: Lean=%q want=%q", i, got, want)
		}
	}
	if cases[2].Want == "ok:6:-128:6:8:slice:true:1,2,3,4,5,6,0,255:" ||
		cases[3].Want == "bounds:-128:6:8:slice:true:1,2,3,4,5,6:0,255" ||
		cases[5].Want == "ok:2:4:0:6:slice:true:5,6:3,4,5,6" {
		log.Fatal("reslice-first, rollback on panic, or hidden reservation escaped detection")
	}
	file, err := os.CreateTemp("", "golean-review-reslice-")
	must(err)
	defer os.Remove(file.Name())
	must(json.NewEncoder(file).Encode(cases))
	must(file.Close())
	must(os.Setenv("GOLEAN_BUFFER_GROW_PRIVATE_CASES", file.Name()))
	const dir = "difftest/bytes/buffer-grow-private/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferGrowPrivate", dir+"bridge_test.go", dir+"grow_test.go"))
	fmt.Printf("Buffer.grow: %d native and Lean fixtures, three incorrect variants distinguished\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
