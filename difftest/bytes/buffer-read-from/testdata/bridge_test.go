package bytes

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
)

var goLeanReadFromError = errors.New("scripted reader error")
var goLeanReadFromPanic = new(int)

type GoLeanReadStep struct {
	Mode, Count, Writes, Value int
}

type goLeanReadFromProbe struct {
	buffer *Buffer
	steps  []GoLeanReadStep
	calls  []string
}

func goLeanReadFromWire(values []byte) string {
	var out []string
	for first := 0; first < len(values); {
		end := first + 1
		for end < len(values) && values[end] == values[first] {
			end++
		}
		out = append(out, fmt.Sprintf("%d.%d", values[first], end-first))
		first = end
	}
	return strings.Join(out, ",")
}

func (r *goLeanReadFromProbe) Read(p []byte) (int, error) {
	if len(r.calls) == len(r.steps) {
		panic("ReadFrom called past end of finite script")
	}
	if len(p) < MinRead || cap(p) != len(p) || r.buffer.lastRead != opInvalid ||
		&p[0] != &r.buffer.buf[:cap(r.buffer.buf)][len(r.buffer.buf)] {
		panic("reader received a non-original or incorrectly sized write slice")
	}
	step := r.steps[len(r.calls)]
	r.calls = append(r.calls, fmt.Sprintf("%d:%s", len(p), goLeanReadFromWire(p)))
	for i := range min(step.Writes, len(p)) {
		p[i] = byte(step.Value)
	}
	switch step.Mode {
	case 1:
		return step.Count, io.EOF
	case 2:
		return step.Count, goLeanReadFromError
	case 3:
		panic(goLeanReadFromPanic)
	default:
		return step.Count, nil
	}
}

func GoLeanReadFrom(tag, length, offset, capacity, flag int, nilReader bool, steps []GoLeanReadStep) (result string) {
	var b *Buffer
	if tag != 0 {
		b = &Buffer{off: offset, lastRead: readOp(flag)}
		if tag == 2 {
			data := make([]byte, capacity)
			for i := range data {
				switch {
				case i < offset:
					data[i] = 11
				case i < length:
					data[i] = 22
				default:
					data[i] = 33
				}
			}
			b.buf = data[:length]
		}
	}
	probe := &goLeanReadFromProbe{buffer: b, steps: steps}
	var reader io.Reader = probe
	if nilReader {
		reader = nil
	}
	status := "ok"
	var count int64
	defer func() {
		if p := recover(); p != nil {
			err, ok := p.(runtime.Error)
			nilFault := ok && err.Error() == "runtime error: invalid memory address or nil pointer dereference"
			switch {
			case b == nil && nilFault:
				result = "nilReceiver"
				return
			case reader == nil && nilFault:
				status = "nilReader"
			case p == errNegativeRead:
				status = "negativeRead"
			case p == ErrTooLarge:
				status = "tooLarge"
			case p == goLeanReadFromPanic:
				status = "readerPanic"
			case ok && strings.HasPrefix(err.Error(), "runtime error: slice bounds out of range"):
				status = "bounds"
			default:
				panic(p)
			}
		}
		tag := "slice"
		if b.buf == nil {
			tag = "nil"
		}
		result = fmt.Sprintf("%s|%d|%s|%d|%d|%s|%s|%s", status, count, tag, b.off, b.lastRead,
			goLeanReadFromWire(b.buf), goLeanReadFromWire(b.buf[len(b.buf):cap(b.buf)]), strings.Join(probe.calls, ";"))
	}()
	var err error
	count, err = b.ReadFrom(reader)
	switch err {
	case nil:
	case goLeanReadFromError:
		status = "readError"
	default:
		panic(err)
	}
	return
}
