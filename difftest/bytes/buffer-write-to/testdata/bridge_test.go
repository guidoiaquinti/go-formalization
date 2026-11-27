package bytes

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

var goLeanWriteToError = errors.New("scripted writer error")

type goLeanWriteToWriter func([]byte) (int, error)

func (w goLeanWriteToWriter) Write(p []byte) (int, error) { return w(p) }

func GoLeanBufferWriteTo(tag, offset, flag int, stored, spare []byte, accepted int, failed bool) (result string) {
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
	call := "-"
	snapshot := func(status string, n int64) string {
		marker := "slice"
		if b.buf == nil {
			marker = "nil"
		}
		return fmt.Sprintf("%s:%d:%s:%s:%d:%d:%s:%s", status, n, call, marker, b.off, b.lastRead, wire(b.buf), wire(b.buf[len(b.buf):cap(b.buf)]))
	}
	defer func() {
		if failure := recover(); failure != nil {
			if b == nil {
				if err, ok := failure.(runtime.Error); ok && err.Error() == "runtime error: invalid memory address or nil pointer dereference" {
					result = "nil-receiver"
					return
				}
			}
			if failure == "bytes.Buffer.WriteTo: invalid Write count" {
				result = snapshot("invalid", 0)
				return
			}
			panic(failure)
		}
	}()
	writer := goLeanWriteToWriter(func(p []byte) (int, error) {
		if call != "-" {
			panic("multiple WriteTo calls")
		}
		call = wire(p) + "/" + wire(p[len(p):cap(p)])
		if failed {
			return accepted, goLeanWriteToError
		}
		return accepted, nil
	})
	n, err := b.WriteTo(writer)
	status := "ok"
	switch err {
	case nil:
	case goLeanWriteToError:
		status = "writer"
	case io.ErrShortWrite:
		status = "short"
	default:
		panic(err)
	}
	return snapshot(status, n)
}
