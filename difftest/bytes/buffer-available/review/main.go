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
		length, capacity, discard int
		action                    func(*bytes.Buffer)
		fixed                     int
	}{
		{0, 0, 0, nil, 0},
		{3, 3, 0, nil, 0},
		{3, 8, 0, nil, 5},
		{3, 8, 2, nil, 5},
		{3, 3, 3, nil, 0},
		{6, 8, 2, func(b *bytes.Buffer) { b.Truncate(1) }, 5},
		{6, 8, 2, func(b *bytes.Buffer) { b.Reset() }, 8},
		{8, 8, 6, func(b *bytes.Buffer) { b.WriteString("xy") }, 4},
		{3, 3, 1, func(b *bytes.Buffer) { b.WriteString(strings.Repeat("x", 70)) }, -1},
		{3, 3, 3, func(b *bytes.Buffer) { b.WriteString("xy") }, 1},
	}
	var rows bytes.Buffer
	var expected []string
	ignoreOffset, ignoreUnread := 0, 0
	for i, c := range cases {
		b := bytes.NewBuffer(make([]byte, c.length, c.capacity))
		b.Next(c.discard)
		if c.action != nil {
			c.action(b)
		}
		unread := b.Bytes()
		want := cap(unread) - len(unread)
		if c.fixed >= 0 && want != c.fixed {
			log.Fatalf("fixture %d: unread spare=%d, expected %d", i, want, c.fixed)
		}
		if got := b.Available(); got != want {
			log.Fatalf("fixture %d: Available=%d, unread spare=%d", i, got, want)
		}
		offset := b.Cap() - cap(unread)
		fmt.Fprintln(&rows, 1, len(unread)+offset, b.Cap())
		expected = append(expected, strconv.Itoa(want))
		if b.Cap()-len(unread) != want {
			ignoreOffset++
		}
		if cap(unread) != want {
			ignoreUnread++
		}
	}
	fmt.Fprintln(&rows, 0)
	expected = append(expected, "panic")
	if !nilPanics() {
		log.Fatal("nil receiver did not trigger its expected panic")
	}
	for i, got := range check.LeanRows("godiff-bytes-buffer-available", len(expected), &rows) {
		if got != expected[i] {
			log.Fatalf("Lean fixture %d returned %q, want %q", i, got, expected[i])
		}
	}
	if ignoreOffset != 3 || ignoreUnread != 7 {
		log.Fatalf("mutation failures offset=%d and unread=%d, want 3 and 7", ignoreOffset, ignoreUnread)
	}
	fmt.Printf("bytes.Buffer.Available passed %d fixed Go and Lean fixtures; ignoring offset and ignoring unread length failed %d and %d.\n",
		len(expected), ignoreOffset, ignoreUnread)
}

func nilPanics() (panicked bool) {
	defer func() {
		if err, ok := recover().(runtime.Error); ok {
			panicked = err.Error() == "runtime error: invalid memory address or nil pointer dereference"
		}
	}()
	var b *bytes.Buffer
	b.Available()
	return false
}
