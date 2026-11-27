package bytes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func goLeanUnreadRuneStatus(b *Buffer) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			err, ok := failure.(runtime.Error)
			if b != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			status = "panic"
		}
	}()
	var duplicate Buffer
	if b != nil {
		duplicate = *b
	}
	err := b.UnreadRune()
	if err == nil {
		return "ok"
	}
	const message = "bytes.Buffer: UnreadRune: previous operation was not a successful ReadRune"
	again := duplicate.UnreadRune()
	if err.Error() != message || again == nil || again.Error() != message || err == again || err == errUnreadByte {
		panic("expected a fresh UnreadRune error")
	}
	return "fresh-error"
}

func GoLeanBufferUnreadRune(tag, offset, flag int, stored, spare []byte) string {
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
	status := goLeanUnreadRuneStatus(b)
	if b == nil {
		return status
	}
	view := b.Bytes()
	marker := "slice"
	if view == nil {
		marker = "nil"
	}
	return fmt.Sprintf("%s:%d:%d:%d:%s:%s:%s:%s", status, b.lastRead, b.off, cap(view), marker,
		unreadRuneWire(view), unreadRuneWire(view[len(view):cap(view)]), unreadRuneWire(b.buf))
}

func unreadRuneWire(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
