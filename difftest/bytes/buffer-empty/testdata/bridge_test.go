package bytes

import (
	"runtime"
	"strconv"
)

func GoLeanBufferEmpty(nilReceiver, nilStorage bool, length, spare, offset int) (out string) {
	var b *Buffer
	if !nilReceiver {
		b = &Buffer{off: offset}
		if !nilStorage {
			b.buf = make([]byte, length, length+spare)
		}
	}
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !nilReceiver || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			out = "panic"
		}
	}()
	return strconv.FormatBool(b.empty())
}
