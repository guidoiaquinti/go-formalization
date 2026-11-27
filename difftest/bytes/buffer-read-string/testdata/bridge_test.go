package bytes

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

func GoLeanBufferReadString(tag, offset, flag int, stored, spare []byte, count int) (result string) {
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
	view, err := b.ReadString(byte(count))
	status := "ok"
	if err == io.EOF {
		status = "eof"
	} else if err != nil {
		panic(err)
	}
	raw := wire([]byte(view))
	result = fmt.Sprintf("%s:%s:%s", status, raw, snapshot())
	if len(view) > 0 {
		b.buf[offset] ^= 0xff
		if wire([]byte(view)) != raw {
			panic("ReadString result aliases buffer backing")
		}
	}
	return result
}
