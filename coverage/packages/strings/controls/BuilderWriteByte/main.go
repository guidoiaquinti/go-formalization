package main

import (
	"fmt"
	"strings"
)

func main() {
	var b strings.Builder
	b.Grow(512)
	for i := range 256 {
		if err := b.WriteByte(byte(i)); err != nil {
			panic(err)
		}
	}
	s := b.String()
	for i := range 256 {
		if s[i] != byte(i) {
			panic("WriteByte did not append exactly the input byte")
		}
	}
	copied := b
	if failure := failedWrite(&copied); failure != "strings: illegal use of non-zero Builder copied by value" ||
		copied.String() != s || b.String() != s {
		panic("copied builder did not panic before changing the contents")
	}
	if len(string(rune(0xff))) == len(s[255:]) {
		panic("rune-encoding mutant unexpectedly matched")
	}
	fmt.Println("Builder.WriteByte: all 256 bytes and copied-builder check passed; rune-encoding mutant rejected")
}

func failedWrite(b *strings.Builder) (failure any) {
	defer func() { failure = recover() }()
	_ = b.WriteByte(1)
	return nil
}
