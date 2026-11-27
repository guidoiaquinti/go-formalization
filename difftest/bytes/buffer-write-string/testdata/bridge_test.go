package bytes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func goLeanWriteStringStatus(b *Buffer, value []byte) (count int, status string) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if b != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			status = "panic"
		}
	}()
	count, err := b.WriteString(string(value))
	if err != nil {
		panic("WriteString returned nonnil error")
	}
	return count, "ok"
}

func GoLeanBufferWriteString(tag, offset, flag int, stored, spare, value []byte) string {
	var b *Buffer
	var storage []byte
	if tag != 0 {
		b = &Buffer{off: offset, lastRead: readOp(flag)}
		if tag == 2 {
			storage = make([]byte, len(stored)+len(spare))
			copy(storage, stored)
			copy(storage[len(stored):], spare)
			b.buf = storage[:len(stored)]
		}
	}
	count, status := goLeanWriteStringStatus(b, value)
	if b == nil {
		return status
	}
	marker := "slice"
	if b.buf == nil {
		marker = "nil"
	}
	reused := len(storage) != 0 && cap(b.buf) != 0 && &storage[0] == &b.buf[:cap(b.buf)][0]
	return fmt.Sprintf("%s:%d:%d:%d:%d:%s:%t:%s:%s", status, count, b.lastRead, b.off, cap(b.buf), marker, reused,
		writeWire(b.buf), writeWire(b.buf[len(b.buf):cap(b.buf)]))
}

func writeWire(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
