package bytes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func goLeanWriteRuneStatus(b *Buffer, value rune) (width int, status string) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if b != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			status = "panic"
		}
	}()
	width, err := b.WriteRune(value)
	if err != nil {
		panic("WriteRune returned nonnil error")
	}
	return width, "ok"
}

func GoLeanBufferWriteRune(tag, offset, flag int, stored, spare []byte, value rune) string {
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
	width, status := goLeanWriteRuneStatus(b, value)
	if b == nil {
		return status
	}
	marker := "slice"
	if b.buf == nil {
		marker = "nil"
	}
	reused := len(storage) != 0 && cap(b.buf) != 0 && &storage[0] == &b.buf[:cap(b.buf)][0]
	return fmt.Sprintf("%s:%d:%d:%d:%d:%s:%t:%s:%s", status, width, b.lastRead, b.off, cap(b.buf), marker, reused,
		writeRuneWire(b.buf), writeRuneWire(b.buf[len(b.buf):cap(b.buf)]))
}

func writeRuneWire(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
