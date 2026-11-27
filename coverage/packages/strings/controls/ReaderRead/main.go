package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	atEnd := afterRune("é")
	if n, err := atEnd.Read(nil); n != 0 || err != io.EOF || atEnd.UnreadRune() != nil {
		panic("EOF changed the previous-rune marker or returned the wrong result")
	}
	inRange := afterRune("éxy")
	if n, err := inRange.Read(nil); n != 0 || err != nil || inRange.UnreadRune() == nil || inRange.Len() != 2 {
		panic("empty in-range read did not clear only the marker")
	}
	short := afterRune("éxy")
	dst := []byte("L???R")
	n, err := short.Read(dst[1:4])
	if n != 2 || err != nil || !bytes.Equal(dst, []byte("Lxy?R")) || short.Len() != 0 {
		panic("short read copied the wrong bytes or returned EOF")
	}

	// Seeking to the current offset clears the marker without advancing.
	clearOnEOF := afterRune("é")
	if _, err := clearOnEOF.Seek(0, io.SeekCurrent); err != nil {
		panic(err)
	}
	_, _ = clearOnEOF.Read(nil)
	if clearOnEOF.UnreadRune() == nil {
		panic("EOF marker-clearing mutant matched")
	}
	if _, err := shortReadEOF(afterRune("éxy"), make([]byte, 3)); err != io.EOF {
		panic("short-read EOF mutant did not differ")
	}
	fmt.Println("Reader.Read: Go passed EOF, empty, and short-copy fixtures; both mutants rejected")
}

func afterRune(input string) *strings.Reader {
	reader := strings.NewReader(input)
	if _, _, err := reader.ReadRune(); err != nil {
		panic(err)
	}
	return reader
}

func shortReadEOF(r *strings.Reader, dst []byte) (int, error) {
	n, err := r.Read(dst)
	if err == nil && n < len(dst) {
		return n, io.EOF
	}
	return n, err
}
