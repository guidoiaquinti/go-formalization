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
	Raw  []byte
	Mode byte
}

type fixture struct {
	input
	Want string
}

func wire(kind string, raw []byte) string {
	fields := make([]string, len(raw))
	for i, b := range raw {
		fields[i] = strconv.Itoa(int(b))
	}
	return kind + ":" + strings.Join(fields, ",")
}

func expected(raw []byte, mode byte) string {
	const hostSafe = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.~!$&'()*+,;=:[]<>\""
	var output []byte
	for len(raw) > 0 {
		b := raw[0]
		if b == '%' {
			n := min(3, len(raw))
			if n != 3 {
				return wire("escape", raw[:n])
			}
			decoded := []byte{0}
			if _, err := hex.Decode(decoded, raw[1:3]); err != nil {
				return wire("escape", raw[:3])
			}
			b = decoded[0]
			// %25 is the only ASCII escape accepted in host mode.
			if mode == 4 && b < 128 && b != '%' || mode == 8 && b != '%' && b != ' ' && !strings.ContainsRune(hostSafe, rune(b)) {
				return wire("escape", raw[:3])
			}
			raw = raw[3:]
		} else {
			if b == '+' && mode == 32 {
				b = ' '
			}
			if (mode == 4 || mode == 8) && b < 128 && !strings.ContainsRune(hostSafe, rune(b)) {
				return wire("host", raw[:1])
			}
			raw = raw[1:]
		}
		output = append(output, b)
	}
	return wire("ok", output)
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
			rows[i] = fixture{c, strings.Join([]string{expected(c.Raw, c.Mode), expected(c.Raw, 32), expected(c.Raw, 2)}, "|")}
		}
		for i, got := range check.LeanRows("godiff-net-url-unescape", len(current), &protocol) {
			if got != rows[i].Want {
				log.Fatalf("%s row %d %+v: got %q want %q", label, first+i, current[i], got, rows[i].Want)
			}
		}
		for _, r := range rows {
			must(encoder.Encode(r))
		}
	}
	fmt.Printf("URL unescaping %s: %d Lean rows agree.\n", label, len(cases))
}

func main() {
	var edges []input
	for i := range 1 << 16 {
		raw := []byte{'x', '%', '2', '5', '+', '%', byte(i >> 8), byte(i), 'z'}
		edges = append(edges, input{raw, []byte{2, 4, 8, 32}[i%4]})
	}
	for mode := range 256 {
		for _, raw := range []string{"", "%", "%a", "%+1x", "%00%ff+", "a/%41", "%41/%0Z", "%20+%25", "\xff%80", "%2b%25"} {
			edges = append(edges, input{[]byte(raw), byte(mode)})
		}
	}
	for v := range 256 {
		for _, mode := range []byte{2, 4, 8, 32} {
			edges = append(edges, input{[]byte{byte(v)}, mode},
				input{[]byte(fmt.Sprintf("%%%02x", v)), mode}, input{[]byte(fmt.Sprintf("%%%02X", v)), mode})
		}
	}
	file, err := os.CreateTemp("", "golean-url-unescape-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	batch("boundaries", edges, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			var raw []byte
			for range rng.Intn(20) {
				switch rng.Intn(4) {
				case 0:
					raw = append(raw, byte(rng.Intn(256)))
				case 1:
					raw = append(raw, fmt.Sprintf("%%%02x", rng.Intn(256))...)
				case 2:
					raw = append(raw, '+')
				case 3:
					raw = append(raw, 'a'+byte(rng.Intn(26)))
				}
			}
			if i%5 == 0 {
				raw = append(raw, '%')
			}
			mode := byte(rng.Intn(256))
			if i%2 == 0 {
				mode = []byte{2, 4, 8, 32}[i/2%4]
			}
			cases[i] = input{raw, mode}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, encoder)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_URL_UNESCAPE_CASES", file.Name()))
	const data = "difftest/net/url/unescape/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanUnescape", data+"bridge_test.go", data+"unescape_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
