package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

var callbackToken = new(int)

type timeoutOnly struct{ calls int }

func (*timeoutOnly) Error() string { return "cause" }
func (e *timeoutOnly) Timeout() bool {
	if e == nil {
		panic(callbackToken)
	}
	e.calls++
	return true
}

type temporaryOnly bool

func (temporaryOnly) Error() string     { return "temporary" }
func (e temporaryOnly) Temporary() bool { return bool(e) }

func main() {
	bad := "\xff\ufffd\u00a0\u2028"
	quoted := "\"\\xff�\\u00a0\\u2028\""
	if got := url.EscapeError(bad).Error(); got != "invalid URL escape "+quoted {
		log.Fatalf("escape diagnostic: %q", got)
	}
	if got := url.InvalidHostError(bad).Error(); got != "invalid character "+quoted+" in host name" {
		log.Fatalf("host diagnostic: %q", got)
	}

	inner := new(timeoutOnly)
	e := &url.Error{Op: "open\xff", URL: bad, Err: inner}
	const operation = "open\xff"
	if got := e.Error(); got != operation+" "+quoted+": cause" {
		log.Fatalf("wrapper diagnostic: %q", got)
	}
	if e.Unwrap() != inner || !e.Timeout() || e.Temporary() || inner.calls != 1 {
		log.Fatal("wrapper identity or capability mismatch")
	}
	wrapped := fmt.Errorf("wrapped: %w", inner)
	e.Err = wrapped
	if e.Unwrap() != wrapped || !errors.Is(e, inner) || e.Timeout() || e.Temporary() || inner.calls != 1 {
		log.Fatal("wrapper traversed nested capability or returned a different error")
	}
	var typedNil *timeoutOnly
	e.Err = typedNil
	if e.Unwrap() == nil || e.Unwrap() != typedNil || e.Temporary() {
		log.Fatal("typed nil or capability identity changed")
	}
	func() {
		defer func() {
			if got := recover(); got != callbackToken {
				log.Fatalf("typed nil capability panic got %v, want identical token", got)
			}
		}()
		e.Timeout()
	}()
	e.Err = temporaryOnly(true)
	if e.Timeout() || !e.Temporary() || e.Unwrap() != temporaryOnly(true) {
		log.Fatal("temporary-only capability was not selected independently")
	}
	var nilWrapper *url.Error
	for _, f := range []func(){
		func() { nilWrapper.Error() }, func() { nilWrapper.Unwrap() },
		func() { nilWrapper.Timeout() }, func() { nilWrapper.Temporary() },
	} {
		func() {
			defer func() {
				if recover() == nil {
					log.Fatal("nil wrapper failed to panic")
				}
			}()
			f()
		}()
	}

	escape := wire("invalid URL escape " + quoted)
	host := wire("invalid character " + quoted + " in host name")
	cases := []struct {
		present, id, timeout, temporary int
		inner, want                     string
	}{
		{1, 11, 2, 0, "cause", strings.Join([]string{wire(operation + " " + quoted + ": cause"), "11", "true", "false", escape, host}, "|")},
		{1, 12, 0, 0, "wrapped: cause", strings.Join([]string{wire(operation + " " + quoted + ": wrapped: cause"), "12", "false", "false", escape, host}, "|")},
		{1, 13, 3, 0, "cause", strings.Join([]string{wire(operation + " " + quoted + ": cause"), "13", "callback77", "false", escape, host}, "|")},
		{0, 0, 0, 0, "", strings.Join([]string{"nilpanic", "nilpanic", "nilpanic", "nilpanic", escape, host}, "|")},
		{1, 14, 0, 2, "temporary", strings.Join([]string{wire(operation + " " + quoted + ": temporary"), "14", "false", "true", escape, host}, "|")},
	}
	var input bytes.Buffer
	for _, c := range cases {
		fmt.Fprint(&input, c.present, " ", c.id, " ", c.timeout, " ", c.temporary)
		for _, s := range []string{operation, bad, c.inner, bad} {
			fmt.Fprint(&input, " ", len(s))
			for _, b := range []byte(s) {
				fmt.Fprint(&input, " ", b)
			}
		}
		fmt.Fprintln(&input)
	}
	rows := check.LeanRows("godiff-net-url-error-methods", len(cases), &input)
	for i, row := range rows {
		if row != cases[i].want {
			log.Fatalf("row %d = %q, want %q", i, row, cases[i].want)
		}
	}
	for _, m := range []struct {
		row, col int
		value    string
	}{
		{0, 0, wire(operation + " " + strings.Replace(quoted, "\\xff", "\ufffd", 1) + ": cause")}, // replace malformed byte
		{1, 1, "11"},    // unwrap beyond the immediate field
		{1, 2, "true"},  // search for Timeout through errors.Is/Unwrap
		{0, 3, "true"},  // infer Temporary from Timeout
		{2, 2, "false"}, // suppress typed nil callback panic
		{3, 3, "false"}, // default a nil receiver to false
		{0, 4, wire("invalid URL escape " + bad)},
		{0, 5, wire("invalid character " + quoted)},
		{4, 3, "false"}, // ignore a present Temporary method
	} {
		if strings.Split(rows[m.row], "|")[m.col] == m.value {
			log.Fatalf("row %d column %d accepted mutant %q", m.row, m.col, m.value)
		}
	}
	log.Print("five independent pinned native and Lean profiles passed and rejected nine wrong observations")
}

func wire(value string) string {
	fields := make([]string, len(value))
	for i, b := range []byte(value) {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}
