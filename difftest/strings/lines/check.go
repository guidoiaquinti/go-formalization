package main

import (
	"log"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type input struct {
	s             string
	first, second int
}

func main() {
	var edges []input
	values := []string{"", "\n", "\n\n", "a", "\na", "a\n", "a\r\nb\r\n", "a\nb\nc", "é\u2028\xff\x00\n", strings.Repeat("a\n", 128)}
	for b := range 256 {
		values = append(values, string([]byte{byte(b), '\n', byte(b)}))
	}
	for _, s := range values {
		for _, budgets := range [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 1}, {len(s), 1}, {len(s) + 1, 1}} {
			edges = append(edges, input{s, budgets[0], budgets[1]})
		}
	}
	check.CompareSamples("godiff-strings-lines", edges, sample, encodeInput, actual, expected)
}

func sample(rng *rand.Rand, _ int) input {
	raw := make([]byte, rng.Intn(129))
	rng.Read(raw)
	for i := range raw {
		if rng.Intn(5) == 0 {
			raw[i] = '\n'
		}
	}
	return input{string(raw), rng.Intn(len(raw) + 2), rng.Intn(len(raw) + 2)}
}

func encodeInput(c input) string {
	row := strconv.Itoa(c.first) + " " + strconv.Itoa(c.second)
	return row + string(appendBytes(nil, c.s, ' '))
}

func appendBytes(row []byte, s string, delimiter byte) []byte {
	for i := range len(s) {
		row = strconv.AppendUint(append(row, delimiter), uint64(s[i]), 10)
	}
	return row
}

func encodeBatches(batches [][]string) string {
	var encoded []string
	for _, batch := range batches {
		var lines []string
		for _, line := range batch {
			lines = append(lines, strings.TrimPrefix(string(appendBytes(nil, line, ',')), ","))
		}
		encoded = append(encoded, strings.Join(lines, ";"))
	}
	return strings.Join(encoded, "|")
}

func actual(c input) string {
	seq := strings.Lines(c.s)
	batches := make([][]string, 4)
	budgets := []int{c.first, c.second, len(c.s) + 1, 1}
	var got []string
	for i, budget := range budgets {
		if budget == 0 {
			continue
		}
		seq(func(line string) bool {
			if line == "" {
				log.Fatalf("strings.Lines(%q) yielded an empty line", c.s)
			}
			batches[i] = append(batches[i], line)
			return len(batches[i]) < budget
		})
		got = append(got, batches[i]...)
	}
	if strings.Join(got, "") != c.s {
		log.Fatalf("strings.Lines(%q) yielded %q across resumed calls", c.s, got)
	}
	return encodeBatches(batches)
}

func expected(c input) string {
	lines := slices.DeleteFunc(strings.SplitAfter(c.s, "\n"), func(line string) bool { return line == "" })
	first := min(c.first, len(lines))
	second := min(first+c.second, len(lines))
	return encodeBatches([][]string{lines[:first], lines[first:second], lines[second:], nil})
}
