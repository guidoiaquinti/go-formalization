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
	stored, spare := []byte{255, 2, 255, 4}, []byte{7, 9}
	cases := []fixture{
		{input{Count: 255}, "nil-receiver"},
		{input{Tag: 1, Flag: 127, Count: 255}, "eof::0:-1::"},
		{input{Tag: 2, Offset: 4, Flag: 4, Stored: stored, Spare: spare, Count: 255}, "eof::4:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare, Count: 255}, "ok:2,255:3:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Offset: 3, Flag: 2, Stored: stored, Spare: spare, Count: 4}, "ok:4:4:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare, Count: 7}, "eof:2,255,4:4:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Stored: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, Count: 9}, "ok:1,2,3,4,5,6,7,8,9:9:-1:1,2,3,4,5,6,7,8,9:"},
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
	for i, got := range check.LeanRows("godiff-bytes-buffer-read-string", len(cases), &model) {
		if got != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].Want)
		}
	}

	backing := []byte{0xff, 0x00, '!', 'C', 'D', '!'}
	buffer := bytes.NewBuffer(backing)
	first, err := buffer.ReadString('!')
	if err != nil || first != "\xff\x00!" {
		log.Fatal("first result did not preserve raw bytes and delimiter")
	}
	backing[0] = 'X'
	buffer.Reset()
	buffer.WriteString("uvwxyz")
	if first != "\xff\x00!" {
		log.Fatal("mutating and then reusing the backing changed a prior string")
	}
	last, err := buffer.ReadString('!')
	if err != io.EOF || last != "uvwxyz" {
		log.Fatal("missing delimiter did not return entire remaining string and EOF")
	}
	last, err = buffer.ReadString('!')
	if err != io.EOF || last != "" || buffer.UnreadByte() != nil || buffer.Len() != 1 {
		log.Fatal("exhausted ReadString lacked EOF, empty value or ordinary read flag")
	}
	if cases[3].Want == "ok:2:2:-1:255,2,255,4:7,9" ||
		cases[2].Want == "eof::4:4:255,2,255,4:7,9" ||
		cases[5].Want == "ok:2,255,4:4:-1:255,2,255,4:7,9" {
		log.Fatal("one incorrect variant escaped")
	}
	f, err := os.CreateTemp("", "golean-read-string-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_BUFFER_READ_STRING_CASES", f.Name()))
	const dir = "difftest/bytes/buffer-read-string/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferReadString", dir+"bridge_test.go", dir+"read_string_test.go"))
	fmt.Printf("Buffer.ReadString: %d native and Lean fixed cases, public backing isolation and unread controls passed\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
