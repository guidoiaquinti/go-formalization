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
		{"", "", -1, "nil"},
		{"", "é", -1, "nil"},
		{"", "", 0, "0:"},
		{"", "a", 4, "4:"},
		{"xxx", "x", 5, "8:"},
		{"baXba", "ab", 4, "9:98,97,88"},
		{"é🚀Xé🚀", "é🚀", 3, "16:195,169,240,159,154,128,88"},
		{"\xff�x\x80", "�", 2, "8:255,239,191,189,120"},
		{"\xffx\xfe", "a", 5, "8:255,120,254"},
		{"éxé", "", 2, "7:195,169,120,195,169"},
		{"éx🚀", "🚀", 4, "11:195,169,120"},
		{"\xe2\x82x\xed\xa0", "\xff", 3, "8:226,130,120"},
		{"aé", "\xa9", 2, "5:97,195,169"},
	}
	var inputs bytes.Buffer
	var emptyFailures, capFailures, byteFailures int
	for _, c := range cases {
		nilInput := 0
		if c.spare < 0 {
			nilInput = 1
		}
		row := fmt.Appendf(nil, "%d %d %d", len(c.cutset), max(0, c.spare), nilInput)
		for _, s := range []string{c.cutset, c.text} {
			for i := range len(s) {
				row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
			}
		}
		fmt.Fprintln(&inputs, string(row))
		var input []byte
		if c.spare >= 0 {
			input = make([]byte, len(c.text), len(c.text)+c.spare)
			copy(input, c.text)
		}
		trimmed := bytes.TrimRight(input, c.cutset)
		if got := observe(trimmed); got != c.want {
			log.Fatalf("TrimRight(%q,%q): Go=%s want=%s", c.text, c.cutset, got, c.want)
		}
		if trimmed != nil && len(trimmed) == 0 && observe(nil) != c.want {
			emptyFailures++
		}
		if trimmed != nil && observe(trimmed[:len(trimmed):len(trimmed)]) != c.want {
			capFailures++
		}
		if observe(trimByteSet(input, c.cutset)) != c.want {
			byteFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-right", len(cases), &inputs) {
		if row != cases[i].want {
			log.Fatalf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	if emptyFailures != 3 || capFailures != 10 || byteFailures != 3 {
		log.Fatalf("nil-over-empty failed %d, clamped cap failed %d, byte-only failed %d",
			emptyFailures, capFailures, byteFailures)
	}
	fmt.Printf("TrimRight passed %d fixtures; nil-over-empty failed %d, clamped cap failed %d, byte-only failed %d\n",
		len(cases), emptyFailures, capFailures, byteFailures)
}

func trimByteSet(input []byte, cutset string) []byte {
	end := len(input)
	for end > 0 && strings.IndexByte(cutset, input[end-1]) >= 0 {
		end--
	}
	return input[:end]
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
