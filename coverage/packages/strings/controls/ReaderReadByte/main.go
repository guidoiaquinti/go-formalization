package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	cases := []struct {
		input string
		want  byte
		err   error
	}{
		{"é\xff", 0xff, nil},
		{"é\x00z", 0, nil},
		{"é", 0, io.EOF},
	}
	var real, preserveAtEOF int
	for _, c := range cases {
		reader := strings.NewReader(c.input)
		if _, _, err := reader.ReadRune(); err != nil {
			panic(err)
		}
		b, err := reader.ReadByte()
		markerCleared := reader.UnreadRune() != nil
		position, seekErr := reader.Seek(0, io.SeekCurrent)
		wantPosition := int64(len("é"))
		if c.err == nil {
			wantPosition++
		}
		if b != c.want || err != c.err || !markerCleared || seekErr != nil || position != wantPosition {
			real++
		}

		mutant := strings.NewReader(c.input)
		if _, _, err := mutant.ReadRune(); err != nil {
			panic(err)
		}
		// Leave the previous-rune marker untouched when returning EOF.
		if mutant.Len() == 0 && mutant.UnreadRune() == nil {
			preserveAtEOF++
		}
	}
	fmt.Printf("Reader.ReadByte: Go=%d/3 preserve-marker-at-EOF=%d/3\n", real, preserveAtEOF)
	if real != 0 || preserveAtEOF != 1 {
		panic("unexpected Reader.ReadByte control counts")
	}
}
