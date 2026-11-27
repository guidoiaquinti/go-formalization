package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	tag, offset   int
	stored, spare []byte
	want          string
}

func main() {
	cases := []fixture{
		{0, 0, nil, nil, "panic"},
		{1, 0, nil, nil, "nil:0::::"},
		{2, 0, nil, nil, "slice:0::::"},
		{2, 1, []byte{10, 20, 30}, []byte{40, 50}, "slice:4:20,30:40,50:21,31:10,21,31,40,50"},
		{2, 2, []byte{1, 2}, []byte{3}, "slice:1::3::1,2,3"},
		{2, 0, []byte{0, 255}, nil, "slice:2:0,255::1,254:1,254"},
		{2, 0, nil, []byte{88, 99}, "slice:2::88,99::88,99"},
	}
	var rows bytes.Buffer
	detached, collapse := 0, 0
	for i, c := range cases {
		if got := native(c, false, false); got != c.want {
			log.Fatalf("Go case %d: %q, want %q", i, got, c.want)
		}
		if got := native(c, true, false); got != c.want {
			detached++
		}
		if got := native(c, false, true); got != c.want {
			collapse++
		}
		fmt.Fprint(&rows, c.tag)
		if c.tag != 0 {
			fmt.Fprint(&rows, " ", len(c.stored), " ", c.offset)
			for _, part := range [][]byte{c.stored, c.spare} {
				for _, b := range part {
					fmt.Fprint(&rows, " ", b)
				}
			}
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-bytes", len(cases), &rows) {
		if got != cases[i].want {
			log.Fatalf("Lean case %d: %q, want %q", i, got, cases[i].want)
		}
	}
	if detached != 2 || collapse != 3 {
		log.Fatalf("mutation counts detached=%d collapse=%d, want 2 and 3", detached, collapse)
	}
	fmt.Printf("bytes.Buffer.Bytes passed %d fixed Go and Lean fixtures; detached and empty-nil mutants failed %d and %d.\n",
		len(cases), detached, collapse)
}

func native(c fixture, detach, collapse bool) (result string) {
	var b *bytes.Buffer
	var backing []byte
	switch c.tag {
	case 1:
		b = new(bytes.Buffer)
	case 2:
		backing = make([]byte, len(c.stored)+len(c.spare))
		copy(backing, c.stored)
		copy(backing[len(c.stored):], c.spare)
		b = bytes.NewBuffer(backing[:len(c.stored)])
		b.Next(c.offset)
	}
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if c.tag != 0 || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			result = "panic"
		}
	}()
	view := b.Bytes()
	if collapse && len(view) == 0 {
		view = nil
	}
	if detach && len(view) > 0 {
		view = bytes.Clone(view[:cap(view)])[:len(view):cap(view)]
	}
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	visible := wire(view)
	spare := wire(view[len(view):cap(view)])
	capacity := cap(view)
	for j := range view {
		view[j] ^= 1
	}
	read := make([]byte, len(view))
	if _, err := io.ReadFull(b, read); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s:%d:%s:%s:%s:%s", marker, capacity, visible, spare, wire(read), wire(backing))
}

func wire(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
