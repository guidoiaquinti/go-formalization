package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	cases := []fixture{
		{input{}, "panic"},
		{input{Tag: 1, Flag: 127}, "ok:0:nil:0:0:::"},
		{input{Tag: 1, Flag: 127, Count: 1}, "eof:0:nil:0:0:167::"},
		{input{Tag: 2, Offset: 3, Flag: 4, Stored: stored, Spare: spare}, "ok:0:slice:0:0:::1,2,255,7,9"},
		{input{Tag: 2, Offset: 3, Flag: 4, Stored: stored, Spare: spare, Count: 2}, "eof:0:slice:0:0:167,198::1,2,255,7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare}, "ok:0:slice:1:0::1,2,255:7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare, Count: 1}, "ok:1:slice:2:-1:2:1,2,255:7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare, Count: 2}, "ok:2:slice:3:-1:2,255:1,2,255:7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare, Count: 3}, "ok:2:slice:3:-1:2,255,229:1,2,255:7,9"},
	}
	var model bytes.Buffer
	for _, f := range cases {
		c := f.Input
		fmt.Fprint(&model, c.Tag)
		if c.Tag != 0 {
			fmt.Fprintf(&model, " %d %d %d %d", len(c.Stored), len(c.Spare), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&model, " %d", b)
				}
			}
			for i := range c.Count {
				fmt.Fprintf(&model, " %d", byte(167+i*31))
			}
		}
		fmt.Fprintln(&model)
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-read", len(cases), &model) {
		if got != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].Want)
		}
	}

	b := bytes.NewBufferString("abc")
	dst := []byte{0, 0, 0, 99}
	if n, err := b.Read(dst); n != 3 || err != nil || string(dst) != "abcc" {
		log.Fatalf("short final read n=%d err=%v bytes=%x", n, err, dst)
	}
	if n, err := b.Read(nil); n != 0 || err != nil || b.UnreadByte() == nil {
		log.Fatal("drained zero-length Read did not reset")
	}
	if n, err := b.Read(dst); n != 0 || err != io.EOF {
		log.Fatal("subsequent nonempty read lacked EOF")
	}
	if cases[3].Want == "ok:0:slice:3:4::1,2,255:7,9" ||
		cases[5].Want == "ok:0:slice:1:2::1,2,255:7,9" ||
		cases[8].Want == "eof:2:slice:3:-1:2,255,229:1,2,255:7,9" {
		log.Fatal("one incorrect variant escaped")
	}
	f, err := os.CreateTemp("", "golean-read-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_BUFFER_READ_CASES", f.Name()))
	const dir = "difftest/bytes/buffer-read/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferRead", dir+"bridge_test.go", dir+"read_test.go"))
	fmt.Printf("Buffer.Read: %d native and Lean fixed cases, public drained-state and partial-read controls passed\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
