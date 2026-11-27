package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

type result struct {
	text, capacity string
}

type fixture struct {
	raw     []byte
	results [3]result
}

func main() {
	mapped := []struct{ input, upper, lower, title rune }{
		{'ǳ', 'Ǳ', 'ǳ', 'ǲ'},
		{'ǲ', 'Ǳ', 'ǳ', 'ǲ'},
		{'K', 'K', 'k', 'K'},
		{'ȿ', 'Ȿ', 'ȿ', 'Ȿ'},
		{'ß', 'ß', 'ß', 'ß'},
	}
	cases := []fixture{
		{nil, [3]result{{"", "0"}, {"", "0"}, {"", "0"}}},
		{[]byte{}, [3]result{{"", "0"}, {"", "0"}, {"", "0"}}},
		{[]byte("Ab\x00z!"), [3]result{{"AB\x00Z!", "5"}, {"ab\x00z!", "5"}, {"AB\x00Z!", "5"}}},
		{[]byte("abc"), [3]result{{"ABC", "3"}, {"abc", "growth"}, {"ABC", "3"}}},
		{[]byte("ABC"), [3]result{{"ABC", "growth"}, {"abc", "3"}, {"ABC", "3"}}},
		{[]byte("ǳǱ"), [3]result{{"ǱǱ", "4"}, {"ǳǳ", "4"}, {"ǲǲ", "4"}}},
		{[]byte("Kİ"), [3]result{{"Kİ", "5"}, {"ki", "5"}, {"Kİ", "5"}}},
		{[]byte("ȿ"), [3]result{{"Ȿ", "growth"}, {"ȿ", "2"}, {"Ȿ", "growth"}}},
		{[]byte{226, 130, 'A'}, [3]result{{"��A", "growth"}, {"��a", "growth"}, {"��A", "growth"}}},
		{[]byte("K\xff"), [3]result{{"K�", "growth"}, {"k�", "4"}, {"K�", "growth"}}},
		{[]byte("ß"), [3]result{{"ß", "2"}, {"ß", "2"}, {"ß", "2"}}},
	}
	var rows bytes.Buffer
	var want []string
	for _, c := range mapped {
		fmt.Fprintln(&rows, 0, c.input)
		if unicode.ToUpper(c.input) != c.upper || unicode.ToLower(c.input) != c.lower || unicode.ToTitle(c.input) != c.title {
			log.Fatalf("native Unicode mapping of %U differs", c.input)
		}
		want = append(want, fmt.Sprintf("%d,%d,%d", c.upper, c.lower, c.title))
	}
	for i, c := range cases {
		fmt.Fprint(&rows, 1)
		for _, b := range c.raw {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
		var parts []string
		for method, f := range []func([]byte) []byte{bytes.ToUpper, bytes.ToLower, bytes.ToTitle} {
			parts = append(parts, observe(i, method, f, c.raw, c.results[method]))
		}
		want = append(want, strings.Join(parts, "|"))
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(rows.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-simple-case")
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
			log.Fatalf("Lean row %d = %q, want %q", i, row, want[i])
		}
	}
	fmt.Printf("bytes simple case reviewer control: %d byte fixtures and %d Unicode mappings passed\n", len(cases), len(mapped))
}

func observe(row, method int, f func([]byte) []byte, raw []byte, want result) string {
	storage := append(bytes.Clone(raw), 0x55, 0xaa)
	var input []byte
	if raw != nil {
		input = storage[:len(raw)]
	}
	before := bytes.Clone(storage)
	output := f(input)
	if output == nil || string(output) != want.text || cap(output) < len(output) ||
		want.capacity != "growth" && fmt.Sprint(cap(output)) != want.capacity {
		log.Fatalf("row %d method %d got nil=%t len=%d cap=%d bytes=%x want %+v", row, method, output == nil, len(output), cap(output), output, want)
	}
	encoded := make([]string, len(output))
	for i, b := range output {
		encoded[i] = fmt.Sprint(b)
		output[i] ^= 0xff
	}
	if !bytes.Equal(storage, before) {
		log.Fatalf("row %d method %d mutated source backing", row, method)
	}
	return "slice:" + want.capacity + ":" + strings.Join(encoded, ",")
}
