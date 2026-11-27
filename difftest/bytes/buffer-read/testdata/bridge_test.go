package bytes

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

func GoLeanBufferRead(tag, offset, flag int, stored, spare []byte, count int) (result string) {
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
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if !ok {
				panic(failure)
			}
			if b == nil && err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
				result = "panic"
				return
			}
			panic(failure)
		}
	}()
	var dst []byte
	if count > 0 {
		dst = make([]byte, count)
		for i := range dst {
			dst[i] = byte(167 + i*31)
		}
	}
	n, err := b.Read(dst)
	status, marker := "ok", "slice"
	if err == io.EOF {
		status = "eof"
	} else if err != nil {
		panic(err)
	}
	if b.buf == nil {
		marker = "nil"
	}
	return fmt.Sprintf("%s:%d:%s:%d:%d:%s:%s:%s", status, n, marker, b.off, b.lastRead, wire(dst), wire(b.buf), wire(b.buf[len(b.buf):cap(b.buf)]))
}
