package bytes

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

func GoLeanReaderRune(tag, offset int, cursor int64, previous int, raw, spare []byte) string {
	var input []byte
	if tag == 2 {
		backing := make([]byte, offset+len(raw)+len(spare)+1)
		backing[0] = 0xa5
		copy(backing[offset:], raw)
		copy(backing[offset+len(raw):], spare)
		input = backing[offset : offset+len(raw) : offset+len(raw)+len(spare)]
	}
	wire := func(raw []byte) string {
		fields := make([]string, len(raw))
		for i, value := range raw {
			fields[i] = strconv.Itoa(int(value))
		}
		return strings.Join(fields, ",")
	}
	describe := func(r *Reader) string {
		if cap(r.s) != cap(input) || len(r.s) != len(input) {
			panic("unexpected reader slice geometry")
		}
		if cap(input) > 0 && &r.s[:cap(r.s)][0] != &input[:cap(input)][0] {
			panic("reader detached from original backing")
		}
		marker := "slice"
		if r.s == nil {
			marker = "nil"
		}
		return fmt.Sprintf("%s:%d:%d:%s:%s", marker, r.i, r.prevRune, wire(r.s), wire(r.s[len(r.s):cap(r.s)]))
	}
	call := func(method func(*Reader) string) (result string) {
		var r *Reader
		if tag != 0 {
			r = &Reader{s: input, i: cursor, prevRune: previous}
		}
		defer func() {
			if value := recover(); value != nil {
				err, ok := value.(runtime.Error)
				if r != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
					panic(value)
				}
				result = "panic"
			}
		}()
		return method(r)
	}
	read := call(func(r *Reader) string {
		value, width, err := r.ReadRune()
		status := "ok"
		if err != nil {
			if err != io.EOF {
				panic("ReadRune returned unexpected error")
			}
			status = "EOF"
		}
		return fmt.Sprintf("%d:%d:%s:%s", value, width, status, describe(r))
	})
	unread := call(func(r *Reader) string {
		err := r.UnreadRune()
		snapshot := describe(r)
		status := "ok"
		if err != nil {
			again := r.UnreadRune()
			if again == nil || again == err || again.Error() != err.Error() {
				panic("UnreadRune did not construct a fresh error")
			}
			status = err.Error()
		}
		return status + ":" + snapshot
	})
	return read + "|" + unread
}
