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

type fixture struct {
	label string
	b     *bytes.Buffer
	want  int
}

func main() {
	newBuffer := func(s string, capacity int) *bytes.Buffer {
		return bytes.NewBuffer(append(make([]byte, 0, capacity), s...))
	}
	view := newBuffer("abcdef", 20)
	view.Next(2)
	truncated := newBuffer("abcdef", 20)
	truncated.Next(2)
	truncated.Truncate(1)
	empty := newBuffer("abc", 16)
	empty.Next(3)
	reset := newBuffer("abc", 16)
	reset.Next(2)
	reset.Reset()
	resliced := newBuffer("abcdef", 20)
	resliced.Next(2)
	resliced.WriteString("gh")
	compact := newBuffer("abcdefghijklmnopqrst", 20)
	compact.Next(17)
	compact.WriteString("xyz")
	grown := newBuffer("abc", 3)
	grown.Next(1)
	grown.WriteString(strings.Repeat("X", 100))
	// The forced-growth setup leaves offset zero and a newly allocated
	// array, so its exact capacity is visible in the unread slice.
	allocated := cap(grown.Bytes())
	if allocated < 102 {
		log.Fatalf("forced-growth capacity = %d, want at least 102", allocated)
	}
	grown.Next(13)
	cases := []fixture{
		{"zero value", new(bytes.Buffer), 0},
		{"spare storage", newBuffer("a", 16), 16},
		{"read prefix", view, 20},
		{"truncated unread data", truncated, 20},
		{"all read", empty, 16},
		{"reset", reset, 16},
		{"resliced append after read", resliced, 20},
		{"compacted append after read", compact, 20},
		{"growing append after read", grown, allocated},
	}
	var requests bytes.Buffer
	var wants []string
	var unreadCapacity, unreadLength, available int
	for _, c := range cases {
		got := c.b.Cap()
		if got != c.want {
			log.Fatalf("%s: Cap=%d, want=%d", c.label, got, c.want)
		}
		if cap(c.b.Bytes()) != c.want {
			unreadCapacity++
		}
		if c.b.Len() != c.want {
			unreadLength++
		}
		if c.b.Available() != c.want {
			available++
		}
		fmt.Fprintln(&requests, 1, c.want)
		wants = append(wants, strconv.Itoa(c.want))
	}
	if !nilPanics() {
		log.Fatal("Cap did not report the nil-pointer runtime panic")
	}
	fmt.Fprintln(&requests, 0)
	wants = append(wants, "panic")
	for i, got := range check.LeanRows("godiff-bytes-buffer-cap", len(wants), &requests) {
		if got != wants[i] {
			log.Fatalf("Lean fixture %d: Cap=%s, want=%s", i, got, wants[i])
		}
	}
	if unreadCapacity != 5 || unreadLength != 8 || available != 7 {
		log.Fatalf("unexpected mutant failures: unread capacity=%d, unread length=%d, available=%d",
			unreadCapacity, unreadLength, available)
	}
	fmt.Printf("bytes.Buffer.Cap passed %d fixed fixtures; unread capacity=%d, unread length=%d, available=%d mutant failures\n",
		len(wants), unreadCapacity, unreadLength, available)
}

func nilPanics() (panicked bool) {
	defer func() {
		if err, ok := recover().(runtime.Error); ok {
			panicked = err.Error() == "runtime error: invalid memory address or nil pointer dereference"
		}
	}()
	var b *bytes.Buffer
	b.Cap()
	return false
}
