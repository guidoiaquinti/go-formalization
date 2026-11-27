package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Tag, Offset, Flag int
	Stored, Spare     []byte
	Rune              rune
}

type fixture struct {
	Input input
	Want  string
}

func result(stored, spare []byte, width, offset int, reused bool) string {
	return fmt.Sprintf("ok:%d:0:%d:%d:slice:%t:%s:%s", width, offset, len(stored)+len(spare), reused, decimal(stored), decimal(spare))
}

func decimal(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func main() {
	full := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	two, replacement, four := []byte{195, 169}, []byte{239, 191, 189}, []byte{240, 144, 128, 128}
	cases := []fixture{
		{input{Rune: 'a'}, "panic"},
		{input{Rune: -1}, "panic"},
		{input{Tag: 1, Flag: -128, Rune: -2147483648}, result(replacement, make([]byte, 61), 3, 0, false)},
		{input{Tag: 2, Flag: 127, Rune: 0xd800}, result(replacement, make([]byte, 5), 3, 0, false)},
		{input{Tag: 2, Flag: 3, Rune: 0x110000}, result(replacement, make([]byte, 5), 3, 0, false)},
		{input{Tag: 2, Stored: full[:2], Spare: full[2:4], Rune: 'é'}, result([]byte{1, 2, 195, 169}, make([]byte, 4), 2, 0, false)},
		{input{Tag: 2, Offset: 1, Flag: 2, Stored: full[:2], Spare: full[2:6], Rune: 'é'}, result([]byte{1, 2, 195, 169}, full[4:6], 2, 1, true)},
		{input{Tag: 2, Offset: 3, Flag: 1, Stored: full[:3], Spare: full[3:5], Rune: 'é'}, result(two, full[2:5], 2, 0, true)},
		{input{Tag: 2, Offset: 2, Flag: 1, Stored: full[:2], Rune: 'é'}, result(two, make([]byte, 6), 2, 0, false)},
		{input{Tag: 2, Offset: 10, Flag: 1, Stored: full, Rune: 0x10000}, result([]byte{11, 12, 240, 144, 128, 128}, full[6:], 4, 0, true)},
		{input{Tag: 2, Offset: 9, Flag: 1, Stored: full, Rune: 0x10000}, result(append([]byte{10, 11, 12}, four...), make([]byte, 9), 4, 0, false)},
		{input{Tag: 2, Offset: 2, Flag: 1, Stored: full[:2], Spare: full[2:3], Rune: 'z'}, result([]byte{1, 2, 'z'}, nil, 1, 2, true)},
	}
	var rows bytes.Buffer
	for _, f := range cases {
		c := f.Input
		fmt.Fprintf(&rows, "%d %d", c.Tag, c.Rune)
		if c.Tag != 0 {
			fmt.Fprintf(&rows, " %d %d %d", len(c.Stored), c.Offset, c.Flag)
			for _, block := range [][]byte{c.Stored, c.Spare} {
				for _, b := range block {
					fmt.Fprintf(&rows, " %d", b)
				}
			}
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-write-rune", len(cases), &rows) {
		if want := cases[i].Want; got != want {
			log.Fatalf("fixture %d: Lean=%q want %q", i, got, want)
		}
	}

	// Capture the live original array through both sides of the four-byte
	// compaction threshold. The two-byte rune still reserves four bytes.
	for _, offset := range []int{9, 10} {
		original := bytes.Clone(full)
		b := bytes.NewBuffer(original[:len(original):len(original)])
		b.Next(offset)
		width, err := b.WriteRune('é')
		if width != 2 || err != nil || b.UnreadByte() == nil {
			log.Fatal("WriteRune returned the wrong width/error or left the read flag valid")
		}
		if offset == 10 {
			if !bytes.Equal(original, []byte{11, 12, 195, 169, 5, 6, 7, 8, 9, 10, 11, 12}) || &b.Bytes()[0] != &original[0] {
				log.Fatal("compaction failed to update the original alias")
			}
		} else if !bytes.Equal(original, full) || &b.Bytes()[0] == &original[0] {
			log.Fatal("reallocation changed the original alias")
		}
	}
	if cases[5].Want == result([]byte{1, 2, 195, 169}, nil, 2, 0, true) ||
		cases[6].Want == result([]byte{1, 2, 195, 169}, full[4:6], 4, 1, true) ||
		cases[2].Want == result([]byte{0}, make([]byte, 63), 1, 0, false) {
		log.Fatal("incorrect reservation, reported width, or negative-rune encoding escaped detection")
	}

	trace, err := os.CreateTemp("", "golean-write-rune-review-")
	must(err)
	defer os.Remove(trace.Name())
	must(json.NewEncoder(trace).Encode(cases))
	must(trace.Close())
	must(os.Setenv("GOLEAN_BUFFER_WRITE_RUNE_CASES", trace.Name()))
	const data = "difftest/bytes/buffer-write-rune/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferWriteRune", data+"bridge_test.go", data+"write_rune_test.go"))
	fmt.Printf("Buffer.WriteRune: %d native and Lean fixtures plus two live alias checks; encoded-width reservation, reported reserved width, and negative ASCII each fail a fixture\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
