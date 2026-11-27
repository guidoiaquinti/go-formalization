package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

//go:noinline
func escaped(s string) *bytes.Buffer { return bytes.NewBufferString(s) }

//go:noinline
func rawByte(s string, offset int) byte { return s[offset] }

func main() {
	repeated := func(n int) string { return strings.Repeat("\xff", n) }
	cases := []struct {
		input    string
		capacity int
	}{
		{"", 0},
		{"\xff\x00\xc0\xaf", 8},
		{repeated(8), 8},
		{repeated(9), 16},
		{repeated(256), 256},
		{repeated(257), 288},
		{repeated(480), 480},
		{repeated(481), 512},
		{repeated(512), 512},
	}
	var rows bytes.Buffer
	var expected []string
	capacityFailures, nilFailures := 0, 0
	for _, c := range cases {
		n := len(c.input)
		spare := make([]byte, c.capacity-n)
		b := escaped(c.input)
		view := b.Bytes()
		if view == nil || b.Len() != n || b.Cap() != c.capacity ||
			cap(view) != c.capacity || string(view) != c.input ||
			!bytes.Equal(view[n:cap(view)], spare) {
			log.Fatalf("native length %d: incorrect geometry or contents", n)
		}
		if b.UnreadByte() == nil || b.UnreadRune() == nil {
			log.Fatal("accepted unread before reading")
		}
		if n > 0 {
			for _, i := range []int{0, n - 1} {
				old := rawByte(c.input, i)
				view[i] ^= 1
				if rawByte(c.input, i) != old || b.Bytes()[i] == old {
					log.Fatal("mutated view lost input isolation")
				}
			}
			fresh := escaped(c.input)
			if fresh.String() != c.input {
				log.Fatal("second constructor observed a previous mutation")
			}
		}
		fmt.Fprint(&rows, "0")
		for _, value := range []byte(c.input) {
			fmt.Fprint(&rows, " ", value)
		}
		fmt.Fprintln(&rows)
		expected = append(expected, fmt.Sprintf("slice:%d:%d:0:0:%s:%s", n, c.capacity, wire([]byte(c.input)), wire(spare)))
		if n != c.capacity {
			capacityFailures++
		}
		if n == 0 {
			nilFailures++
		}
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-new-buffer-string")
	cmd.Stdin = &rows
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	actual := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(actual) != len(expected) {
		log.Fatalf("got %d model fixtures want %d", len(actual), len(expected))
	}
	for i, got := range actual {
		if got != expected[i] {
			log.Fatalf("fixture %d: Lean=%q want %q", i, got, expected[i])
		}
	}
	if capacityFailures != 4 || nilFailures != 1 {
		log.Fatalf("mutation counts capacity=%d nilness=%d", capacityFailures, nilFailures)
	}
	fmt.Printf("NewBufferString reviewer control: %d native and Lean fixtures passed; exact-length and nil-empty alternatives fail %d/%d.\n", len(actual), capacityFailures, nilFailures)
}

func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}
