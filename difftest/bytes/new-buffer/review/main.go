package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	raw          []byte
	length, lead int
	want         string
}

func main() {
	cases := []fixture{
		{nil, 0, 0, "nil:0:0:0:0::"},
		{[]byte{}, 0, 1, "slice:0:0:0:0::"},
		{[]byte{255, 0, 9}, 0, 3, "slice:0:3:0:0::255,0,9"},
		{[]byte{255, 0, 9}, 2, 3, "slice:2:3:0:0:255,0:9"},
		{[]byte{255, 0, 9}, 3, 3, "slice:3:3:0:0:255,0,9:"},
		{[]byte{0x80, 0xc0, 0, 3, 4}, 2, 1, "slice:2:5:0:0:128,192:0,3,4"},
	}
	var rows bytes.Buffer
	for _, f := range cases {
		marker := 0
		if f.raw != nil {
			marker = 1
		}
		fmt.Fprintf(&rows, "%d %d", marker, f.length)
		for _, b := range f.raw {
			fmt.Fprintf(&rows, " %d", b)
		}
		fmt.Fprintln(&rows)
	}
	for i, got := range check.LeanRows("godiff-bytes-new-buffer", len(cases), &rows) {
		f := cases[i]
		var src []byte
		if f.raw != nil {
			mem := bytes.Repeat([]byte{78}, f.lead+len(f.raw)+5)
			copy(mem[f.lead:], f.raw)
			src = mem[f.lead : f.lead+f.length : f.lead+len(f.raw)]
		}
		var address *byte
		if cap(src) != 0 {
			address = &src[:cap(src)][0]
		}
		b := bytes.NewBuffer(src)
		unread := b.Bytes()
		if b.Cap() > 0 && &unread[:cap(unread)][0] != address {
			log.Fatal("constructor detached the supplied backing")
		}
		if b.UnreadByte() == nil || b.UnreadRune() == nil {
			log.Fatal("constructor made a prior read observable")
		}
		storage := "slice"
		if unread == nil {
			storage = "nil"
		}
		native := fmt.Sprintf("%s:%d:%d:%d:0:%s:%s", storage, b.Len(), b.Cap(), b.Cap()-cap(unread), wire(unread), wire(unread[len(unread):cap(unread)]))
		if got != f.want || native != f.want {
			log.Fatalf("fixture %d: Go=%s Lean=%s want=%s", i, native, got, f.want)
		}
	}
	backing := []byte{11, 22, 33, 44}
	buffer := bytes.NewBuffer(backing[1:1:3])
	first := &buffer.Bytes()[:buffer.Cap()][0]
	if err := buffer.WriteByte(0xff); err != nil || &buffer.Bytes()[0] != first ||
		!bytes.Equal(buffer.Bytes()[:buffer.Cap()], []byte{0xff, 33}) {
		log.Fatal("initial empty slice did not adopt capacity and spare contents")
	}
	if cases[1].want == cases[0].want || cases[3].want == "slice:2:2:0:0:255,0:" ||
		cases[2].want == "slice:0:3:0:0::0,0,0" {
		log.Fatal("one incorrect variant escaped")
	}
	fmt.Printf("NewBuffer: %d native and Lean fixtures passed, plus first-write reuse\n", len(cases))
}

func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}
