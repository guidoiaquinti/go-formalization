package bytesplit

import (
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s, sep       string
	limit, spare int
}

type Splitter func([]byte, []byte, int) [][]byte

func Compare(executable string, split Splitter, unbounded, retainSeparator bool) {
	var edges []input
	add := func(s, sep string) {
		if unbounded {
			for _, spare := range []int{0, 3} {
				edges = append(edges, input{s, sep, -1, spare})
			}
			return
		}
		for _, n := range []int{math.MinInt, -4, -1, 0, 1, 2, 3, len(s), len(s) + 1, math.MaxInt} {
			edges = append(edges, input{s, sep, n, len(s)%3 + 3})
		}
	}
	for _, s := range []string{"", "a", "aaaa", "baaaabaa", "éé\U0001F680", "\xff\xc0\x80\xed\xa0\x80\xe2\x80", "a\x00a"} {
		for _, sep := range []string{"", "a", "aa", "aaa", "é", "\xa9", "\xff", "\x00", "missing"} {
			add(s, sep)
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		add(raw+"a"+raw, "")
		add(raw+"a"+raw, raw)
	}
	check.CompareSamples(executable, edges, func(rng *rand.Rand, i int) input {
		c := sample(rng, i)
		if unbounded {
			c.limit = -1
		}
		return c
	}, encodeInput, func(c input) string { return actual(c, split) }, func(c input) string { return expected(c, retainSeparator) })
}

func sample(rng *rand.Rand, i int) input {
	raw := make([]byte, rng.Intn(65))
	rng.Read(raw)
	start := rng.Intn(len(raw) + 1)
	end := start + rng.Intn(min(6, len(raw)-start)+1)
	sep := string(raw[start:end])
	if rng.Intn(3) == 0 {
		sep = strings.Repeat("a", rng.Intn(4))
	}
	if rng.Intn(3) == 0 {
		raw = []byte(sep + string(raw) + sep + sep)
	}
	limits := [...]int{-1, 0, 1, 2, rng.Intn(len(raw) + 3), math.MaxInt}
	return input{string(raw), sep, limits[i%len(limits)], rng.Intn(8)}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(c.limit), 10)
	row = strconv.AppendInt(append(row, ' '), int64(len(c.sep)), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.spare), 10)
	for _, s := range []string{c.sep, c.s} {
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
	}
	return string(row)
}

type fragment struct {
	text     string
	capacity int
}

func describe(parts []fragment) string {
	if parts == nil {
		return "nil"
	}
	fields := make([]string, len(parts))
	for i, part := range parts {
		data := make([]string, len(part.text))
		for j := range len(part.text) {
			data[j] = strconv.Itoa(int(part.text[j]))
		}
		fields[i] = "f" + strconv.Itoa(part.capacity) + ":" + strings.Join(data, ",")
	}
	return "slice" + strings.Join(fields, ";")
}

func actual(c input, split Splitter) string {
	src := make([]byte, len(c.s), len(c.s)+c.spare)
	copy(src, c.s)
	fragments := split(src, []byte(c.sep), c.limit)
	if fragments == nil {
		return "nil"
	}
	parts := make([]fragment, len(fragments))
	for i, p := range fragments {
		parts[i] = fragment{string(p), cap(p)}
	}
	return describe(parts)
}

func expected(c input, retainSeparator bool) string {
	if c.limit == 0 {
		return "nil"
	}
	split := strings.Split
	joiner := c.sep
	if retainSeparator {
		split = strings.SplitAfter
		joiner = ""
	}
	pieces := split(c.s, c.sep)
	if c.limit > 0 && c.limit < len(pieces) {
		pieces[c.limit-1] = strings.Join(pieces[c.limit-1:], joiner)
		pieces = pieces[:c.limit]
	}
	parts := make([]fragment, len(pieces))
	for i, p := range pieces {
		parts[i] = fragment{p, len(p)}
	}
	slots := len(c.s)
	if c.limit > 0 {
		slots = min(slots, c.limit)
	}
	if len(parts) > 0 && (c.sep != "" || len(parts) == slots) {
		parts[len(parts)-1].capacity += c.spare
	}
	return describe(parts)
}
