package main

import (
	"bufio"
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
	pairs        []string
	text         string
	stringWriter bool
	script       [][2]int
}

func main() {
	file, err := os.CreateTemp("", "golean-generic-replacer-write-string-")
	must(err)
	defer os.Remove(file.Name())
	buffer := bufio.NewWriter(file)
	fixture := json.NewEncoder(buffer)
	compare("boundaries", boundaries(), fixture)
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5000)
			for i := range cases {
				cases[i] = sample(rng, i)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases, fixture)
		}
	}
	must(buffer.Flush())
	must(file.Close())
	must(os.Setenv("GOLEAN_GENERIC_REPLACER_WRITE_STRING_CASES", file.Name()))
	const data = "difftest/strings/generic-replacer-write-string/testdata/"
	must(check.TestInPackage("strings", "TestGoLeanGenericReplacerWriteString", data+"bridge_test.go", data+"write_string_test.go"))
}

func boundaries() []input {
	var cases []input
	for _, pairs := range [][]string{
		nil, {"", ""}, {"", "|", "a", "X"}, {"a", "X", "", "|"}, {"", "first", "", "second"},
		{"ab", "L", "a", "S"}, {"a", "", "ab", "L"}, {"a", "first", "a", "second"},
		{"\x80", "\xff\x00", "\xff", "", "\xe2\x82", "!"}, {"ǆ", "z", "", "."},
	} {
		for _, text := range []string{"", "abc", "abba", "xabab", "aaa", "z", "ǆ🚀", "\xff\x80\xe2\x82", "\x00\xff"} {
			for _, script := range [][][2]int{nil, {{0, 0}}, {{-1, 1}}, {{0, 0}, {1, 2}}, {{-1, 0}, {-1, 0}, {0, 1}}, {{0, 0}, {1, 0}, {1, 0}}} {
				for _, stringWriter := range []bool{false, true} {
					cases = append(cases, input{pairs, text, stringWriter, script})
				}
			}
		}
	}
	for b := range 256 {
		key := string([]byte{byte(b)})
		pairs := []string{"", "", key + key, "\x00\xff", key, ""}
		cases = append(cases, input{pairs, key + "!" + key + key, false, nil},
			input{pairs, key + key, true, [][2]int{{0, 0}, {-1, 0}, {-1, 0}, {1, 2}}})
	}
	return cases
}

func sample(rng *rand.Rand, index int) input {
	alphabet := [...]int{2, 8, 256}[index%3]
	raw := func(limit int) string {
		out := make([]byte, rng.Intn(limit+1))
		for i := range out {
			out[i] = byte(rng.Intn(alphabet))
		}
		return string(out)
	}
	var pairs []string
	for range rng.Intn(12) {
		key := raw(12)
		if len(pairs) != 0 {
			base := pairs[2*rng.Intn(len(pairs)/2)]
			switch rng.Intn(4) {
			case 0:
				key = base[:rng.Intn(len(base)+1)]
			case 1:
				key = base + key
			}
		}
		pairs = append(pairs, key, raw(16))
	}
	text := raw(36)
	if len(pairs) != 0 && index%3 != 0 {
		key := pairs[2*rng.Intn(len(pairs)/2)]
		position := rng.Intn(len(text) + 1)
		text = text[:position] + key + text[position:] + key
	}
	script := make([][2]int, rng.Intn(12))
	for i := range script {
		script[i][0] = rng.Intn(7) - 1
		if index%5 != 0 && rng.Intn(4) == 0 {
			script[i][1] = 1 + rng.Intn(2)
		}
	}
	return input{pairs, text, rng.Intn(2) == 0, script}
}

// At each byte position, enumerate the first eligible pair. If it is
// empty, emit it and then consider the first nonempty pair at that same
// byte position before advancing.
func expected(c input) string {
	find := func(text string, nonempty bool) int {
		for i := 0; i < len(c.pairs); i += 2 {
			key := c.pairs[i]
			if (!nonempty || key != "") && strings.HasPrefix(text, key) {
				return i
			}
		}
		return -1
	}
	var payloads []string
	last := 0
	for at := 0; at <= len(c.text); {
		text := c.text[at:]
		match := find(text, false)
		if match >= 0 && c.pairs[match] == "" {
			payloads = append(payloads, c.text[last:at], c.pairs[match+1])
			last = at
			match = find(text, true)
		}
		if match < 0 {
			at++
			continue
		}
		payloads = append(payloads, c.text[last:at], c.pairs[match+1])
		at += len(c.pairs[match])
		last = at
	}
	if last < len(c.text) {
		payloads = append(payloads, c.text[last:])
	}
	calls := make([]string, 0, len(payloads))
	mode := byte('B')
	if c.stringWriter {
		mode = 'S'
	}
	count, failure := 0, 0
	for i, payload := range payloads {
		call := []byte{mode}
		for j := range len(payload) {
			call = strconv.AppendUint(append(call, ','), uint64(payload[j]), 10)
		}
		calls = append(calls, string(call))
		accepted := len(payload)
		if i < len(c.script) {
			if limit := c.script[i][0]; limit >= 0 {
				accepted = min(accepted, limit)
			}
			failure = c.script[i][1]
		}
		count += accepted
		if failure != 0 {
			break
		}
	}
	return fmt.Sprintf("%d:%d:%s", count, failure, strings.Join(calls, ";"))
}

func encode(c input) string {
	mode := 0
	if c.stringWriter {
		mode = 1
	}
	row := fmt.Appendf(nil, "%d %d %d", mode, len(c.script), len(c.pairs)/2)
	for _, reply := range c.script {
		row = fmt.Appendf(row, " %d %d", reply[0], reply[1])
	}
	add := func(raw string) {
		for i := range len(raw) {
			row = strconv.AppendUint(append(row, ' '), uint64(raw[i]), 10)
		}
	}
	for i := 0; i < len(c.pairs); i += 2 {
		key, value := c.pairs[i], c.pairs[i+1]
		row = fmt.Appendf(row, " %d %d", len(key), len(value))
		add(key)
		add(value)
	}
	add(c.text)
	return string(row)
}

func compare(label string, cases []input, fixture *json.Encoder) {
	var model bytes.Buffer
	for _, c := range cases {
		fmt.Fprintln(&model, encode(c))
	}
	for i, got := range check.LeanRows("godiff-strings-generic-replacer-write-string", len(cases), &model) {
		c := cases[i]
		want := expected(c)
		if got != want {
			log.Fatalf("%s row %d pairs=%q text=%x: Lean=%s want=%s", label, i, c.pairs, c.text, got, want)
		}
		pairs := make([][]byte, len(c.pairs))
		for j, field := range c.pairs {
			pairs[j] = []byte(field)
		}
		must(fixture.Encode(struct {
			Pairs        [][]byte
			Text         []byte
			StringWriter bool
			Script       [][2]int
			Want         string
		}{pairs, []byte(c.text), c.stringWriter, c.script, want}))
	}
	fmt.Printf("genericReplacer.WriteString %s: %d Lean traces agree.\n", label, len(cases))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
