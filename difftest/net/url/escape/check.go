package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Raw  []byte
	Mode byte
}

type fixture struct {
	input
	Want string
}

const (
	segment = 2
	query   = 32
)

var specials = [...]string{"$&+,/:;=@", "$&+:=@", "!$&'()*+,;=:[]<>\"",
	"!$&'()*+,;=:[]<>\"", "$&+,;=", "", "$&+,/:;=?@!()*"}

func classify(c byte) (mask byte, isHex bool, nibble byte) {
	const base = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.~"
	for i, suffix := range specials {
		if strings.IndexByte(base, c) >= 0 || strings.IndexByte(suffix, c) >= 0 {
			mask |= 1 << i
		}
	}
	value, err := strconv.ParseUint(string([]byte{c}), 16, 8)
	if err == nil {
		return mask | 128, true, byte(value)
	}
	// Direct unhex is called on every byte. For non-hex inputs its
	// result is the source arithmetic, outside that helper's contract.
	return mask, false, 9*(c/64) + c%16
}

func encodeBytes(raw []byte, mode byte) string {
	var out strings.Builder
	for _, c := range raw {
		mask, _, _ := classify(c)
		switch {
		case mode == query && c == ' ':
			out.WriteByte('+')
		case mode&mask == 0:
			fmt.Fprintf(&out, "%%%02X", c)
		default:
			out.WriteByte(c)
		}
	}
	return out.String()
}

func wire(raw string) string {
	fields := make([]string, len(raw))
	for i := range len(raw) {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return strings.Join(fields, ",")
}

func expected(c input) string {
	fields := make([]string, len(c.Raw))
	for i, b := range c.Raw {
		mask, hex, nibble := classify(b)
		fields[i] = fmt.Sprintf("%t:%t:%d", mask&c.Mode == 0, hex, nibble)
	}
	return strings.Join([]string{strings.Join(fields, ","), wire(encodeBytes(c.Raw, c.Mode)),
		wire(encodeBytes(c.Raw, query)), wire(encodeBytes(c.Raw, segment))}, "|")
}

func batch(label string, cases []input, encoder *json.Encoder) {
	const shard = 1024
	for first := 0; first < len(cases); first += shard {
		current := cases[first:min(first+shard, len(cases))]
		var protocol bytes.Buffer
		rows := make([]fixture, len(current))
		for i, c := range current {
			fmt.Fprint(&protocol, c.Mode)
			for _, b := range c.Raw {
				fmt.Fprint(&protocol, " ", b)
			}
			fmt.Fprintln(&protocol)
			rows[i] = fixture{c, expected(c)}
		}
		for i, got := range check.LeanRows("godiff-net-url-escape", len(current), &protocol) {
			if got != rows[i].Want {
				log.Fatalf("%s row %d %+v: got %q want %q", label, first+i, current[i], got, rows[i].Want)
			}
		}
		for _, r := range rows {
			must(encoder.Encode(r))
		}
	}
	fmt.Printf("URL escaping %s: %d Lean rows agree.\n", label, len(cases))
}

func main() {
	var edges []input
	for mode := range 256 {
		for c := range 256 {
			edges = append(edges, input{[]byte{byte(c)}, byte(mode)})
		}
	}
	for _, mode := range []byte{0, 1, 2, 4, 8, 16, 32, 64, 128, 255} {
		for _, raw := range []string{"", " ", "+%/%\x00", "a&b-c?d#e", "\xff\xe2\x82K",
			strings.Repeat(" ", 64), strings.Repeat(" ", 65), strings.Repeat("x", 61) + "!x",
			strings.Repeat("x", 62) + "!x", strings.Repeat("x", 63) + "!x"} {
			edges = append(edges, input{[]byte(raw), mode})
		}
	}
	file, err := os.CreateTemp("", "golean-url-escape-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	batch("boundaries", edges, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			raw := make([]byte, rng.Intn(57))
			if i%3 == 0 {
				const candidates = "aAZ019~?!+-_*/\x00\xff"
				for j := range raw {
					raw[j] = candidates[rng.Intn(len(candidates))]
				}
			} else {
				rng.Read(raw)
			}
			mode := byte(rng.Intn(256))
			if i%2 == 0 {
				mode = 1 << rng.Intn(7)
			}
			cases[i] = input{raw, mode}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, encoder)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_URL_ESCAPE_CASES", file.Name()))
	const data = "difftest/net/url/escape/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanEscape", data+"bridge_test.go", data+"escape_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
