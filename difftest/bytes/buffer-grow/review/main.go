package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type fixture struct {
	name, status        string
	tag, length, offset int
	backing             []byte
	count, maxInt       int
	resultOffset, flag  int
	stored, spare       []byte
	reused, leanOnly    bool
}

func main() {
	maxInt := int(^uint(0) >> 1)
	old := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	cases := []fixture{
		{name: "negative precedes nil panic", status: "negative", count: -1},
		{name: "nil receiver", status: "nil-receiver"},
		{name: "nil empty retains nil", status: "ok", tag: 1},
		{name: "nil reserve small", status: "ok", tag: 1, count: 64, spare: make([]byte, 64)},
		{name: "fit preserves read", status: "ok", tag: 2, length: 5, offset: 2, count: 3,
			backing: old, resultOffset: 2, flag: -1, stored: old[:5], spare: old[5:], reused: true},
		{name: "consumed reset even with zero", status: "ok", tag: 2, length: 5, offset: 5,
			backing: old, spare: old, reused: true},
		{name: "compaction at half threshold", status: "ok", tag: 2, length: 8, offset: 6, count: 2,
			backing: old, flag: -1, stored: old[6:], spare: old[2:], reused: true},
		{name: "allocation above half threshold", status: "ok", tag: 2, length: 8, offset: 5, count: 2,
			backing: old, flag: -1, stored: old[5:], spare: make([]byte, 13)},
		{name: "overflow after consumed reset", status: "too-large", tag: 2, length: 5, offset: 5,
			backing: old, count: maxInt, spare: old, reused: true},
		{name: "overflow keeps partial read", status: "too-large", tag: 2, length: 5, offset: 2,
			count: maxInt, backing: old, resultOffset: 2, flag: -1, stored: old[:5], spare: old[5:], reused: true},
		{name: "overflow threshold equality allocates", status: "ok", tag: 2, length: 8, offset: 4,
			count: 3, maxInt: 19, backing: old, flag: -1, stored: old[4:], spare: make([]byte, 12)},
		{name: "overflow threshold plus one panics", status: "too-large", tag: 2, length: 8, offset: 4,
			count: 4, maxInt: 19, backing: old, resultOffset: 4, flag: -1, stored: old, reused: true, leanOnly: true},
	}
	var rows bytes.Buffer
	var native int
	for _, c := range cases {
		if !c.leanOnly {
			c.checkGo()
			native++
		}
		c.write(&rows, maxInt)
	}
	cmd := exec.Command(filepath.Join("lean", ".lake", "build", "bin", "godiff-bytes-buffer-grow"))
	cmd.Stdin = &rows
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Lean control: %v: %s", err, output)
	}
	results := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(results) != len(cases) {
		log.Fatalf("got %d results for %d fixtures", len(results), len(cases))
	}
	for i, c := range cases {
		if want := c.expected(); results[i] != want {
			log.Fatalf("%s: Lean=%q want %q", c.name, results[i], want)
		}
	}
	fmt.Printf("Buffer.Grow reviewer control: %d native and %d Lean fixtures agree.\n", native, len(cases))
}

func grow(b *bytes.Buffer, n int) (status string) {
	defer func() {
		if failure := recover(); failure != nil {
			switch {
			case failure == bytes.ErrTooLarge:
				status = "too-large"
			case failure == "bytes.Buffer.Grow: negative count":
				status = "negative"
			default:
				err, ok := failure.(runtime.Error)
				if !ok || b != nil || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
					panic(failure)
				}
				status = "nil-receiver"
			}
		}
	}()
	b.Grow(n)
	return "ok"
}

func (c fixture) checkGo() {
	backing := bytes.Clone(c.backing)
	var b *bytes.Buffer
	if c.tag == 1 {
		b = new(bytes.Buffer)
	} else if c.tag == 2 {
		b = bytes.NewBuffer(backing[:c.length:len(backing)])
	}
	if c.offset > 0 {
		b.Next(c.offset)
	}
	if got := grow(b, c.count); got != c.status {
		log.Fatalf("%s: native status=%s want %s", c.name, got, c.status)
	}
	if b == nil {
		return
	}
	capacity := len(c.stored) + len(c.spare)
	if b.Cap() != capacity || (b.Bytes() == nil) != (c.tag == 1 && capacity == 0) {
		log.Fatalf("%s: capacity or nilness mismatch", c.name)
	}
	if !bytes.Equal(b.Bytes(), c.stored[c.resultOffset:]) {
		log.Fatalf("%s: contents=%v want %v", c.name, b.Bytes(), c.stored[c.resultOffset:])
	}
	current := b.Bytes()[:capacity-c.resultOffset]
	want := append(bytes.Clone(c.stored), c.spare...)
	if !bytes.Equal(current, want[c.resultOffset:]) {
		log.Fatalf("%s: accessible backing=%v want %v", c.name, current, want[c.resultOffset:])
	}
	if c.reused {
		if !bytes.Equal(backing, want) || (len(current) > 0 && &current[0] != &backing[c.resultOffset]) {
			log.Fatalf("%s: backing reuse mismatch", c.name)
		}
	} else if !bytes.Equal(backing, c.backing) || (len(current) > 0 && len(backing) > 0 && &current[0] == &backing[0]) {
		log.Fatalf("%s: allocation changed or retained original backing", c.name)
	}
	if (b.UnreadByte() == nil) != (c.flag != 0) {
		log.Fatalf("%s: read flag mismatch", c.name)
	}
}

func (c fixture) write(rows *bytes.Buffer, maxInt int) {
	if c.maxInt != 0 {
		maxInt = c.maxInt
	}
	fmt.Fprintf(rows, "%d %d %d", c.tag, c.count, maxInt)
	if c.tag != 0 {
		flag := 0
		if c.offset > 0 {
			flag = -1
		}
		fmt.Fprintf(rows, " %d %d %d", c.length, c.offset, flag)
		for _, value := range c.backing {
			fmt.Fprintf(rows, " %d", value)
		}
	}
	fmt.Fprintln(rows)
}

func (c fixture) expected() string {
	if c.tag == 0 {
		return c.status
	}
	capacity := len(c.stored) + len(c.spare)
	marker := "slice"
	if c.tag == 1 && capacity == 0 {
		marker = "nil"
	}
	return fmt.Sprintf("%s:%d:%d:%d:%s:%t:%s:%s", c.status, c.flag, c.resultOffset,
		capacity, marker, c.reused, wire(c.stored), wire(c.spare))
}

func wire(block []byte) string {
	fields := make([]string, len(block))
	for i, value := range block {
		fields[i] = strconv.Itoa(int(value))
	}
	return strings.Join(fields, ",")
}
