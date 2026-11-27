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
	name, stored, spare string
	tag, offset, read   int
	change              bool
	want                byte
	eof                 bool
}

func (c fixture) setup() *bytes.Buffer {
	if c.tag == 1 {
		return new(bytes.Buffer)
	}
	raw := c.stored + c.spare
	storage := make([]byte, len(raw))
	copy(storage, raw)
	b := bytes.NewBuffer(storage[:len(c.stored)])
	switch {
	case c.read < 0:
		b.Next(c.offset)
	case c.read > 0:
		if _, width, err := b.ReadRune(); err != nil || width != c.read {
			log.Fatalf("%s: ReadRune width=%d, error=%v", c.name, width, err)
		}
	}
	if c.change {
		// A fresh aliased unread view can be changed before ReadByte.
		b.Bytes()[0] = c.want
	}
	return b
}

func main() {
	tests := []fixture{
		{name: "zero value", tag: 1, eof: true},
		{name: "nonnil empty", tag: 2, eof: true},
		{name: "spare-only", tag: 2, spare: "\x00\xff", eof: true},
		{name: "exhausted regular", tag: 2, stored: "abc", spare: "\x09", offset: 3, read: -1, eof: true},
		{name: "exhausted rune", tag: 2, stored: "🚀", spare: "\xff", offset: 4, read: 4, eof: true},
		{name: "NUL", tag: 2, stored: "\x00z", want: 0},
		{name: "high byte", tag: 2, stored: "\xffz", spare: "\x80", want: 0xff},
		{name: "following regular read", tag: 2, stored: "ab\x80z", spare: "\x00", offset: 2, read: -1, want: 0x80},
		{name: "following rune", tag: 2, stored: "éx", spare: "\x99", offset: 2, read: 2, want: 'x'},
		{name: "changed view", tag: 2, stored: "abz", spare: "\xff", offset: 1, read: -1, change: true, want: 0xa5},
	}
	var requests bytes.Buffer
	var wants []string
	var omitReset, retainEOF, initialByte, runeFlag int
	for _, c := range tests {
		b := c.setup()
		got, err := b.ReadByte()
		var wantError error
		if c.eof {
			wantError = io.EOF
		}
		if got != c.want || err != wantError {
			log.Fatalf("%s: ReadByte=(%d,%v), want (%d,%v)", c.name, got, err, c.want, wantError)
		}
		stored := c.stored
		if c.change {
			stored = stored[:c.offset] + string([]byte{c.want}) + stored[c.offset+1:]
		}
		marker := "slice"
		if c.tag == 1 {
			marker = "nil"
		}
		status, offset, unread, spare, retry := "ok", c.offset+1, "", c.spare, "rejected"
		if c.eof {
			status, offset, spare = "eof", 0, stored+c.spare
			if c.offset > 0 {
				omitReset++
			}
			old := c.setup()
			if old.UnreadByte() == nil {
				retainEOF++
			}
		} else {
			unread = stored[offset:]
			retry = strconv.Itoa(int(c.want))
			if stored[0] != c.want {
				initialByte++
			}
			runeFlag++
		}
		view := b.Bytes()
		capacity := len(unread) + len(spare)
		if (view == nil) != (c.tag == 1) || string(view) != unread ||
			cap(view) != capacity || b.Cap()-cap(view) != offset || string(view[len(view):cap(view)]) != spare {
			log.Fatalf("%s: retained geometry differs after ReadByte", c.name)
		}
		if err := b.UnreadByte(); (err == nil) == c.eof {
			log.Fatalf("%s: unexpected UnreadByte result %v", c.name, err)
		}
		if !c.eof {
			if again, err := b.ReadByte(); err != nil || again != c.want {
				log.Fatalf("%s: reread=%d, error=%v", c.name, again, err)
			}
		}
		runeBuffer := c.setup()
		runeBuffer.ReadByte()
		if runeBuffer.UnreadRune() == nil {
			log.Fatalf("%s: ReadByte permitted UnreadRune", c.name)
		}
		fmt.Fprintf(&requests, "%d %d %d %d", c.tag, len(stored), c.offset, c.read)
		for _, value := range []byte(stored + c.spare) {
			fmt.Fprintf(&requests, " %d", value)
		}
		fmt.Fprintln(&requests)
		wants = append(wants, fmt.Sprintf("%s:%d:%s:%d:%d:%s:%s:%s:true",
			status, c.want, marker, offset, capacity, encode(unread), encode(spare), retry))
	}
	if !nilPanics() {
		log.Fatal("ReadByte did not report the nil-pointer runtime panic")
	}
	fmt.Fprintln(&requests, 0)
	wants = append(wants, "nil-receiver")
	for i, got := range check.LeanRows("godiff-bytes-buffer-read-byte", len(wants), &requests) {
		if got != wants[i] {
			log.Fatalf("Lean ReadByte fixture %d: %s want %s", i, got, wants[i])
		}
	}
	if omitReset != 2 || retainEOF != 2 || initialByte != 3 || runeFlag != 5 {
		log.Fatalf("unexpected controls: omit reset=%d, EOF flag=%d, first byte=%d, rune flag=%d",
			omitReset, retainEOF, initialByte, runeFlag)
	}
	fmt.Printf("Buffer.ReadByte: %d native and Lean fixtures; omit reset=%d, EOF flag=%d, first byte=%d, rune flag=%d mutant failures\n",
		len(wants), omitReset, retainEOF, initialByte, runeFlag)
}

func encode(raw string) string {
	fields := make([]string, len(raw))
	for i := range len(raw) {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(fields, ",")
}

func nilPanics() (panicked bool) {
	defer func() {
		if err, ok := recover().(runtime.Error); ok {
			panicked = err.Error() == "runtime error: invalid memory address or nil pointer dereference"
		}
	}()
	var b *bytes.Buffer
	b.ReadByte()
	return false
}
