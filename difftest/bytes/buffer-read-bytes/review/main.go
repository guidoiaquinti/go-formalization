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
		{input{Tag: 1, Flag: 127, Count: 255}, "eof:nil:0:-1::"},
		{input{Tag: 2, Offset: 4, Flag: 4, Stored: stored, Spare: spare, Count: 255}, "eof:nil:4:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare, Count: 255}, "ok:slice:8:2,255:0,0,0,0,0,0:3:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Offset: 3, Flag: 2, Stored: stored, Spare: spare, Count: 4}, "ok:slice:8:4:0,0,0,0,0,0,0:4:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: stored, Spare: spare, Count: 7}, "eof:slice:8:2,255,4:0,0,0,0,0:4:-1:255,2,255,4:7,9"},
		{input{Tag: 2, Stored: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, Count: 9}, "ok:slice:16:1,2,3,4,5,6,7,8,9:0,0,0,0,0,0,0:9:-1:1,2,3,4,5,6,7,8,9:"},
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
	for i, got := range check.LeanRows("godiff-bytes-buffer-read-bytes", len(cases), &model) {
		if got != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].Want)
		}
	}

	backing := []byte("ab!cd!")
	buffer := bytes.NewBuffer(backing)
	first, err := buffer.ReadBytes('!')
	if err != nil || string(first) != "ab!" {
		log.Fatal("first result omitted delimiter")
	}
	first[0] = 'X'
	if string(backing) != "ab!cd!" {
		log.Fatal("mutating returned bytes changed original backing")
	}
	buffer.Reset()
	buffer.WriteString("uvwxyz")
	if string(first) != "Xb!" {
		log.Fatal("buffer reuse changed prior result")
	}
	last, err := buffer.ReadBytes('!')
	if err != io.EOF || string(last) != "uvwxyz" {
		log.Fatal("missing delimiter did not return the full remainder with EOF")
	}
	last, err = buffer.ReadBytes('!')
	if err != io.EOF || last != nil || buffer.UnreadByte() != nil || buffer.Len() != 1 {
		log.Fatal("exhausted ReadBytes did not return nil EOF and install ordinary read flag")
	}
	if cases[3].Want == "ok:slice:8:2:0,0,0,0,0,0,0:2:-1:255,2,255,4:7,9" ||
		cases[2].Want == "eof:nil:4:4:255,2,255,4:7,9" ||
		cases[6].Want == "ok:slice:9:1,2,3,4,5,6,7,8,9::9:-1:1,2,3,4,5,6,7,8,9:" {
		log.Fatal("one incorrect variant escaped")
	}
	f, err := os.CreateTemp("", "golean-read-bytes-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_BUFFER_READ_BYTES_CASES", f.Name()))
	const dir = "difftest/bytes/buffer-read-bytes/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferReadBytes", dir+"bridge_test.go", dir+"read_bytes_test.go"))
	fmt.Printf("Buffer.ReadBytes: %d native and Lean fixed cases, public backing isolation and unread controls passed\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
