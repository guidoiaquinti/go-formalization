package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		text, cutset string
		spare        int
		want         string
	}{
		{"", "", 0, "nil"},
		{"", "a", 4, "nil"},
		{"a", "a", 2, "nil"},
		{"abc", "abc", 0, "nil"},
		{"\xff\x80", "�", 2, "nil"},
		{"aaXaa", "a", 4, "7:88,97,97"},
		{"baXba", "ab", 4, "7:88,98,97"},
		{"é🚀Xé🚀", "é🚀", 3, "10:88,195,169,240,159,154,128"},
		{"\xff�x\x80", "�", 2, "4:120,128"},
		{"\xffx\xfe", "a", 5, "8:255,120,254"},
		{"éxé", "", 2, "7:195,169,120,195,169"},
		{"\xe2\x82x\xed\xa0", "\xff", 3, "6:120,237,160"},
		{"éa", "\xc3", 2, "5:195,169,97"},
	}
	var inputs bytes.Buffer
	var nilFailures, capFailures, byteFailures, directionFailures int
	for _, c := range cases {
		row := fmt.Appendf(nil, "%d %d", len(c.cutset), c.spare)
		for _, text := range []string{c.cutset, c.text} {
			for i := range len(text) {
				row = strconv.AppendUint(append(row, ' '), uint64(text[i]), 10)
			}
		}
		fmt.Fprintln(&inputs, string(row))
		var input []byte
		if len(c.text)+c.spare != 0 {
			input = make([]byte, len(c.text), len(c.text)+c.spare)
			copy(input, c.text)
		}
		result := bytes.TrimLeft(input, c.cutset)
		if got := observe(result); got != c.want {
			log.Fatalf("TrimLeft(%q,%q): Go=%s want=%s", c.text, c.cutset, got, c.want)
		}
		if result == nil && observe([]byte{}) != c.want {
			nilFailures++
		}
		if result != nil && observe(result[:len(result):len(result)]) != c.want {
			capFailures++
		}
		if observe(trimLeftBytes(input, c.cutset)) != c.want {
			byteFailures++
		}
		if observe(bytes.Trim(input, c.cutset)) != c.want {
			directionFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-left", len(cases), &inputs) {
		if row != cases[i].want {
			log.Fatalf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	if nilFailures != 5 || capFailures != 8 || byteFailures != 4 || directionFailures != 5 {
		log.Fatalf("nonnil-empty=%d capped=%d byte-only=%d two-sided=%d", nilFailures, capFailures, byteFailures, directionFailures)
	}
	fmt.Printf("TrimLeft passed %d fixtures; mutants failed nonnil-empty=%d capped=%d byte-only=%d two-sided=%d\n",
		len(cases), nilFailures, capFailures, byteFailures, directionFailures)
}

func trimLeftBytes(input []byte, cutset string) []byte {
	for len(input) != 0 && strings.IndexByte(cutset, input[0]) >= 0 {
		input = input[1:]
	}
	if len(input) == 0 {
		return nil
	}
	return input
}

func observe(result []byte) string {
	if result == nil {
		return "nil"
	}
	row := fmt.Appendf(nil, "%d:", cap(result))
	for i, value := range result {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(value), 10)
	}
	return string(row)
}
