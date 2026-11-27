package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	cases := []struct {
		text  string
		spare int
		want  string
	}{
		{"", 0, "nil"},
		{"", 4, "nil"},
		{" \t\r\n", 3, "nil"},
		{" \u2000\u00a0\t", 2, "nil"},
		{" \tX  ", 4, "7:88"},
		{" \u00a0X\v", 3, "5:88"},
		{" \tX\u3000 \r", 4, "10:88"},
		{"\t\xffX \t", 2, "6:255,88"},
		{" \u200bX\ufeff ", 1, "9:226,128,139,88,239,187,191"},
		{"\tX\u0085", 2, "5:88"},
		{"\x85X ", 1, "4:133,88"},
		{" \ta\va \t", 2, "7:97,11,97"},
	}
	var inputs bytes.Buffer
	var emptyFailures, capFailures, asciiFailures int
	for _, c := range cases {
		row := strconv.AppendInt(nil, int64(c.spare), 10)
		for i := range len(c.text) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
		}
		fmt.Fprintln(&inputs, string(row))
		var input []byte
		if len(c.text)+c.spare != 0 {
			input = make([]byte, len(c.text), len(c.text)+c.spare)
			copy(input, c.text)
		}
		trimmed := bytes.TrimSpace(input)
		if got := observe(trimmed); got != c.want {
			log.Fatalf("TrimSpace(%q): Go=%s want=%s", c.text, got, c.want)
		}
		if trimmed == nil && observe([]byte{}) != c.want {
			emptyFailures++
		}
		if trimmed != nil && observe(trimmed[:len(trimmed):len(trimmed)]) != c.want {
			capFailures++
		}
		if observe(bytes.Trim(input, " \t\n\v\f\r")) != c.want {
			asciiFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-space", len(cases), &inputs) {
		if row != cases[i].want {
			log.Fatalf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	if emptyFailures != 4 || capFailures != 8 || asciiFailures != 4 {
		log.Fatalf("nonnil-empty failed %d, clamped capacity failed %d, ASCII-only failed %d",
			emptyFailures, capFailures, asciiFailures)
	}
	fmt.Printf("TrimSpace passed %d fixtures; nonnil-empty failed %d, clamped capacity failed %d, ASCII-only failed %d\n",
		len(cases), emptyFailures, capFailures, asciiFailures)
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
