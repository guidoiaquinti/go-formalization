package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	name, stored, spare string
	tag, offset, read   int
	want                rune
	width               int
}

func (c fixture) setup() *bytes.Buffer {
	if c.tag == 1 {
		return new(bytes.Buffer)
	}
	all := c.stored + c.spare
	storage := make([]byte, len(all))
	copy(storage, all)
	b := bytes.NewBuffer(storage[:len(c.stored)])
	switch {
	case c.read < 0:
		b.Next(c.offset)
	case c.read > 0:
		if _, width, err := b.ReadRune(); err != nil || width != c.read {
			log.Fatalf("%s: ReadRune setup width=%d error=%v", c.name, width, err)
		}
	}
	return b
}

func main() {
	tests := []fixture{
		{name: "zero value", tag: 1},
		{name: "empty backing", tag: 2},
		{name: "consumed ordinary", stored: "abc", spare: "\x00\xff", tag: 2, offset: 3, read: -1},
		{name: "consumed rune", stored: "🚀", spare: "\x80", tag: 2, offset: 4, read: 4},
		{name: "NUL", stored: "\x00x", tag: 2, want: 0, width: 1},
		{name: "two byte", stored: "aéz", spare: "\xff", tag: 2, offset: 1, read: -1, want: 'é', width: 2},
		{name: "encoded RuneError", stored: "�z", tag: 2, want: utf8.RuneError, width: 3},
		{name: "four byte", stored: "é🚀z", tag: 2, offset: 2, read: 2, want: '🚀', width: 4},
		{name: "maximum scalar", stored: "\xf4\x8f\xbf\xbfz", tag: 2, want: utf8.MaxRune, width: 4},
		{name: "overlong", stored: "\xc0\xaf", tag: 2, want: utf8.RuneError, width: 1},
		{name: "truncated", stored: "\xf0\x9f\x92", tag: 2, want: utf8.RuneError, width: 1},
		{name: "surrogate", stored: "\xed\xa0\x80", tag: 2, want: utf8.RuneError, width: 1},
		{name: "above Unicode maximum", stored: "\xf4\x90\x80\x80", tag: 2, want: utf8.RuneError, width: 1},
		{name: "partial prior byte", stored: "éz", tag: 2, offset: 1, read: -1, want: utf8.RuneError, width: 1},
	}
	var input bytes.Buffer
	var expected []string
	var noReset, consumeThree, sameRewind int
	for _, c := range tests {
		b := c.setup()
		got, width, err := b.ReadRune()
		status, marker := "ok", "slice"
		offset := c.offset + c.width
		remaining, spare := c.stored[offset:], c.spare
		var wantErr error
		if c.tag == 1 {
			marker = "nil"
		}
		if c.width == 0 {
			status, offset, remaining, spare = "eof", 0, "", c.stored+c.spare
			wantErr = io.EOF
			if c.offset != 0 {
				noReset++
			}
		}
		if got != c.want || width != c.width || err != wantErr {
			log.Fatalf("%s: ReadRune=%U,%d,%v want=%U,%d,%v", c.name, got, width, err, c.want, c.width, wantErr)
		}
		view, capacity := b.Bytes(), len(remaining)+len(spare)
		if (view == nil) != (c.tag == 1) || string(view) != remaining || cap(view) != capacity ||
			b.Cap()-cap(view) != offset || string(view[len(view):cap(view)]) != spare {
			log.Fatalf("%s: changed buffer geometry differs", c.name)
		}
		runeBuffer := c.setup()
		runeBuffer.ReadRune()
		byteBack, runeBack := "rejected", "rejected"
		if c.width == 0 {
			if b.UnreadByte() == nil || runeBuffer.UnreadRune() == nil {
				log.Fatalf("%s: EOF retained an unread flag", c.name)
			}
		} else {
			if err := b.UnreadByte(); err != nil {
				log.Fatalf("%s: UnreadByte failed: %v", c.name, err)
			}
			back, err := b.ReadByte()
			if err != nil || back != c.stored[offset-1] {
				log.Fatalf("%s: UnreadByte retry did not read last original byte", c.name)
			}
			byteBack = strconv.Itoa(int(back))
			if err := runeBuffer.UnreadRune(); err != nil {
				log.Fatalf("%s: UnreadRune failed: %v", c.name, err)
			}
			start := runeBuffer.Cap() - cap(runeBuffer.Bytes())
			retry, retryWidth, err := runeBuffer.ReadRune()
			if err != nil || start != c.offset || retry != c.want || retryWidth != c.width {
				log.Fatalf("%s: UnreadRune retry did not rewind full width", c.name)
			}
			runeBack = fmt.Sprintf("%d,%d,%d", start, retry, retryWidth)
			if c.want == utf8.RuneError && c.width != 3 {
				consumeThree++
			}
			if c.width > 1 {
				sameRewind++
			}
		}
		fmt.Fprintf(&input, "%d %d %d %d", c.tag, len(c.stored), c.offset, c.read)
		for _, octet := range []byte(c.stored + c.spare) {
			fmt.Fprintf(&input, " %d", octet)
		}
		fmt.Fprintln(&input)
		expected = append(expected, fmt.Sprintf("%s:%d:%d:%s:%d:%d:%s:%s:%s:%s", status, c.want, c.width, marker,
			offset, capacity, decimal(remaining), decimal(spare), byteBack, runeBack))
	}
	if !nilPanics() {
		log.Fatal("ReadRune did not report the nil-pointer runtime panic")
	}
	fmt.Fprintln(&input, 0)
	expected = append(expected, "nil-receiver")
	for i, got := range check.LeanRows("godiff-bytes-buffer-read-rune", len(expected), &input) {
		if got != expected[i] {
			log.Fatalf("Lean fixture %d: got=%s want=%s", i, got, expected[i])
		}
	}
	if noReset != 2 || consumeThree != 5 || sameRewind != 4 {
		log.Fatalf("controls: no reset=%d, RuneError width three=%d, same rewind=%d", noReset, consumeThree, sameRewind)
	}
	fmt.Printf("Buffer.ReadRune: %d native and Lean fixtures; no reset=%d, RuneError width three=%d, same rewind=%d mutant failures\n",
		len(expected), noReset, consumeThree, sameRewind)
}

func decimal(raw string) string {
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
	b.ReadRune()
	return false
}
