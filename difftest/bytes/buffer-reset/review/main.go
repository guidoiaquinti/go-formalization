package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type history struct {
	name, stored, spare string
	nilStorage, rune    bool
	offset              int
}

func (h history) prepare() (*bytes.Buffer, []byte, int) {
	if h.nilStorage {
		return new(bytes.Buffer), nil, 0
	}
	raw := h.stored + h.spare
	storage := make([]byte, len(raw))
	copy(storage, raw)
	b := bytes.NewBuffer(storage[:len(h.stored)])
	if h.rune {
		_, width, err := b.ReadRune()
		if err != nil || width != h.offset {
			log.Fatalf("%s: setup ReadRune width=%d error=%v", h.name, width, err)
		}
		return b, storage, width
	}
	if h.offset != 0 {
		b.Next(h.offset)
		return b, storage, -1
	}
	return b, storage, 0
}

func main() {
	tests := []history{
		{name: "zero value", nilStorage: true},
		{name: "nonnil empty"},
		{name: "spare-only", spare: "\x18\xef"},
		{name: "unread", stored: "abc", spare: "\x09"},
		{name: "Next partial", stored: "abcdef", spare: "\x19\xfe", offset: 3},
		{name: "Next full", stored: "abc", offset: 3},
		{name: "ASCII rune", stored: "x!", spare: "\xff", rune: true, offset: 1},
		{name: "two-byte rune", stored: "éx", rune: true, offset: 2},
		{name: "three-byte RuneError", stored: "�x", spare: "\x80", rune: true, offset: 3},
		{name: "four-byte rune", stored: "🚀x", rune: true, offset: 4},
		{name: "malformed byte", stored: "\xffx", spare: "\x00", rune: true, offset: 1},
	}
	var lean bytes.Buffer
	var wants []string
	var preservedOffset, retainedByteTag, retainedRuneTag, discardedStorage int
	for _, h := range tests {
		beforeByte, _, tag := h.prepare()
		beforeRune, _, _ := h.prepare()
		byteWasValid := beforeByte.UnreadByte() == nil
		runeWasValid := beforeRune.UnreadRune() == nil
		if byteWasValid != (tag != 0) || runeWasValid != (tag > 0) {
			log.Fatalf("%s: unreachable tag %d (UnreadByte=%t UnreadRune=%t)", h.name, tag, byteWasValid, runeWasValid)
		}
		if byteWasValid {
			retainedByteTag++
		}
		if runeWasValid {
			retainedRuneTag++
		}
		if h.offset != 0 {
			preservedOffset++
		}
		if h.stored+h.spare != "" {
			discardedStorage++
		}
		b, storage, _ := h.prepare()
		b.Reset()
		view := b.Bytes()
		full := h.stored + h.spare
		if len(view) != 0 || cap(view) != len(full) || string(view[:cap(view)]) != full ||
			(view == nil) != h.nilStorage || b.UnreadByte() == nil {
			log.Fatalf("%s: Reset kept incorrect bytes, nilness, length, capacity, or byte flag", h.name)
		}
		br, _, _ := h.prepare()
		br.Reset()
		if br.UnreadRune() == nil {
			log.Fatalf("%s: Reset retained the rune flag", h.name)
		}
		reused := ""
		if len(full) > 0 {
			b.WriteByte(0xa5)
			reused = "\xa5" + full[1:]
			if string(storage) != reused || b.Cap() != len(full) || !bytes.Equal(b.Bytes(), []byte{0xa5}) {
				log.Fatalf("%s: fitting append did not reuse original backing at zero", h.name)
			}
		}
		marker, inputTag := "slice", 2
		if h.nilStorage {
			marker, inputTag = "nil", 1
		}
		fmt.Fprintf(&lean, "%d %d %d %d", inputTag, len(h.stored), h.offset, tag)
		for i := range len(full) {
			fmt.Fprintf(&lean, " %d", full[i])
		}
		fmt.Fprintln(&lean)
		wants = append(wants, fmt.Sprintf("%s:0:%d:%s:true:true:%s", marker, len(full), decimal(full), decimal(reused)))
	}
	if !nilPanics() {
		log.Fatal("Reset did not report its nil-pointer panic")
	}
	fmt.Fprintln(&lean, 0)
	wants = append(wants, "panic")
	for i, got := range check.LeanRows("godiff-bytes-buffer-reset", len(wants), &lean) {
		if want := wants[i]; got != want {
			log.Fatalf("Lean reset fixture %d: got=%s want=%s", i, got, want)
		}
	}
	if preservedOffset != 7 || retainedByteTag != 7 || retainedRuneTag != 5 || discardedStorage != 9 {
		log.Fatalf("unexpected mutant controls: offset=%d, byte flag=%d, rune flag=%d, discarded=%d",
			preservedOffset, retainedByteTag, retainedRuneTag, discardedStorage)
	}
	fmt.Printf("bytes.Buffer.Reset: %d native and Lean fixtures; offset=%d, byte flag=%d, rune flag=%d, discarded=%d mutant failures\n",
		len(wants), preservedOffset, retainedByteTag, retainedRuneTag, discardedStorage)
}

func decimal(raw string) string {
	parts := make([]string, len(raw))
	for i := range parts {
		parts[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(parts, ",")
}

func nilPanics() (panicked bool) {
	defer func() {
		if err, ok := recover().(runtime.Error); ok {
			panicked = err.Error() == "runtime error: invalid memory address or nil pointer dereference"
		}
	}()
	var b *bytes.Buffer
	b.Reset()
	return false
}
