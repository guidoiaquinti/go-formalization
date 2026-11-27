package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type errorFunc func() string

func (f errorFunc) Error() string { return f() }

type nilError struct{}

var nilCalls int

func (*nilError) Error() string {
	nilCalls++
	return "\x00\xff"
}

type fixture struct {
	present bool
	offset  int64
	decimal string
	mode    int
	raw     string
}

func native(c fixture, read bool) (out string) {
	// Separate callback and panic markers ensure the recovered identity is checked.
	marker := new(int)
	calls, before := 0, nilCalls
	defer func() {
		if p := recover(); p != nil {
			calls += nilCalls - before
			if p == marker {
				out = fmt.Sprintf("panic:%d", calls)
				return
			}
			if e, ok := p.(runtime.Error); ok && strings.Contains(e.Error(), "nil pointer") {
				out = fmt.Sprintf("nil:%d", calls)
				return
			}
			panic(p)
		}
	}()
	var e error
	switch c.mode {
	case 1:
		e = errorFunc(func() string { calls++; return c.raw })
	case 2:
		e = errorFunc(func() string { calls++; panic(marker) })
	case 3:
		e = (*nilError)(nil)
	}
	var got string
	if read {
		var outer *flate.ReadError
		if c.present {
			outer = &flate.ReadError{Offset: c.offset, Err: e}
		}
		got = outer.Error()
	} else {
		var outer *flate.WriteError
		if c.present {
			outer = &flate.WriteError{Offset: c.offset, Err: e}
		}
		got = outer.Error()
	}
	return fmt.Sprintf("ok:%d:%s", calls+nilCalls-before, wire(got))
}

func wire(raw string) string {
	parts := make([]string, len(raw))
	for i := range len(raw) {
		parts[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(parts, ",")
}

func main() {
	fixed := []fixture{
		{true, math.MinInt64, "-9223372036854775808", 1, "\xff\x00\x80"},
		{true, -71, "-71", 1, "\r\n💡"},
		{true, 0, "0", 1, ""},
		{true, math.MaxInt64, "9223372036854775807", 2, ""},
		{false, -1, "-1", 2, ""},
		{true, 99, "99", 0, ""},
		{true, -4, "-4", 3, "\x00\xff"},
		{false, 1, "1", 3, "\x00\xff"},
		{false, 2, "2", 1, "body"},
	}
	var input bytes.Buffer
	want := make([]string, 0, len(fixed))
	var mutations [3]int
	for _, c := range fixed {
		present := 0
		if c.present {
			present = 1
		}
		fmt.Fprint(&input, present, " ", c.offset, " ", c.mode)
		for _, b := range []byte(c.raw) {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
		corrupt := "flate: corrupt input before offset " + c.decimal
		internal := "flate: internal error: " + c.raw
		fields := []string{wire(corrupt), wire(internal)}
		if flate.CorruptInputError(c.offset).Error() != corrupt || flate.InternalError(c.raw).Error() != internal {
			log.Fatal("native value error failed literal formatting")
		}
		for i, kind := range []string{"read", "write"} {
			expected := "nil:0"
			if c.present && c.mode == 2 {
				expected = "panic:1"
			} else if c.present && c.mode != 0 {
				expected = "ok:1:" + wire("flate: "+kind+" error at offset "+c.decimal+": "+c.raw)
			}
			if got := native(c, i == 0); got != expected {
				log.Fatalf("native %s at offset %d: got %q, want %q", kind, c.offset, got, expected)
			}
			fields = append(fields, expected)
		}
		want = append(want, strings.Join(fields, "|"))
		if strconv.FormatUint(uint64(c.offset), 10) != c.decimal {
			mutations[0]++
		}
		if strings.ToValidUTF8(c.raw, "�") != c.raw {
			mutations[1]++
		}
		if c.present && c.mode == 2 {
			wrapped := errorFunc(func() string { panic(new(int)) })
			if got := fmt.Sprint(wrapped); strings.HasPrefix(got, "%!v(PANIC=") {
				mutations[2]++
			}
		}
	}
	for i, got := range check.LeanRows("godiff-flate-error-text", len(fixed), &input) {
		if got != want[i] {
			log.Fatalf("control row %d: Lean got %q want %q", i, got, want[i])
		}
	}
	if expected := [3]int{4, 3, 1}; mutations != expected {
		log.Fatalf("mutation disagreements %v, want %v", mutations, expected)
	}
	fmt.Printf("%d literal cases for all four methods; mutation disagreements %v\n", len(fixed), mutations)
}
