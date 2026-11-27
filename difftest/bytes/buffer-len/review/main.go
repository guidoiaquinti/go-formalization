package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		initial string
		discard int
		action  func(*bytes.Buffer)
		want    int
	}{
		{"", 0, nil, 0},
		{"abc", 0, nil, 3},
		{"abcdef", 2, nil, 4},
		{"abc", 3, nil, 0},
		{"abcdef", 2, func(b *bytes.Buffer) { b.Truncate(1) }, 1},
		{"abcdef", 2, func(b *bytes.Buffer) { b.Reset() }, 0},
		{"abcdefgh", 6, func(b *bytes.Buffer) { b.WriteString("xy") }, 4},
		{"abc", 1, func(b *bytes.Buffer) { b.WriteString(strings.Repeat("x", 70)) }, 72},
		{"abc", 3, func(b *bytes.Buffer) { b.WriteString("xy") }, 2},
	}
	var rows bytes.Buffer
	var expected []string
	retainedLengthFailures := 0
	for i, c := range cases {
		b := bytes.NewBuffer(append(make([]byte, 0, len(c.initial)), c.initial...))
		b.Next(c.discard)
		if c.action != nil {
			c.action(b)
		}
		if got := b.Len(); got != c.want || len(b.Bytes()) != c.want {
			log.Fatalf("fixture %d: Len=%d, Bytes length=%d, want %d", i, got, len(b.Bytes()), c.want)
		}
		length := b.Cap() - b.Available()
		offset := length - len(b.Bytes())
		fmt.Fprintln(&rows, 1, length, offset)
		expected = append(expected, strconv.Itoa(c.want))
		if length != c.want {
			retainedLengthFailures++
		}
	}
	fmt.Fprintln(&rows, 0)
	expected = append(expected, "panic")
	if !nilPanics() {
		log.Fatal("nil receiver did not trigger its expected panic")
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-len", len(expected), &rows) {
		if got != expected[i] {
			log.Fatalf("Lean fixture %d returned %q, want %q", i, got, expected[i])
		}
	}
	if retainedLengthFailures != 3 {
		log.Fatalf("retained-length mutant failed %d fixtures, want 3", retainedLengthFailures)
	}
	fmt.Printf("bytes.Buffer.Len passed %d fixed Go and Lean fixtures; retained-length mutant failed %d.\n",
		len(expected), retainedLengthFailures)
}

func nilPanics() (panicked bool) {
	defer func() {
		if err, ok := recover().(runtime.Error); ok {
			panicked = err.Error() == "runtime error: invalid memory address or nil pointer dereference"
		}
	}()
	var b *bytes.Buffer
	b.Len()
	return false
}
