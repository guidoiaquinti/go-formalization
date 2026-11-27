package main

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type writer struct {
	reader *strings.Reader
	count  int
	err    error
	calls  int
}

func (w *writer) Write([]byte) (int, error) {
	panic("Write called instead of WriteString")
}

func (w *writer) WriteString(s string) (int, error) {
	w.calls++
	if s != "xyz" || reflect.ValueOf(w.reader).Elem().FieldByName("prevRune").Int() != -1 {
		panic("wrong suffix or rune marker at callback")
	}
	return w.count, w.err
}

func main() {
	token := errors.New("writer failed")
	for _, c := range []struct {
		count int
		err   error
		want  error
	}{
		{1, nil, io.ErrShortWrite},
		{1, token, token},
		{-1, nil, io.ErrShortWrite},
		{3, nil, nil},
	} {
		r := preparedReader()
		w := &writer{reader: r, count: c.count, err: c.err}
		n, err := r.WriteTo(w)
		if n != int64(c.count) || err != c.want || w.calls != 1 || r.Len() != 3-c.count {
			panic("incorrect writer outcome, error identity, or cursor")
		}
		if c.err == token && shortWriteFirst(c.count, 3, c.err) == err {
			panic("short-write-before-writer-error mutant matched")
		}
	}
	r := preparedReader()
	w := &writer{reader: r, count: 4}
	if failure := writePanic(r, w); failure != "strings.Reader.WriteTo: invalid WriteString count" ||
		r.Len() != 3 || w.calls != 1 {
		panic("excessive count changed cursor or did not panic")
	}
	if n, err := strings.NewReader("").WriteTo(nil); n != 0 || err != nil {
		panic("EOF required a writer")
	}
	fmt.Println("Reader.WriteTo: Go passed six fixtures; Write dispatch and short-write precedence mutants rejected")
}

func preparedReader() *strings.Reader {
	r := strings.NewReader("éxyz")
	if _, _, err := r.ReadRune(); err != nil {
		panic(err)
	}
	return r
}

func shortWriteFirst(count, length int, writerError error) error {
	if count != length {
		return io.ErrShortWrite
	}
	return writerError
}

func writePanic(r *strings.Reader, w io.Writer) (failure any) {
	defer func() { failure = recover() }()
	_, _ = r.WriteTo(w)
	return nil
}
