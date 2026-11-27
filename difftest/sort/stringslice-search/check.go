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

type input struct {
	xs  []string
	key string
}

func main() {
	edges := []input{{nil, ""}, {[]string{}, "x"}, {[]string{""}, ""}, {[]string{"a"}, ""}, {[]string{"a"}, "ab"}, {[]string{"a", "a", "ab"}, "a"}, {[]string{"a", "ab", "b"}, "aa"}, {[]string{"\x00", "\x80", "\xff"}, "\x80"}}

	check.CompareSamples("godiff-stringslice-search", edges, func(r *rand.Rand, i int) input {
		xs := make([]string, r.Intn(65))
		for j := range xs {
			if i%2 == 0 {
				xs[j] = string([]byte{byte(r.Intn(5))})
			} else {
				raw := make([]byte, r.Intn(9))
				r.Read(raw)
				xs[j] = string(raw)
			}
		}
		slices.Sort(xs)
		raw := make([]byte, r.Intn(9))
		r.Read(raw)
		key := string(raw)
		if i%2 == 0 {
			key = string([]byte{byte(r.Intn(7))})
		} else if i%4 == 1 && len(xs) != 0 {
			key = xs[r.Intn(len(xs))]
		}
		return input{xs, key}
	}, func(c input) string {
		var row strings.Builder
		row.WriteByte('0')
		for _, v := range append([]string{c.key}, c.xs...) {
			row.WriteString(" -1")
			for i := range len(v) {
				fmt.Fprintf(&row, " %d", v[i])
			}
		}
		return row.String()
	}, func(c input) string {
		return strconv.Itoa(sort.StringSlice(c.xs).Search(c.key))
	}, func(c input) string {
		for i, v := range c.xs {
			if bytes.Compare([]byte(v), []byte(c.key)) >= 0 {
				return strconv.Itoa(i)
			}
		}
		return strconv.Itoa(len(c.xs))
	})
}
