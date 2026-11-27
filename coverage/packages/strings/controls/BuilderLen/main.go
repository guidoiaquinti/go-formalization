package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	var b strings.Builder
	checkLen(&b, 0)
	b.Grow(64)
	checkLen(&b, 0)
	b.WriteString("é")
	b.WriteByte(0xff)
	checkLen(&b, 3)
	b.WriteRune('🙂')
	checkLen(&b, 7)
	if b.Cap() == b.Len() || utf8.RuneCountInString(b.String()) == b.Len() {
		panic("capacity or rune-count mutant matched the byte length")
	}
	b.Reset()
	checkLen(&b, 0)
	if !nilPanics() {
		panic("nil receiver returned normally")
	}
	fmt.Println("Builder.Len: Go passed six state and nil checks; capacity and rune-count mutants rejected")
}

func checkLen(b *strings.Builder, want int) {
	if got := b.Len(); got != want {
		panic(fmt.Sprintf("Len=%d, want %d", got, want))
	}
}

func nilPanics() (panicked bool) {
	defer func() { panicked = recover() != nil }()
	var b *strings.Builder
	b.Len()
	return false
}
