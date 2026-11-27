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
	text     string
	kind     int
	spare    int
	nilInput bool
	want     string
}

func main() {
	cases := []fixture{
		{"", 0, 0, true, "nil|"},
		{"", 0, 0, false, "0:|"},
		{"", 0, 3, false, "3:|"},
		{"a", 0, 0, false, "panic"},
		{"é\xe0", 2, 3, false, "6:|65533,233"},
		{"x\xffé", 1, 2, false, "6:120,255,195,169|233"},
		{"x🚀\xff�", 3, 3, false, "12:120,240,159,154,128|65533,65533,128640"},
		{"é\x80x🚀é", 4, 2, false, "12:195,169,128,120|233,128640,120"},
		{"🚀é", 4, 0, false, "6:|233,128640"},
		{"z\x80é", 4, 0, false, "4:122,128|233,65533"},
	}
	var inputs bytes.Buffer
	var nilFailures, capFailures, widthFailures int
	for _, c := range cases {
		nilInput := 0
		if c.nilInput {
			nilInput = 1
		}
		row := fmt.Appendf(nil, "%d %d %d", c.kind, c.spare, nilInput)
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
		if result != nil && len(result) == 0 && observe(nil, calls) != c.want {
			nilFailures++
		}
		if result != nil && observe(result[:len(result):len(result)], calls) != c.want {
			capFailures++
		}
		if len(result) > 0 {
			_, width := utf8.DecodeLastRune(result)
			end := len(result) - width + 1
			if observe(result[:end], calls) != c.want {
				widthFailures++
			}
		}
	}
	for i, row := range check.LeanRows("godiff-bytes-trim-right-func", len(cases), &inputs) {
		if row != cases[i].want {
			log.Fatalf("case %d: Lean=%s want=%s", i, row, cases[i].want)
		}
	}
	if nilFailures != 4 || capFailures != 7 || widthFailures != 2 {
		log.Fatalf("nil-over-empty failed %d, clamped capacity failed %d, one-byte width failed %d",
			nilFailures, capFailures, widthFailures)
	}
	fmt.Printf("TrimRightFunc passed %d fixtures; nil-over-empty failed %d, clamped capacity failed %d, one-byte width failed %d\n",
		len(cases), nilFailures, capFailures, widthFailures)
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
	if !c.nilInput {
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
	result = bytes.TrimRightFunc(input, f)
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
