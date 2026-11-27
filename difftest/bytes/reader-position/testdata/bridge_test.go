package bytes

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

var goLeanWriteError = errors.New("scripted writer failure")
var goLeanWritePanic = new(int)

type goLeanPositionWriter struct {
	reader        *Reader
	script, count int
	trace         string
	cursor        int64
}

func (w *goLeanPositionWriter) WriteString(string) (int, error) {
	panic("unexpected dispatch to WriteString")
}

func (w *goLeanPositionWriter) Write(data []byte) (int, error) {
	if w.trace != "none" || len(data) == 0 || w.reader.prevRune != -1 || w.reader.i != w.cursor {
		panic("unexpected writer call or previous-rune marker")
	}
	if w.reader.i >= int64(len(w.reader.s)) || &data[0] != &w.reader.s[w.reader.i] {
		panic("writer did not receive the original unread slice")
	}
	fields := make([]string, len(data))
	for i, b := range data {
		fields[i] = strconv.Itoa(int(b))
	}
	w.trace = fmt.Sprintf("%d:%s", cap(data), strings.Join(fields, ","))
	switch w.script {
	case 1:
		return w.count, goLeanWriteError
	case 2:
		panic(goLeanWritePanic)
	default:
		return w.count, nil
	}
}

func GoLeanReaderPosition(tag, previous int, cursor, offset int64, whence, script, count int, raw, spare []byte) string {
	newReader := func() *Reader {
		if tag == 0 {
			return nil
		}
		r := &Reader{i: cursor, prevRune: previous}
		if tag == 2 {
			storage := make([]byte, len(raw)+len(spare))
			copy(storage, raw)
			copy(storage[len(raw):], spare)
			r.s = storage[:len(raw)]
		}
		return r
	}
	assertStorage := func(r *Reader) {
		if !Equal(r.s, raw) || !Equal(r.s[len(r.s):cap(r.s)], spare) || (r.s == nil) != (tag == 1) {
			panic("changed source storage")
		}
	}
	seek := func() (result string) {
		r := newReader()
		defer func() {
			if p := recover(); p != nil {
				err, ok := p.(runtime.Error)
				if r != nil || !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
					panic(p)
				}
				result = "nilReader"
			}
		}()
		pos, err := r.Seek(offset, whence)
		assertStorage(r)
		status := "ok"
		if err != nil {
			switch err.Error() {
			case "bytes.Reader.Seek: invalid whence":
				status = "whence"
			case "bytes.Reader.Seek: negative position":
				status = "negative"
			default:
				panic(err)
			}
		}
		return fmt.Sprintf("%s:%d:%d:%d", status, pos, r.i, r.prevRune)
	}()
	write := func() (result string) {
		r := newReader()
		probe := &goLeanPositionWriter{r, script, count, "none", cursor}
		var writer io.Writer = probe
		if script == 3 {
			writer = nil
		}
		var n int64
		status := "ok"
		defer func() {
			if p := recover(); p != nil {
				err, runtimeFault := p.(runtime.Error)
				nilFault := runtimeFault && err.Error() == "runtime error: invalid memory address or nil pointer dereference"
				switch {
				case r == nil && nilFault:
					result = "nilReader"
					return
				case writer == nil && nilFault:
					status = "nilWriter"
				case p == goLeanWritePanic:
					status = "writerPanic"
				case p == "bytes.Reader.WriteTo: invalid Write count":
					status = "invalidCount"
				default:
					panic(p)
				}
			}
			assertStorage(r)
			result = fmt.Sprintf("%s:%d:%d:%d:%s", status, n, r.i, r.prevRune, probe.trace)
		}()
		var err error
		n, err = r.WriteTo(writer)
		switch err {
		case nil:
		case io.ErrShortWrite:
			status = "shortWrite"
		case goLeanWriteError:
			status = "writerError"
		default:
			panic(err)
		}
		return
	}()
	return seek + "|" + write
}
