package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
)

type fixture struct {
	raw                []byte
	count              int
	repeatPanic, title string
}

func main() {
	large := []byte(strings.Repeat("-", 8192) + "a")
	cases := []fixture{
		{nil, -1, "bytes: negative Repeat count", ""},
		{[]byte{}, math.MaxInt, "", ""},
		{[]byte("ab"), math.MaxInt, "bytes: Repeat output length overflow", "Ab"},
		{[]byte("abc"), math.MaxInt, "bytes: Repeat output length overflow", "Abc"},
		{[]byte("x"), 0, "", "X"},
		{[]byte{255, 0, 13}, 3, "", "�\x00\r"},
		{[]byte("aB-"), 9000, "", "AB-"},
		{large, 3, "", strings.Repeat("-", 8192) + "A"},
		{[]byte("ab-cd_ef9gh"), 1, "", "Ab-Cd_ef9gh"},
		{[]byte("ab—cd\u00a0ef"), 1, "", "Ab—cd\u00a0Ef"},
		{[]byte{0xe2, 0x82, 'a', ' ', 'b'}, 1, "", "��a B"},
		{[]byte("ǳ Ǳǳ"), 1, "", "ǲ ǲǳ"},
	}
	var rows bytes.Buffer
	var want []string
	for i, c := range cases {
		fmt.Fprint(&rows, math.MaxInt, " ", c.count)
		for _, b := range c.raw {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
		repeat := assessRepeat(i, c)
		storage := append(bytes.Clone(c.raw), 0x55, 0xaa)
		var input []byte
		if c.raw != nil {
			input = storage[:len(c.raw)]
		}
		before := bytes.Clone(storage)
		title := bytes.Title(input)
		if string(title) != c.title {
			log.Fatalf("fixture %d Title=%x, want %x", i, title, c.title)
		}
		capacity := "growth"
		if len(c.title) <= len(c.raw) {
			capacity = fmt.Sprint(len(c.raw))
		}
		titleWire := observe(i, title, storage, before, capacity)
		want = append(want, repeat+"|"+titleWire)
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(rows.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-repeat-title")
	cmd.Stdin = &rows
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(got) != len(want) {
		log.Fatalf("got %d Lean rows, want %d", len(got), len(want))
	}
	for i, row := range got {
		if row != want[i] {
			log.Fatalf("Lean row %d differs from independent observation", i)
		}
	}
	fmt.Printf("bytes repeat and title review: %d native and Lean fixtures passed\n", len(cases))
}

func assessRepeat(row int, c fixture) (result string) {
	defer func() {
		p := recover()
		if p == nil && c.repeatPanic == "" {
			return
		}
		if c.repeatPanic == "" || p != c.repeatPanic {
			log.Fatalf("row %d Repeat panic=%v want %q", row, p, c.repeatPanic)
		}
		result = c.repeatPanic
	}()
	storage := append(bytes.Clone(c.raw), 0x55, 0xaa)
	var input []byte
	if c.raw != nil {
		input = storage[:len(c.raw)]
	}
	before := bytes.Clone(storage)
	output := bytes.Repeat(input, c.count)
	if c.repeatPanic != "" {
		log.Fatalf("row %d Repeat returned, expected panic %q", row, c.repeatPanic)
	}
	if len(output) != len(c.raw)*c.count {
		log.Fatalf("row %d Repeat length %d", row, len(output))
	}
	for i, b := range output {
		if b != c.raw[i%len(c.raw)] {
			log.Fatalf("row %d Repeat byte %d lost alignment", row, i)
		}
	}
	return observe(row, output, storage, before, fmt.Sprint(len(output)))
}

func observe(row int, output, storage, before []byte, capacity string) string {
	if output == nil || cap(output) < len(output) ||
		capacity != "growth" && fmt.Sprint(cap(output)) != capacity {
		log.Fatalf("row %d invalid output cap %d", row, cap(output))
	}
	encoded := make([]string, len(output))
	for i, b := range output {
		encoded[i] = fmt.Sprint(b)
		output[i] ^= 0xff
	}
	if !bytes.Equal(storage, before) {
		log.Fatalf("row %d output shares source backing", row)
	}
	return "slice:" + capacity + ":" + strings.Join(encoded, ",")
}
