package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Offset, Previous int
	Cursor                int64
	Raw, Spare            []byte
}

type fixture struct {
	Input input
	Want  string
}

const atStart = "bytes.Reader.UnreadByte: at beginning of slice"

func main() {
	raw, spare := []byte{0, 255}, []byte{77}
	cases := []fixture{
		{input{}, "panic|panic"},
		{input{Tag: 1, Previous: 9}, "0:EOF:nil:0:-1::|" + atStart + ":nil:0:9::"},
		{input{Tag: 1, Previous: 9, Cursor: 1}, "0:EOF:nil:1:-1::|ok:nil:0:-1::"},
		{input{Tag: 2, Offset: 1, Previous: 8, Raw: raw, Spare: spare}, "0:ok:slice:1:-1:0,255:77|" + atStart + ":slice:0:8:0,255:77"},
		{input{Tag: 2, Offset: 2, Previous: 8, Cursor: 1, Raw: raw, Spare: spare}, "255:ok:slice:2:-1:0,255:77|ok:slice:0:-1:0,255:77"},
		{input{Tag: 2, Offset: 3, Previous: 8, Cursor: 2, Raw: raw, Spare: spare}, "0:EOF:slice:2:-1:0,255:77|ok:slice:1:-1:0,255:77"},
		{input{Tag: 2, Offset: 3, Previous: 8, Cursor: 3, Raw: raw, Spare: spare}, "0:EOF:slice:3:-1:0,255:77|ok:slice:2:-1:0,255:77"},
		{input{Tag: 2, Offset: 3, Previous: -1, Cursor: math.MaxInt64, Raw: raw, Spare: spare}, "0:EOF:slice:9223372036854775807:-1:0,255:77|ok:slice:9223372036854775806:-1:0,255:77"},
		{input{Tag: 2, Offset: 3, Previous: 0, Spare: spare}, "0:EOF:slice:0:-1::77|" + atStart + ":slice:0:0::77"},
	}
	var model bytes.Buffer
	var rejects [3]int
	for _, f := range cases {
		c := f.Input
		fmt.Fprintf(&model, "%d %d %d %d", c.Tag, len(c.Raw), c.Cursor, c.Previous)
		for _, part := range [][]byte{c.Raw, c.Spare} {
			for _, b := range part {
				fmt.Fprintf(&model, " %d", b)
			}
		}
		fmt.Fprintln(&model)
		if c.Tag == 0 {
			continue
		}
		read, unread, _ := strings.Cut(f.Want, "|")
		if strings.HasPrefix(read, "0:EOF:") && c.Previous != -1 {
			rejects[0]++
		}
		if c.Cursor == 0 && c.Previous != -1 {
			rejects[1]++
		}
		if c.Cursor > int64(len(c.Raw)) && strings.HasPrefix(unread, "ok:") {
			rejects[2]++
		}
	}
	for i, got := range check.LeanRows("godiff-bytes-reader-byte", len(cases), &model) {
		if got != cases[i].Want {
			log.Fatalf("fixture %d: Lean=%s want=%s", i, got, cases[i].Want)
		}
	}
	if rejects != [3]int{5, 3, 3} {
		log.Fatalf("incorrect EOF/failed-unread/clamp variants have unexpected mismatch counts %v", rejects)
	}

	r := bytes.NewReader([]byte("é"))
	if _, width, err := r.ReadRune(); width != 2 || err != nil {
		log.Fatal("rune setup")
	}
	if value, err := r.ReadByte(); value != 0 || err != io.EOF || r.UnreadRune() == nil {
		log.Fatal("ReadByte EOF did not invalidate previous rune")
	}
	if r.UnreadByte() != nil {
		log.Fatal("UnreadByte after EOF did not rewind")
	}
	if value, err := r.ReadByte(); value != 0xa9 || err != nil {
		log.Fatal("UnreadByte rewound more than one byte")
	}
	r = bytes.NewReader([]byte("é"))
	r.ReadRune()
	if r.UnreadByte() != nil || r.UnreadByte() != nil {
		log.Fatal("UnreadByte requires prior successful read")
	}
	failure1, failure2 := r.UnreadByte(), r.UnreadByte()
	if failure1 == nil || failure2 == nil || failure1 == failure2 || failure1.Error() != atStart || failure2.Error() != atStart {
		log.Fatal("at-start errors were not distinct with identical message")
	}
	f, err := os.CreateTemp("", "golean-reader-byte-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_READER_BYTE_CASES", f.Name()))
	const data = "difftest/bytes/reader-byte/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReaderByte", data+"bridge_test.go", data+"reader_byte_test.go"))
	fmt.Printf("Reader byte: %d native/Lean fixed rows, two public sequences and incorrect variants %v\n", len(cases), rejects)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
