package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	text  string
	kind  int
	spare int
	want  string
}

func main() {
	cases := []fixture{
		{"", 0, 0, "nil|"},
		{"", 0, 5, "nil|"},
		{"a", 0, 0, "panic"},
		{"\xffa", 1, 2, "4:255,97|65533"},
		{"é\xe0", 2, 2, "nil|233,65533"},
		{"\xff�🚀\xc3", 3, 5, "10:240,159,154,128,195|65533,65533,128640"},
		{"é🚀x\xffé", 4, 4, "8:120,255,195,169|233,128640,120"},
		{"é🚀", 4, 2, "nil|233,128640"},
		{"🚀z", 4, 0, "1:122|128640,122"},
	}
	var inputs bytes.Buffer
	var nilFailures, capFailures, traceFailures int
	for _, c := range cases {
		row := fmt.Appendf(nil, "%d %d", c.kind, c.spare)
		for i := range len(c.text) {
			row = strconv.AppendUint(append(row, ' '), uint64(c.text[i]), 10)
		}
		fmt.Fprintln(&inputs, string(row))
		result, calls, panicked := run(c)
		if panicked {
			if c.want != "panic" {
				log.Fatalf("input %q kind %d unexpectedly panicked", c.text, c.kind)
			}
			continue
		}
		if got := observe(result, calls); got != c.want {
			log.Fatalf("input %q kind %d: Go=%s want=%s", c.text, c.kind, got, c.want)
		}
		if result == nil && observe([]byte{}, calls) != c.want {
			nilFailures++
		}
		if result != nil && observe(result[:len(result):len(result)], calls) != c.want {
			capFailures++
		}
		if c.kind != 0 && observe(result, []rune(c.text)) != c.want {
			traceFailures++
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-left-func", len(cases), &inputs) {
		if row != cases[i].want {
			log.Fatalf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	if nilFailures != 4 || capFailures != 3 || traceFailures != 3 {
		log.Fatalf("empty-over-nil failed %d, clamped capacity failed %d, full tracing failed %d",
			nilFailures, capFailures, traceFailures)
	}
	fmt.Printf("TrimLeftFunc passed %d fixtures; empty-over-nil failed %d, clamped capacity failed %d, full tracing failed %d\n",
		len(cases), nilFailures, capFailures, traceFailures)
}

func run(c fixture) (result []byte, calls []rune, panicked bool) {
	defer func() {
		if value := recover(); value != nil {
			err, ok := value.(runtime.Error)
			if !ok || err.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			panicked = true
		}
	}()
	var input []byte
	if len(c.text)+c.spare > 0 {
		input = make([]byte, len(c.text), len(c.text)+c.spare)
		copy(input, c.text)
	}
	var f func(rune) bool
	if c.kind != 0 {
		f = func(r rune) bool {
			calls = append(calls, r)
			switch c.kind {
			case 1:
				return false
			case 2:
				return true
			case 3:
				return r == utf8.RuneError
			case 4:
				return r == 'é' || r == '🚀'
			default:
				panic("invalid predicate")
			}
		}
	}
	result = bytes.TrimLeftFunc(input, f)
	return
}

func observe(result []byte, calls []rune) string {
	var row []byte
	if result == nil {
		row = []byte("nil")
	} else {
		row = fmt.Appendf(row, "%d:", cap(result))
	}
	for i, b := range result {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendUint(row, uint64(b), 10)
	}
	row = append(row, '|')
	for i, r := range calls {
		if i != 0 {
			row = append(row, ',')
		}
		row = strconv.AppendInt(row, int64(r), 10)
	}
	return string(row)
}
