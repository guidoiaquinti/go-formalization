package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type fixture struct {
	name                    string
	tag, offset, length     int
	backing, detached       []byte
	aliasStart, aliasLength int
	resultOffset            int
	result, spare           []byte
	reused                  bool
}

func main() {
	ten := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	cases := []fixture{
		{name: "nil receiver"},
		{name: "nil empty", tag: 1, aliasStart: -1},
		{name: "nil small allocation", tag: 1, detached: []byte{255, 0}, aliasStart: -1,
			result: []byte{255, 0}, spare: make([]byte, 62)},
		{name: "nonnil zero capacity", tag: 2, backing: []byte{}, aliasStart: -1},
		{name: "zero write preserves consumed prefix", tag: 2, length: 6, offset: 6,
			backing: ten[:6], aliasStart: -1, resultOffset: 6, result: ten[:6], reused: true},
		{name: "overlapping fast append from spare", tag: 2, length: 3, offset: 1,
			backing: ten[:8], aliasStart: 4, aliasLength: 3, resultOffset: 1,
			result: []byte{1, 2, 3, 5, 6, 7}, spare: []byte{7, 8}, reused: true},
		{name: "consumed reset and overlap", tag: 2, length: 6, offset: 6,
			backing: ten[:6], aliasStart: 2, aliasLength: 4,
			result: []byte{3, 4, 5, 6}, spare: []byte{5, 6}, reused: true},
		{name: "compaction changes aliased source", tag: 2, length: 10, offset: 8,
			backing: ten, aliasStart: 0, aliasLength: 3,
			result: []byte{9, 10, 9, 10, 3}, spare: ten[5:], reused: true},
		{name: "above compaction cutoff keeps old source", tag: 2, length: 10, offset: 7,
			backing: ten, aliasStart: 0, aliasLength: 3,
			result: []byte{8, 9, 10, 1, 2, 3}, spare: make([]byte, 10)},
	}
	var rows bytes.Buffer
	for _, c := range cases {
		c.checkGo()
		c.write(&rows)
	}
	cmd := exec.Command(filepath.Join("lean", ".lake", "build", "bin", "godiff-bytes-buffer-write"))
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
	fmt.Printf("Buffer.Write reviewer control: %d native and Lean fixtures agree.\n", len(cases))
}

func (c fixture) count() int {
	if c.aliasStart >= 0 {
		return c.aliasLength
	}
	return len(c.detached)
}

func (c fixture) checkGo() {
	if c.tag == 0 {
		defer func() {
			if recover() == nil {
				log.Fatal("nil receiver did not panic")
			}
		}()
		var receiver *bytes.Buffer
		receiver.Write(nil)
		return
	}
	backing := bytes.Clone(c.backing)
	b := new(bytes.Buffer)
	if c.tag == 2 {
		b = bytes.NewBuffer(backing[:c.length:len(backing)])
	}
	if c.offset > 0 {
		b.Next(c.offset)
	}
	input := c.detached
	if c.aliasStart >= 0 {
		input = backing[c.aliasStart : c.aliasStart+c.aliasLength]
	}
	n, err := b.Write(input)
	if n != c.count() || err != nil || b.UnreadByte() == nil {
		log.Fatalf("%s: write result or read flag mismatch", c.name)
	}
	capacity := len(c.result) + len(c.spare)
	if b.Cap() != capacity || (b.Bytes() == nil) != (c.tag == 1 && n == 0) {
		log.Fatalf("%s: capacity or nilness mismatch", c.name)
	}
	current := b.Bytes()[:b.Cap()-c.resultOffset]
	want := append(bytes.Clone(c.result), c.spare...)
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
}

func (c fixture) write(rows *bytes.Buffer) {
	if c.tag == 0 {
		fmt.Fprintln(rows, "0")
		return
	}
	flag := 0
	if c.offset > 0 {
		flag = -1
	}
	fmt.Fprintf(rows, "%d %d %d %d %d %d %d", c.tag, c.length, len(c.backing)-c.length,
		c.offset, flag, c.aliasStart, c.aliasLength)
	for _, block := range [][]byte{c.backing, c.detached} {
		for _, value := range block {
			fmt.Fprintf(rows, " %d", value)
		}
	}
	fmt.Fprintln(rows)
}

func (c fixture) expected() string {
	if c.tag == 0 {
		return "panic"
	}
	marker := "slice"
	if c.tag == 1 && c.count() == 0 {
		marker = "nil"
	}
	return fmt.Sprintf("ok:%d:0:%d:%d:%s:%t:%s:%s", c.count(), c.resultOffset,
		len(c.result)+len(c.spare), marker, c.reused, wire(c.result), wire(c.spare))
}

func wire(block []byte) string {
	fields := make([]string, len(block))
	for i, value := range block {
		fields[i] = strconv.Itoa(int(value))
	}
	return strings.Join(fields, ",")
}
