package main

import (
	"fmt"
	"maps"
	"math/rand"
	"net/url"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type entry struct {
	key   string
	words []string
}

type input struct {
	entries    []entry
	key, value string
}

func text(raw string) string {
	fields := make([]string, len(raw))
	for i := range len(raw) {
		fields[i] = strconv.Itoa(int(raw[i]))
	}
	return fmt.Sprintf("%d:%s", len(raw), strings.Join(fields, ","))
}

func render(entries []entry) string {
	if entries == nil {
		return "nil"
	}
	frozen := slices.Clone(entries)
	slices.SortFunc(frozen, func(a, b entry) int { return strings.Compare(a.key, b.key) })
	fields := make([]string, len(frozen))
	for i, e := range frozen {
		payload := "nil"
		if e.words != nil {
			items := make([]string, len(e.words))
			for j, word := range e.words {
				items[j] = text(word)
			}
			payload = "list:" + strings.Join(items, ";")
		}
		fields[i] = text(e.key) + "=" + payload
	}
	return "map:" + strings.Join(fields, "/")
}

func nativeEntries(m url.Values) []entry {
	if m == nil {
		return nil
	}
	entries := make([]entry, 0, len(m))
	for key, words := range m {
		entries = append(entries, entry{key, words})
	}
	return entries
}

func native(c input) string {
	construct := func() url.Values {
		if c.entries == nil {
			return nil
		}
		m := make(url.Values, len(c.entries))
		for _, e := range c.entries {
			m[e.key] = slices.Clone(e.words)
		}
		return m
	}
	m := construct()
	call := func(action func(url.Values)) (out string) {
		before := construct()
		defer func() {
			if p := recover(); p != nil {
				err, ok := p.(runtime.Error)
				if before != nil || !ok || err.Error() != "assignment to entry in nil map" {
					panic(p)
				}
				out = "panic"
			}
		}()
		action(before)
		return render(nativeEntries(before))
	}
	set := call(func(v url.Values) { v.Set(c.key, c.value) })
	add := call(func(v url.Values) { v.Add(c.key, c.value) })
	del := call(func(v url.Values) { v.Del(c.key) })
	cloned := m.Clone()
	cloneWire := render(nativeEntries(cloned))
	encoded := m.Encode()
	if cloned != nil {
		// Mutating copied value slots and entries must leave the original map intact.
		for key, words := range cloned {
			for i := range words {
				words[i] += "\x00review"
			}
			delete(cloned, key)
		}
		cloned["\x00review"] = []string{"probe"}
		if !maps.EqualFunc(m, construct(), func(a, b []string) bool {
			return (a == nil) == (b == nil) && slices.Equal(a, b)
		}) {
			panic("Values.Clone aliases original data")
		}
	}
	return strings.Join([]string{text(m.Get(c.key)), strconv.FormatBool(m.Has(c.key)), set, add, del, cloneWire, text(encoded)}, "|")
}

func replace(entries []entry, key string, words []string) []entry {
	retained := make([]entry, 0, len(entries)+1)
	for _, e := range entries {
		if e.key != key {
			retained = append(retained, e)
		}
	}
	if words != nil {
		retained = append(retained, entry{key, words})
	}
	return retained
}

func quote(raw string) string {
	const safe = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.~"
	var out strings.Builder
	for i := range len(raw) {
		b := raw[i]
		switch {
		case b == ' ':
			out.WriteByte('+')
		case strings.IndexByte(safe, b) >= 0:
			out.WriteByte(b)
		default:
			fmt.Fprintf(&out, "%%%02X", b)
		}
	}
	return out.String()
}

func expected(c input) string {
	var first string
	var existing []string
	found := false
	for _, e := range c.entries {
		if e.key == c.key {
			existing, found = e.words, true
			if len(e.words) > 0 {
				first = e.words[0]
			}
			break
		}
	}
	if c.entries == nil {
		return text(first) + "|false|panic|panic|nil|nil|" + text("")
	}
	keys := slices.Clone(c.entries)
	slices.SortFunc(keys, func(a, b entry) int { return strings.Compare(a.key, b.key) })
	var pairs []string
	for _, e := range keys {
		for _, value := range e.words {
			pairs = append(pairs, quote(e.key)+"="+quote(value))
		}
	}
	return strings.Join([]string{text(first), strconv.FormatBool(found),
		render(replace(c.entries, c.key, []string{c.value})),
		render(replace(c.entries, c.key, append(slices.Clone(existing), c.value))),
		render(replace(c.entries, c.key, nil)), render(c.entries), text(strings.Join(pairs, "&"))}, "|")
}

func encode(c input) string {
	tag := 0
	if c.entries != nil {
		tag = 1
	}
	row := fmt.Appendf(nil, "%d %d", tag, len(c.entries))
	appendString := func(s string) {
		row = fmt.Appendf(row, " %d", len(s))
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
	}
	appendString(c.key)
	appendString(c.value)
	for _, e := range c.entries {
		appendString(e.key)
		if e.words == nil {
			row = append(row, " -1"...)
		} else {
			row = fmt.Appendf(row, " %d", len(e.words))
			for _, word := range e.words {
				appendString(word)
			}
		}
	}
	return string(row)
}

func main() {
	words := []string{"", "a", "A", "a\x00", "a b", "+", "/", "&", "=", "%", "%ff", "é", "\xff", "\x80", "�"}
	var edges []input
	for _, key := range words {
		for _, value := range words {
			for _, data := range [][]entry{nil, {}, {{key, nil}}, {{key, []string{}}},
				{{key, []string{"", value, "+"}}}, {{key + "\x00", []string{value}}, {key, []string{""}}, {"\xff" + key, []string{"é"}}}} {
				edges = append(edges, input{data, key, value})
			}
		}
	}
	check.CompareSamples("godiff-net-url-values", edges, sample, encode, native, expected)
}

func sample(rng *rand.Rand, index int) input {
	word := func() string {
		raw := make([]byte, rng.Intn(12))
		if index%2 == 0 {
			const alphabet = "aAbB+-& =\x00\xff"
			for i := range raw {
				raw[i] = alphabet[rng.Intn(len(alphabet))]
			}
		} else {
			rng.Read(raw)
		}
		return string(raw)
	}
	c := input{key: word(), value: word()}
	if index%13 == 0 {
		return c
	}
	c.entries = make([]entry, 0)
	seen := map[string]bool{}
	for i := range rng.Intn(9) {
		key := word()
		if i == 0 && index%3 != 0 {
			key = c.key
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		var values []string
		if index%5 != 0 {
			values = make([]string, rng.Intn(5))
			for j := range values {
				values[j] = word()
			}
		}
		c.entries = append(c.entries, entry{key, values})
	}
	return c
}
