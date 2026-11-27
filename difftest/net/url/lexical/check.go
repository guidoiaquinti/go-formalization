package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Present bool
	Raw     []byte
}
type fixture struct {
	input
	Want string
}

var portOnly = regexp.MustCompile("^:[0-9]*$")
var schemePrefix = regexp.MustCompile("(?s)^([a-zA-Z][a-zA-Z0-9+.-]*):(.*)$")

func wire(raw string) string {
	fields := make([]string, len(raw))
	for i := range fields {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return fmt.Sprintf("%d:%s", len(raw), strings.Join(fields, ","))
}

func expected(c input) string {
	raw := string(c.Raw)
	host, port := raw, ""
	for i := len(raw) - 1; i >= 0; i-- {
		if raw[i] != ':' {
			continue
		}
		if portOnly.MatchString(raw[i:]) {
			host, port = raw[:i], raw[i+1:]
		}
		break
	}
	if len(host) >= 2 && host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1 : len(host)-1]
	}
	scheme := wire("") + "/" + wire(raw)
	if strings.HasPrefix(raw, ":") {
		scheme = "missing"
	} else if match := schemePrefix.FindStringSubmatch(raw); match != nil {
		scheme = wire(match[1]) + "/" + wire(match[2])
	}
	h, p, a := "panic", "panic", "panic"
	if c.Present {
		h, p, a = wire(host), wire(port), strconv.FormatBool(raw != "")
	}
	control := strings.IndexAny(raw, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x7f") >= 0
	return strings.Join([]string{strconv.FormatBool(raw == "" || portOnly.MatchString(raw)),
		wire(host) + "/" + wire(port), h, p, a, scheme, strconv.FormatBool(control)}, "|")
}

func batch(label string, cases []input, encoder *json.Encoder) {
	const shard = 1024
	for first := 0; first < len(cases); first += shard {
		current := cases[first:min(first+shard, len(cases))]
		var protocol bytes.Buffer
		rows := make([]fixture, len(current))
		for i, c := range current {
			tag := 0
			if c.Present {
				tag = 1
			}
			fmt.Fprint(&protocol, tag)
			for _, b := range c.Raw {
				fmt.Fprint(&protocol, " ", b)
			}
			fmt.Fprintln(&protocol)
			rows[i] = fixture{c, expected(c)}
		}
		for i, got := range check.LeanRows("godiff-net-url-lexical", len(current), &protocol) {
			if got != rows[i].Want {
				log.Fatalf("%s row %d %+v got %q want %q", label, first+i, current[i], got, rows[i].Want)
			}
		}
		for _, row := range rows {
			must(encoder.Encode(row))
		}
	}
	fmt.Printf("URL lexical %s: %d Lean rows agree.\n", label, len(cases))
}

func main() {
	var edges []input
	for i := range 1 << 16 {
		raw := []byte{byte(i >> 8), byte(i)}
		edges = append(edges, input{i%3 != 0, raw})
	}
	for value := range 256 {
		b := byte(value)
		edges = append(edges, input{true, []byte{b}}, input{true, []byte{'A', b, ':', 'x'}},
			input{true, []byte{'[', b, ']', ':', '1'}}, input{true, []byte{':', b}})
	}
	for _, raw := range []string{"", "[]", "[:]:", "[::1]", "[::1]:080", "[::1]:x", "[[::1]]", "x:123456789123456789",
		"a:b:", "a:b:9", "x:%31", "x:１２", "x:+1", "x:-1", "HTTP://x", "a+.-1:\xff", ":", "\xff",
		"\x00", "\x7f", "\x80", "a\x00:", "a%00:", "a\n:"} {
		edges = append(edges, input{false, []byte(raw)}, input{true, []byte(raw)})
	}
	f, err := os.CreateTemp("", "golean-url-lexical-")
	must(err)
	defer os.Remove(f.Name())
	encoder := json.NewEncoder(f)
	batch("boundaries", edges, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			raw := make([]byte, rng.Intn(42))
			if i%2 == 0 {
				const alphabet = "0123456789AZaz[]:+.-\x00\xff%/?"
				for j := range raw {
					raw[j] = alphabet[rng.Intn(len(alphabet))]
				}
			} else {
				rng.Read(raw)
			}
			switch i % 5 {
			case 0:
				raw = append(raw, ':', '0', byte('0'+rng.Intn(10)))
			case 1:
				raw = append([]byte("ABC+.-9:"), raw...)
			case 2:
				raw = append(append([]byte{'['}, raw...), ']')
			}
			cases[i] = input{i%11 != 0, raw}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, encoder)
	}
	must(f.Close())
	must(os.Setenv("GOLEAN_URL_LEXICAL_CASES", f.Name()))
	const data = "difftest/net/url/lexical/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanURLLexical", data+"bridge_test.go", data+"lexical_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
