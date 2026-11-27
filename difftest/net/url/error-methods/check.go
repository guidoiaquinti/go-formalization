package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	kind                int
	present, flag       bool
	op, raw, inner, bad string
}

type nested struct {
	text  string
	flag  bool
	calls [2]int
}

type timeoutOnly struct{ *nested }

func (e timeoutOnly) Error() string { return e.text }
func (e timeoutOnly) Timeout() bool { e.calls[0]++; return e.flag }

type temporaryOnly struct{ *nested }

func (e temporaryOnly) Error() string   { return e.text }
func (e temporaryOnly) Temporary() bool { e.calls[1]++; return e.flag }

type formattingPanic struct{}

func (formattingPanic) Error() string { panic("inner message") }

type both struct{ *nested }

var callbackToken = errors.New("callback marker")

func (e *both) Error() string { return e.text }
func (e *both) Timeout() bool {
	if e == nil {
		panic(callbackToken)
	}
	e.calls[0]++
	return e.flag
}
func (e *both) Temporary() bool {
	if e == nil {
		panic(callbackToken)
	}
	e.calls[1]++
	return !e.flag
}

func main() {
	var fixed []fixture
	for _, raw := range []string{"", "quote\"\\", "\x00\a\b\f\n\r\t\v\x7f", "\xff\xc2A", "é日😀", "\u00a0\u200d\U0010ffff", "\ufffd"} {
		for kind := range 7 {
			for _, flag := range []bool{false, true} {
				fixed = append(fixed, fixture{kind, true, flag, "\xffop", raw, "bad\x00\xff", raw})
			}
		}
		fixed = append(fixed, fixture{raw: raw, bad: raw})
	}
	for raw := range 256 {
		data := string([]byte{byte(raw)})
		fixed = append(fixed, fixture{4, true, true, "op", data, "bad", data})
	}
	for r := rune(1); r <= utf8.MaxRune; r++ {
		if strconv.IsPrint(r) == strconv.IsPrint(r-1) {
			continue
		}
		for _, edge := range []rune{r - 1, r} {
			raw := string(edge)
			fixed = append(fixed, fixture{4, true, false, "op", raw, "bad", raw})
		}
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			rows := make([]fixture, 5000)
			for i := range rows {
				rows[i] = fixture{rng.Intn(7), i%17 != 0, rng.Intn(2) == 0, sample(rng), sample(rng), sample(rng), sample(rng)}
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), rows)
		}
	}
}

func sample(rng *rand.Rand) string {
	if rng.Intn(3) == 0 {
		var b strings.Builder
		for range rng.Intn(10) {
			r := rune(rng.Intn(utf8.MaxRune + 1))
			b.WriteRune(r)
		}
		return b.String()
	}
	b := make([]byte, rng.Intn(28))
	rng.Read(b)
	return string(b)
}

func compare(label string, rows []fixture) {
	var input bytes.Buffer
	for _, f := range rows {
		timeout, temporary, inner := expect(f)
		present := 0
		if f.present {
			present = 1
		}
		fmt.Fprint(&input, present, " ", f.kind, " ", timeout, " ", temporary)
		writeBytes(&input, f.op)
		writeBytes(&input, f.raw)
		writeBytes(&input, inner)
		writeBytes(&input, f.bad)
		fmt.Fprintln(&input)
	}
	for i, lean := range check.LeanRows("godiff-net-url-error-methods", len(rows), &input) {
		f := rows[i]
		want := oracle(f)
		got := native(f)
		if lean != want || got != want {
			log.Fatalf("%s row %d: Lean %q Go %q want %q fixture %+v", label, i, lean, got, want, f)
		}
	}
	fmt.Printf("URL error methods %s: %d pinned native and Lean rows\n", label, len(rows))
}

func writeBytes(b *bytes.Buffer, s string) {
	fmt.Fprint(b, " ", len(s))
	for _, c := range []byte(s) {
		fmt.Fprint(b, " ", c)
	}
}

func encoded(s string) string {
	var buf strings.Builder
	for i, c := range []byte(s) {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.Itoa(int(c)))
	}
	return buf.String()
}

func expect(f fixture) (timeout, temporary int, inner string) {
	inner = f.inner
	code := 1
	if f.flag {
		code = 2
	}
	switch f.kind {
	case 0:
		inner = "%!s(<nil>)"
	case 2:
		timeout = code
	case 3:
		temporary = code
	case 4:
		timeout, temporary = code, 3-code
	case 5:
		timeout, temporary, inner = 3, 3, "<nil>"
	case 6:
		inner = "%!s(PANIC=Error method: inner message)"
	}
	return
}

func oracle(f fixture) string {
	t, p, inner := expect(f)
	result := []string{"nilpanic", "nilpanic", "nilpanic", "nilpanic"}
	values := []string{"false", "false", "true", "callback77"}
	if f.present {
		result = []string{encoded(f.op + " " + strconv.Quote(f.raw) + ": " + inner), strconv.Itoa(f.kind), values[t], values[p]}
	}
	return strings.Join(append(result, encoded("invalid URL escape "+strconv.Quote(f.bad)),
		encoded("invalid character "+strconv.Quote(f.bad)+" in host name")), "|")
}

func native(f fixture) string {
	state := &nested{text: f.inner, flag: f.flag}
	var inner error
	switch f.kind {
	case 1:
		inner = errors.New(f.inner)
	case 2:
		inner = timeoutOnly{state}
	case 3:
		inner = temporaryOnly{state}
	case 4:
		inner = &both{state}
	case 5:
		inner = (*both)(nil)
	case 6:
		inner = formattingPanic{}
	}
	var wrapped *url.Error
	if f.present {
		wrapped = &url.Error{Op: f.op, URL: f.raw, Err: inner}
	}
	result := []string{
		invoke(func() string { return encoded(wrapped.Error()) }),
		invoke(func() string {
			if v := wrapped.Unwrap(); v != inner {
				return "wrong identity"
			}
			return strconv.Itoa(f.kind)
		}),
		invoke(func() string { return strconv.FormatBool(wrapped.Timeout()) }),
		invoke(func() string { return strconv.FormatBool(wrapped.Temporary()) }),
		encoded(url.EscapeError(f.bad).Error()),
		encoded(url.InvalidHostError(f.bad).Error()),
	}
	want := [2]int{}
	if f.present {
		if f.kind == 2 || f.kind == 4 {
			want[0] = 1
		}
		if f.kind == 3 || f.kind == 4 {
			want[1] = 1
		}
	}
	if state.calls != want {
		log.Fatalf("callback counts %v want %v", state.calls, want)
	}
	return strings.Join(result, "|")
}

func invoke(call func() string) (result string) {
	defer func() {
		switch p := recover(); p {
		case nil:
		case callbackToken:
			result = "callback77"
		default:
			if _, ok := p.(runtime.Error); !ok {
				panic(p)
			}
			result = "nilpanic"
		}
	}()
	return call()
}
