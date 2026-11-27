package main

import (
	"fmt"
	"io"
	"math"
	"strings"
)

func main() {
	start := strings.NewReader("é")
	if err := start.UnreadByte(); err == nil || err.Error() != "strings.Reader.UnreadByte: at beginning of string" {
		panic("unexpected beginning-of-string result")
	}
	if position, err := start.Seek(0, io.SeekCurrent); err != nil || position != 0 {
		panic("failed UnreadByte changed cursor")
	}
	if _, _, err := start.ReadRune(); err != nil {
		panic(err)
	}
	if err := start.UnreadByte(); err != nil {
		panic(err)
	}
	if start.UnreadRune() == nil {
		panic("UnreadByte retained the previous-rune marker")
	}
	if b, err := start.ReadByte(); err != nil || b != 0xa9 {
		panic("UnreadByte did not return to the second byte of the rune")
	}

	pastEndRejected := 0
	for _, offset := range []int64{3, math.MaxInt64} {
		r := strings.NewReader("é")
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			panic(err)
		}
		if err := r.UnreadByte(); err != nil {
			panic(err)
		}
		if position, err := r.Seek(0, io.SeekCurrent); err != nil || position != offset-1 {
			panic("incorrect cursor beyond end of input")
		}
		if offset > r.Size() {
			pastEndRejected++
		}
	}
	if pastEndRejected != 2 {
		panic("unexpected failures for mutant rejecting positions past the end")
	}
	fmt.Printf("Reader.UnreadByte: Go passed 4 fixtures, past-end rejection mutant failed %d\n", pastEndRejected)
}
