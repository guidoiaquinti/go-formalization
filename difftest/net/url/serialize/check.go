package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type fixture struct {
	present     bool
	flags, kind int
	parts       [10]string
	backing     []byte
	length      int
}

func main() {
	var fixed []fixture
	for _, parts := range [][10]string{
		{},
		{"https", "", "", "", "example.test", "recipes/pasta", "", "day=2", "week 1", ""},
		{"https", "", "ingredient", "olive oil", "example.test", "/recipes/pasta", "", "", "week 1", "week%201"},
		{"https", "", "", "", "", "/folder", "", "", "", ""},
		{"mailto", "recipes@example.test", "", "", "", "", "", "", "heading", ""},
		{"example", "//example.test/inbox", "", "", "", "", "", "q=1", "", ""},
		{"", "", "", "", "", "a:b/c", "", "", "", ""},
		{"", "", "", "", "", "*", "%2a", "", "", ""},
		{"https", "", "", "", "café.example", "/foo/bar", "/foo%2fbar", "", "a/b", "a%2fb"},
		{"", "", "", "", "", "/", "", "", "x", ""},
	} {
		for flags := range 4 {
			for kind := range 3 {
				for _, cap := range []int{0, 1, 10, 32, 128} {
					backing := bytes.Repeat([]byte{0x80}, cap)
					fixed = append(fixed, fixture{true, flags, kind, parts, backing, min(cap, 3)})
				}
			}
		}
	}
	for _, cap := range []int{0, 1, 8} {
		fixed = append(fixed, fixture{backing: bytes.Repeat([]byte{0xfe}, cap), length: min(cap, 2)})
	}
	compare("fixed", fixed)
	for _, seed := range []int64{17, 2718, 65537} {
		r := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]fixture, 5000)
			for i := range cases {
				u := [10]string{
					choose(r, "", "https", "example"), choose(r, "", "", "", "contact", "//example.test/inbox"), choose(r, "", "reader", "S Alice"),
					choose(r, "", "red apples", "word"), choose(r, "", "example.test", "localhost:8080", "[::1]:443", "café.example"),
					choose(r, "", "/", "/one", "a:b", "sub/folder", "*", "/tail space", "/one/café"), "",
					choose(r, "", "a=1", "key=x+y", "letter=%C3%A9"), choose(r, "", "later", "αβ", "first page", "a/b"), "",
				}
				if u[5] == "/one" && r.Intn(3) == 0 {
					u[6] = "/%6fne"
				}
				if u[8] == "a/b" && r.Intn(3) == 0 {
					u[9] = "a%2fb"
				}
				backing := make([]byte, r.Intn(85))
				_, err := r.Read(backing)
				if err != nil {
					log.Fatal(err)
				}
				cases[i] = fixture{r.Intn(20) != 0, r.Intn(4), r.Intn(3), u, backing, r.Intn(len(backing) + 1)}
			}
			compare(fmt.Sprintf("seed %d block %d", seed, block), cases)
		}
	}
}

func choose(r *rand.Rand, values ...string) string { return values[r.Intn(len(values))] }

func compare(label string, cases []fixture) {
	var input bytes.Buffer
	for _, c := range cases {
		present := 0
		if c.present {
			present = 1
		}
		fmt.Fprint(&input, present, " ", c.flags, " ", c.kind, " ", c.length)
		for _, part := range c.parts {
			fmt.Fprint(&input, " ", len(part))
		}
		for _, part := range c.parts {
			for i := range len(part) {
				fmt.Fprint(&input, " ", part[i])
			}
		}
		for _, b := range c.backing {
			fmt.Fprint(&input, " ", b)
		}
		fmt.Fprintln(&input)
	}
	for i, lean := range check.LeanRows("godiff-net-url-serialize", len(cases), &input) {
		c := cases[i]
		got, want := native(c), oracle(c)
		if got != want || got != lean {
			log.Fatalf("%s row %d: fixture %+v; Go %q Lean %q oracle %q", label, i, c, got, lean, want)
		}
	}
	fmt.Printf("%s: %d pinned Go and Lean serialization rows agree\n", label, len(cases))
}

func (c fixture) url() *url.URL {
	if !c.present {
		return nil
	}
	p := c.parts
	var user *url.Userinfo
	if c.kind == 1 {
		user = url.User(p[2])
	} else if c.kind == 2 {
		user = url.UserPassword(p[2], p[3])
	}
	return &url.URL{Scheme: p[0], Opaque: p[1], User: user, Host: p[4], Path: p[5], RawPath: p[6], RawQuery: p[7], Fragment: p[8], RawFragment: p[9],
		OmitHost: c.flags&1 != 0, ForceQuery: c.flags&2 != 0}
}

func wire(raw string) string {
	cells := make([]string, len(raw))
	for i := range len(raw) {
		cells[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(cells, ",")
}

func observe(call func() string) (result string) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(runtime.Error); ok && strings.Contains(err.Error(), "nil pointer") {
				result = "nil"
			} else {
				panic(p)
			}
		}
	}()
	return call()
}

func native(c fixture) string {
	u := c.url()
	var snapshot url.URL
	if u != nil {
		snapshot = *u
	}
	stringified := observe(func() string { return wire(u.String()) })
	request := observe(func() string { return wire(u.RequestURI()) })
	redacted := wire(u.Redacted())
	marshaled := observe(func() string {
		data, err := u.MarshalBinary()
		if err != nil || len(data) == 0 && data != nil {
			log.Fatal("unexpected marshal state")
		}
		return wire(string(data))
	})
	old := make([]byte, len(c.backing))
	copy(old, c.backing)
	appended := observe(func() string {
		data, err := u.AppendBinary(old[:c.length:len(old)])
		if err != nil {
			log.Fatal(err)
		}
		reuse := cap(data) == 0 && len(old) == 0
		if cap(data) > 0 && len(old) > 0 {
			reuse = &data[:cap(data)][0] == &old[0]
		}
		return fmt.Sprintf("%s:%s:%t", wire(string(data)), wire(string(old)), reuse)
	})
	if appended == "nil" {
		appended += ":" + wire(string(old))
	}
	if u != nil && *u != snapshot {
		log.Fatal("URL rendering changed receiver")
	}
	return strings.Join([]string{stringified, request, redacted, marshaled, appended}, "|")
}

func render(u *url.URL) string {
	var prefix string
	if u.Scheme != "" {
		prefix = u.Scheme + ":"
	}
	if u.Opaque != "" {
		prefix += u.Opaque
	} else {
		path := u.EscapedPath()
		authority := u.Host != "" || u.User != nil
		if authority || u.Scheme != "" && u.Path != "" && !u.OmitHost {
			prefix += "//"
		}
		if u.User != nil {
			prefix += u.User.String() + "@"
		}
		for i := range len(u.Host) {
			b := u.Host[i]
			if b >= 128 {
				prefix += fmt.Sprintf("%%%02X", b)
			} else {
				prefix += string(b)
			}
		}
		switch {
		case u.OmitHost && !authority && strings.HasPrefix(path, "//"):
			prefix += "%2F"
			path = path[1:]
		case u.Host != "" && path != "" && !strings.HasPrefix(path, "/"):
			prefix += "/"
		case prefix == "":
			first, _, _ := strings.Cut(path, "/")
			if strings.Contains(first, ":") {
				prefix += "./"
			}
		}
		prefix += path
	}
	if u.ForceQuery || u.RawQuery != "" {
		prefix += "?" + u.RawQuery
	}
	if u.Fragment != "" {
		prefix += "#" + u.EscapedFragment()
	}
	return prefix
}

func oracle(c fixture) string {
	u := c.url()
	if u == nil {
		return "nil|nil||nil|nil:" + wire(string(c.backing))
	}
	text := render(u)
	request := u.Opaque
	if request == "" {
		request = u.EscapedPath()
		if request == "" {
			request = "/"
		}
	} else if strings.HasPrefix(request, "//") {
		request = u.Scheme + ":" + request
	}
	if u.ForceQuery || u.RawQuery != "" {
		request += "?" + u.RawQuery
	}
	copyURL := *u
	if _, set := u.User.Password(); set {
		copyURL.User = url.UserPassword(u.User.Username(), "xxxxx")
	}
	old := slices.Clone(c.backing)
	result := string(old[:c.length]) + text
	reuse := len(result) <= len(old)
	if reuse {
		copy(old, result)
	}
	return fmt.Sprintf("%s|%s|%s|%s|%s:%s:%t", wire(text), wire(request), wire(render(&copyURL)), wire(text), wire(result), wire(string(old)), reuse)
}
