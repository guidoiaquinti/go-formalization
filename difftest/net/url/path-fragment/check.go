package main

import (
	"bytes"
	"encoding/hex"
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
	Present              bool
	Mode                 byte
	Decoded, Raw, Update []byte
}

type fixture struct {
	input
	Want string
}

var extra = [...]string{"$&+,/:;=@", "$&+:=@", "!$&'()*+,;=:[]<>\"",
	"!$&'()*+,;=:[]<>\"", "$&+,;=", "", "$&+,/:;=?@!()*", "0123456789ABCDEFabcdef"}

func allowed(c byte, mode byte) bool {
	for i, suffix := range extra {
		if mode&(1<<i) == 0 {
			continue
		}
		if i != 7 && strings.IndexByte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_.~", c) >= 0 ||
			strings.IndexByte(suffix, c) >= 0 {
			return true
		}
	}
	return false
}

func valid(raw []byte, mode byte) bool {
	for _, b := range raw {
		if strings.IndexByte("!$&'()*+,;=:@[]%", b) < 0 && !allowed(b, mode) {
			return false
		}
	}
	return true
}

func encodePart(raw []byte, mode byte) []byte {
	var out bytes.Buffer
	for _, b := range raw {
		if allowed(b, mode) {
			out.WriteByte(b)
		} else {
			fmt.Fprintf(&out, "%%%02X", b)
		}
	}
	return out.Bytes()
}

func decodePart(raw []byte) (value, invalid []byte) {
	for len(raw) > 0 {
		before, after, found := bytes.Cut(raw, []byte{'%'})
		value = append(value, before...)
		if !found {
			return value, nil
		}
		if len(after) < 2 {
			return nil, raw[len(before):]
		}
		var one [1]byte
		if _, err := hex.Decode(one[:], after[:2]); err != nil {
			return nil, raw[len(before) : len(before)+3]
		}
		value = append(value, one[0])
		raw = after[2:]
	}
	return value, nil
}

func wire(raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return strings.Join(fields, ",")
}

func expected(c input) string {
	escaped := func(mode byte, star bool) string {
		if !c.Present {
			return "panic"
		}
		if len(c.Raw) > 0 && valid(c.Raw, mode) {
			decoded, fail := decodePart(c.Raw)
			if fail == nil && bytes.Equal(decoded, c.Decoded) {
				return "ok:" + wire(c.Raw)
			}
		}
		if star && bytes.Equal(c.Decoded, []byte{'*'}) {
			return "ok:42"
		}
		return "ok:" + wire(encodePart(c.Decoded, mode))
	}
	update := func(mode byte) string {
		out, bad := decodePart(c.Update)
		if bad != nil {
			state := "nil"
			if c.Present {
				state = wire(c.Decoded) + "/" + wire(c.Raw)
			}
			return "escape:" + wire(bad) + "/" + state
		}
		if !c.Present {
			return "panic"
		}
		raw := c.Update
		if bytes.Equal(encodePart(out, mode), raw) {
			raw = nil
		}
		return "ok/" + wire(out) + "/" + wire(raw)
	}
	return strings.Join([]string{strconv.FormatBool(valid(c.Raw, c.Mode)), escaped(1, true), escaped(64, false), update(1), update(64)}, "|")
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
			fmt.Fprintf(&protocol, "%d %d %d %d", tag, c.Mode, len(c.Decoded), len(c.Raw))
			for _, xs := range [][]byte{c.Decoded, c.Raw, c.Update} {
				for _, b := range xs {
					fmt.Fprint(&protocol, " ", b)
				}
			}
			fmt.Fprintln(&protocol)
			rows[i] = fixture{c, expected(c)}
		}
		for i, got := range check.LeanRows("godiff-net-url-path-fragment", len(current), &protocol) {
			if got != rows[i].Want {
				log.Fatalf("%s row %d %+v: got %q want %q", label, first+i, current[i], got, rows[i].Want)
			}
		}
		for _, r := range rows {
			must(encoder.Encode(r))
		}
	}
	fmt.Printf("URL path and fragment %s: %d Lean rows agree.\n", label, len(cases))
}

func main() {
	var edges []input
	for mode := range 256 {
		for b := range 256 {
			raw := []byte{byte(b)}
			edges = append(edges, input{true, byte(mode), raw, raw, raw})
		}
	}
	words := []string{"", "*", "%2a", "%2A", "a/b", "a%2fb", "x%3f", "x?", "%", "a%g1", "a%1", "%25", "%FF", "\xff", "[]!", "%20", "%2B", "+", "#", "é"}
	for _, decoded := range words {
		for _, raw := range words {
			for _, present := range []bool{false, true} {
				edges = append(edges, input{present, 1, []byte(decoded), []byte(raw), []byte(raw)})
			}
		}
	}
	file, err := os.CreateTemp("", "golean-url-path-fragment-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	batch("boundaries", edges, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		randomBytes := func() []byte {
			raw := make([]byte, rng.Intn(40))
			rng.Read(raw)
			return raw
		}
		cases := make([]input, 20000)
		for i := range cases {
			decoded, raw, update := randomBytes(), randomBytes(), randomBytes()
			switch i % 4 {
			case 0:
				raw = encodePart(decoded, 1)
			case 1:
				raw = encodePart(decoded, 64)
			case 2:
				raw = append(encodePart(decoded, 1), '%')
			}
			if i%2 == 0 {
				update = raw
			}
			cases[i] = input{i%13 != 0, byte(rng.Intn(256)), decoded, raw, update}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, encoder)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_URL_PATH_FRAGMENT_CASES", file.Name()))
	const data = "difftest/net/url/path-fragment/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanPathFragment", data+"bridge_test.go", data+"path_fragment_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
