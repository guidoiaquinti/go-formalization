package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	var real, retainedCursor int
	for _, replacement := range []string{"", "xyz", "\xffé"} {
		r := strings.NewReader("éold")
		_, oldCursor, err := r.ReadRune()
		if err != nil || oldCursor != 2 {
			panic("unexpected starting cursor")
		}
		r.Reset(replacement)
		if r.UnreadRune() == nil {
			panic("Reset retained the previous-rune marker")
		}
		position, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			panic(err)
		}
		result, err := io.ReadAll(r)
		if err != nil {
			panic(err)
		}
		if position != 0 || string(result) != replacement {
			real++
		}
		if oldCursor != 0 {
			retainedCursor++
		}
	}
	fmt.Printf("Reader.Reset: Go=%d/3 retained-cursor=%d/3\n", real, retainedCursor)
	if real != 0 || retainedCursor != 3 {
		panic("unexpected Reader.Reset control counts")
	}
}
