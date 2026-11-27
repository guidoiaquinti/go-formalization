package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"slices"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

type namedSlice []int

func row(prefix string, values []int) string {
	var out strings.Builder
	out.WriteString(prefix)
	for _, v := range values {
		fmt.Fprintf(&out, " %d", v)
	}
	return out.String()
}

func checkOverflow() {
	defer func() {
		if got := recover(); got != "len out of range" {
			log.Fatalf("Concat overflow panic = %v, expected len out of range", got)
		}
		fmt.Println("Concat: zero-sized elements exercise native length overflow.")
	}()
	slices.Concat(make([]struct{}, math.MaxInt), []struct{}{{}})
}

func main() {
	checkOverflow()
	if got := check.LeanRows("godiff-slices-concat", 1, strings.NewReader("3 -1 1 2 -1 3 4\n"))[0]; got != "len out of range" {
		log.Fatalf("modeled overflow = %q, expected len out of range", got)
	}
	edges := [][]namedSlice{nil, {}, {nil}, {{}}, {nil, {}}, {{1}, nil, {2, 3}}, {{0, 255}, {}, {7}}}
	check.CompareSamples("godiff-slices-concat", edges, func(r *rand.Rand, _ int) []namedSlice {
		parts := make([]namedSlice, r.Intn(9))
		for i := range parts {
			if r.Intn(5) == 0 {
				continue
			}
			parts[i] = make(namedSlice, r.Intn(17))
			for j := range parts[i] {
				parts[i][j] = r.Intn(256)
			}
		}
		return parts
	}, func(parts []namedSlice) string {
		var out strings.Builder
		fmt.Fprint(&out, math.MaxInt)
		for _, part := range parts {
			out.WriteString(row(" -1", part))
		}
		return out.String()
	}, func(parts []namedSlice) string {
		got := slices.Concat(parts...)
		if got == nil {
			return "nil"
		}
		offset := 0
		for _, part := range parts {
			if len(part) > 0 && offset < len(got) {
				old := part[0]
				got[offset] ^= 1
				if part[0] != old {
					log.Fatal("Concat aliases the first element of an input part")
				}
				got[offset] ^= 1
			}
			offset += len(part)
		}
		return row("ok", got)
	}, func(parts []namedSlice) string {
		var values []int
		for _, part := range parts {
			for _, v := range part {
				values = append(values, v)
			}
		}
		if len(values) == 0 {
			return "nil"
		}
		return row("ok", values)
	})
}
