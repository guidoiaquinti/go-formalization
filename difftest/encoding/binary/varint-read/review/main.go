package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"slices"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type reader struct {
	bytes []byte
	end   error
	calls int
	limit int
}

func (r *reader) ReadByte() (byte, error) {
	index := r.calls
	r.calls++
	if r.calls > r.limit {
		panic("extra ReadByte call")
	}
	if index >= len(r.bytes) {
		return 0xff, r.end
	}
	return r.bytes[index], nil
}

type fixture struct {
	name  string
	bytes []byte
	end   error
	tag   int
	u     uint64
	s     int64
	err   error
	label string
	calls int
}

func ten(first, last byte) []byte {
	return append(append([]byte{first}, slices.Repeat([]byte{128}, 8)...), last, 42)
}

func main() {
	custom := errors.New("EOF")
	wrapped := fmt.Errorf("reader: %w", io.EOF)
	probe := &reader{bytes: ten(128, 3), limit: 10}
	_, overflow := binary.ReadUvarint(probe)
	if overflow == nil || overflow.Error() != "binary: varint overflows a 64-bit integer" {
		log.Fatalf("unchecked overflow sentinel: %v", overflow)
	}
	cases := []fixture{
		{"empty", nil, io.EOF, 0, 0, 0, io.EOF, "eof", 1},
		{"partial negative", []byte{129}, io.EOF, 0, 1, -1, io.ErrUnexpectedEOF, "unexpected", 2},
		{"partial positive", []byte{134}, io.ErrUnexpectedEOF, 1, 6, 3, io.ErrUnexpectedEOF, "unexpected", 2},
		{"EOF lookalike", []byte{133}, custom, 2, 5, -3, custom, "other", 2},
		{"wrapped EOF", []byte{131}, wrapped, 3, 3, -2, wrapped, "other", 2},
		{"invalid terminating tenth payload", ten(129, 3), io.EOF, 0, 1, -1, overflow, "overflow", 10},
		{"tenth continuation with low bit set", ten(129, 129), io.EOF, 0, 1<<63 | 1, -(1 << 62) - 1, overflow, "overflow", 10},
		{"tenth continuation with low bit clear", ten(128, 254), io.EOF, 0, 0, 0, overflow, "overflow", 10},
		{"valid tenth payload", append(slices.Repeat([]byte{255}, 9), 1, 42), io.EOF, 0, math.MaxUint64, math.MinInt64, nil, "ok", 10},
		{"stop after first byte", []byte{2, 255}, io.EOF, 0, 2, 1, nil, "ok", 1},
	}
	var wire bytes.Buffer
	var rejects [5]int
	for _, c := range cases {
		first := &reader{bytes: c.bytes, end: c.end, limit: c.calls}
		second := &reader{bytes: c.bytes, end: c.end, limit: c.calls}
		u, err := binary.ReadUvarint(first)
		s, serr := binary.ReadVarint(second)
		if u != c.u || s != c.s || err != c.err || serr != c.err || first.calls != c.calls || second.calls != c.calls {
			log.Fatalf("%s: native (%d,%v,%d) (%d,%v,%d)", c.name, u, err, first.calls, s, serr, second.calls)
		}
		fmt.Fprint(&wire, c.tag)
		for _, b := range c.bytes {
			fmt.Fprintf(&wire, " %d", b)
		}
		fmt.Fprintln(&wire)
		// Count differences from incorrect status, value and read-budget predictions.
		if errors.Is(c.end, io.EOF) && c.end != io.EOF && c.err == c.end {
			rejects[0]++
		}
		if c.err != nil && c.s != 0 {
			rejects[1]++
		}
		if c.err == overflow && c.bytes[9] < 128 && u|(uint64(c.bytes[9])<<63) != u {
			rejects[2]++
		}
		if c.err == overflow && c.bytes[9] >= 128 && u&^(uint64(1)<<63) != u {
			rejects[3]++
		}
		if c.err == overflow && c.bytes[9] >= 128 && c.calls != 11 {
			rejects[4]++
		}
	}
	for i, got := range check.LeanRows("godiff-binary-varint-read", len(cases), &wire) {
		c := cases[i]
		want := fmt.Sprintf("%d:%s:%d|%d:%s:%d", c.u, c.label, c.calls, c.s, c.label, c.calls)
		if got != want {
			log.Fatalf("%s: Lean %s, want %s", c.name, got, want)
		}
	}
	if rejects != [5]int{1, 6, 1, 1, 2} {
		log.Fatalf("unexpected mutant rejection counts %v", rejects)
	}
	fmt.Printf("Varint reader: %d independent native and Lean fixtures, mutant rejections %v\n", len(cases), rejects)
}
