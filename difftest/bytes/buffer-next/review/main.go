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
	tag, start, read    int
	count, size         int
	status              string
}

func (f fixture) setup() (*bytes.Buffer, []byte) {
	if f.tag == 0 {
		return nil, nil
	}
	if f.tag == 1 {
		return new(bytes.Buffer), nil
	}
	full := f.stored + f.spare
	storage := make([]byte, len(full))
	copy(storage, full)
	b := bytes.NewBuffer(storage[:len(f.stored)])
	switch {
	case f.read < 0:
		b.Next(f.start)
	case f.read > 0:
		if _, width, err := b.ReadRune(); err != nil || width != f.read {
			log.Fatalf("%s: ReadRune width=%d error=%v", f.name, width, err)
		}
	}
	return b, storage
}

func safeNext(b *bytes.Buffer, n int) (view []byte, status string) {
	defer func() {
		if value := recover(); value != nil {
			fault, ok := value.(runtime.Error)
			if !ok {
				panic(value)
			}
			switch {
			case b == nil && fault.Error() == "runtime error: invalid memory address or nil pointer dereference":
				status = "nil-receiver"
			case b != nil && strings.HasPrefix(fault.Error(), "runtime error: slice bounds out of range"):
				status = "bounds"
			default:
				panic(value)
			}
		}
	}()
	return b.Next(n), "ok"
}

func main() {
	tests := []fixture{
		{name: "nil negative", count: -1, status: "nil-receiver"},
		{name: "nil zero", status: "nil-receiver"},
		{name: "nil positive", count: math.MaxInt, status: "nil-receiver"},
		{name: "zero value negative", tag: 1, count: -1, status: "bounds"},
		{name: "zero value zero", tag: 1, status: "ok"},
		{name: "zero value oversized", tag: 1, count: math.MaxInt, status: "ok"},
		{name: "spare-only oversized", tag: 2, spare: "\x00\xff", count: math.MaxInt, status: "ok"},
		{"regular negative", "abcdef", "\x00\xff", 2, 2, -1, -1, 0, "bounds"},
		{"regular minimum", "abcdef", "\x00\xff", 2, 2, -1, math.MinInt, 0, "bounds"},
		{"regular zero", "abcdef", "\x00\xff", 2, 2, -1, 0, 0, "ok"},
		{"below unread", "abcdef", "\x00\xff", 2, 2, -1, 3, 3, "ok"},
		{"exact unread", "abcdef", "\x00\xff", 2, 2, -1, 4, 4, "ok"},
		{"unread plus one", "abcdef", "\x00\xff", 2, 2, -1, 5, 4, "ok"},
		{"maximum", "abcdef", "\x00\xff", 2, 2, -1, math.MaxInt, 4, "ok"},
		{"exhausted with spare", "abc", "\x00\xff", 2, 3, -1, math.MaxInt, 0, "ok"},
		{"prior rune negative", "🚀ab", "\x00", 2, 4, 4, -1, 0, "bounds"},
		{"prior rune zero", "🚀ab", "\x00", 2, 4, 4, 0, 0, "ok"},
		{"prior rune positive", "🚀ab", "\x00", 2, 4, 4, 1, 1, "ok"},
	}
	var input bytes.Buffer
	var expected []string
	var copyFailures, noClampFailures, retainedByteFailures, retainedRuneFailures int
	for _, f := range tests {
		b, storage := f.setup()
		view, status := safeNext(b, f.count)
		if status != f.status {
			log.Fatalf("%s: Next status=%s want=%s", f.name, status, f.status)
		}
		if b == nil {
			fmt.Fprintln(&input, 0, f.count)
			expected = append(expected, status)
			continue
		}
		beforeByte, _ := f.setup()
		beforeRune, _ := f.setup()
		if byteValid, runeValid := beforeByte.UnreadByte() == nil, beforeRune.UnreadRune() == nil; byteValid != (f.read != 0) || runeValid != (f.read > 0) {
			log.Fatalf("%s: wrong setup read state", f.name)
		}
		raw := f.stored + f.spare
		offset := f.start + f.size
		visible, remaining := f.stored[f.start:offset], f.stored[offset:]
		capacity, marker := 0, "nil"
		if status == "ok" {
			capacity = len(raw) - f.start
			if f.tag == 2 {
				marker = "slice"
			}
		}
		if string(view) != visible || cap(view) != capacity || (view == nil) != (marker == "nil") {
			log.Fatalf("%s: wrong view bytes, capacity or nilness", f.name)
		}
		backing := []byte(raw)
		for i := range view {
			view[i] ^= 1
			backing[f.start+i] ^= 1
		}
		if string(storage) != string(backing) || string(b.Bytes()) != remaining ||
			b.Cap()-cap(b.Bytes()) != offset {
			log.Fatalf("%s: Next view mutation did not match original backing or cursor", f.name)
		}
		unreturned := "rejected"
		if f.size > 0 {
			if err := b.UnreadByte(); err != nil {
				log.Fatalf("%s: positive Next must allow UnreadByte", f.name)
			}
			byte, err := b.ReadByte()
			if err != nil || byte != backing[offset-1] {
				log.Fatalf("%s: reread byte=%d error=%v", f.name, byte, err)
			}
			unreturned = strconv.Itoa(int(byte))
			copyFailures++
		} else if b.UnreadByte() == nil {
			log.Fatalf("%s: empty or failed Next retained UnreadByte", f.name)
		} else if f.read != 0 {
			retainedByteFailures++
		}
		runeBuffer, _ := f.setup()
		safeNext(runeBuffer, f.count)
		if runeBuffer.UnreadRune() == nil {
			log.Fatalf("%s: Next retained UnreadRune", f.name)
		}
		if f.read > 0 {
			retainedRuneFailures++
		}
		if status == "ok" && f.count > f.size {
			noClampFailures++
		}
		fmt.Fprintf(&input, "%d %d %d %d %d", f.tag, len(f.stored), f.start, f.read, f.count)
		for _, byte := range []byte(raw) {
			fmt.Fprintf(&input, " %d", byte)
		}
		fmt.Fprintln(&input)
		expected = append(expected, fmt.Sprintf("%s:%s:%d:%s:%d:%s:%s:%s:true", status, marker,
			capacity, encode(visible), offset, encode(remaining), encode(string(backing)), unreturned))
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-next", len(expected), &input) {
		if got != expected[i] {
			log.Fatalf("Lean fixture %s: got=%s want=%s", tests[i].name, got, expected[i])
		}
	}
	if copyFailures != 5 || noClampFailures != 5 || retainedByteFailures != 6 || retainedRuneFailures != 3 {
		log.Fatalf("mutants: copy=%d no clamp=%d retained byte=%d retained rune=%d",
			copyFailures, noClampFailures, retainedByteFailures, retainedRuneFailures)
	}
	fmt.Printf("Buffer.Next: %d native and Lean fixtures; copy=%d, no clamp=%d, retained byte=%d, retained rune=%d mutant failures\n",
		len(expected), copyFailures, noClampFailures, retainedByteFailures, retainedRuneFailures)
}

func encode(raw string) string {
	fields := make([]string, len(raw))
	for i := range len(raw) {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(fields, ",")
}
