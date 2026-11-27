package main

import (
	"fmt"
	"io"
	"math"
	"strings"
)

func main() {
	for _, c := range []struct {
		offset int64
		whence int
		err    string
	}{
		{0, 7, "strings.Reader.Seek: invalid whence"},
		{-4, io.SeekCurrent, "strings.Reader.Seek: negative position"},
		{math.MaxInt64, io.SeekCurrent, "strings.Reader.Seek: negative position"},
		{math.MaxInt64, io.SeekEnd, "strings.Reader.Seek: negative position"},
	} {
		r := preparedReader()
		got, err := r.Seek(c.offset, c.whence)
		if got != 0 || err == nil || err.Error() != c.err ||
			r.Len() != 1 || r.UnreadRune() == nil {
			panic("Seek error changed the cursor or retained the previous-rune marker")
		}
	}
	if preparedReader().UnreadRune() != nil {
		panic("delayed marker clear mutant did not retain an observable marker")
	}

	r := preparedReader()
	if pos, err := r.Seek(math.MaxInt64, io.SeekStart); pos != math.MaxInt64 || err != nil {
		panic("SeekStart rejected the maximum int64 offset")
	}
	if pos, err := r.Seek(1, io.SeekCurrent); pos != 0 || err == nil ||
		err.Error() != "strings.Reader.Seek: negative position" {
		panic("wrapped relative seek did not differ from saturation")
	}
	if pos, err := r.Seek(-1, io.SeekCurrent); pos != math.MaxInt64-1 || err != nil {
		panic("failed relative seek changed the cursor")
	}
	fmt.Println("Reader.Seek: Go passed error ordering and wraparound fixtures; delayed clearing and saturation rejected")
}

func preparedReader() *strings.Reader {
	r := strings.NewReader("xéy")
	if _, err := r.ReadByte(); err != nil {
		panic(err)
	}
	if _, _, err := r.ReadRune(); err != nil {
		panic(err)
	}
	return r
}
