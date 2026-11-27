package bytes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func goLeanReaderSetupPanic(f func()) (status string) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			status = "panic"
		}
	}()
	f()
	return "returned"
}

func GoLeanReaderSetup(tag, offset int, cursor int64, previous int, raw, spare []byte) string {
	var input []byte
	if tag == 1 {
		backing := make([]byte, offset+len(raw)+len(spare)+7)
		copy(backing[offset:], raw)
		copy(backing[offset+len(raw):], spare)
		input = backing[offset : offset+len(raw) : offset+len(raw)+len(spare)]
	}
	wire := func(raw []byte) string {
		fields := make([]string, len(raw))
		for i, b := range raw {
			fields[i] = strconv.Itoa(int(b))
		}
		return strings.Join(fields, ",")
	}
	describe := func(r *Reader) string {
		marker := "slice"
		if r.s == nil {
			marker = "nil"
		}
		if cap(r.s) > 0 && &r.s[:cap(r.s)][0] != &input[:cap(input)][0] {
			panic("reader missed original backing")
		}
		return fmt.Sprintf("%s:%d:%d:%s:%s", marker, r.i, r.prevRune, wire(r.s), wire(r.s[len(r.s):cap(r.s)]))
	}
	constructor := describe(NewReader(input))
	original := Reader{s: []byte{111, 112, 113}, i: cursor, prevRune: previous}
	original.Reset(input)
	reset := describe(&original)
	sized := Reader{s: input, i: cursor, prevRune: previous}
	size := sized.Size()
	if sized.i != cursor || sized.prevRune != previous {
		panic("Size changed reader state")
	}
	var missing *Reader
	nilReset := goLeanReaderSetupPanic(func() { missing.Reset(input) })
	nilSize := goLeanReaderSetupPanic(func() { missing.Size() })
	return fmt.Sprintf("%s|%s|%s|%d|%s", constructor, reset, nilReset, size, nilSize)
}
