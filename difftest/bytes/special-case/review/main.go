package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

type fixture struct {
	raw   []byte
	table unicode.SpecialCase
	want  [3]string
}

func main() {
	cases := []fixture{
		{nil, nil, [3]string{"", "", ""}},
		{[]byte{}, nil, [3]string{"", "", ""}},
		{[]byte("İiIı"), unicode.TurkishCase, [3]string{"İİII", "iiıı", "İİII"}},
		{[]byte("Aa"), unicode.SpecialCase{{Lo: 'A', Hi: 'A', Delta: [3]rune{0, 0, 0}}}, [3]string{"AA", "Aa", "AA"}},
		{[]byte("AB"), unicode.SpecialCase{{Lo: 'A', Hi: 'A', Delta: [3]rune{-66, 1, 2}}}, [3]string{"B", "Bb", "CB"}},
		{[]byte{0xff, 0xe2, 0x82, 'A'}, unicode.SpecialCase{{Lo: 0xfffd, Hi: 0xfffd,
			Delta: [3]rune{'X' - 0xfffd, 'y' - 0xfffd, -0xfffe}}}, [3]string{"XXXA", "yyya", "A"}},
		{[]byte("\U0010ffff"), unicode.SpecialCase{{Lo: 0x10ffff, Hi: 0x10ffff, Delta: [3]rune{1, 0, -1}}}, [3]string{"�", "\U0010ffff", "\U0010fffe"}},
		{[]byte("AB"), unicode.SpecialCase{{Lo: 0x80000001, Hi: 'B', Delta: [3]rune{math.MaxInt32, math.MaxInt32, math.MaxInt32}}}, [3]string{"AA", "BB", "AA"}},
		{[]byte("aA"), unicode.SpecialCase{{Lo: 'a', Hi: 'a', Delta: [3]rune{1, 2, 3}},
			{Lo: 'A', Hi: 'z', Delta: [3]rune{0, 1, 2}}}, [3]string{"aA", "bB", "cC"}},
		{[]byte("A"), unicode.SpecialCase{{Lo: 'A', Hi: 'A', Delta: [3]rune{math.MinInt32, math.MinInt32, math.MinInt32}}}, [3]string{"", "", ""}},
	}
	var rows bytes.Buffer
	var expected []string
	for i, c := range cases {
		fmt.Fprint(&rows, len(c.table))
		for _, entry := range c.table {
			fmt.Fprint(&rows, " ", int32(entry.Lo), " ", int32(entry.Hi), " ", entry.Delta[0], " ", entry.Delta[1], " ", entry.Delta[2])
		}
		for _, b := range c.raw {
			fmt.Fprint(&rows, " ", b)
		}
		fmt.Fprintln(&rows)
		var parts []string
		for method, f := range []func(unicode.SpecialCase, []byte) []byte{bytes.ToUpperSpecial, bytes.ToLowerSpecial, bytes.ToTitleSpecial} {
			parts = append(parts, observe(i, method, f, c))
		}
		expected = append(expected, strings.Join(parts, "|"))
	}
	if len(os.Args) == 2 && os.Args[1] == "rows" {
		_, _ = os.Stdout.Write(rows.Bytes())
		return
	}
	cmd := exec.Command("lean/.lake/build/bin/godiff-bytes-special-case")
	cmd.Stdin = &rows
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	got := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(got) != len(expected) {
		log.Fatalf("got %d Lean rows, want %d", len(got), len(expected))
	}
	for i, row := range got {
		if row != expected[i] {
			log.Fatalf("Lean row %d got %q want %q", i, row, expected[i])
		}
	}
	fmt.Printf("bytes special-case review: %d fixed fixtures passed on all three entries\n", len(cases))
}

func observe(row, method int, f func(unicode.SpecialCase, []byte) []byte, c fixture) string {
	storage := append(bytes.Clone(c.raw), 0x55, 0xaa)
	var input []byte
	if c.raw != nil {
		input = storage[:len(c.raw)]
	}
	before := bytes.Clone(storage)
	out := f(c.table, input)
	want := c.want[method]
	capacity := "growth"
	if len(want) <= len(c.raw) {
		capacity = fmt.Sprint(len(c.raw))
	}
	if out == nil || string(out) != want || cap(out) < len(out) ||
		capacity != "growth" && fmt.Sprint(cap(out)) != capacity {
		log.Fatalf("row %d method %d got %x cap %d, want %x cap %s", row, method, out, cap(out), want, capacity)
	}
	encoded := make([]string, len(out))
	for i, b := range out {
		encoded[i] = fmt.Sprint(b)
		out[i] ^= 0xff
	}
	if !bytes.Equal(storage, before) {
		log.Fatalf("row %d method %d aliases input backing", row, method)
	}
	return "slice:" + capacity + ":" + strings.Join(encoded, ",")
}
