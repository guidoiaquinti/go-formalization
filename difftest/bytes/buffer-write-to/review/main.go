package main

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Failed            bool
}

type fixture struct {
	Input input
	Want  string
}

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(raw []byte) (int, error) { return f(raw) }

func main() {
	stored, spare := []byte{1, 2, 255}, []byte{7, 9}
	makeInput := func(count int, fail bool) input {
		return input{2, 1, 4, stored, spare, count, fail}
	}
	cases := []fixture{
		{input{}, "nil-receiver"},
		{input{Tag: 1, Flag: 127, Count: 99, Failed: true}, "ok:0:-:nil:0:0::"},
		{input{Tag: 2, Offset: 3, Flag: 4, Stored: stored, Spare: spare, Count: 99, Failed: true}, "ok:0:-:slice:0:0::1,2,255,7,9"},
		{makeInput(0, false), "short:0:2,255/7,9:slice:1:0:1,2,255:7,9"},
		{makeInput(1, false), "short:1:2,255/7,9:slice:2:0:1,2,255:7,9"},
		{makeInput(1, true), "writer:1:2,255/7,9:slice:2:0:1,2,255:7,9"},
		{makeInput(2, true), "writer:2:2,255/7,9:slice:3:0:1,2,255:7,9"},
		{makeInput(2, false), "ok:2:2,255/7,9:slice:0:0::1,2,255,7,9"},
		{makeInput(3, true), "invalid:0:2,255/7,9:slice:1:0:1,2,255:7,9"},
	}
	var model bytes.Buffer
	for _, f := range cases {
		c := f.Input
		fmt.Fprint(&model, c.Tag)
		if c.Tag != 0 {
			failed := 0
			if c.Failed {
				failed = 1
			}
			fmt.Fprintf(&model, " %d %d %d %d %d", c.Count, failed, len(c.Stored), c.Offset, c.Flag)
			for _, part := range [][]byte{c.Stored, c.Spare} {
				for _, b := range part {
					fmt.Fprintf(&model, " %d", b)
				}
			}
		}
		fmt.Fprintln(&model)
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-write-to", len(cases), &model) {
		if got != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].Want)
		}
	}
	for _, failure := range []error{nil, errors.New("observer failure")} {
		backing := []byte{1, 2, 3, 7, 9}
		b := bytes.NewBuffer(backing[:3])
		b.Next(1)
		calls := 0
		n, err := b.WriteTo(writerFunc(func(raw []byte) (int, error) {
			calls++
			if !bytes.Equal(raw, []byte{2, 3}) || !bytes.Equal(raw[:cap(raw)], []byte{2, 3, 7, 9}) || b.UnreadByte() == nil {
				log.Fatal("wrong offered view or read flag not cleared before callback")
			}
			return len(raw), failure
		}))
		if calls != 1 || n != 2 || err != failure || b.Len() != 0 {
			log.Fatal("write count or error identity changed")
		}
		wantCap := b.Cap()
		if failure != nil {
			wantCap -= 3
		}
		if cap(b.Bytes()) != wantCap {
			log.Fatal("full failed write must not reset")
		}
		n, err = b.WriteTo(nil)
		if n != 0 || err != nil || cap(b.Bytes()) != b.Cap() {
			log.Fatal("exhausted WriteTo attempted a callback")
		}
	}
	if cases[4].Want == cases[5].Want || cases[6].Want == cases[7].Want || cases[8].Want == cases[5].Want {
		log.Fatal("one incorrect reply ordering escaped")
	}
	f, err := os.CreateTemp("", "golean-write-to-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_BUFFER_WRITE_TO_CASES", f.Name()))
	const data = "difftest/bytes/buffer-write-to/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferWriteTo", data+"bridge_test.go", data+"write_to_test.go"))
	short := bytes.NewBufferString("raw")
	if n, err := short.WriteTo(writerFunc(func([]byte) (int, error) { return 1, nil })); n != 1 || err != io.ErrShortWrite || short.String() != "aw" {
		log.Fatalf("short write returned %d,%v with remainder %q", n, err, short.String())
	}
	fmt.Printf("Buffer.WriteTo: %d native and Lean fixtures plus three public callback sequences passed\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
