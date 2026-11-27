package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	present                 bool
	flags, kind             int
	parts                   [10]string
	backing                 []byte
	length                  int
	text, request, redacted string
}

func main() {
	cases := []fixture{
		{true, 1, 0, [10]string{"https", "", "", "", "", "//notes/today", "", "v=1"}, bytes.Repeat([]byte{'@'}, 72), 1,
			"https:%2F/notes/today?v=1", "//notes/today?v=1", "https:%2F/notes/today?v=1"},
		{true, 2, 2, [10]string{"https", "", "a", "b", "example.test", "/note", "/%6eote", "", "x/y", "x%2fy"}, []byte{1, 2}, 1,
			"https://a:b@example.test/%6eote?#x%2fy", "/%6eote?", "https://a:xxxxx@example.test/%6eote?#x%2fy"},
		{true, 2, 0, [10]string{"mail", "//example.test/msg", "", "", "", "", "", "", "x"}, bytes.Repeat([]byte{0x80}, 50), 2,
			"mail://example.test/msg?#x", "mail://example.test/msg?", "mail://example.test/msg?#x"},
		{true, 0, 0, [10]string{}, nil, 0, "", "/", ""},
		{false, 0, 0, [10]string{}, []byte{0xff}, 0, "", "", ""},
	}
	var input bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.present {
			present = 1
		}
		fmt.Fprint(&input, present, " ", c.flags, " ", c.kind, " ", c.length)
		for _, p := range c.parts {
			fmt.Fprint(&input, " ", len(p))
		}
		for _, p := range c.parts {
			for _, b := range []byte(p) {
				fmt.Fprint(&input, " ", b)
			}
		}
		for _, b := range c.backing {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	var rejected [5]int
	for i, got := range check.LeanRows("godiff-net-url-serialize", len(cases), &input) {
		c := cases[i]
		want := reference(c)
		if native := native(c); got != want || native != want {
			log.Fatalf("row %d got Lean %q, Go %q, want %q", i, got, native, want)
		}
		parts := strings.Split(got, "|")
		wrong := [5]string{parts[0], parts[1], parts[2], parts[3], parts[4]}
		if c.flags&1 != 0 {
			wrong[0] = wire("https://notes/today?v=1")
		}
		if c.parts[8] != "" {
			wrong[1] = wire(c.request + "#" + c.parts[8])
		}
		if c.kind == 2 {
			wrong[2] = wire(c.text)
		}
		if c.present {
			wrong[3] = wire(c.request)
		}
		if c.present && c.length+len(c.text) <= len(c.backing) {
			wrong[4] = appendWire(c, c.text, c.backing, false)
		}
		for j, w := range wrong {
			if w != parts[j] {
				rejected[j]++
			}
		}
	}
	if want := [5]int{1, 2, 1, 4, 3}; rejected != want {
		log.Fatalf("mutant counts %v want %v", rejected, want)
	}
	fmt.Printf("URL serialization review: %d pinned Go and Lean cases, mutation counts %v\n", len(cases), rejected)
}

func wire(s string) string {
	parts := make([]string, len(s))
	for i, b := range []byte(s) {
		parts[i] = strconv.Itoa(int(b))
	}
	return strings.Join(parts, ",")
}

func appendWire(c fixture, text string, old []byte, reused bool) string {
	return fmt.Sprintf("%s:%s:%t", wire(string(c.backing[:c.length])+text), wire(string(old)), reused)
}

func reference(c fixture) string {
	if !c.present {
		return "nil|nil||nil|nil:" + wire(string(c.backing))
	}
	old := bytes.Clone(c.backing)
	reuse := len(c.text)+c.length <= len(old)
	if reuse {
		copy(old[c.length:], c.text)
	}
	return strings.Join([]string{wire(c.text), wire(c.request), wire(c.redacted), wire(c.text), appendWire(c, c.text, old, reuse)}, "|")
}

func observe(f func() string) (result string) {
	defer func() {
		if x := recover(); x != nil {
			if _, ok := x.(runtime.Error); !ok {
				panic(x)
			}
			result = "nil"
		}
	}()
	return f()
}

func native(c fixture) string {
	var u *url.URL
	if c.present {
		p := c.parts
		u = &url.URL{Scheme: p[0], Opaque: p[1], Host: p[4], Path: p[5], RawPath: p[6], RawQuery: p[7],
			Fragment: p[8], RawFragment: p[9], OmitHost: c.flags&1 != 0, ForceQuery: c.flags&2 != 0}
		if c.kind == 2 {
			u.User = url.UserPassword(p[2], p[3])
		}
	}
	marshal := observe(func() string {
		b, err := u.MarshalBinary()
		if err != nil || len(b) == 0 && b != nil {
			log.Fatal("unexpected marshal result")
		}
		return wire(string(b))
	})
	old := bytes.Clone(c.backing)
	appended := observe(func() string {
		b, err := u.AppendBinary(old[:c.length])
		if err != nil {
			log.Fatal(err)
		}
		reused := cap(b) == 0 && len(old) == 0
		if len(old) > 0 && cap(b) > 0 {
			reused = &b[:cap(b)][0] == &old[0]
		}
		return fmt.Sprintf("%s:%s:%t", wire(string(b)), wire(string(old)), reused)
	})
	if appended == "nil" {
		appended += ":" + wire(string(old))
	}
	return strings.Join([]string{observe(func() string { return wire(u.String()) }),
		observe(func() string { return wire(u.RequestURI()) }), wire(u.Redacted()), marshal, appended}, "|")
}
