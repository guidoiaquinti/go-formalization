package bytes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func goLeanWriteByteStatus(b *Buffer, value byte) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if b != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			status = "panic"
		}
	}()
	if err := b.WriteByte(value); err != nil {
		panic("WriteByte returned nonnil error")
	}
	return "ok"
}

func GoLeanBufferWriteByte(tag, offset, flag int, stored, spare []byte, value byte) string {
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
	status := goLeanWriteByteStatus(b, value)
	if b == nil {
		return status
	}
	marker := "slice"
	if b.buf == nil {
		marker = "nil"
	}
	reused := len(storage) != 0 && cap(b.buf) != 0 && &storage[0] == &b.buf[:cap(b.buf)][0]
	return fmt.Sprintf("%s:%d:%d:%d:%s:%t:%s:%s", status, b.lastRead, b.off, cap(b.buf), marker, reused,
		writeByteWire(b.buf), writeByteWire(b.buf[len(b.buf):cap(b.buf)]))
}

func writeByteWire(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
