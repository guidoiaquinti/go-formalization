package bytes

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

func GoLeanReaderCopy(tag, previous int, cursor, offset int64, sharing, size int, raw, spare, separate []byte) string {
	wire := func(data []byte) string {
		fields := make([]string, len(data))
		for i, b := range data {
			fields[i] = strconv.Itoa(int(b))
		}
		return strings.Join(fields, ",")
	}
	run := func(positioned bool) (result string) {
		var r *Reader
		var source []byte
		if tag != 0 {
			r = &Reader{i: cursor, prevRune: previous}
			if tag == 2 {
				source = make([]byte, len(raw)+len(spare))
				copy(source, raw)
				copy(source[len(raw):], spare)
				r.s = source[:len(raw)]
			}
		}

		var dst []byte
		if sharing >= 0 {
			dst = source[sharing : sharing+size]
		} else if separate != nil {
			dst = make([]byte, len(separate)+1)
			dst[len(separate)] = 0xf1
			copy(dst, separate)
			dst = dst[:len(separate)]
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
		var n int
		var err error
		if positioned {
			n, err = r.ReadAt(dst, offset)
		} else {
			n, err = r.Read(dst)
		}
		status := "ok"
		if err == io.EOF {
			status = "EOF"
		} else if err != nil {
			if !positioned || offset >= 0 {
				panic("unexpected Reader error")
			}
			_, again := r.ReadAt(nil, offset)
			if again == nil || again == err || again.Error() != err.Error() {
				panic("ReadAt did not construct a fresh negative-offset error")
			}
			status = err.Error()
		}
		if sharing < 0 && separate != nil && dst[:cap(dst)][len(dst)] != 0xf1 {
			panic("changed spare disjoint destination byte")
		}
		state := "receiver-nil"
		if r != nil {
			if len(r.s) != len(raw) || cap(r.s) != len(source) ||
				(cap(r.s) > 0 && &r.s[:cap(r.s)][0] != &source[0]) {
				panic("reader slice geometry changed")
			}
			tag := "slice"
			if r.s == nil {
				tag = "nil"
			}
			state = fmt.Sprintf("%s:%d:%d:%s:%s", tag, r.i, r.prevRune, wire(r.s), wire(r.s[len(r.s):cap(r.s)]))
		}
		return fmt.Sprintf("%d:%s:%s:%s", n, status, wire(dst), state)
	}
	return run(false) + "|" + run(true)
}
