package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	Config        string
	Raw           []byte
	Params, Start int
	Present       bool
}

type fixture struct {
	input
	Want string
}

func limit(params int, config string) (bool, bool) {
	standard := params <= 10_000
	custom, err := strconv.Atoi(config)
	if err != nil {
		return standard, false
	}
	allowed := custom == 0 || params < custom
	return allowed, allowed != standard
}

func wire(s string) string {
	fields := make([]string, len(s))
	for i := range fields {
		fields[i] = strconv.Itoa(int(s[i]))
	}
	return fmt.Sprintf("%d:%s", len(s), strings.Join(fields, ","))
}

func state(values map[string][]string) string {
	if values == nil {
		return "nil"
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	var parts []string
	for _, key := range keys {
		payload := "nil"
		if values[key] != nil {
			items := make([]string, len(values[key]))
			for i, w := range values[key] {
				items[i] = wire(w)
			}
			payload = "list:" + strings.Join(items, ";")
		}
		parts = append(parts, wire(key)+"="+payload)
	}
	return "map:" + strings.Join(parts, "/")
}

func decode(raw string) (string, string) {
	var out strings.Builder
	for i := 0; i < len(raw); i++ {
		switch raw[i] {
		case '+':
			out.WriteByte(' ')
		case '%':
			triple := raw[i:min(i+3, len(raw))]
			if len(triple) < 3 {
				return "", triple
			}
			var value [1]byte
			if _, err := hex.Decode(value[:], []byte(triple[1:])); err != nil {
				return "", triple
			}
			out.WriteByte(value[0])
			i += 2
		default:
			out.WriteByte(raw[i])
		}
	}
	return out.String(), ""
}

func parse(values map[string][]string, raw, config string) (string, string) {
	parts := strings.Split(raw, "&")
	if allowed, _ := limit(len(parts), config); !allowed {
		return "limit", state(values)
	}
	status := "ok"
	for _, field := range parts {
		if field == "" {
			continue
		}
		if strings.IndexByte(field, ';') >= 0 {
			status = "semicolon"
			continue
		}
		key, value, _ := strings.Cut(field, "=")
		key, bad := decode(key)
		if bad == "" {
			value, bad = decode(value)
		}
		if bad != "" {
			if status == "ok" {
				status = "escape:" + wire(bad)
			}
			continue
		}
		if values == nil {
			return "panic", ""
		}
		values[key] = append(values[key], value)
	}
	return status, state(values)
}

func expected(c input) string {
	m := map[string][]string(nil)
	if c.Start == 1 {
		m = map[string][]string{}
	}
	if c.Start == 2 {
		m = map[string][]string{"a": {"old"}, "z": nil}
	}
	config, raw := c.Config, string(c.Raw)
	ok, changed := limit(c.Params, config)
	private, privateState := parse(m, raw, config)
	if private != "panic" {
		private += "/" + privateState
	}
	public, publicState := parse(map[string][]string{}, raw, config)
	result := "panic"
	if c.Present {
		result = publicState
	}
	return fmt.Sprintf("%t:%t|%s|%s/%s|%s", ok, changed, private, public, publicState, result)
}

func batch(label string, cases []input, encoder *json.Encoder) {
	const shard = 512
	for first := 0; first < len(cases); first += shard {
		current := cases[first:min(first+shard, len(cases))]
		var protocol bytes.Buffer
		rows := make([]fixture, len(current))
		for i, c := range current {
			tag := 0
			if c.Present {
				tag = 1
			}
			fmt.Fprintf(&protocol, "%d %d %d %d %d", strconv.IntSize, c.Params, c.Start, tag, len(c.Config))
			for j := range len(c.Config) {
				fmt.Fprint(&protocol, " ", c.Config[j])
			}
			for _, b := range c.Raw {
				fmt.Fprint(&protocol, " ", b)
			}
			fmt.Fprintln(&protocol)
			rows[i] = fixture{c, expected(c)}
		}
		for i, got := range check.LeanRows("godiff-net-url-query", len(current), &protocol) {
			if got != rows[i].Want {
				log.Fatalf("%s row %d got %q want %q", label, first+i, got, rows[i].Want)
			}
		}
		for _, r := range rows {
			must(encoder.Encode(r))
		}
	}
	fmt.Printf("URL query %s: %d Lean rows agree.\n", label, len(cases))
}

func main() {
	settings := []string{"", "0", "1", "2", "5", "10000", "10001", "-1", "+3", "-0", "00", "abc", " 2", "1_0",
		strconv.Itoa(math.MaxInt), strconv.Itoa(math.MinInt), strconv.Itoa(math.MaxInt) + "0"}
	words := []string{"", "a", "a=1&a=2", "a+b=+%2B", "x=%00%ff", "%0x&good=1", "%GG;v=%XX&%", "%&;x&good",
		"a=b=c&=&=+&bare", "&&a&", "%61=old&a=new", "a=%2G&x=%%%&z", "z=\xff", "a=%3B;b=c", "a=%3B&semi=%3b", "%G1=v%12"}
	var edges []input
	for _, setting := range settings {
		for i, raw := range words {
			for _, n := range []int{math.MinInt, -1, 0, 1, 2, 9999, 10000, 10001, math.MaxInt} {
				edges = append(edges, input{setting, []byte(raw), n, i % 3, i%4 != 0})
			}
		}
	}
	for _, n := range []int{9999, 10000, 10001} {
		for _, setting := range []string{"", "0", "10000", "10001", "bad"} {
			edges = append(edges, input{setting, []byte(strings.Repeat("&", n-1) + "a=x"), n, 2, true})
		}
	}
	file, err := os.CreateTemp("", "golean-url-query-")
	must(err)
	defer os.Remove(file.Name())
	encoder := json.NewEncoder(file)
	batch("boundaries", edges, encoder)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		cases := make([]input, 20000)
		for i := range cases {
			cfg := settings[i%len(settings)]
			var parts []string
			for range rng.Intn(12) {
				if rng.Intn(2) == 0 {
					parts = append(parts, words[rng.Intn(len(words))])
				} else {
					raw := make([]byte, rng.Intn(9))
					rng.Read(raw)
					parts = append(parts, string(raw))
				}
			}
			n := rng.Intn(20) - 3
			if i%3 == 0 {
				n = 9996 + rng.Intn(10)
			}
			if i%7 == 0 {
				n = int(rng.Uint64())
			}
			cases[i] = input{cfg, []byte(strings.Join(parts, "&")), n, i % 3, i%11 != 0}
		}
		batch(fmt.Sprintf("seed=%d", seed), cases, encoder)
	}
	must(file.Close())
	must(os.Setenv("GOLEAN_URL_QUERY_CASES", file.Name()))
	const data = "difftest/net/url/query/testdata/"
	must(check.TestInPackage("net/url", "TestGoLeanURLQuery", data+"bridge_test.go", data+"query_test.go"))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
