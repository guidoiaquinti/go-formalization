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
	Byte              byte
}

type fixture struct {
	Input input
	Want  string
}

func result(stored, spare []byte, offset int, reused bool) string {
	return fmt.Sprintf("ok:0:%d:%d:slice:%t:%s:%s", offset, len(stored)+len(spare), reused, decimal(stored), decimal(spare))
}

func decimal(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func main() {
	const added = 0xff
	full := []byte{1, 2, 3, 4, 5, 6}
	large := bytes.Repeat([]byte{0x51}, 256)
	cases := []fixture{
		{input{Byte: added}, "panic"},
		{input{Tag: 1, Flag: -128, Byte: added}, result([]byte{added}, make([]byte, 63), 0, false)},
		{input{Tag: 2, Flag: 127, Byte: added}, result([]byte{added}, make([]byte, 7), 0, false)},
		{input{Tag: 2, Offset: 6, Flag: -1, Stored: full, Spare: []byte{99, 98}, Byte: added}, result([]byte{1, 2, 3, 4, 5, 6, added}, []byte{98}, 6, true)},
		{input{Tag: 2, Offset: 6, Flag: -1, Stored: full, Byte: added}, result([]byte{added}, full[1:], 0, true)},
		{input{Tag: 2, Offset: 4, Flag: 4, Stored: full, Byte: added}, result([]byte{5, 6, added}, full[3:], 0, true)},
		{input{Tag: 2, Offset: 3, Flag: 3, Stored: full, Byte: added}, result([]byte{4, 5, 6, added}, make([]byte, 4), 0, false)},
		{input{Tag: 2, Flag: -1, Stored: large, Byte: added}, result(append(bytes.Clone(large), added), make([]byte, 255), 0, false)},
		{input{Tag: 2, Offset: 4, Flag: 4, Stored: full[:5], Byte: added}, result([]byte{5, added}, full[2:5], 0, true)},
		{input{Tag: 2, Offset: 3, Flag: 3, Stored: full[:5], Byte: added}, result([]byte{4, 5, added}, make([]byte, 5), 0, false)},
	}
	var rows bytes.Buffer
	for _, f := range cases {
		c := f.Input
		fmt.Fprintf(&rows, "%d %d", c.Tag, c.Byte)
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
	for i, got := range check.LeanRows("godiff-bytes-buffer-write-byte", len(cases), &rows) {
		if want := cases[i].Want; got != want {
			log.Fatalf("fixture %d: Lean=%q want %q", i, got, want)
		}
	}

	// A live original slice outside the receiver observes in-place copies.
	for _, c := range []struct {
		offset  int
		compact bool
	}{{4, true}, {3, false}} {
		original := bytes.Clone(full)
		b := bytes.NewBuffer(original[:len(original):len(original)])
		b.Next(c.offset)
		if err := b.WriteByte(added); err != nil {
			log.Fatal(err)
		}
		if err := b.UnreadByte(); err == nil {
			log.Fatal("WriteByte left the previous read valid")
		}
		if c.compact {
			if !bytes.Equal(original, []byte{5, 6, added, 4, 5, 6}) || &b.Bytes()[0] != &original[0] {
				log.Fatal("compaction failed to change the live alias")
			}
		} else if !bytes.Equal(original, full) || &b.Bytes()[0] == &original[0] {
			log.Fatal("reallocation changed the live alias")
		}
	}

	var strictHalf, alwaysCompact, nilEqualsEmpty int
	if got := result([]byte{5, 6, added}, make([]byte, 5), 0, false); got != cases[5].Want {
		strictHalf++
	}
	if got := result([]byte{4, 5, 6, added}, full[4:], 0, true); got != cases[6].Want {
		alwaysCompact++
	}
	if cases[1].Want != cases[2].Want {
		nilEqualsEmpty++
	}
	if strictHalf != 1 || alwaysCompact != 1 || nilEqualsEmpty != 1 {
		log.Fatal("incorrect variants were not all detected")
	}
	trace, err := os.CreateTemp("", "golean-write-byte-review-")
	must(err)
	defer os.Remove(trace.Name())
	must(json.NewEncoder(trace).Encode(cases))
	must(trace.Close())
	must(os.Setenv("GOLEAN_BUFFER_WRITE_BYTE_CASES", trace.Name()))
	const data = "difftest/bytes/buffer-write-byte/testdata/"
	must(check.TestInPackage("bytes", "TestGoLeanBufferWriteByte", data+"bridge_test.go", data+"write_byte_test.go"))
	fmt.Printf("Buffer.WriteByte: %d native and Lean fixtures plus two live alias checks; strict half, always compact and nil equivalence each fail a fixture\n", len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
