package main

import (
	"fmt"
	"io"
	"math"
	"strings"
)

func main() {
	var reader strings.Reader
	if reader.Size() != 0 {
		panic("zero Reader size is nonzero")
	}
	var real, unreadLength int
	for _, text := range []string{"", "abc", "é🙂"} {
		reader.Reset(text)
		for _, seek := range []int64{0, 1, int64(len(text)), math.MaxInt64} {
			if _, err := reader.Seek(seek, io.SeekStart); err != nil {
				panic(err)
			}
			want := int64(len(text))
			if reader.Size() != want {
				real++
			}
			if int64(reader.Len()) != want {
				unreadLength++
			}
		}
	}
	fmt.Printf("Reader.Size: Go=%d/12 unread-length=%d/12\n", real, unreadLength)
	if real != 0 || unreadLength != 6 {
		panic("unexpected Reader.Size control counts")
	}
}
