package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	edges := [][]string{nil, {}, {""}, {"a", "a"}, {"a", "ab"}, {"ab", "a"}, {"\x00", "\xff"}, {"\xff", "\x00"}, {"a", "b", ""}}

	check.CompareSamples("godiff-strings-are-sorted", edges, func(r *rand.Rand, i int) []string {
		xs := make([]string, r.Intn(33))
		for j := range xs {
			if i%3 == 0 {
				xs[j] = string([]byte{byte(r.Intn(5))})
			} else {
				raw := make([]byte, r.Intn(9))
				r.Read(raw)
				xs[j] = string(raw)
			}
		}
		if i%3 != 2 {
			slices.Sort(xs)
			if i%3 == 1 && len(xs) > 1 {
				j := r.Intn(len(xs) - 1)
				xs[j], xs[j+1] = xs[j+1], xs[j]
			}
		}
		return xs
	}, func(xs []string) string {
		var row strings.Builder
		row.WriteByte('0')
		for _, x := range xs {
			row.WriteString(" -1")
			for i := range len(x) {
				fmt.Fprintf(&row, " %d", x[i])
			}
		}
		return row.String()
	}, func(xs []string) string {
		return strconv.FormatBool(sort.StringsAreSorted(xs))
	}, func(xs []string) string {
		for i, x := range xs {
			for _, y := range xs[i+1:] {
				if bytes.Compare([]byte(x), []byte(y)) > 0 {
					return "false"
				}
			}
		}
		return "true"
	})
}
