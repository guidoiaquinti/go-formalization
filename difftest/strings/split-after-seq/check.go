package main

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s, sep        string
	first, second int
}

func main() {
	edges := []input{}
	add := func(s, sep string) {
		for _, budgets := range [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 1}, {len(s) + 2, 1}} {
			edges = append(edges, input{s, sep, budgets[0], budgets[1]})
		}
	}
	for _, s := range []string{"", "a", "aa", "aaa", "baaaabaa", "éé\U0001F680", "\xff\xc0\x80\xed\xa0\x80\xe2\x80", "a\x00a"} {
		for _, sep := range []string{"", "a", "aa", "aaa", "é", "\xa9", "\xff", "\x00", "missing"} {
			add(s, sep)
		}
	}
	for b := range 256 {
		raw := string([]byte{byte(b)})
		add(raw+"a"+raw, "")
		add(raw+"a"+raw, raw)
	}
	check.CompareSamples("godiff-strings-split-after-seq", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
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
	return input{string(raw), sep, rng.Intn(len(raw) + 2), rng.Intn(len(raw) + 2)}
}

func encodeInput(c input) string {
	row := strconv.AppendInt(nil, int64(len(c.sep)), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.first), 10)
	row = strconv.AppendInt(append(row, ' '), int64(c.second), 10)
	for _, s := range []string{c.sep, c.s} {
		for i := range len(s) {
			row = strconv.AppendUint(append(row, ' '), uint64(s[i]), 10)
		}
	}
	return string(row)
}

func describe(batches [][]string) string {
	encoded := make([]string, len(batches))
	for i, batch := range batches {
		fields := make([]string, len(batch))
		for j, field := range batch {
			bytes := make([]string, len(field))
			for k := range len(field) {
				bytes[k] = strconv.Itoa(int(field[k]))
			}
			fields[j] = "f" + strings.Join(bytes, ",")
		}
		encoded[i] = strings.Join(fields, ";")
	}
	return strings.Join(encoded, "|")
}

func actual(c input) string {
	seq := strings.SplitAfterSeq(c.s, c.sep)
	budgets := []int{c.first, c.second, len(c.s) + 2, 1}
	batches := make([][]string, len(budgets))
	for i, budget := range budgets {
		if budget > 0 {
			seq(func(part string) bool {
				batches[i] = append(batches[i], part)
				return len(batches[i]) < budget
			})
		}
	}
	return describe(batches)
}

func expected(c input) string {
	fragments := strings.SplitAfter(c.s, c.sep)
	budgets := []int{c.first, c.second, len(c.s) + 2, 1}
	batches := make([][]string, len(budgets))
	for i, budget := range budgets {
		yielded := min(budget, len(fragments))
		batches[i] = fragments[:yielded]
		advance := yielded
		if budget > 0 && (yielded == budget || (c.sep != "" && yielded == len(fragments))) {
			// A false result or the final nonempty-separator yield leaves
			// the most recent fragment at the current cursor.
			advance--
		}
		fragments = fragments[advance:]
	}
	return describe(batches)
}
