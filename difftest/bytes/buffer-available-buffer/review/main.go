package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	tag, offset            int
	stored, spare, payload []byte
	want                   string
}

func main() {
	cases := []fixture{
		{0, 0, nil, nil, nil, "panic"},
		{1, 0, nil, nil, nil, "nil:0:0::::0:"},
		{1, 0, nil, nil, []byte{0}, "nil:0:0::::1:0"},
		{2, 0, []byte{}, nil, nil, "slice:0:0::::0:"},
		{2, 1, []byte{32}, nil, []byte{42}, "slice:0:0::::1:42"},
		{2, 1, []byte{1, 2, 3}, []byte{9, 8, 7}, []byte{255}, "slice:0:3:9,8,7:255,8,7:2,3:1:2,3,255"},
		{2, 2, []byte{1, 2}, []byte{8, 7}, []byte{0, 1}, "slice:0:2:8,7:0,1::2:0,1"},
		{2, 1, []byte{11, 22}, []byte{33}, []byte{44, 55}, "slice:0:1:33:33:22:2:22,44,55"},
		{2, 0, []byte{3}, []byte{4, 5}, nil, "slice:0:2:4,5:4,5:3:0:3"},
	}
	var input bytes.Buffer
	copyTail, nilZero := 0, 0
	for i, c := range cases {
		fmt.Fprint(&input, c.tag)
		if c.tag != 0 {
			fmt.Fprint(&input, " ", len(c.stored), " ", len(c.spare), " ", c.offset)
			for _, seq := range [][]byte{c.stored, c.spare, c.payload} {
				for _, value := range seq {
					fmt.Fprint(&input, " ", value)
				}
			}
		}
		fmt.Fprintln(&input)
		if got := observe(c, false, false); got != c.want {
			log.Fatalf("Go fixture %d returned %q, want %q", i, got, c.want)
		}
		if observe(c, true, false) != c.want {
			copyTail++
		}
		if observe(c, false, true) != c.want {
			nilZero++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-buffer-available-buffer", len(cases), &input) {
		if row != cases[i].want {
			log.Fatalf("Lean fixture %d returned %q, want %q", i, row, cases[i].want)
		}
	}
	if copyTail != 2 || nilZero != 2 {
		log.Fatalf("mutation counts copy=%d nil-zero-capacity=%d, want 2 each", copyTail, nilZero)
	}
	fmt.Printf("Buffer.AvailableBuffer passed %d fixed Go and Lean fixtures; copied tail and nil zero capacity mutants failed %d and %d.\n",
		len(cases), copyTail, nilZero)
}

func observe(c fixture, detached, nilZero bool) (output string) {
	var b *bytes.Buffer
	switch c.tag {
	case 1:
		b = new(bytes.Buffer)
	case 2:
		initial := make([]byte, len(c.stored)+len(c.spare))
		copy(initial, c.stored)
		copy(initial[len(c.stored):], c.spare)
		b = bytes.NewBuffer(initial[:len(c.stored)])
		b.Next(c.offset)
	}
	defer func() {
		if value := recover(); value != nil {
			if _, ok := value.(runtime.Error); !ok || c.tag != 0 {
				panic(value)
			}
			output = "panic"
		}
	}()
	view := b.AvailableBuffer()
	if detached && cap(view) > 0 {
		view = bytes.Clone(view[:cap(view)])[0:0:cap(view)]
	}
	if nilZero && cap(view) == 0 {
		view = nil
	}
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	before := wire(view[:cap(view)])
	data := append(view, c.payload...)
	tail := b.AvailableBuffer()
	after := wire(tail[:cap(tail)])
	unread := wire(b.Bytes())
	written, err := b.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s:%d:%d:%s:%s:%s:%d:%s", marker, len(view), cap(view), before,
		after, unread, written, wire(b.Bytes()))
}

func wire(seq []byte) string {
	values := make([]string, len(seq))
	for i, value := range seq {
		values[i] = fmt.Sprint(value)
	}
	return strings.Join(values, ",")
}
