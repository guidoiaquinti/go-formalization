package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	stored, spare := []byte{1, 2, 255}, []byte{7, 9}
	state := func(count int) input { return input{2, 1, -128, stored, spare, count} }
	cases := []fixture{
		{input{Count: -1}, "nil-receiver"},
		{input{Tag: 1, Flag: 127}, "ok:nil:0:127::"},
		{input{Tag: 1, Flag: 127, Count: 1}, "eof:nil:0:127::"},
		{input{Tag: 2, Offset: 3, Flag: 4, Stored: stored, Spare: spare}, "ok:slice:2::7,9:3:4:1,2,255:7,9"},
		{state(0), "ok:slice:4::2,255,7,9:1:-128:1,2,255:7,9"},
		{state(1), "ok:slice:4:2:255,7,9:1:-128:1,2,255:7,9"},
		{state(2), "ok:slice:4:2,255:7,9:1:-128:1,2,255:7,9"},
		{state(math.MaxInt), "eof:slice:4:2,255:7,9:1:-128:1,2,255:7,9"},
		{state(math.MinInt), "bounds:1:-128:1,2,255:7,9"},
	}
	var model bytes.Buffer
	for _, f := range cases {
		c := f.Input
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
	for i, got := range check.LeanRows("godiff-bytes-buffer-peek", len(cases), &model) {
		if got != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].Want)
		}
	}

	raw := []byte{1, 2, 3, 4, 99}
	b := bytes.NewBuffer(raw[:4])
	b.Next(1)
	peek, err := b.Peek(1)
	if err != nil || len(peek) != 1 || cap(peek) != 4 {
		log.Fatal("short Peek capacity or error changed")
	}
	peek[:cap(peek)][3] = 55
	if raw[4] != 55 || b.UnreadByte() != nil {
		log.Fatal("short Peek did not alias spare or invalidated prior unread")
	}
	peek, err = b.Peek(math.MaxInt)
	if err != io.EOF || len(peek) != 4 || &peek[0] != &raw[0] {
		log.Fatal("overlong Peek did not return original remainder with EOF")
	}
	if cases[6].Want == cases[7].Want || cases[5].Want == "ok:slice:1:2::1:-128:1,2,255:7,9" ||
		cases[8].Want == cases[4].Want {
		log.Fatal("one incorrect variant escaped")
	}
	f, err := os.CreateTemp("", "golean-peek-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_BUFFER_PEEK_CASES", f.Name()))
	const dir = "difftest/bytes/buffer-peek/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferPeek", dir+"bridge_test.go", dir+"peek_test.go"))
	fmt.Printf("Buffer.Peek: %d native and Lean fixed cases, public alias and unread controls passed\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
