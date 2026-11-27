package main

import (
	"fmt"
	"strings"
)

func main() {
	var b strings.Builder
	if b.Cap() != 0 {
		panic("zero Builder has capacity")
	}
	b.Grow(64)
	reserve := b.Cap()
	if reserve < 64 || b.Len() != 0 {
		panic("Grow did not reserve empty backing storage")
	}
	b.WriteString("é")
	b.WriteByte(0xff)
	if b.Cap() != reserve || b.Len() != 3 {
		panic("write within reserve changed capacity or byte length")
	}
	if b.Len() == b.Cap() || b.Cap()-b.Len() == b.Cap() {
		panic("length or remaining-capacity mutant matched")
	}
	b.Reset()
	if b.Cap() != 0 || b.Len() != 0 {
		panic("Reset retained backing storage")
	}
	if !nilPanics() {
		panic("nil receiver returned normally")
	}
	fmt.Println("Builder.Cap: Go passed zero, reserve, write, reset and nil checks; both mutants rejected")
}

func nilPanics() (panicked bool) {
	defer func() { panicked = recover() != nil }()
	var b *strings.Builder
	b.Cap()
	return false
}
