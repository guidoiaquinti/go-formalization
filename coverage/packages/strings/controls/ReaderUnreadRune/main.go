package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	const beginning = "strings.Reader.UnreadRune: at beginning of string"
	const missing = "strings.Reader.UnreadRune: previous operation was not ReadRune"
	r := strings.NewReader("xé")
	checkError(r.UnreadRune(), beginning)
	if _, err := r.ReadByte(); err != nil {
		panic(err)
	}
	checkError(r.UnreadRune(), missing)
	decoded, width, err := r.ReadRune()
	if err != nil || decoded != 'é' || width != 2 {
		panic("unexpected initial rune")
	}
	if err := r.UnreadRune(); err != nil {
		panic(err)
	}
	checkError(r.UnreadRune(), missing)
	decoded, width, err = r.ReadRune()
	if err != nil || decoded != 'é' || width != 2 {
		panic("rewind did not recover the entire rune")
	}
	if _, err := r.Seek(-1, io.SeekCurrent); err != nil {
		panic(err)
	}
	decoded, _, err = r.ReadRune()
	if err != nil || decoded == 'é' {
		panic("single-byte rewind mutant matched")
	}
	fmt.Println("Reader.UnreadRune: Go passed the error, rewind, and repeat checks; single-byte rewind rejected")
}

func checkError(err error, want string) {
	if err == nil || err.Error() != want {
		panic(fmt.Sprintf("got %v, want %q", err, want))
	}
}
