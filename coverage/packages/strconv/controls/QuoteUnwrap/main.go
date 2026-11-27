package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

type pointerError struct{}

func (*pointerError) Error() string { panic("unexpected Error call") }

func unwrap(e *strconv.NumError) (got error, failure any) {
	defer func() { failure = recover() }()
	return e.Unwrap(), nil
}

func main() {
	cases := []struct {
		s    string
		want bool
	}{
		{"", true},
		{"a\tb", true},
		{"a\rb", false},
		{"a\nb", false},
		{"a`b", false},
		{"a\ufeffb", false},
		{"\x7f", false},
		{"\u0085", true},
		{"\ufffd", true},
		{"\xef\xbf", false},
		{"\xed\xa0\x80", false},
		{"\xff", false},
	}
	var replaceMalformed, rejectRuneError int
	for _, c := range cases {
		if got := strconv.CanBackquote(c.s); got != c.want {
			log.Fatalf("CanBackquote(%q)=%t want %t", c.s, got, c.want)
		}
		if strconv.CanBackquote(strings.ToValidUTF8(c.s, "\ufffd")) != c.want {
			replaceMalformed++
		}
		if (strconv.CanBackquote(c.s) && !strings.ContainsRune(c.s, utf8.RuneError)) != c.want {
			rejectRuneError++
		}
	}
	var typedNil *pointerError
	sameA, sameB := errors.New("same"), errors.New("same")
	values := []error{nil, typedNil, sameA, sameB, strconv.ErrSyntax}
	var collapse, canonicalize int
	for _, want := range values {
		got, p := unwrap(&strconv.NumError{Func: "x", Num: "\xff", Err: want})
		if p != nil || got != want {
			log.Fatalf("Unwrap did not preserve error of dynamic type %T", want)
		}
		mutant := got
		if ptr, ok := got.(*pointerError); ok && ptr == nil {
			mutant = nil
		}
		if mutant != want {
			collapse++
		}
		mutant = got
		if got == sameB {
			mutant = sameA
		}
		if mutant != want {
			canonicalize++
		}
	}
	_, failure := unwrap(nil)
	if _, isRuntime := failure.(runtime.Error); !isRuntime {
		log.Fatalf("nil receiver recovered %T, want runtime.Error", failure)
	}
	if replaceMalformed != 3 || rejectRuneError != 1 || collapse != 1 || canonicalize != 1 {
		log.Fatalf("mutation failures malformed=%d RuneError=%d typed-nil=%d canonical=%d",
			replaceMalformed, rejectRuneError, collapse, canonicalize)
	}
	fmt.Printf("CanBackquote: %d fixtures pass, mutations fail %d/%d; Unwrap: %d fixtures pass, mutations fail %d/%d\n",
		len(cases), replaceMalformed, rejectRuneError, len(values)+1, collapse, canonicalize)
}
