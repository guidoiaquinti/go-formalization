package bytes

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func goLeanGrowPrivateStatus(b *Buffer, count int) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			if failure == ErrTooLarge {
				status = "too-large"
				return
			}
			err, ok := failure.(runtime.Error)
			if ok && b == nil && err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
				status = "nil-receiver"
				return
			}
			if ok && b != nil && strings.HasPrefix(err.Error(), "runtime error: slice bounds out of range") {
				status = "bounds"
				return
			}
			panic(failure)
		}
	}()
	return fmt.Sprintf("ok:%d", b.grow(count))
}

func GoLeanBufferGrowPrivate(tag, offset, flag int, stored, spare []byte, count int) string {
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
	status := goLeanGrowPrivateStatus(b, count)
	if b == nil {
		return status
	}
	marker := "slice"
	if b.buf == nil {
		marker = "nil"
	}
	reused := len(storage) != 0 && cap(b.buf) != 0 && &storage[0] == &b.buf[:cap(b.buf)][0]
	return fmt.Sprintf("%s:%d:%d:%d:%s:%t:%s:%s", status, b.lastRead, b.off, cap(b.buf), marker, reused,
		writeWire(b.buf), writeWire(b.buf[len(b.buf):cap(b.buf)]))
}

func writeWire(raw []byte) string {
	out := make([]string, len(raw))
	for i, b := range raw {
		out[i] = strconv.Itoa(int(b))
	}
	return strings.Join(out, ",")
}
