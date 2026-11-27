package bytes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func goLeanUnreadByteStatus(b *Buffer) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if b != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			status = "panic"
		}
	}()
	err := b.UnreadByte()
	if err == nil {
		return "ok"
	}
	const message = "bytes.Buffer: UnreadByte: previous operation was not a successful read"
	if err != errUnreadByte || err.Error() != message {
		panic("expected sentinel errUnreadByte")
	}
	return "sentinel-error"
}

func GoLeanBufferUnreadByte(tag, offset, flag int, stored, spare []byte) string {
	var b *Buffer
	if tag != 0 {
		b = &Buffer{off: offset, lastRead: readOp(flag)}
		if tag == 2 {
			storage := make([]byte, len(stored)+len(spare))
			copy(storage, stored)
			copy(storage[len(stored):], spare)
			b.buf = storage[:len(stored)]
		}
	}
	status := goLeanUnreadByteStatus(b)
	if b == nil {
		return status
	}
	view := b.Bytes()
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	return fmt.Sprintf("%s:%d:%d:%d:%s:%s:%s:%s", status, b.lastRead, b.off, cap(view), marker,
		unreadByteWire(view), unreadByteWire(view[len(view):cap(view)]), unreadByteWire(b.buf))
}

func unreadByteWire(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
