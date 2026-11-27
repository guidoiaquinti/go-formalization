package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	name, stored, spare string
	start, count        int
	tag, read           int
	end, offset         int
	status              string
}

func (c fixture) prepare() (*bytes.Buffer, []byte) {
	if c.tag == 0 {
		return nil, nil
	}
	if c.tag == 1 {
		return new(bytes.Buffer), nil
	}
	raw := c.stored + c.spare
	storage := make([]byte, len(raw))
	copy(storage, raw)
	b := bytes.NewBuffer(storage[:len(c.stored)])
	if c.read > 0 {
		if _, width, err := b.ReadRune(); err != nil || width != c.read {
			log.Fatalf("%s: setup rune width=%d error=%v", c.name, width, err)
		}
	} else if c.read == -1 {
		b.Next(c.start)
	}
	return b, storage
}

func truncate(b *bytes.Buffer, count int) (status string) {
	defer func() {
		if value := recover(); value != nil {
			if b != nil && value == "bytes.Buffer: truncation out of range" {
				status = "bounds"
				return
			}
			if err, ok := value.(runtime.Error); b == nil && ok &&
				err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
				status = "nil-receiver"
				return
			}
			panic(value)
		}
	}()
	b.Truncate(count)
	return "ok"
}

func main() {
	tests := []fixture{
		{name: "nil zero", count: 0, status: "nil-receiver"},
		{name: "nil negative", count: -1, status: "nil-receiver"},
		{name: "nil positive", count: 1, status: "nil-receiver"},
		{name: "zero-value reset", tag: 1, status: "ok"},
		{name: "zero-value oversize", tag: 1, count: 1, status: "bounds"},
		{"regular reset", "abcdef", "\x09\xff", 2, 0, 2, -1, 0, 0, "ok"},
		{"regular retain one", "abcdef", "\x09\xff", 2, 1, 2, -1, 3, 2, "ok"},
		{"regular retain all", "abcdef", "\x09\xff", 2, 4, 2, -1, 6, 2, "ok"},
		{"regular reject extra", "abcdef", "\x09\xff", 2, 5, 2, -1, 6, 2, "bounds"},
		{"regular reject negative", "abcdef", "\x09\xff", 2, math.MinInt, 2, -1, 6, 2, "bounds"},
		{"rune reset", "🚀xy", "\x80", 4, 0, 2, 4, 0, 0, "ok"},
		{"rune retain one", "🚀xy", "\x80", 4, 1, 2, 4, 5, 4, "ok"},
		{"rune reject extra", "🚀xy", "\x80", 4, math.MaxInt, 2, 4, 6, 4, "bounds"},
		{"rune reject negative", "🚀xy", "\x80", 4, -1, 2, 4, 6, 4, "bounds"},
	}
	var requests bytes.Buffer
	var wants []string
	var rejectedByte, rejectedRune, zeroReslicing int
	for _, c := range tests {
		b, storage := c.prepare()
		if got := truncate(b, c.count); got != c.status {
			log.Fatalf("%s: status=%s want=%s", c.name, got, c.status)
		}
		if b == nil {
			fmt.Fprintln(&requests, 0, c.count)
			wants = append(wants, c.status)
			continue
		}
		beforeByte, _ := c.prepare()
		beforeRune, _ := c.prepare()
		if byteValid, runeValid := beforeByte.UnreadByte() == nil, beforeRune.UnreadRune() == nil; byteValid != (c.read != 0) || runeValid != (c.read > 0) {
			log.Fatalf("%s: setup did not establish the expected read flag", c.name)
		}
		unread, spare := c.stored[c.offset:c.end], c.stored[c.end:]+c.spare
		view := b.Bytes()
		if (view == nil) != (c.tag == 1) || string(view) != unread || cap(view) != len(unread)+len(spare) ||
			b.Cap()-cap(view) != c.offset || string(view[len(view):cap(view)]) != spare || b.UnreadByte() == nil {
			log.Fatalf("%s: unexpected bytes, offset, capacity, nilness or unread-byte state", c.name)
		}
		rb, _ := c.prepare()
		truncate(rb, c.count)
		if rb.UnreadRune() == nil {
			log.Fatalf("%s: truncate retained the rune-unread flag", c.name)
		}
		if c.status == "bounds" && c.read != 0 {
			rejectedByte++
			if c.read > 0 {
				rejectedRune++
			}
		}
		if c.count == 0 && c.start > 0 {
			zeroReslicing++
			b.WriteByte(0xa5)
			if string(storage) != "\xa5"+(c.stored + c.spare)[1:] {
				log.Fatalf("%s: reset append did not reuse the original backing at zero", c.name)
			}
		}
		marker := "slice"
		if c.tag == 1 {
			marker = "nil"
		}
		fmt.Fprintf(&requests, "%d %d %d %d %d", c.tag, len(c.stored), c.start, c.read, c.count)
		for _, value := range []byte(c.stored + c.spare) {
			fmt.Fprintf(&requests, " %d", value)
		}
		fmt.Fprintln(&requests)
		wants = append(wants, fmt.Sprintf("%s:%s:%d:%d:%s:%s:true:true", c.status, marker,
			c.offset, len(unread)+len(spare), encode(unread), encode(spare)))
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-truncate", len(wants), &requests) {
		if got != wants[i] {
			log.Fatalf("Lean %s: Truncate=%s want=%s", tests[i].name, got, wants[i])
		}
	}
	if rejectedByte != 4 || rejectedRune != 2 || zeroReslicing != 2 {
		log.Fatalf("unexpected controls: retained byte flag=%d, retained rune flag=%d, zero reslice=%d",
			rejectedByte, rejectedRune, zeroReslicing)
	}
	fmt.Printf("Buffer.Truncate: %d native and Lean fixtures; retained byte flag=%d, retained rune flag=%d, zero reslice=%d mutant failures\n",
		len(wants), rejectedByte, rejectedRune, zeroReslicing)
}

func encode(raw string) string {
	parts := make([]string, len(raw))
	for i := range len(raw) {
		parts[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(parts, ",")
}
