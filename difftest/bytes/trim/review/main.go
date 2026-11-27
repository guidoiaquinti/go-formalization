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
		{"xxx", "x", 5, "nil"},
		{"baXba", "ab", 4, "7:88"},
		{"é🚀Xé🚀", "é🚀", 3, "10:88"},
		{"\xff�x\x80", "�", 2, "4:120"},
		{"\xffx\xfe", "a", 5, "8:255,120,254"},
		{"éxé", "", 2, "7:195,169,120,195,169"},
		{"éx🚀", "🚀", 4, "11:195,169,120"},
		{"\xe2\x82x\xed\xa0", "\xff", 3, "6:120"},
		{"éa", "\xc3", 2, "5:195,169,97"},
	}
	var inputs bytes.Buffer
	var emptyFailures, capFailures, byteFailures int
	for _, c := range cases {
		row := fmt.Appendf(nil, "%d %d", len(c.cutset), c.spare)
		for _, s := range []string{c.cutset, c.text} {
			for i := range len(s) {
				row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
			}
		}
		fmt.Fprintln(&inputs, string(row))
		var input []byte
		if len(c.text)+c.spare != 0 {
			input = make([]byte, len(c.text), len(c.text)+c.spare)
			copy(input, c.text)
		}
		trimmed := bytes.Trim(input, c.cutset)
		if got := observe(trimmed); got != c.want {
			log.Fatalf("Trim(%q,%q): Go=%s want=%s", c.text, c.cutset, got, c.want)
		}
		if trimmed == nil && observe([]byte{}) != c.want {
			emptyFailures++
		}
		if trimmed != nil && observe(trimmed[:len(trimmed):len(trimmed)]) != c.want {
			capFailures++
		}
		if observe(trimByteSet(input, c.cutset)) != c.want {
			byteFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim", len(cases), &inputs) {
		if row != cases[i].want {
			log.Fatalf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	if emptyFailures != 3 || capFailures != 8 || byteFailures != 3 {
		log.Fatalf("nonnil-empty failed %d, clamped cap failed %d, byte-only failed %d",
			emptyFailures, capFailures, byteFailures)
	}
	fmt.Printf("Trim passed %d fixtures; nonnil-empty failed %d, clamped cap failed %d, byte-only failed %d\n",
		len(cases), emptyFailures, capFailures, byteFailures)
}

func trimByteSet(input []byte, cutset string) []byte {
	left, right := 0, len(input)
	for right > left && strings.IndexByte(cutset, input[right-1]) >= 0 {
		right--
	}
	for left < right && strings.IndexByte(cutset, input[left]) >= 0 {
		left++
	}
	if left == right {
		return nil
	}
	return input[left:right]
}

func observe(result []byte) string {
	if result == nil {
		return "nil"
	}
	row := fmt.Appendf(nil, "%d:", cap(result))
	for i, b := range result {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(b), 10)
	}
	return string(row)
}
