package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
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

const start = "bytes.Reader.UnreadRune: at beginning of slice"
const previous = "bytes.Reader.UnreadRune: previous operation was not ReadRune"

func describe(c input, cursor int64, marker int) string {
	kind := "slice"
	if c.Tag == 1 {
		kind = "nil"
	}
	wire := func(raw []byte) string {
		fields := make([]string, len(raw))
		for i, v := range raw {
			fields[i] = strconv.Itoa(int(v))
		}
		return strings.Join(fields, ",")
	}
	return fmt.Sprintf("%s:%d:%d:%s:%s", kind, cursor, marker, wire(c.Raw), wire(c.Spare))
}
func read(c input, value, width, marker int, status string) string {
	return fmt.Sprintf("%d:%d:%s:%s", value, width, status, describe(c, c.Cursor+int64(width), marker))
}
func unread(c input, status string, cursor int64, marker int) string {
	return status + ":" + describe(c, cursor, marker)
}

func main() {
	literal := input{2, 1, -7, 1, []byte{9, 239, 191, 189, 33}, []byte{255}}
	malformed := input{2, 2, 0, 1, []byte{9, 237, 160, 128}, []byte{80}}
	truncated := input{2, 1, 1, 2, []byte{0, 1, 226, 130}, []byte{172}}
	four := input{2, 1, 2, 3, []byte{1, 2, 3, 240, 159, 145, 139}, []byte{10}}
	exhausted := input{2, 3, 8, math.MaxInt64, []byte{0}, []byte{65}}
	atStart := input{2, 1, 99, 0, []byte{127}, nil}
	empty := input{1, 0, -8, 0, nil, nil}
	forward := input{2, 1, 14, 1, []byte{255}, []byte{94}}
	cases := []fixture{
		{input{}, "panic|panic"},
		{literal, read(literal, 65533, 3, 1, "ok") + "|" + unread(literal, previous, 1, -7)},
		{malformed, read(malformed, 65533, 1, 1, "ok") + "|" + unread(malformed, "ok", 0, -1)},
		{truncated, read(truncated, 65533, 1, 2, "ok") + "|" + unread(truncated, "ok", 1, -1)},
		{four, read(four, 128075, 4, 3, "ok") + "|" + unread(four, "ok", 2, -1)},
		{exhausted, read(exhausted, 0, 0, -1, "EOF") + "|" + unread(exhausted, "ok", 8, -1)},
		{atStart, read(atStart, 127, 1, 0, "ok") + "|" + unread(atStart, start, 0, 99)},
		{empty, read(empty, 0, 0, -1, "EOF") + "|" + unread(empty, start, 0, -8)},
		{forward, read(forward, 0, 0, -1, "EOF") + "|" + unread(forward, "ok", 14, -1)},
	}

	var wire bytes.Buffer
	for _, f := range cases {
		c := f.Input
		fmt.Fprintf(&wire, "%d %d %d %d", c.Tag, len(c.Raw), c.Cursor, c.Previous)
		for _, part := range [][]byte{c.Raw, c.Spare} {
			for _, v := range part {
				fmt.Fprintf(&wire, " %d", v)
			}
		}
		fmt.Fprintln(&wire)
	}
	var rejects [4]int
	for i, got := range check.LeanRows("godiff-bytes-reader-rune", len(cases), &wire) {
		f := cases[i]
		if got != f.Want {
			log.Fatalf("row %d: Lean %q, want %q", i, got, f.Want)
		}
		c := f.Input
		gotRead, gotUnread, _ := strings.Cut(got, "|")
		if strings.HasPrefix(gotRead, "65533:") && read(c, 65533, 1, int(c.Cursor), "ok") != gotRead {
			rejects[0]++
		}
		if strings.HasPrefix(gotRead, "0:0:EOF:") && read(c, 0, 0, c.Previous, "EOF") != gotRead {
			rejects[1]++
		}
		if strings.HasPrefix(gotUnread, "ok:") && unread(c, "ok", c.Cursor-1, -1) != gotUnread {
			rejects[2]++
		}
		if c.Tag != 0 && c.Cursor == 0 {
			mutant := unread(c, previous, c.Cursor, c.Previous)
			if c.Previous >= 0 {
				mutant = unread(c, "ok", int64(c.Previous), -1)
			}
			if mutant != gotUnread {
				rejects[3]++
			}
		}
	}
	if rejects != [4]int{1, 3, 2, 2} {
		log.Fatalf("unexpected mutant rejection counts: %v", rejects)
	}
	r := bytes.NewReader([]byte{9, 239, 191, 189, 255})
	r.ReadByte()
	if v, w, e := r.ReadRune(); v != 65533 || w != 3 || e != nil {
		log.Fatal("valid replacement rune")
	}
	if r.UnreadRune() != nil || r.UnreadRune() == nil {
		log.Fatal("UnreadRune did not consume the marker")
	}
	if v, w, e := r.ReadRune(); v != 65533 || w != 3 || e != nil {
		log.Fatal("valid replacement rune after rewind")
	}
	if v, w, e := r.ReadRune(); v != 65533 || w != 1 || e != nil {
		log.Fatal("malformed byte")
	}
	if v, w, e := r.ReadRune(); v != 0 || w != 0 || e != io.EOF || r.UnreadRune() == nil {
		log.Fatal("EOF did not clear the marker")
	}

	f, err := os.CreateTemp("", "golean-reader-rune-review-")
	must(err)
	defer os.Remove(f.Name())
	must(json.NewEncoder(f).Encode(cases))
	must(f.Close())
	must(os.Setenv("GOLEAN_READER_RUNE_CASES", f.Name()))
	const data = "difftest/bytes/reader-rune/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanReaderRune", data+"bridge_test.go", data+"reader_rune_test.go"))
	fmt.Printf("Reader rune: %d independent native and Lean rows, one public sequence, mutant rejections %v\n", len(cases), rejects)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
