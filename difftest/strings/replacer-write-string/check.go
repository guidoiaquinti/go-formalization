package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	pairs        []string
	text, second string
	stringWriter bool
	script       [][2]int
	mode         int
}

func main() {
	compare("boundaries", boundaries())
	for _, seed := range []int64{17, 2718, 65537} {
		rng := rand.New(rand.NewSource(seed))
		for block := range 4 {
			cases := make([]input, 5000)
			for i := range cases {
				cases[i] = sample(rng, i)
			}
			compare(fmt.Sprintf("seed=%d block=%d", seed, block), cases)
		}
	}
}

func boundaries() []input {
	cases := []input{{mode: 0}, {text: "abc", second: "def", mode: 0}, {text: "abc", second: "", mode: 1}, {text: "", second: "a", mode: 1}}
	for _, pairs := range [][]string{
		nil, {"a", "b"}, {"a", "ab"}, {"ab", "X"}, {"a", ""}, {"a", "a", "b", "c"}, {"", ""}, {"", "|", "a", "X"}, {"a", "X", "", "|"}, {"", "first", "", "second"},
		{"ab", "L", "a", "S"}, {"a", "", "ab", "L"}, {"a", "first", "a", "second"},
		{"\x80", "\xff\x00", "\xff", "", "\xe2\x82", "!"}, {"ǆ", "z", "", "."},
	} {
		for _, text := range []string{"", "abc", "abba", "xabab", "aaa", "z", "ǆ🚀", "\xff\x80\xe2\x82", "\x00\xff"} {
			for _, script := range [][][2]int{nil, {{0, 0}}, {{-1, 1}}, {{0, 0}, {1, 2}}, {{-1, 0}, {-1, 0}, {0, 1}}, {{0, 0}, {1, 0}, {1, 0}}} {
				for _, stringWriter := range []bool{false, true} {
					cases = append(cases, input{pairs, text, text + "z", stringWriter, script, 2})
				}
			}
		}
	}
	for b := range 256 {
		key := string([]byte{byte(b)})
		pairs := []string{"", "", key + key, "\x00\xff", key, ""}
		cases = append(cases, input{pairs, key + "!" + key + key, key, false, nil, 2},
			input{pairs, key + key, key + "a", true, [][2]int{{0, 0}, {-1, 0}, {-1, 0}, {1, 2}}, 2})
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
	return input{pairs, text, text + "\xff", rng.Intn(2) == 0, script, 2}
}

func genericPayloads(c input) []string {
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
	return payloads
}

type event struct {
	mode    byte
	payload string
}

func payloadEvents(c input) []event {
	mode := byte('B')
	if c.stringWriter {
		mode = 'S'
	}
	var out []event
	if len(c.pairs) == 2 && len(c.pairs[0]) > 1 {
		key, value := c.pairs[0], c.pairs[1]
		rest := c.text
		for {
			at := strings.Index(rest, key)
			if at < 0 {
				return append(out, event{mode, rest})
			}
			out = append(out, event{mode, rest[:at]}, event{mode, value})
			rest = rest[at+len(key):]
		}
	}
	var values [256]*string
	allByteValues := true
	for i := 0; i < len(c.pairs); i += 2 {
		key, value := c.pairs[i], c.pairs[i+1]
		if len(key) != 1 {
			for _, payload := range genericPayloads(c) {
				out = append(out, event{mode, payload})
			}
			return out
		}
		if values[key[0]] == nil {
			values[key[0]] = &c.pairs[i+1]
		}
		allByteValues = allByteValues && len(value) == 1
	}
	last := 0
	for i := range len(c.text) {
		original := c.text[i]
		if value := values[original]; value != nil && (!allByteValues || (*value)[0] != original) {
			if last < i {
				out = append(out, event{mode, c.text[last:i]})
			}
			out = append(out, event{'B', *value})
			last = i + 1
		}
	}
	if last < len(c.text) {
		out = append(out, event{mode, c.text[last:]})
	}
	return out
}

var failures = [...]error{nil, errors.New("first"), errors.New("second")}

type scripted struct {
	script [][2]int
	calls  []string
}

type stringScripted struct{ *scripted }

func (w *scripted) Write(raw []byte) (int, error) {
	if raw == nil {
		return w.reply(event{'N', ""})
	}
	return w.reply(event{'B', string(raw)})
}

func (w stringScripted) WriteString(raw string) (int, error) {
	return w.reply(event{'S', raw})
}

func (w *scripted) reply(c event) (int, error) {
	reply := [2]int{-1, 0}
	if i := len(w.calls); i < len(w.script) {
		reply = w.script[i]
	}
	wire := []byte{c.mode}
	for i := range len(c.payload) {
		wire = strconv.AppendUint(append(wire, ','), uint64(c.payload[i]), 10)
	}
	w.calls = append(w.calls, string(wire))
	count := len(c.payload)
	if reply[0] >= 0 {
		count = min(count, reply[0])
	}
	return count, failures[reply[1]]
}

func show(count int, failure error, calls []string) string {
	code := -1
	for i, value := range failures {
		if failure == value {
			code = i
			break
		}
	}
	return fmt.Sprintf("%d:%d:%s", count, code, strings.Join(calls, ";"))
}

func second(c input) input {
	return input{c.pairs, c.second, "", !c.stringWriter, nil, c.mode}
}

func expected(c input) string {
	if c.mode == 0 {
		return "panic"
	}
	one := func(c input) string {
		writer := &scripted{script: c.script}
		count := 0
		var failure error
		for _, event := range payloadEvents(c) {
			n, err := writer.reply(event)
			count += n
			failure = err
			if err != nil {
				break
			}
		}
		return show(count, failure, writer.calls)
	}
	return one(c) + "|" + one(second(c))
}

func observe(c input) (result string) {
	var replacer *strings.Replacer
	switch c.mode {
	case 1:
		replacer = new(strings.Replacer)
	case 2:
		replacer = strings.NewReplacer(c.pairs...)
	}
	defer func() {
		if value := recover(); value != nil {
			fault, ok := value.(runtime.Error)
			if c.mode != 0 || !ok || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(value)
			}
			result = "panic"
		}
	}()
	one := func(c input) string {
		writer := &scripted{script: c.script}
		var target io.Writer = writer
		if c.stringWriter {
			target = stringScripted{writer}
		}
		count, failure := replacer.WriteString(target, c.text)
		return show(count, failure, writer.calls)
	}
	return one(c) + "|" + one(second(c))
}

func compare(label string, cases []input) {
	var input bytes.Buffer
	for _, c := range cases {
		mode := 0
		if c.stringWriter {
			mode = 1
		}
		fmt.Fprintf(&input, "%d %d %d %d", c.mode, mode, len(c.script), len(c.pairs)/2)
		for _, reply := range c.script {
			fmt.Fprintf(&input, " %d %d", reply[0], reply[1])
		}
		add := func(raw string) {
			for i := range len(raw) {
				fmt.Fprintf(&input, " %d", raw[i])
			}
		}
		for i := 0; i < len(c.pairs); i += 2 {
			fmt.Fprintf(&input, " %d %d", len(c.pairs[i]), len(c.pairs[i+1]))
			add(c.pairs[i])
			add(c.pairs[i+1])
		}
		fmt.Fprintf(&input, " %d", len(c.text))
		add(c.text)
		add(c.second)
		fmt.Fprintln(&input)
	}
	for i, model := range check.LeanRows("godiff-strings-replacer-write-string", len(cases), &input) {
		c := cases[i]
		want := expected(c)
		if native := observe(c); native != want || model != want {
			log.Fatalf("%s row %d pairs=%q inputs=%x,%x: Go=%q Lean=%q want=%q", label, i, c.pairs, c.text, c.second, native, model, want)
		}
	}
	fmt.Printf("strings.Replacer.WriteString %s: %d first and repeated trace pairs agree.\n", label, len(cases))
}
