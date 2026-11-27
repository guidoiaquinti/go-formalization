package main

import (
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf8"
)

func main() {
	cases := []struct {
		s      string
		cursor int64
	}{
		{"", 0},
		{"abcd", 0},
		{"abcd", 3},
		{"abcd", 4},
		{"abcd", math.MaxInt64},
		{"aé🙂", 1},
		{"aé🙂", 2},
	}
	var real, originalLength, runeCount int
	for _, c := range cases {
		reader := strings.NewReader(c.s)
		if _, err := reader.Seek(c.cursor, io.SeekStart); err != nil {
			panic(err)
		}
		length := reader.Len()
		unread, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		if length != len(unread) || reader.Len() != 0 {
			real++
		}
		if len(c.s) != len(unread) {
			originalLength++
		}
		if utf8.RuneCount(unread) != len(unread) {
			runeCount++
		}
	}
	fmt.Printf("Reader.Len: Go=%d/7 original-length=%d/7 rune-count=%d/7\n", real, originalLength, runeCount)
	if real != 0 || originalLength != 5 || runeCount != 2 {
		panic("unexpected Reader.Len control counts")
	}
}
