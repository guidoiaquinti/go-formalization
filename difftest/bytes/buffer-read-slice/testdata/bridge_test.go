package bytes

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

func GoLeanBufferReadSlice(tag, offset, flag int, stored, spare []byte, count int) (result string) {
	var b *Buffer
	if tag != 0 {
		b = &Buffer{off: offset, lastRead: readOp(flag)}
		if tag == 2 {
			backing := make([]byte, len(stored)+len(spare))
			copy(backing, stored)
			copy(backing[len(stored):], spare)
			b.buf = backing[:len(stored)]
		}
	}
	wire := func(raw []byte) string {
		fields := make([]string, len(raw))
		for i, value := range raw {
			fields[i] = strconv.Itoa(int(value))
		}
		return strings.Join(fields, ",")
	}
	snapshot := func() string {
		return fmt.Sprintf("%d:%d:%s:%s", b.off, b.lastRead, wire(b.buf), wire(b.buf[len(b.buf):cap(b.buf)]))
	}
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if !ok {
				panic(failure)
			}
			if b == nil && err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
				result = "nil-receiver"
				return
			}

			panic(failure)
		}
	}()
	view, err := b.readSlice(byte(count))
	status := "ok"
	if err == io.EOF {
		status = "eof"
	} else if err != nil {
		panic(err)
	}
	if view == nil {
		return status + ":nil:" + snapshot()
	}
	result = fmt.Sprintf("%s:slice:%d:%s:%s:%s", status, cap(view), wire(view), wire(view[len(view):cap(view)]), snapshot())
	if cap(view) > 0 {
		complete := view[:cap(view)]
		original := complete[0]
		complete[0] ^= 0xff
		if b.buf[:cap(b.buf)][offset] != complete[0] {
			panic("readSlice result lost backing alias")
		}
		complete[0] = original
	}
	return result
}
