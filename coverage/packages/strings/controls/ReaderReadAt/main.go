package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("éabc")
	if _, _, err := r.ReadRune(); err != nil {
		panic(err)
	}
	dst := []byte{'L', '?', '?', '?', 'R'}
	n, err := r.ReadAt(dst[1:4], 3)
	if n != 2 || err != io.EOF || !bytes.Equal(dst, []byte("Lbc?R")) {
		panic("incorrect short read, EOF, or destination")
	}
	endCount, endError := r.ReadAt(nil, r.Size())
	if endCount != 0 || endError != io.EOF {
		panic("empty destination at EOF did not return EOF")
	}
	if n, err := r.ReadAt(nil, 0); n != 0 || err != nil {
		panic("empty destination at a valid offset returned an error")
	}
	if err := r.UnreadRune(); err != nil || r.Len() != len("éabc") {
		panic("ReadAt changed reader state")
	}
	var nilReader *strings.Reader
	if n, err := nilReader.ReadAt(dst, -1); n != 0 || err == nil ||
		err.Error() != "strings.Reader.ReadAt: negative offset" || !bytes.Equal(dst, []byte("Lbc?R")) {
		panic("negative offset was not rejected before receiver access")
	}
	if !panicsOnReadAt(nilReader) {
		panic("nil receiver with nonnegative offset did not panic")
	}
	_, emptyFirstError := readAtEmptyFirst(r, nil, r.Size())
	if emptyFirstError == endError || bytes.Equal([]byte{'L', 'b', 'c', 0, 'R'}, dst) {
		panic("empty-first or zero-filled-tail mutant matched")
	}
	fmt.Println("Reader.ReadAt: short copy, empty requests, reader state, and nil receiver ordering passed")
}

func panicsOnReadAt(r *strings.Reader) (panicked bool) {
	defer func() { panicked = recover() != nil }()
	r.ReadAt(nil, 0)
	return false
}

func readAtEmptyFirst(r *strings.Reader, dst []byte, offset int64) (int, error) {
	if len(dst) == 0 {
		return 0, nil
	}
	return r.ReadAt(dst, offset)
}
